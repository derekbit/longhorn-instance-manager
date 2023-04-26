package disk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/golang/protobuf/ptypes/empty"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	issiutil "github.com/longhorn/go-iscsi-helper/util"
	spdknvme "github.com/longhorn/go-spdk-helper/pkg/nvme"
	spdkclient "github.com/longhorn/go-spdk-helper/pkg/spdk/client"
	spdktypes "github.com/longhorn/go-spdk-helper/pkg/spdk/types"
	spdkutil "github.com/longhorn/go-spdk-helper/pkg/types"

	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/meta"
	"github.com/longhorn/longhorn-instance-manager/pkg/types"
	"github.com/longhorn/longhorn-instance-manager/pkg/util/broadcaster"
)

const (
	defaultClusterSize = 4 * 1024 * 1024 // 4MB

	devPath = "/dev/longhorn"
)

type Server struct {
	shutdownCh    chan error
	HealthChecker HealthChecker

	spdkEnabled bool
	spdkClient  *spdkclient.Client
	sync.Mutex

	replicaBroadcaster *broadcaster.Broadcaster
	replicaBroadcastCh chan interface{}
	replicaUpdateCh    chan interface{}

	engineBroadcaster *broadcaster.Broadcaster
	engineBroadcastCh chan interface{}
	engineUpdateCh    chan interface{}
}

func NewServer(spdkEnabled bool, shutdownCh chan error) (*Server, error) {
	s := &Server{
		spdkEnabled:   spdkEnabled,
		shutdownCh:    shutdownCh,
		HealthChecker: &GRPCHealthChecker{},

		replicaBroadcaster: &broadcaster.Broadcaster{},
		replicaBroadcastCh: make(chan interface{}),
		replicaUpdateCh:    make(chan interface{}),

		engineBroadcaster: &broadcaster.Broadcaster{},
		engineBroadcastCh: make(chan interface{}),
		engineUpdateCh:    make(chan interface{}),
	}

	c, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := s.replicaBroadcaster.Subscribe(c, s.replicaBroadcastConnector); err != nil {
		return nil, err
	}

	if _, err := s.replicaBroadcaster.Subscribe(c, s.engineBroadcastConnector); err != nil {
		return nil, err
	}
	go s.startMonitoring()

	return s, nil
}

func (s *Server) getSpdkClient() (*spdkclient.Client, error) {
	if s.spdkEnabled {
		if s.spdkClient == nil {
			spdkClient, err := spdkclient.NewClient()
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create spdk client")
			}

			s.spdkClient = spdkClient
		}
	}

	return s.spdkClient, nil
}

func (s *Server) startMonitoring() {
	done := false
	for {
		select {
		case <-s.shutdownCh:
			logrus.Info("Disk Server is shutting down")
			done = true
			break
		case r := <-s.replicaUpdateCh:
			//pm.lock.RLock()
			// Modify response to indicate deletion.

			//pm.lock.RUnlock()
			s.replicaBroadcastCh <- interface{}(r)
		case e := <-s.engineUpdateCh:
			//pm.lock.RLock()
			// Modify response to indicate deletion.
			s.engineBroadcastCh <- interface{}(e)
		}
		if done {
			break
		}
	}
}

func (s *Server) bdevLvolGetLvstore(lvsName, uuid string) (*spdktypes.LvstoreInfo, error) {
	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get spdk client")
	}

	if spdkClient == nil {
		return nil, grpcstatus.Error(grpccodes.FailedPrecondition, "SPDK is not enabled")
	}

	lvstoreInfos, err := spdkClient.BdevLvolGetLvstore(lvsName, uuid)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get lvstore info")
	}

	if len(lvstoreInfos) != 1 {
		return nil, fmt.Errorf("number of lvstores with UUID %v is not 1", uuid)
	}

	return &lvstoreInfos[0], nil
}

