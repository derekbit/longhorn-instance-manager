package disk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/ptypes/empty"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	spdkclient "github.com/longhorn/go-spdk-helper/pkg/spdk/client"
	spdktypes "github.com/longhorn/go-spdk-helper/pkg/spdk/types"

	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/types"
)

const (
	defaultClusterSize = 4 * 1024 * 1024 // 4MB
)

type Server struct {
	shutdownCh    chan error
	HealthChecker HealthChecker
}

func NewServer(shutdownCh chan error) (*Server, error) {
	s := &Server{
		shutdownCh:    shutdownCh,
		HealthChecker: &GRPCHealthChecker{},
	}

	go s.startMonitoring()

	return s, nil
}

func (s *Server) startMonitoring() {
	done := false
	for {
		select {
		case <-s.shutdownCh:
			logrus.Info("Disk Server is shutting down")
			done = true
			break
		}
		if done {
			break
		}
	}
}

func (s *Server) DiskCreate(ctx context.Context, req *rpc.DiskCreateRequest) (*rpc.DiskCreateResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"diskName": req.DiskName,
		"diskPath": req.DiskPath,
	})

	log.Info("Creating disk")

	spdkCli, err := spdkclient.NewClient()
	if err != nil {
		log.WithError(err).Error("Failed to create SPDK client")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	diskPath := getDiskPath(req.DiskPath)
	bdevName, err := spdkCli.BdevAioCreate(diskPath, req.DiskName, 4096)
	if err != nil {
		log.WithError(err).Error("Failed to create AIO bdev")
		return nil, grpcstatus.Error(grpccodes.NotFound, err.Error())
	}

	uuid, err := spdkCli.BdevLvolCreateLvstore(bdevName, req.DiskName, defaultClusterSize)
	if err != nil {
		log.WithError(err).Error("Failed to create lvstore")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	lvstoreInfos, err := spdkCli.BdevLvolGetLvstore("", uuid)
	if err != nil {
		log.WithError(err).Error("Failed to get lvstore info")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}
	if len(lvstoreInfos) != 1 {
		log.WithError(err).Error("Number of lvstore info is not 1")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	log.Info("Created disk")

	diskID, err := getDeviceNumber(diskPath)
	if err != nil {
		log.WithError(err).Error("Failed to get disk ID")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	info := lvstoreInfos[0]

	return &rpc.DiskCreateResponse{
		DiskInfo: &rpc.DiskInfo{
			Id:          diskID,
			Uuid:        info.UUID,
			Path:        req.DiskPath,
			Type:        "block",
			TotalSize:   int64(info.TotalDataClusters * info.ClusterSize),
			FreeSize:    int64(info.FreeClusters * info.ClusterSize),
			TotalBlocks: int64(info.TotalDataClusters * info.ClusterSize / info.BlockSize),
			FreeBlocks:  int64(info.FreeClusters * info.ClusterSize / info.BlockSize),
			BlockSize:   int64(info.BlockSize),
			ClusterSize: int64(info.ClusterSize),
			Readonly:    false,
		},
	}, nil
}

func (s *Server) DiskInfo(ctx context.Context, req *rpc.DiskInfoRequest) (*rpc.DiskInfoResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"diskPath": req.DiskPath,
	})

	log.Info("Getting disk info")

	spdkCli, err := spdkclient.NewClient()
	if err != nil {
		log.WithError(err).Error("Failed to create SPDK client")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	bdevInfos, err := spdkCli.BdevGetBdevs("", 30)
	if err != nil {
		log.WithError(err).Error("Failed to get bdev infos")
		return nil, grpcstatus.Error(grpccodes.NotFound, err.Error())
	}

	diskPath := getDiskPath(req.DiskPath)
	bdevName := ""
	for _, info := range bdevInfos {
		if info.DriverSpecific != nil &&
			info.DriverSpecific.Aio != nil &&
			info.DriverSpecific.Aio.FileName == diskPath {
			bdevName = info.Name
		}
	}
	if bdevName == "" {
		log.WithError(err).Error("Failed to find bdev name")
		return nil, grpcstatus.Errorf(grpccodes.NotFound, "failed to find bdev name for disk path %v", diskPath)
	}

	lvstoreName := bdevName
	lvstoreInfos, err := spdkCli.BdevLvolGetLvstore(lvstoreName, "")
	if err != nil {
		log.WithError(err).Error("Failed to get lvstore info")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	if len(lvstoreInfos) != 1 {
		log.WithError(err).Error("Number of lvstore info is not 1")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	log.Info("Got disk info")

	diskID, err := getDeviceNumber(diskPath)
	if err != nil {
		log.WithError(err).Error("Failed to get disk ID")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	info := lvstoreInfos[0]

	return &rpc.DiskInfoResponse{
		DiskInfo: &rpc.DiskInfo{
			Id:          diskID,
			Uuid:        info.UUID,
			Path:        req.DiskPath,
			Type:        "block",
			TotalSize:   int64(info.TotalDataClusters * info.ClusterSize),
			FreeSize:    int64(info.FreeClusters * info.ClusterSize),
			TotalBlocks: int64(info.TotalDataClusters * info.ClusterSize / info.BlockSize),
			FreeBlocks:  int64(info.FreeClusters * info.ClusterSize / info.BlockSize),
			BlockSize:   int64(info.BlockSize),
			ClusterSize: int64(info.ClusterSize),
			Readonly:    false,
		},
	}, nil
}

func (s *Server) ReplicaCreate(ctx context.Context, req *rpc.ReplicaCreateRequest) (*rpc.ReplicaCreateResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"name":        req.Name,
		"lvstoreUUID": req.LvstoreUuid,
		"size":        req.Size,
	})

	log.Info("Creating replica")

	spdkCli, err := spdkclient.NewClient()
	if err != nil {
		log.WithError(err).Error("Failed to create SPDK client")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	lvstoreUUID := req.LvstoreUuid
	lvstoreInfos, err := spdkCli.BdevLvolGetLvstore("", lvstoreUUID)
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

	uuid, err := spdkCli.BdevLvolCreate(lvstoreInfo.Name, req.Name, "", sizeInMib, spdktypes.BdevLvolClearMethodUnmap, true)
	if err != nil {
		log.WithError(err).Error("Failed to create lvol")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	log.Info("Created replica")

	lvolInfos, err := spdkCli.BdevLvolGet(uuid, 30)
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

	return &rpc.ReplicaCreateResponse{
		ReplicaInfo: &rpc.ReplicaInfo{
			Name:        name,
			Uuid:        lvolInfo.UUID,
			BdevName:    lvolInfo.DriverSpecific.Lvol.BaseBdev,
			LvstoreUuid: lvolInfo.DriverSpecific.Lvol.LvolStoreUUID,

			TotalSize:   int64(lvolInfo.NumBlocks * uint64(lvolInfo.BlockSize)),
			TotalBlocks: int64(lvolInfo.NumBlocks),
			BlockSize:   int64(lvolInfo.BlockSize),

			ThinProvision: lvolInfo.DriverSpecific.Lvol.ThinProvision,

			State: types.ProcessStateRunning,
		},
	}, nil
}

func (s *Server) ReplicaDelete(ctx context.Context, req *rpc.ReplicaDeleteRequest) (*empty.Empty, error) {
	log := logrus.WithFields(logrus.Fields{
		"name":        req.Name,
		"lvstoreUUID": req.LvstoreUuid,
	})

	log.Info("Deleting replica")

	spdkCli, err := spdkclient.NewClient()
	if err != nil {
		log.WithError(err).Error("Failed to create SPDK client")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	lvstoreInfos, err := spdkCli.BdevLvolGetLvstore("", req.LvstoreUuid)
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
	deleted, err := spdkCli.BdevLvolDelete(aliasName)
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

func (s *Server) ReplicaInfo(ctx context.Context, req *rpc.ReplicaInfoRequest) (*rpc.ReplicaInfoResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"name":        req.Name,
		"lvstoreUUID": req.LvstoreUuid,
	})

	log.Info("Getting replica info")

	spdkCli, err := spdkclient.NewClient()
	if err != nil {
		log.WithError(err).Error("Failed to create SPDK client")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	lvstoreInfos, err := spdkCli.BdevLvolGetLvstore("", req.LvstoreUuid)
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
	lvolInfos, err := spdkCli.BdevLvolGet(aliasName, 30)
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

	return &rpc.ReplicaInfoResponse{
		ReplicaInfo: &rpc.ReplicaInfo{
			Name:        name,
			Uuid:        lvolInfo.UUID,
			BdevName:    lvolInfo.DriverSpecific.Lvol.BaseBdev,
			LvstoreUuid: lvolInfo.DriverSpecific.Lvol.BaseBdev,

			TotalSize:   int64(lvolInfo.NumBlocks * uint64(lvolInfo.BlockSize)),
			TotalBlocks: int64(lvolInfo.NumBlocks),
			BlockSize:   int64(lvolInfo.BlockSize),

			ThinProvision: lvolInfo.DriverSpecific.Lvol.ThinProvision,

			State: types.ProcessStateRunning,
		},
	}, nil
}

func (s *Server) ReplicaList(ctx context.Context, req *empty.Empty) (*rpc.ReplicaListResponse, error) {
	log := logrus.WithFields(logrus.Fields{})

	log.Info("Getting replica list")

	spdkCli, err := spdkclient.NewClient()
	if err != nil {
		log.WithError(err).Error("Failed to create SPDK client")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	lvolInfos, err := spdkCli.BdevLvolGet("", 30)
	if err != nil {
		log.WithError(err).Error("Failed to get lvol infos")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	log.Info("Got replica list")

	resp := &rpc.ReplicaListResponse{
		ReplicaInfos: map[string]*rpc.ReplicaInfo{},
	}
	for _, info := range lvolInfos {
		name := getNameFromAlias(info.Aliases)
		if name == "" {
			return nil, grpcstatus.Error(grpccodes.Internal, "Failed to get name from alias")
		}

		resp.ReplicaInfos[info.Name] = &rpc.ReplicaInfo{
			Name:        name,
			Uuid:        info.UUID,
			BdevName:    info.DriverSpecific.Lvol.BaseBdev,
			LvstoreUuid: info.DriverSpecific.Lvol.BaseBdev,

			TotalSize:   int64(info.NumBlocks * uint64(info.BlockSize)),
			TotalBlocks: int64(info.NumBlocks),
			BlockSize:   int64(info.BlockSize),

			ThinProvision: info.DriverSpecific.Lvol.ThinProvision,

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
