package spdk

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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
	"github.com/longhorn/longhorn-instance-manager/pkg/util"
	"github.com/longhorn/longhorn-instance-manager/pkg/util/broadcaster"
)

const (
	devPath = "/dev/longhorn"

	defaultEnginePortCount  = 1
	defaultReplicaPortCount = 1
)

type Server struct {
	sync.Mutex

	shutdownCh    chan error
	HealthChecker HealthChecker

	spdkEnabled bool

	portRangeMin   int32
	portRangeMax   int32
	availablePorts *util.Bitmap

	replicaBroadcaster *broadcaster.Broadcaster
	replicaBroadcastCh chan interface{}
	replicaUpdateCh    chan interface{}

	engineBroadcaster *broadcaster.Broadcaster
	engineBroadcastCh chan interface{}
	engineUpdateCh    chan interface{}
}

func NewServer(portRange string, spdkEnabled bool, shutdownCh chan error) (*Server, error) {
	start, end, err := parsePortRange(portRange)
	if err != nil {
		return nil, err
	}

	s := &Server{
		shutdownCh:    shutdownCh,
		HealthChecker: &GRPCHealthChecker{},

		spdkEnabled: spdkEnabled,

		portRangeMin:   start,
		portRangeMax:   end,
		availablePorts: util.NewBitmap(start, end),

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

func (s *Server) newSPDKClient() (*spdkclient.Client, error) {
	if s.spdkEnabled {
		return spdkclient.NewClient()
	}
	return nil, fmt.Errorf("SPDK is not enabled")
}

func (s *Server) startMonitoring() {
	done := false
	for {
		select {
		case <-s.shutdownCh:
			logrus.Info("SPDK Server is shutting down")
			done = true
		case r := <-s.replicaUpdateCh:
			s.replicaBroadcastCh <- interface{}(r)
		case e := <-s.engineUpdateCh:
			s.engineBroadcastCh <- interface{}(e)
		}
		if done {
			break
		}
	}
}

func (s *Server) ReplicaCreate(ctx context.Context, req *rpc.ReplicaCreateRequest) (*rpc.Replica, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name":           req.Name,
		"lvsUUID":        req.LvsUuid,
		"specSize":       req.SpecSize,
		"exposeRequired": req.ExposeRequired,
	})

	log.Info("Creating replica")

	spdkClient, err := s.newSPDKClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get SPDK client").Error())
	}
	defer spdkClient.Close()

	lvstoreInfo, err := bdevLvolGetLvstore(spdkClient, log, req.LvsUuid)
	if err != nil {
		log.WithError(err).Error("Failed to get lvstore")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	sizeInMib := uint64(req.SpecSize / 1024 / 1024)

	uuid, err := spdkClient.BdevLvolCreate(lvstoreInfo.Name, req.Name, "", sizeInMib, spdktypes.BdevLvolClearMethodUnmap, true)
	if err != nil {
		log.WithError(err).Error("Failed to create lvol")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	localPort := int32(0)
	localIP, err := issiutil.GetIPToHost()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get local IP")
	}

	if req.ExposeRequired {
		portStart, _, err := s.allocateInstancePorts(req.Name, defaultReplicaPortCount)
		if err != nil {
			log.WithError(err).Error("Failed to allocate ports")
			return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to allocate ports").Error())
		}
		localPort = portStart

		nqn := spdkutil.GetNQN(req.Name)
		err = spdkClient.StartExposeBdev(nqn, lvstoreInfo.Name+"/"+req.Name, localIP, strconv.Itoa(int(localPort)))
		if err != nil {
			log.WithError(err).Error("Failed to start exposing bdev")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}
	}

	log.Info("Created replica")

	lvolInfo, err := bdevLvolGet(spdkClient, log, uuid)
	if err != nil {
		log.WithError(err).Error("Failed to get lvol")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	r := &rpc.Replica{
		Name:       req.Name,
		Uuid:       lvolInfo.UUID,
		LvsName:    lvstoreInfo.Name,
		LvsUuid:    lvstoreInfo.UUID,
		SpecSize:   uint64(lvolInfo.NumBlocks * uint64(lvolInfo.BlockSize)),
		ActualSize: uint64(lvolInfo.NumBlocks * uint64(lvolInfo.BlockSize)),
		Snapshots:  map[string]*rpc.Lvol{},
		Rebuilding: false,
		Port:       int32(localPort),
	}

	s.replicaUpdateCh <- interface{}(r)

	return r, nil
}

func (s *Server) ReplicaDelete(ctx context.Context, req *rpc.ReplicaDeleteRequest) (*empty.Empty, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name": req.Name,
	})

	log.Info("Deleting replica")

	spdkClient, err := s.newSPDKClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get SPDK client").Error())
	}
	defer spdkClient.Close()

	deleted, err := spdkClient.BdevLvolDelete(req.Name)
	if err != nil {
		log.WithError(err).Error("Failed to delete lvol")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	if !deleted {
		log.WithError(err).Error("Failed to delete lvol")
		return nil, grpcstatus.Errorf(grpccodes.Internal, "failed to delete lvol %v", req.Name)
	}

	log.Info("Deleted replica")

	return &empty.Empty{}, nil
}