func (s *Server) getBdevNameFromDiskPath(diskPath string) (string, error) {
	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get spdk client")
	}

	bdevInfos, err := spdkClient.BdevGetBdevs("", 30)
	if err != nil {
		return "", grpcstatus.Error(grpccodes.NotFound, errors.Wrapf(err, "failed to get bdevs").Error())
	}

	bdevName := ""
	for _, info := range bdevInfos {
		if info.DriverSpecific != nil &&
			info.DriverSpecific.Aio != nil &&
			info.DriverSpecific.Aio.FileName == diskPath {
			bdevName = info.Name
		}
	}

	if bdevName == "" {
		return "", grpcstatus.Error(grpccodes.NotFound, errors.Errorf("failed to find bdev with disk path %v", diskPath).Error())
	}

	return bdevName, nil
}

func (s *Server) DiskCreate(ctx context.Context, req *rpc.DiskCreateRequest) (*rpc.Disk, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"diskName": req.DiskName,
		"diskPath": req.DiskPath,
	})

	log.Info("Creating disk")

	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get spdk client")
	}

	diskPath := getDiskPath(req.DiskPath)

	bdevName, err := spdkClient.BdevAioCreate(diskPath, req.DiskName, 4096)
	if err != nil {
		resp, parseErr := parseErrorMessage(err.Error())
		if parseErr != nil || !strings.EqualFold(resp.Message, syscall.Errno(syscall.EEXIST).Error()) {
			log.WithError(err).Error("Failed to create AIO bdev")
			return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to create AIO bdev").Error())
		}
	}

	uuid, err := spdkClient.BdevLvolCreateLvstore(bdevName, req.DiskName, defaultClusterSize)
	if err != nil {
		log.WithError(err).Error("Failed to create lvstore")
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to create lvstore").Error())
	}

	lvstoreInfo, err := s.bdevLvolGetLvstore("", uuid)
	if err != nil {
		log.WithError(err).Errorf("Failed to get lvstore with name %v and UUID %v", req.DiskName, uuid)
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to get lvstore with name %v and UUID %v", req.DiskName, uuid).Error())
	}

	log.Info("Created disk")

	// A disk does not have a fsid, so we use the device number as the disk ID
	diskID, err := getDeviceNumber(diskPath)
	if err != nil {
		log.WithError(err).Error("Failed to get disk ID")
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to get disk ID").Error())
	}

	return &rpc.Disk{
		Id:          diskID,
		Uuid:        lvstoreInfo.UUID,
		Path:        req.DiskPath,
		Type:        DiskTypeBlock,
		TotalSize:   int64(lvstoreInfo.TotalDataClusters * lvstoreInfo.ClusterSize),
		FreeSize:    int64(lvstoreInfo.FreeClusters * lvstoreInfo.ClusterSize),
		TotalBlocks: int64(lvstoreInfo.TotalDataClusters * lvstoreInfo.ClusterSize / lvstoreInfo.BlockSize),
		FreeBlocks:  int64(lvstoreInfo.FreeClusters * lvstoreInfo.ClusterSize / lvstoreInfo.BlockSize),
		BlockSize:   int64(lvstoreInfo.BlockSize),
		ClusterSize: int64(lvstoreInfo.ClusterSize),
		Readonly:    false,
	}, nil
}

func (s *Server) DiskGet(ctx context.Context, req *rpc.DiskGetRequest) (*rpc.Disk, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"diskPath": req.DiskPath,
	})

	log.Info("Getting disk info")

	diskPath := getDiskPath(req.DiskPath)

	bdevName, err := s.getBdevNameFromDiskPath(diskPath)
	if err != nil {
		log.WithError(err).Errorf("Failed to get bdev name for disk path %v", diskPath)
		return nil, grpcstatus.Error(grpccodes.NotFound, errors.Wrapf(err, "failed to get bdev name for disk path %v", diskPath).Error())
	}

	lvstoreName := bdevName
	lvstoreInfo, err := s.bdevLvolGetLvstore(lvstoreName, "")
	if err != nil {
		log.WithError(err).Errorf("Failed to get %v lvstore", lvstoreName)
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to get %v lvstore", lvstoreName).Error())
	}

	log.Info("Got disk info")

	diskID, err := getDeviceNumber(diskPath)
	if err != nil {
		log.WithError(err).Error("Failed to get disk ID")
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to get disk ID").Error())
	}

	return &rpc.Disk{
		Id:          diskID,
		Uuid:        lvstoreInfo.UUID,
		Path:        req.DiskPath,
		Type:        DiskTypeBlock,
		TotalSize:   int64(lvstoreInfo.TotalDataClusters * lvstoreInfo.ClusterSize),
		FreeSize:    int64(lvstoreInfo.FreeClusters * lvstoreInfo.ClusterSize),
		TotalBlocks: int64(lvstoreInfo.TotalDataClusters * lvstoreInfo.ClusterSize / lvstoreInfo.BlockSize),
		FreeBlocks:  int64(lvstoreInfo.FreeClusters * lvstoreInfo.ClusterSize / lvstoreInfo.BlockSize),
		BlockSize:   int64(lvstoreInfo.BlockSize),
		ClusterSize: int64(lvstoreInfo.ClusterSize),
		Readonly:    false,
	}, nil
}

