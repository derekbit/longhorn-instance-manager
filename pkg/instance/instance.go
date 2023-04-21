package instance

import (
	"context"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/longhorn/longhorn-instance-manager/pkg/client"
	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/meta"
)

const (
	BackendStoreDriverTypeLonghorn = "longhorn"
	BackendStoreDriverTypeSpdkAio  = "spdk-aio"
)

type Server struct {
	logsDir                      string
	processManagerServiceAddress string
	diskServiceAddress           string
	shutdownCh                   chan error
	HealthChecker                HealthChecker
}

func NewServer(logsDir, processManagerServiceAddress, diskServiceAddress string, shutdownCh chan error) (*Server, error) {
	s := &Server{
		logsDir:                      logsDir,
		processManagerServiceAddress: processManagerServiceAddress,
		diskServiceAddress:           diskServiceAddress,
		shutdownCh:                   shutdownCh,
		HealthChecker:                &GRPCHealthChecker{},
	}

	go s.startMonitoring()

	return s, nil
}

func (s *Server) startMonitoring() {
	done := false
	for {
		select {
		case <-s.shutdownCh:
			logrus.Info("Instance Server is shutting down")
			done = true
			break
		}
		if done {
			break
		}
	}
}

func (s *Server) VersionGet(ctx context.Context, req *empty.Empty) (*rpc.InstanceVersionResponse, error) {
	v := meta.GetVersion()
	return &rpc.InstanceVersionResponse{
		Version:   v.Version,
		GitCommit: v.GitCommit,
		BuildDate: v.BuildDate,

		InstanceManagerAPIVersion:    int64(v.InstanceManagerAPIVersion),
		InstanceManagerAPIMinVersion: int64(v.InstanceManagerAPIMinVersion),

		InstanceManagerProxyAPIVersion:    int64(v.InstanceManagerProxyAPIVersion),
		InstanceManagerProxyAPIMinVersion: int64(v.InstanceManagerProxyAPIMinVersion),
	}, nil
}

func (s *Server) InstanceCreate(ctx context.Context, req *rpc.InstanceCreateRequest) (*rpc.InstanceResponse, error) {
	logrus.WithFields(logrus.Fields{
		"name":               req.Spec.Name,
		"type":               req.Spec.Type,
		"backendStoreDriver": req.Spec.BackendStoreDriver,
	}).Info("Creating instance")

	pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
	if err != nil {
		return nil, err
	}
	defer pmClient.Close()

	switch req.Spec.BackendStoreDriver {
	case BackendStoreDriverTypeLonghorn:
		if req.Spec.Process == nil {
			return nil, fmt.Errorf("process is required for longhorn backend store driver")
		}
		process, err := pmClient.ProcessCreate(req.Spec.Name, req.Spec.Process.Binary, int(req.Spec.PortCount), req.Spec.Process.Args, req.Spec.PortArgs)
		if err != nil {
			return nil, err
		}

		return &rpc.InstanceResponse{
			Spec: &rpc.InstanceSpec{
				Name: process.Spec.Name,
				Process: &rpc.Process{
					Binary: process.Spec.Binary,
					Args:   process.Spec.Args,
				},
				PortCount: int32(process.Spec.PortCount),
				PortArgs:  process.Spec.PortArgs,
			},
			Status: &rpc.InstanceStatus{
				State:     process.Status.State,
				PortStart: process.Status.PortStart,
				PortEnd:   process.Status.PortEnd,
				ErrorMsg:  process.Status.ErrorMsg,
			},
			Deleted: process.Deleted,
		}, nil
	case BackendStoreDriverTypeSpdkAio:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.Spec.BackendStoreDriver)
	}
}

