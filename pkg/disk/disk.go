package disk

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/ptypes/empty"

	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
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

	return &rpc.DiskCreateResponse{}, nil
}

func (s *Server) DiskInfo(ctx context.Context, req *rpc.DiskInfoRequest) (*rpc.DiskInfoResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"diskPath": req.DiskPath,
	})

	log.Info("Getting disk info")

	return &rpc.DiskInfoResponse{}, nil
}

func (s *Server) ReplicaCreate(ctx context.Context, req *rpc.ReplicaCreateRequest) (*rpc.ReplicaCreateResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"name":        req.Name,
		"lvstoreUUID": req.LvstoreUuid,
		"size":        req.Size,
	})

	log.Info("Creating replica")

	return &rpc.ReplicaCreateResponse{}, nil
}

func (s *Server) ReplicaDelete(ctx context.Context, req *rpc.ReplicaDeleteRequest) (*empty.Empty, error) {
	log := logrus.WithFields(logrus.Fields{
		"name":        req.Name,
		"lvstoreUUID": req.LvstoreUuid,
	})

	log.Info("Deleting replica")

	return &empty.Empty{}, nil
}

func (s *Server) ReplicaInfo(ctx context.Context, req *rpc.ReplicaInfoRequest) (*rpc.ReplicaInfoResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"name":        req.Name,
		"lvstoreUUID": req.LvstoreUuid,
	})

	log.Info("Getting replica info")

	return &rpc.ReplicaInfoResponse{}, nil
}

func (s *Server) ReplicaList(ctx context.Context, req *empty.Empty) (*rpc.ReplicaListResponse, error) {
	log := logrus.WithFields(logrus.Fields{})

	log.Info("Getting replica list")

	return &rpc.ReplicaListResponse{}, nil
}