func (s *Server) ReplicaGet(ctx context.Context, req *rpc.ReplicaGetRequest) (*rpc.Replica, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name": req.Name,
	})

	log.Info("Getting replica info")

	spdkClient, err := s.newSPDKClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get SPDK client").Error())
	}
	defer spdkClient.Close()

	lvol, err := bdevLvolGet(spdkClient, log, req.Name)
	if err != nil {
		log.WithError(err).Error("Failed to get lvol")
		return nil, grpcstatus.Errorf(grpccodes.Internal, errors.Wrapf(err, "failed to get lvol %v", req.Name).Error())
	}

	log.Info("Got replica info")

	localPort := 0
	nqn := spdkutil.GetNQN(req.Name)
	listeners, err := spdkClient.NvmfSubsystemGetListeners(nqn, "")
	if err != nil {
		resp, parseErr := parseErrorMessage(err.Error())
		if parseErr != nil || !strings.Contains(resp.Message, "Invalid parameters") {
			log.WithError(err).Error("Failed to get listener")
			return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get listener").Error())
		}
	} else {
		listener := listeners[0]
		localPort, err = strconv.Atoi(listener.Address.Trsvcid)
		if err != nil {
			log.WithError(err).Error("Failed to parse port")
			return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to parse port").Error())
		}
	}

	return &rpc.Replica{
		Name:       req.Name,
		Uuid:       lvol.UUID,
		LvsName:    lvol.Name,
		LvsUuid:    lvol.UUID,
		SpecSize:   uint64(lvol.NumBlocks * uint64(lvol.BlockSize)),
		ActualSize: uint64(lvol.NumBlocks * uint64(lvol.BlockSize)),
		Snapshots:  map[string]*rpc.Lvol{},
		Rebuilding: false,
		Port:       int32(localPort),
	}, nil
}

func (s *Server) ReplicaList(ctx context.Context, req *empty.Empty) (*rpc.ReplicaListResponse, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{})

	log.Info("Getting replica list")

	spdkClient, err := s.newSPDKClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get SPDK client").Error())
	}
	defer spdkClient.Close()

	lvolInfos, err := spdkClient.BdevLvolGet("", 3000)
	if err != nil {
		log.WithError(err).Error("Failed to get lvol infos")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	resp := &rpc.ReplicaListResponse{
		Replicas: map[string]*rpc.Replica{},
	}
	for _, info := range lvolInfos {
		name := getNameFromAlias(info.Aliases)
		if name == "" {
			return nil, grpcstatus.Error(grpccodes.Internal, "failed to get name from alias")
		}

		localPort := 0
		nqn := spdkutil.GetNQN(name)
		listeners, err := spdkClient.NvmfSubsystemGetListeners(nqn, "")
		if err != nil {
			resp, parseErr := parseErrorMessage(err.Error())
			if parseErr != nil || !strings.Contains(resp.Message, "Invalid parameters") {
				log.WithError(err).Error("Failed to get listener")
				return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get listener").Error())
			}
		} else {
			listener := listeners[0]
			localPort, err = strconv.Atoi(listener.Address.Trsvcid)
			if err != nil {
				log.WithError(err).Error("Failed to parse port")
				return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to parse port").Error())
			}
		}

		resp.Replicas[info.Name] = &rpc.Replica{
			Name:       name,
			Uuid:       info.UUID,
			LvsName:    info.Name,
			LvsUuid:    info.UUID,
			SpecSize:   uint64(info.NumBlocks * uint64(info.BlockSize)),
			ActualSize: uint64(info.NumBlocks * uint64(info.BlockSize)),
			Snapshots:  map[string]*rpc.Lvol{},
			Rebuilding: false,
			Port:       int32(localPort),
		}
	}

	log.Info("Got replica list")

	return resp, nil
}