func (s *Server) InstanceDelete(ctx context.Context, req *rpc.InstanceDeleteRequest) (*rpc.InstanceResponse, error) {
	logrus.WithFields(logrus.Fields{
		"name":               req.Name,
		"type":               req.Type,
		"backendStoreDriver": req.BackendStoreDriver,
	}).Info("Deleting instance")

	pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
	if err != nil {
		return nil, err
	}
	defer pmClient.Close()

	switch req.BackendStoreDriver {
	case BackendStoreDriverTypeLonghorn:
		process, err := pmClient.ProcessDelete(req.Name)
		if err != nil {
			return nil, err
		}

		return &rpc.InstanceResponse{
			Spec: &rpc.InstanceSpec{
				Name: process.Spec.Name,
				Process: &rpc.Process{
					Binary: process.Spec.Binary,
					Args:   process.Spec.Args,
				},
				PortCount: int32(process.Spec.PortCount),
				PortArgs:  process.Spec.PortArgs,
			},
			Status: &rpc.InstanceStatus{
				State:     process.Status.State,
				PortStart: process.Status.PortStart,
				PortEnd:   process.Status.PortEnd,
				ErrorMsg:  process.Status.ErrorMsg,
			},
			Deleted: process.Deleted,
		}, nil
	case BackendStoreDriverTypeSpdkAio:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (s *Server) InstanceGet(ctx context.Context, req *rpc.InstanceGetRequest) (*rpc.InstanceResponse, error) {
	logrus.WithFields(logrus.Fields{
		"name":               req.Name,
		"type":               req.Type,
		"backendStoreDriver": req.BackendStoreDriver,
	}).Info("Getting instance")

	pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
	if err != nil {
		return nil, err
	}
	defer pmClient.Close()

	switch req.BackendStoreDriver {
	case BackendStoreDriverTypeLonghorn:
		process, err := pmClient.ProcessGet(req.Name)
		if err != nil {
			return nil, err
		}

		return &rpc.InstanceResponse{
			Spec: &rpc.InstanceSpec{
				Name: process.Spec.Name,
				Process: &rpc.Process{
					Binary: process.Spec.Binary,
					Args:   process.Spec.Args,
				},
				PortCount: int32(process.Spec.PortCount),
				PortArgs:  process.Spec.PortArgs,
			},
			Status: &rpc.InstanceStatus{
				State:     process.Status.State,
				PortStart: process.Status.PortStart,
				PortEnd:   process.Status.PortEnd,
				ErrorMsg:  process.Status.ErrorMsg,
			},
			Deleted: process.Deleted,
		}, nil
	case BackendStoreDriverTypeSpdkAio:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (s *Server) InstanceList(ctx context.Context, req *empty.Empty) (*rpc.InstanceListResponse, error) {
	logrus.WithFields(logrus.Fields{}).Info("Listing instances")

	pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
	if err != nil {
		return nil, err
	}
	defer pmClient.Close()

	processes, err := pmClient.ProcessList()
	if err != nil {
		return nil, err
	}

	instances := map[string]*rpc.InstanceResponse{}

	for _, process := range processes {
		instances[process.Spec.Name] = &rpc.InstanceResponse{
			Spec: &rpc.InstanceSpec{
				Name: process.Spec.Name,
				Process: &rpc.Process{
					Binary: process.Spec.Binary,
					Args:   process.Spec.Args,
				},
				PortCount: int32(process.Spec.PortCount),
				PortArgs:  process.Spec.PortArgs,
			},
			Status: &rpc.InstanceStatus{
				State:     process.Status.State,
				PortStart: process.Status.PortStart,
				PortEnd:   process.Status.PortEnd,
				ErrorMsg:  process.Status.ErrorMsg,
			},
			Deleted: process.Deleted,
		}
	}

	return &rpc.InstanceListResponse{
		Instances: instances,
	}, nil
}

func (s *Server) InstanceReplace(ctx context.Context, req *rpc.InstanceReplaceRequest) (*rpc.InstanceResponse, error) {
	logrus.WithFields(logrus.Fields{
		"name":               req.Spec.Name,
		"type":               req.Spec.Type,
		"backendStoreDriver": req.Spec.BackendStoreDriver,
	}).Info("Replacing instance")

	pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
	if err != nil {
		return nil, err
	}
	defer pmClient.Close()

	switch req.Spec.BackendStoreDriver {
	case BackendStoreDriverTypeLonghorn:
		if req.Spec.Process == nil {
			return nil, fmt.Errorf("process is required for longhorn backend store driver")
		}

		process, err := pmClient.ProcessReplace(req.Spec.Name, req.Spec.Process.Binary, int(req.Spec.PortCount), req.Spec.Process.Args, req.Spec.PortArgs, req.TerminateSignal)
		if err != nil {
			return nil, err
		}

		return &rpc.InstanceResponse{
			Spec: &rpc.InstanceSpec{
				Name: process.Spec.Name,
				Process: &rpc.Process{
					Binary: process.Spec.Binary,
					Args:   process.Spec.Args,
				},
				PortCount: int32(process.Spec.PortCount),
				PortArgs:  process.Spec.PortArgs,
			},
			Status: &rpc.InstanceStatus{
				State:     process.Status.State,
				PortStart: process.Status.PortStart,
				PortEnd:   process.Status.PortEnd,
				ErrorMsg:  process.Status.ErrorMsg,
			},
			Deleted: process.Deleted,
		}, nil
	case BackendStoreDriverTypeSpdkAio:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.Spec.BackendStoreDriver)
	}
}

func (s *Server) InstanceLog(req *rpc.InstanceLogRequest, srv rpc.InstanceService_InstanceLogServer) error {
	logrus.WithFields(logrus.Fields{
		"name":               req.Name,
		"type":               req.Type,
		"backendStoreDriver": req.BackendStoreDriver,
	}).Info("Getting instance log")

	pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
	if err != nil {
		return err
	}
	defer pmClient.Close()

	switch req.BackendStoreDriver {
	case BackendStoreDriverTypeLonghorn:
		stream, err := pmClient.ProcessLog(context.Background(), req.Name)
		if err != nil {
			return err
		}
		for {
			line, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err := srv.Send(&rpc.LogResponse{Line: line}); err != nil {
				return err
			}
		}
		return nil
	case BackendStoreDriverTypeSpdkAio:
		return nil
	default:
		return fmt.Errorf("unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (s *Server) InstanceWatch(req *empty.Empty, srv rpc.InstanceService_InstanceWatchServer) error {
	return nil
}
