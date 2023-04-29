package spdk

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	spdkclient "github.com/longhorn/go-spdk-helper/pkg/spdk/client"

	"github.com/longhorn/longhorn-instance-manager/pkg/util/broadcaster"
)

const (
	defaultClusterSize = 4 * 1024 * 1024 // 4MB
	defaultBlockSize   = 4096            // 4KB

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
		return spdkclient.NewClient()
	}
	return nil, fmt.Errorf("spdk is not enabled")
}

func (s *Server) Close() error {
	if s.spdkClient == nil {
		return nil
	}
	return s.spdkClient.Close()
}

func (s *Server) startMonitoring() {
	done := false
	for {
		select {
		case <-s.shutdownCh:
			logrus.Info("SPDK Server is shutting down")
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

/*
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

func (s *Server) DiskCreate(ctx context.Context, req *rpc.DiskCreateRequest) (*rpc.Disk, error) {
	s.Lock()
	defer s.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"diskName":  req.DiskName,
		"diskPath":  req.DiskPath,
		"blockSize": req.BlockSize,
	})

	log.Info("Creating disk")

	if req.DiskName == "" || req.DiskPath == "" {
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "disk name and path are required")
	}

	blockSize := uint64(defaultBlockSize)
	if req.BlockSize > 0 {
		blockSize = uint64(req.BlockSize)
	}

	ok, err := util.IsBlockDevice(req.DiskPath)
	if err != nil {
		log.WithError(err).Error("Failed to check if disk is block device")
		return nil, grpcstatus.Error(grpccodes.FailedPrecondition, errors.Wrap(err, "failed to check if disk is block device").Error())
	}
	if !ok {
		log.Errorf("Disk %v is not a block device", req.DiskPath)
		return nil, grpcstatus.Errorf(grpccodes.FailedPrecondition, "disk %v is not a block device", req.DiskPath)
	}

	size, err := util.GetDiskDeviceSize(req.DiskPath)
	if err != nil {
		log.WithError(err).Error("Failed to get disk size")
		return nil, grpcstatus.Error(grpccodes.FailedPrecondition, errors.Wrap(err, "failed to get disk size").Error())
	}
	if size == 0 {
		return nil, grpcstatus.Errorf(grpccodes.FailedPrecondition, "disk %v size is 0", req.DiskPath)
	}

	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get spdk client").Error())
	}
	defer spdkClient.Close()

	diskPath := getDiskPath(req.DiskPath)

	bdevName, err := spdkClient.BdevAioCreate(diskPath, req.DiskName, blockSize)
	if err != nil {
		resp, parseErr := parseErrorMessage(err.Error())
		if parseErr != nil || !strings.EqualFold(resp.Message, syscall.Errno(syscall.EEXIST).Error()) {
			log.WithError(err).Error("Failed to create AIO bdev")
			return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to create AIO bdev").Error())
		}
	}

	lvstoreName := bdevName
	lvstoreInfos, err := spdkClient.BdevLvolGetLvstore("", "")
	if err != nil {
		log.WithError(err).Errorf("Failed to get lvstore with name %v", lvstoreName)
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to get lvstore with name %v", lvstoreName).Error())
	}

	uuid := ""
	switch len(lvstoreInfos) {
	case 0:
		uuid, err = spdkClient.BdevLvolCreateLvstore(lvstoreName, req.DiskName, defaultClusterSize)
		if err != nil {
			log.WithError(err).Error("Failed to create lvstore")
			return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to create lvstore").Error())
		}
	case 1:
		log.Info("Found an existing lvstore %v", lvstoreInfos[0].Name)
		lvstoreInfo := &lvstoreInfos[0]
		if lvstoreInfo.Name != lvstoreName {
			log.Info("Renaming the existing lvstore %v to %v", lvstoreInfo.Name, lvstoreName)
			renamed, err := spdkClient.BdevLvolRenameLvstore(lvstoreInfo.Name, lvstoreName)
			if err != nil {
				log.WithError(err).Errorf("Failed to rename lvstore from %v to %v", lvstoreInfo.Name, lvstoreName)
				return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to rename lvstore from %v to %v", lvstoreInfo.Name, lvstoreName).Error())
			}
			if !renamed {
				log.Errorf("Failed to rename lvstore from %v to %v", lvstoreInfo.Name, lvstoreName)
				return nil, grpcstatus.Errorf(grpccodes.Internal, "failed to rename lvstore from %v to %v", lvstoreInfo.Name, lvstoreName)
			}
		}
		uuid = lvstoreInfo.UUID
	default:
		log.Error("Found more than one lvstore")
		return nil, grpcstatus.Error(grpccodes.Internal, "found more than one lvstore")
	}

	lvstoreInfo, err := s.bdevLvolGetLvstore(spdkClient, "", uuid)
	if err != nil {
		log.WithError(err).Errorf("Failed to get lvstore with name %v and UUID %v", req.DiskName, uuid)
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to get lvstore with name %v and UUID %v", req.DiskName, uuid).Error())
	}

	log.Info("Created disk")

	// A disk does not have a fsid, so we use the device number as the disk ID
	diskID, err := util.GetDeviceNumber(diskPath)
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
		"diskName": req.DiskName,
		"diskPath": req.DiskPath,
	})

	log.Info("Getting disk info")

	if req.DiskName == "" || req.DiskPath == "" {
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "disk path is required")
	}

	spdkClient, err := s.getSpdkClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to get spdk client").Error())
	}
	defer spdkClient.Close()

	// Check if the disk exists
	bdevInfos, err := spdkClient.BdevGetBdevs(req.DiskName, 3000)
	if err != nil {
		resp, parseErr := parseErrorMessage(err.Error())
		if parseErr != nil || !strings.EqualFold(resp.Message, syscall.Errno(syscall.ENODEV).Error()) {
			log.WithError(err).Errorf("Failed to get bdev with name %v", req.DiskName)
			return nil, grpcstatus.Errorf(grpccodes.Internal, errors.Wrapf(err, "failed to get bdev with name %v", req.DiskName).Error())
		}
	}
	if len(bdevInfos) == 0 {
		log.WithError(err).Errorf("Cannot find bdev with name %v", req.DiskName)
		return nil, grpcstatus.Errorf(grpccodes.NotFound, "cannot find bdev with name %v", req.DiskName)
	}

	diskPath := getDiskPath(req.DiskPath)

	var bdevInfo *spdktypes.BdevInfo
	for i, info := range bdevInfos {
		if info.DriverSpecific != nil ||
			info.DriverSpecific.Aio != nil ||
			info.DriverSpecific.Aio.FileName == diskPath {
			bdevInfo = &bdevInfos[i]
		}
	}
	if bdevInfo == nil {
		log.WithError(err).Errorf("Failed to get bdev name for disk path %v", diskPath)
		return nil, grpcstatus.Errorf(grpccodes.NotFound, errors.Wrapf(err, "failed to get bdev name for disk path %v", diskPath).Error())
	}

	// Get the disk information from the lvstore
	lvstoreName := req.DiskName
	lvstoreInfo, err := s.bdevLvolGetLvstore(spdkClient, lvstoreName, "")
	if err != nil {
		log.WithError(err).Errorf("Failed to get %v lvstore", lvstoreName)
		return nil, grpcstatus.Error(grpccodes.NotFound, errors.Wrapf(err, "failed to get %v lvstore", lvstoreName).Error())
	}

	log.Info("Got disk info")

	diskID, err := util.GetDeviceNumber(diskPath)
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

func (s *Server) bdevLvolGetLvstore(spdkClient *spdkclient.Client, lvsName, uuid string) (*spdktypes.LvstoreInfo, error) {
	lvstoreInfos, err := spdkClient.BdevLvolGetLvstore(lvsName, uuid)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get lvstore info")
	}

	if len(lvstoreInfos) == 0 {
		return nil, errors.New(syscall.ENOENT.Error())
	}

	if len(lvstoreInfos) != 1 {
		return nil, fmt.Errorf("number of lvstores with name %v is not 1", lvsName)
	}

	return &lvstoreInfos[0], nil
}

func getDiskPath(path string) string {
	return filepath.Join("/host", path)
}
*/