func (s *Server) ReplicaWatch(req *empty.Empty, srv rpc.SPDKService_ReplicaWatchServer) error {
	responseChan, err := s.Subscribe(types.InstanceTypeReplica)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			logrus.WithError(err).Error("SPDK service update watch errored out")
		} else {
			logrus.Info("SPDK service update watch ended successfully")
		}
	}()
	logrus.Info("Started new SPDK service update watch")

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

func (s *Server) EngineCreate(ctx context.Context, req *rpc.EngineCreateRequest) (*rpc.Engine, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name":              req.Name,
		"specSize":          req.SpecSize,
		"replicaAddressMap": req.ReplicaAddressMap,
	})

	log.Info("Creating engine")

	spdkClient, err := s.newSPDKClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get SPDK client").Error())
	}
	defer spdkClient.Close()

	localbdevs := map[string]struct{}{}
	bdevs := []string{}

	// Get all local bdevs
	lvolInfos, err := spdkClient.BdevLvolGet("", 3000)
	if err != nil {
		log.WithError(err).Error("Failed to get lvol infos")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}
	for _, info := range lvolInfos {
		replicaName := getNameFromAlias(info.Aliases)
		if _, ok := req.ReplicaAddressMap[replicaName]; ok {
			bdevs = append(bdevs, info.Aliases[0])
			localbdevs[replicaName] = struct{}{}
		}
	}

	// Attach remote NVMe controllers
	for replicaName, replicaAddr := range req.ReplicaAddressMap {
		if _, ok := localbdevs[replicaName]; ok {
			continue
		}

		log.Infof("Attaching NVMe controller for replica %v at %v", replicaName, replicaAddr)

		nqn := spdkutil.GetNQN(replicaName)
		ip, port, err := splitHostPort(replicaAddr)
		if err != nil {
			log.WithError(err).Error("Failed to split host port")
			return nil, grpcstatus.Error(grpccodes.InvalidArgument, err.Error())
		}

		bdevNameList, err := spdkClient.BdevNvmeAttachController(replicaName, nqn, ip, port, spdktypes.NvmeTransportTypeTCP, spdktypes.NvmeAddressFamilyIPv4)
		if err != nil {
			log.WithError(err).Error("Failed to attach NVMe controller")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}
		bdevs = append(bdevs, bdevNameList...)
	}

	log.Infof("Creating RAID with bdevs %v", bdevs)

	localIP, err := issiutil.GetIPToHost()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	portStart, _, err := s.allocateInstancePorts(req.Name, defaultReplicaPortCount)
	if err != nil {
		log.WithError(err).Error("Failed to allocate ports")
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to allocate ports").Error())
	}
	localPort := portStart

	created, err := spdkClient.BdevRaidCreate(req.Name, spdktypes.BdevRaidLevelRaid1, 0, bdevs)
	if err != nil {
		log.WithError(err).Error("Failed to create RAID")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	if created {
		log.Infof("Created RAID with bdevs %v", bdevs)
	}

	raidNQN := spdkutil.GetNQN(req.Name)
	err = spdkClient.StartExposeBdev(raidNQN, req.Name, localIP, strconv.Itoa(int(localPort)))
	if err != nil {
		log.WithError(err).Error("Failed to start exposing bdev")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	initiator, err := spdknvme.NewInitiator(raidNQN, localIP, strconv.Itoa(int(localPort)))
	if err != nil {
		log.WithError(err).Error("Failed to create NVMe initiator")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	err = initiator.StartInitiator()
	if err != nil {
		log.WithError(err).Error("Failed to start NVMe initiator")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	err = waitForDeviceReady(filepath.Join("/dev", initiator.ControllerName)+"n1", 10)
	if err != nil {
		log.WithError(err).Error("Failed to wait for longhorn device ready")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	volumeName := getVolumeName(req.Name)
	err = createLonghornDevice(filepath.Join("/dev", initiator.ControllerName)+"n1", volumeName)
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

	engine, err := getEngine(spdkClient, req.Name, log)
	if err != nil {
		log.WithError(err).Error("Failed to get engine")
		return nil, err
	}

	s.engineUpdateCh <- interface{}(engine)
	return engine, nil
}

func (s *Server) EngineDelete(ctx context.Context, req *rpc.EngineDeleteRequest) (*empty.Empty, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name": req.Name,
	})

	log.Info("Deleting engine")

	spdkClient, err := s.newSPDKClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get SPDK client").Error())
	}
	defer spdkClient.Close()

	engine, err := getEngine(spdkClient, req.Name, log)
	if err != nil {
		log.WithError(err).Error("Failed to get engine")
		return nil, err
	}

	volumeName := getVolumeName(req.Name)
	err = deleteLonghornDevice(volumeName)
	if err != nil {
		log.WithError(err).Error("Failed to delete longhorn device")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	raidNQN := spdkutil.GetNQN(req.Name)
	localIP, err := issiutil.GetIPToHost()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get local IP")
	}

	listeners, err := spdkClient.NvmfSubsystemGetListeners(raidNQN, "")
	if err != nil {
		log.WithError(err).Error("Failed to get NVMe subsystem listeners")
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to get NVMe subsystem listeners").Error())
	}
	listener := listeners[0]
	localPort := listener.Address.Trsvcid

	nvmeCli, err := spdknvme.NewInitiator(raidNQN, localIP, localPort)
	if err != nil {
		log.WithError(err).Error("Failed to create NVMe initiator")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	err = nvmeCli.StopInitiator()
	if err != nil {
		log.WithError(err).Error("Failed to stop NVMe initiator")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	err = spdkClient.StopExposeBdev(raidNQN)
	if err != nil {
		log.WithError(err).Error("Failed to stop expose bdev")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	portStart, err := strconv.Atoi(localPort)
	if err != nil {
		log.WithError(err).Error("Failed to convert local port to int")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}
	portEnd := portStart
	s.releaseInstancePorts(req.Name, int32(portStart), int32(portEnd))

	_, err = spdkClient.BdevRaidDelete(req.Name)
	if err != nil {
		log.WithError(err).Error("Failed to delete bdev raid")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	s.engineUpdateCh <- interface{}(engine)

	log.Info("Deleted engine")

	return &empty.Empty{}, nil
}

func (s *Server) EngineGet(ctx context.Context, req *rpc.EngineGetRequest) (*rpc.Engine, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"name": req.Name,
	})

	log.Info("Getting engine info")

	spdkClient, err := s.newSPDKClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get SPDK client").Error())
	}
	defer spdkClient.Close()

	engine, err := getEngine(spdkClient, req.Name, log)
	if err != nil {
		log.WithError(err).Error("Failed to get engine")
		return nil, err
	}

	return engine, nil
}

func (s *Server) EngineList(ctx context.Context, req *empty.Empty) (*rpc.EngineListResponse, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{})

	log.Info("Getting engine list")

	spdkClient, err := s.newSPDKClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get SPDK client").Error())
	}
	defer spdkClient.Close()

	resp := &rpc.EngineListResponse{
		Engines: map[string]*rpc.Engine{},
	}
	bdevRaidInfos, err := spdkClient.BdevRaidGetInfoByCategory(spdktypes.BdevRaidCategoryAll)
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
			SpecSize:          0,
			ActualSize:        0,
			Ip:                listenerList[0].Address.Traddr,
			Port:              int32(port),
			ReplicaAddressMap: map[string]string{},
			ReplicaModeMap:    map[string]rpc.ReplicaMode{},
			Endpoint:          "",
		}

		resp.Engines[info.Name] = engine
	}

	log.Info("Got engine list")

	return resp, nil
}

func (s *Server) EngineWatch(req *empty.Empty, srv rpc.SPDKService_EngineWatchServer) error {
	responseChan, err := s.Subscribe(types.InstanceTypeEngine)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			logrus.WithError(err).Error("SPDK service update watch errored out")
		} else {
			logrus.Info("SPDK service update watch ended successfully")
		}
	}()
	logrus.Info("Started new SPDK service update watch")

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
