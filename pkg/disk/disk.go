package disk

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	spdkclient "github.com/longhorn/go-spdk-helper/pkg/spdk/client"
	spdktypes "github.com/longhorn/go-spdk-helper/pkg/spdk/types"

	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/meta"
)

const (
	defaultClusterSize = 4 * 1024 * 1024 // 4MB
	defaultBlockSize   = 4096            // 4KB

	hostPrefix = "/host"
)

type Server struct {
	sync.Mutex

	shutdownCh    chan error
	HealthChecker HealthChecker

	spdkEnabled bool
}

func NewServer(spdkEnabled bool, shutdownCh chan error) (*Server, error) {
	s := &Server{
		spdkEnabled:   spdkEnabled,
		shutdownCh:    shutdownCh,
		HealthChecker: &GRPCHealthChecker{},
	}

	go s.startMonitoring()

	return s, nil
}

func (s *Server) newSPDKClient() (*spdkclient.Client, error) {
	if !s.spdkEnabled {
		return nil, fmt.Errorf("SPDK is not enabled")
	}
	return spdkclient.NewClient()
}

func (s *Server) startMonitoring() {
	done := false
	for {
		select {
		case <-s.shutdownCh:
			logrus.Info("Disk gRPC Server is shutting down")
			done = true
		}
		if done {
			break
		}
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
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "disk name and disk path are required")
	}

	blockSize := uint64(defaultBlockSize)
	if req.BlockSize > 0 {
		log.Infof("Using custom block size %v", req.BlockSize)
		blockSize = uint64(req.BlockSize)
	}

	ok, err := isBlockDevice(req.DiskPath)
	if err != nil {
		log.WithError(err).Error("Failed to check if disk is a block device")
		return nil, grpcstatus.Error(grpccodes.FailedPrecondition, errors.Wrap(err, "failed to check if disk is a block device").Error())
	}
	if !ok {
		log.Errorf("Disk %v is not a block device", req.DiskPath)
		return nil, grpcstatus.Errorf(grpccodes.FailedPrecondition, "disk %v is not a block device", req.DiskPath)
	}

	size, err := getDiskDeviceSize(req.DiskPath)
	if err != nil {
		log.WithError(err).Error("Failed to get disk size")
		return nil, grpcstatus.Error(grpccodes.FailedPrecondition, errors.Wrap(err, "failed to get disk size").Error())
	}
	if size == 0 {
		return nil, grpcstatus.Errorf(grpccodes.FailedPrecondition, "disk %v size is 0", req.DiskPath)
	}

	spdkClient, err := s.newSPDKClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to get SPDK client").Error())
	}
	defer spdkClient.Close()

	diskPath := getDiskPath(req.DiskPath)

	bdevName, err := spdkClient.BdevAioCreate(diskPath, req.DiskName, blockSize)
	if err != nil {
		resp, parseErr := parseErrorMessage(err.Error())
		if parseErr != nil || !isFileExists(resp.Message) {
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
		log.Infof("Found an existing lvstore %v", lvstoreInfos[0].Name)
		lvstoreInfo := &lvstoreInfos[0]
		if lvstoreInfo.Name != lvstoreName {
			log.Infof("Renaming the existing lvstore %v to %v", lvstoreInfo.Name, lvstoreName)
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
		"diskName": req.DiskName,
		"diskPath": req.DiskPath,
	})

	log.Info("Getting disk info")

	if req.DiskName == "" || req.DiskPath == "" {
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "disk path is required")
	}

	spdkClient, err := s.newSPDKClient()
	if err != nil {
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrapf(err, "failed to get SPDK client").Error())
	}
	defer spdkClient.Close()

	// Check if the disk exists
	bdevInfos, err := spdkClient.BdevGetBdevs(req.DiskName, 3000)
	if err != nil {
		resp, parseErr := parseErrorMessage(err.Error())
		if parseErr != nil || !isNoSuchDevice(resp.Message) {
			log.WithError(err).Errorf("Failed to get AIO bdev with name %v", req.DiskName)
			return nil, grpcstatus.Errorf(grpccodes.Internal, errors.Wrapf(err, "failed to get AIO bdev with name %v", req.DiskName).Error())
		}
	}
	if len(bdevInfos) == 0 {
		log.WithError(err).Errorf("Cannot find AIO bdev with name %v", req.DiskName)
		return nil, grpcstatus.Errorf(grpccodes.NotFound, "cannot find AIO bdev with name %v", req.DiskName)
	}

	diskPath := getDiskPath(req.DiskPath)

	var bdevInfo *spdktypes.BdevInfo
	for i, info := range bdevInfos {
		if info.DriverSpecific != nil ||
			info.DriverSpecific.Aio != nil ||
			info.DriverSpecific.Aio.FileName == diskPath {
			bdevInfo = &bdevInfos[i]
			break
		}
	}
	if bdevInfo == nil {
		log.WithError(err).Errorf("Failed to get AIO bdev name for disk path %v", diskPath)
		return nil, grpcstatus.Errorf(grpccodes.NotFound, errors.Wrapf(err, "failed to get AIO bdev name for disk path %v", diskPath).Error())
	}

	// Get the disk information from the lvstore
	lvstoreName := req.DiskName
	lvstoreInfo, err := s.bdevLvolGetLvstore(spdkClient, lvstoreName, "")
	if err != nil {
		log.WithError(err).Errorf("Failed to get %v lvstore", lvstoreName)
		return nil, grpcstatus.Error(grpccodes.NotFound, errors.Wrapf(err, "failed to get %v lvstore", lvstoreName).Error())
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
	return filepath.Join(hostPrefix, path)
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

func isBlockDevice(fullPath string) (bool, error) {
	var st unix.Stat_t
	err := unix.Stat(fullPath, &st)
	if err != nil {
		return false, err
	}

	return (st.Mode & unix.S_IFMT) == unix.S_IFBLK, nil
}

func getDiskDeviceSize(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to open %s", path)
	}
	defer file.Close()

	pos, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to seek %s", path)
	}
	return pos, nil
}
