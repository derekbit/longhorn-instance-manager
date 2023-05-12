package disk

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	spdkhelpertypes "github.com/longhorn/go-spdk-helper/pkg/types"
	spdkclient "github.com/longhorn/longhorn-spdk-engine/pkg/client"
	spdkrpc "github.com/longhorn/longhorn-spdk-engine/proto/spdkrpc"

	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/meta"
)

const (
	spdkTgtReadinessProbeTimeout = 60 * time.Second
)

type Server struct {
	sync.RWMutex

	spdkServiceAddress string
	shutdownCh         chan error
	HealthChecker      HealthChecker

	spdkEnabled bool

	spdkClient *spdkclient.SPDKClient
}

func NewServer(spdkEnabled bool, spdkServiceAddress string, shutdownCh chan error) (*Server, error) {
	s := &Server{
		spdkEnabled:        spdkEnabled,
		spdkServiceAddress: spdkServiceAddress,
		shutdownCh:         shutdownCh,
		HealthChecker:      &GRPCHealthChecker{},
	}

	if s.spdkEnabled {
		logrus.Info("Creating SPDK client since SPDK is enabled")

		if !isSPDKTgtReady(spdkTgtReadinessProbeTimeout) {
			return nil, fmt.Errorf("spdk_tgt is not ready in %v", spdkTgtReadinessProbeTimeout)
		}

		spdkClient, err := spdkclient.NewSPDKClient(s.spdkServiceAddress)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create SPDK client")
		}
		s.spdkClient = spdkClient
	}

	go s.startMonitoring()

	return s, nil
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
	log := logrus.WithFields(logrus.Fields{
		"diskType":  req.DiskType,
		"diskName":  req.DiskName,
		"diskPath":  req.DiskPath,
		"blockSize": req.BlockSize,
	})

	log.Info("Creating disk")

	if req.DiskName == "" || req.DiskPath == "" {
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "disk name and disk path are required")
	}

	s.Lock()
	defer s.Unlock()

	switch req.DiskType {
	case rpc.DiskType_block:
		ret, err := s.spdkClient.DiskCreate(req.DiskName, req.DiskPath, req.BlockSize)
		if err != nil {
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}
		return spdkDiskToDisk(ret), nil
	case rpc.DiskType_filesystem:
		// TODO: implement filesystem disk type
		fallthrough
	default:
		return nil, grpcstatus.Errorf(grpccodes.Unimplemented, "unsupported disk type %v", req.DiskType)
	}
}

func (s *Server) DiskDelete(ctx context.Context, req *rpc.DiskDeleteRequest) (*emptypb.Empty, error) {
	log := logrus.WithFields(logrus.Fields{
		"diskName": req.DiskName,
		"diskUUID": req.DiskUuid,
	})

	log.Info("Deleting disk")

	if req.DiskName == "" || req.DiskUuid == "" {
		return &empty.Empty{}, grpcstatus.Error(grpccodes.InvalidArgument, "disk name and disk UUID are required")
	}

	s.Lock()
	defer s.Unlock()

	if s.spdkEnabled {
		return &empty.Empty{}, s.spdkClient.DiskDelete(req.DiskName, req.DiskUuid)
	}
	return &empty.Empty{}, nil
}

func (s *Server) DiskGet(ctx context.Context, req *rpc.DiskGetRequest) (*rpc.Disk, error) {
	log := logrus.WithFields(logrus.Fields{
		"diskType": req.DiskType,
		"diskName": req.DiskName,
		"diskPath": req.DiskPath,
	})

	log.Info("Getting disk info")

	if req.DiskName == "" || req.DiskPath == "" {
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "disk name and disk path are required")
	}

	s.RLock()
	defer s.RUnlock()

	switch req.DiskType {
	case rpc.DiskType_block:
		ret, err := s.spdkClient.DiskGet(req.DiskName, req.DiskPath)
		if err != nil {
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}
		return spdkDiskToDisk(ret), nil
	case rpc.DiskType_filesystem:
		// TODO: implement filesystem disk type
		fallthrough
	default:
		return nil, grpcstatus.Errorf(grpccodes.Unimplemented, "unsupported disk type %v", req.DiskType)
	}
}

func isSPDKTgtReady(timeout time.Duration) bool {
	for i := 0; i < int(timeout.Seconds()); i++ {
		conn, err := net.DialTimeout(spdkhelpertypes.DefaultJSONServerNetwork, spdkhelpertypes.DefaultUnixDomainSocketPath, 1*time.Second)
		if err == nil {
			conn.Close()
			return true
		}
		time.Sleep(time.Second)
	}
	return false
}

func spdkDiskToDisk(disk *spdkrpc.Disk) *rpc.Disk {
	return &rpc.Disk{
		Id:          disk.Id,
		Uuid:        disk.Uuid,
		Path:        disk.Path,
		Type:        disk.Type,
		TotalSize:   disk.TotalSize,
		FreeSize:    disk.FreeSize,
		TotalBlocks: disk.TotalBlocks,
		FreeBlocks:  disk.FreeBlocks,
		BlockSize:   disk.BlockSize,
		ClusterSize: disk.ClusterSize,
	}
}