func (s *Server) ReplicaCreate(ctx context.Context, req *rpc.ReplicaCreateRequest) (*rpc.Replica, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name":        req.Name,
		"lvstoreUUID": req.LvstoreUuid,
		"size":        req.Size,
		"address":     req.Address,
	})

	log.Info("Creating replica")

	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get spdk client")
	}

	lvstoreUUID := req.LvstoreUuid
	lvstoreInfos, err := spdkClient.BdevLvolGetLvstore("", lvstoreUUID)
	if err != nil {
		log.WithError(err).Error("Failed to get lvstore info")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	if len(lvstoreInfos) != 1 {
		log.WithError(err).Error("Number of lvstore info is not 1")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	lvstoreInfo := lvstoreInfos[0]

	sizeInMib := uint64(req.Size / 1024 / 1024)

	uuid, err := spdkClient.BdevLvolCreate(lvstoreInfo.Name, req.Name, "", sizeInMib, spdktypes.BdevLvolClearMethodUnmap, true)
	if err != nil {
		log.WithError(err).Error("Failed to create lvol")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	////
	/*
		nqn := spdkutil.GetNQN(req.Name)
		err = spdkClient.StartExposeBdev(nqn, lvstoreInfo.Name+"/"+req.Name, req.Address, "4420")
		if err != nil {
			log.WithError(err).Error("Failed to start exposing bdev")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}
	*/
	////

	log.Info("Created replica")

	lvolInfos, err := spdkClient.BdevLvolGet(uuid, 30)
	if err != nil {
		log.WithError(err).Error("Failed to get lvol info")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	if len(lvolInfos) != 1 {
		log.WithError(err).Error("Number of lvol info is not 1")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	lvolInfo := lvolInfos[0]

	name := getNameFromAlias(lvolInfo.Aliases)
	if name == "" {
		return nil, grpcstatus.Error(grpccodes.NotFound, "failed to get name from alias")
	}

	r := &rpc.Replica{
		Name:        name,
		Uuid:        lvolInfo.UUID,
		BdevName:    lvolInfo.DriverSpecific.Lvol.BaseBdev,
		LvstoreUuid: lvolInfo.DriverSpecific.Lvol.LvolStoreUUID,

		TotalSize:   int64(lvolInfo.NumBlocks * uint64(lvolInfo.BlockSize)),
		TotalBlocks: int64(lvolInfo.NumBlocks),
		BlockSize:   int64(lvolInfo.BlockSize),

		Port: 4420,

		State: types.ProcessStateRunning,
	}

	s.replicaUpdateCh <- interface{}(r)

	return r, nil
}

func (s *Server) ReplicaDelete(ctx context.Context, req *rpc.ReplicaDeleteRequest) (*empty.Empty, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name":        req.Name,
		"lvstoreUUID": req.LvstoreUuid,
	})

	log.Info("Deleting replica")

	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get spdk client")
	}

	lvstoreInfos, err := spdkClient.BdevLvolGetLvstore("", req.LvstoreUuid)
	if err != nil {
		log.WithError(err).Error("Failed to list lvstore info")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	lvstoreName := getLvolNameFromUUID(lvstoreInfos, req.LvstoreUuid)
	if lvstoreName == "" {
		log.WithError(err).Error("Failed to get lvstore name")
		return nil, grpcstatus.Error(grpccodes.NotFound, err.Error())
	}

	aliasName := lvstoreName + "/" + req.Name
	deleted, err := spdkClient.BdevLvolDelete(aliasName)
	if err != nil {
		log.WithError(err).Error("Failed to delete lvol")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	if !deleted {
		log.WithError(err).Error("Failed to delete lvol")
		return nil, grpcstatus.Errorf(grpccodes.Internal, "failed to delete lvol %v", aliasName)
	}

	log.Info("Deleted replica")

	return &empty.Empty{}, nil
}

func (s *Server) ReplicaGet(ctx context.Context, req *rpc.ReplicaGetRequest) (*rpc.Replica, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name":        req.Name,
		"lvstoreUUID": req.LvstoreUuid,
	})

	log.Info("Getting replica info")

	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get spdk client")
	}

	lvstoreInfos, err := spdkClient.BdevLvolGetLvstore("", req.LvstoreUuid)
	if err != nil {
		log.WithError(err).Error("Failed to list lvstore info")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	lvstoreName := getLvolNameFromUUID(lvstoreInfos, req.LvstoreUuid)
	if lvstoreName == "" {
		log.WithError(err).Error("Failed to get lvstore name")
		return nil, grpcstatus.Error(grpccodes.NotFound, err.Error())
	}

	aliasName := lvstoreName + "/" + req.Name
	lvolInfos, err := spdkClient.BdevLvolGet(aliasName, 30)
	if err != nil {
		log.WithError(err).Error("Failed to get lvol info")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	if len(lvolInfos) != 1 {
		log.WithError(err).Error("Number of lvol info is not 1")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	log.Info("Got replica info")

	lvolInfo := lvolInfos[0]

	name := getNameFromAlias(lvolInfo.Aliases)
	if name == "" {
		return nil, grpcstatus.Error(grpccodes.NotFound, "failed to get name from alias")
	}

	return &rpc.Replica{
		Name:        name,
		Uuid:        lvolInfo.UUID,
		BdevName:    lvolInfo.DriverSpecific.Lvol.BaseBdev,
		LvstoreUuid: lvolInfo.DriverSpecific.Lvol.BaseBdev,

		TotalSize:   int64(lvolInfo.NumBlocks * uint64(lvolInfo.BlockSize)),
		TotalBlocks: int64(lvolInfo.NumBlocks),
		BlockSize:   int64(lvolInfo.BlockSize),

		Port: 4420,

		State: types.ProcessStateRunning,
	}, nil
}

func (s *Server) ReplicaList(ctx context.Context, req *empty.Empty) (*rpc.ReplicaListResponse, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{})

	log.Info("Getting replica list")

	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get spdk client")
	}

	lvolInfos, err := spdkClient.BdevLvolGet("", 30)
	if err != nil {
		log.WithError(err).Error("Failed to get lvol infos")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	log.Info("Got replica list")

	resp := &rpc.ReplicaListResponse{
		Replicas: map[string]*rpc.Replica{},
	}
	for _, info := range lvolInfos {
		name := getNameFromAlias(info.Aliases)
		if name == "" {
			return nil, grpcstatus.Error(grpccodes.Internal, "Failed to get name from alias")
		}

		resp.Replicas[info.Name] = &rpc.Replica{
			Name:        name,
			Uuid:        info.UUID,
			BdevName:    info.DriverSpecific.Lvol.BaseBdev,
			LvstoreUuid: info.DriverSpecific.Lvol.BaseBdev,

			TotalSize:   int64(info.NumBlocks * uint64(info.BlockSize)),
			TotalBlocks: int64(info.NumBlocks),
			BlockSize:   int64(info.BlockSize),

			Port: 4420,

			State: types.ProcessStateRunning,
		}
	}

	return resp, nil
}

func getDeviceNumber(filename string) (string, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return "", err
	}
	dev := fi.Sys().(*syscall.Stat_t).Dev
	major := int(dev >> 8)
	minor := int(dev & 0xff)
	return fmt.Sprintf("%d-%d", major, minor), nil
}

func getNameFromAlias(alias []string) string {
	if len(alias) != 1 {
		return ""
	}

	splitName := strings.Split(alias[0], "/")
	return splitName[1]
}

func getDiskPath(path string) string {
	return filepath.Join("/host", path)
}

func getLvolNameFromUUID(lvstoreInfos []spdktypes.LvstoreInfo, lvstoreUUID string) string {
	lvstoreName := ""
	for _, info := range lvstoreInfos {
		if info.UUID == lvstoreUUID {
			lvstoreName = info.Name
			break
		}
	}
	return lvstoreName
}

func (s *Server) VersionGet(ctx context.Context, req *empty.Empty) (*rpc.VersionResponse, error) {
	v := meta.GetVersion()
	return &rpc.VersionResponse{
		Version:   v.Version,
		GitCommit: v.GitCommit,
		BuildDate: v.BuildDate,

		InstanceManagerAPIVersion:    int64(v.InstanceManagerAPIVersion),
		InstanceManagerAPIMinVersion: int64(v.InstanceManagerAPIMinVersion),

		InstanceManagerProxyAPIVersion:    int64(v.InstanceManagerProxyAPIVersion),
		InstanceManagerProxyAPIMinVersion: int64(v.InstanceManagerProxyAPIMinVersion),

		InstanceManagerDiskServiceAPIVersion:    int64(v.InstanceManagerDiskServiceAPIVersion),
		InstanceManagerDiskServiceAPIMinVersion: int64(v.InstanceManagerDiskServiceAPIMinVersion),
	}, nil
}

func (s *Server) EngineCreate(ctx context.Context, req *rpc.EngineCreateRequest) (*rpc.Engine, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name":              req.Name,
		"volumeName":        req.VolumeName,
		"address":           req.Address,
		"replicaAddressMap": req.ReplicaAddressMap,
		"frontend":          req.Frontend,
	})

	log.Info("Creating engine")

	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get spdk client")
	}

	bdevs := []string{}

	lvolInfos, err := spdkClient.BdevLvolGet("", 30)
	if err != nil {
		log.WithError(err).Error("Failed to get lvol infos")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}
	for _, info := range lvolInfos {
		bdevs = append(bdevs, info.Aliases[0])
	}

	localIP, err := issiutil.GetIPToHost()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get local IP")
	}
	localPort := "4422"

	/*
		for name, addr := range req.ReplicaAddressMap {
			nqn := spdkutil.GetNQN(name)
			addressComponents := strings.Split(addr, ":")
			if len(addressComponents) != 2 {
				return nil, grpcstatus.Errorf(grpccodes.InvalidArgument, "Invalid address %v", addr)
			}

			bdevNameList, err := spdkCli.BdevNvmeAttachController(name, nqn, addressComponents[0], addressComponents[1], spdktypes.NvmeTransportTypeTCP, spdktypes.NvmeAddressFamilyIPv4)
			if err != nil {
				log.WithError(err).Error("Failed to attach NVMe controller")
				return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
			}
			bdevs = append(bdevs, bdevNameList...)
		}
	*/

	created, err := spdkClient.BdevRaidCreate(req.Name, spdktypes.BdevRaidLevelRaid1, 0, bdevs)
	if err != nil {
		log.WithError(err).Error("Failed to create RAID")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	if created {
		log.Info("Created RAID")
	}

	raidNQN := spdkutil.GetNQN(req.Name)
	err = spdkClient.StartExposeBdev(raidNQN, req.Name, localIP, localPort)
	if err != nil {
		log.WithError(err).Error("Failed to start exposing bdev")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	nvmeCli, err := spdknvme.NewInitiator(raidNQN, localIP, localPort)
	if err != nil {
		log.WithError(err).Error("Failed to create NVMe initiator")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	err = nvmeCli.StartInitiator()
	if err != nil {
		log.WithError(err).Error("Failed to start NVMe initiator")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	time.Sleep(5 * time.Second)

	err = createLonghornDevice(filepath.Join("/dev", nvmeCli.ControllerName)+"n1", req.VolumeName)
	if err != nil {
		log.WithError(err).Error("Failed to create longhorn device")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	listenerList, err := spdkClient.NvmfSubsystemGetListeners(raidNQN, "")
	if err != nil {
		log.WithError(err).Error("Failed to get NVMe subsystem listeners")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	if len(listenerList) != 1 {
		log.WithError(err).Error("Invalid NVMe subsystem listeners")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	port, err := strconv.Atoi(listenerList[0].Address.Trsvcid)
	if err != nil {
		log.WithError(err).Error("Failed to convert port")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}
	e := &rpc.Engine{
		Name:              req.Name,
		Address:           req.Address,
		ReplicaAddressMap: req.ReplicaAddressMap,
		Frontend:          req.Frontend,
		Endpoint:          localIP,
		Ip:                listenerList[0].Address.Traddr,
		Port:              int32(port),
	}

	s.engineUpdateCh <- interface{}(e)

	return e, nil
}

func (s *Server) EngineDelete(ctx context.Context, req *rpc.EngineDeleteRequest) (*empty.Empty, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name": req.Name,
	})

	log.Info("Deleting engine")
	log.Info("Deleted engine")

	return nil, grpcstatus.Error(grpccodes.Unimplemented, "")
}

func (s *Server) EngineGet(ctx context.Context, req *rpc.EngineGetRequest) (*rpc.Engine, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name": req.Name,
	})

	log.Info("Getting engine info")

	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get spdk client")
	}

	bdevRaidInfos, err := spdkClient.BdevRaidGetBdevs(spdktypes.BdevRaidCategoryAll)
	if err != nil {
		log.WithError(err).Error("Failed to get bdev raid infos")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	for _, info := range bdevRaidInfos {
		if info.Name != req.Name {
			continue
		}

		raidNQN := spdkutil.GetNQN(req.Name)
		listenerList, err := spdkClient.NvmfSubsystemGetListeners(raidNQN, "")
		if err != nil {
			log.WithError(err).Error("Failed to get NVMe subsystem listeners")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}

		if len(listenerList) != 1 {
			log.WithError(err).Error("Invalid NVMe subsystem listeners")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}

		port, err := strconv.Atoi(listenerList[0].Address.Trsvcid)
		if err != nil {
			log.WithError(err).Error("Failed to convert port")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}

		engine := &rpc.Engine{
			Name:              info.Name,
			Uuid:              "",
			Size:              0,
			Address:           "",
			ReplicaAddressMap: map[string]string{},
			ReplicaModeMap:    map[string]rpc.ReplicaMode{},
			Endpoint:          "",
			Frontend:          "",
			Ip:                listenerList[0].Address.Traddr,
			Port:              int32(port),
		}

		if info.State == spdktypes.BdevRaidCategoryOffline {
			engine.FrontendState = "down"
		} else {
			engine.FrontendState = "up"
		}

		return engine, nil
	}

	return nil, grpcstatus.Error(grpccodes.NotFound, "")
}

func (s *Server) EngineList(ctx context.Context, req *empty.Empty) (*rpc.EngineListResponse, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{})

	log.Info("Getting engine list")

	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get spdk client")
	}

	resp := &rpc.EngineListResponse{
		Engines: map[string]*rpc.Engine{},
	}
	bdevRaidInfos, err := spdkClient.BdevRaidGetBdevs(spdktypes.BdevRaidCategoryAll)
	if err != nil {
		log.WithError(err).Error("Failed to get bdev raid infos")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	for _, info := range bdevRaidInfos {
		raidNQN := spdkutil.GetNQN(info.Name)
		listenerList, err := spdkClient.NvmfSubsystemGetListeners(raidNQN, "")
		if err != nil {
			log.WithError(err).Error("Failed to get NVMe subsystem listeners")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}

		if len(listenerList) != 1 {
			log.WithError(err).Error("Invalid NVMe subsystem listeners")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}

		port, err := strconv.Atoi(listenerList[0].Address.Trsvcid)
		if err != nil {
			log.WithError(err).Error("Failed to convert port")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}

		engine := &rpc.Engine{
			Name:              info.Name,
			Uuid:              "",
			Size:              0,
			Address:           "",
			ReplicaAddressMap: map[string]string{},
			ReplicaModeMap:    map[string]rpc.ReplicaMode{},
			Endpoint:          "",
			Frontend:          "",
			Ip:                listenerList[0].Address.Traddr,
			Port:              int32(port),
		}

		if info.State == spdktypes.BdevRaidCategoryOffline {
			engine.FrontendState = "down"
		} else {
			engine.FrontendState = "up"
		}

		resp.Engines[info.Name] = engine
	}

	log.Info("Got engine list")

	return resp, nil
}

func createLonghornDevice(devicePath, name string) error {
	logrus.Infof("Creating longhorn device: devicePath=%s, name=%s", devicePath, name)

	if _, err := os.Stat(devPath); os.IsNotExist(err) {
		if err := os.MkdirAll(devPath, 0755); err != nil {
			logrus.Fatalf("device %v: Cannot create directory %v", name, devPath)
		}
	}

	// Get the major and minor numbers of the NVMe device.
	major, minor, err := getDeviceNumbers(devicePath)
	if err != nil {
		return err
	}

	longhornDevPath := filepath.Join(devPath, name)

	return duplicateDevice(major, minor, longhornDevPath)
}

func getDeviceNumbers(devicePath string) (major, minor uint32, err error) {
	fileInfo, err := os.Stat(devicePath)
	if err != nil {
		return 0, 0, err
	}

	statT := fileInfo.Sys().(*syscall.Stat_t)
	major = uint32(int(statT.Rdev) >> 8)
	minor = uint32(int(statT.Rdev) & 0xFF)
	return major, minor, nil
}

func duplicateDevice(major, minor uint32, dest string) error {
	if err := mknod(dest, major, minor); err != nil {
		return fmt.Errorf("couldn't create device %s: %w", dest, err)
	}
	if err := os.Chmod(dest, 0660); err != nil {
		return fmt.Errorf("couldn't change permission of the device %s: %w", dest, err)
	}
	return nil
}

func mknod(device string, major, minor uint32) error {
	var fileMode os.FileMode = 0660
	fileMode |= unix.S_IFBLK
	dev := int(unix.Mkdev(uint32(major), uint32(minor)))

	logrus.Infof("Creating device %s %d:%d", device, major, minor)
	return unix.Mknod(device, uint32(fileMode), dev)
}

func (s *Server) replicaBroadcastConnector() (chan interface{}, error) {
	return s.replicaBroadcastCh, nil
}

func (s *Server) engineBroadcastConnector() (chan interface{}, error) {
	return s.engineBroadcastCh, nil
}

func (s *Server) Subscribe(instanceType string) (<-chan interface{}, error) {
	switch instanceType {
	case types.InstanceTypeEngine:
		return s.engineBroadcaster.Subscribe(context.TODO(), s.engineBroadcastConnector)
	case types.InstanceTypeReplica:
		return s.replicaBroadcaster.Subscribe(context.TODO(), s.replicaBroadcastConnector)
	default:
		return nil, fmt.Errorf("unknown instance type %v", instanceType)
	}
}

func (s *Server) ReplicaWatch(req *empty.Empty, srv rpc.DiskService_ReplicaWatchServer) error {
	responseChan, err := s.Subscribe(types.InstanceTypeReplica)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			logrus.WithError(err).Error("Disk service update watch errored out")
		} else {
			logrus.Info("Disk service update watch ended successfully")
		}
	}()
	logrus.Info("Started new disk service update watch")

	for resp := range responseChan {
		r, ok := resp.(*rpc.Replica)
		if !ok {
			return fmt.Errorf("BUG: cannot get Replica from channel")
		}
		if err := srv.Send(r); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) EngineWatch(req *empty.Empty, srv rpc.DiskService_EngineWatchServer) error {
	responseChan, err := s.Subscribe(types.InstanceTypeEngine)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			logrus.WithError(err).Error("Disk service update watch errored out")
		} else {
			logrus.Info("Disk service update watch ended successfully")
		}
	}()
	logrus.Info("Started new disk service update watch")

	for resp := range responseChan {
		e, ok := resp.(*rpc.Engine)
		if !ok {
			return fmt.Errorf("BUG: cannot get Engine from channel")
		}
		if err := srv.Send(e); err != nil {
			return err
		}
	}

	return nil
}
