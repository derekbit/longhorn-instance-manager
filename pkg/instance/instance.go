package instance

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/longhorn/longhorn-instance-manager/pkg/client"
	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/meta"
	"github.com/longhorn/longhorn-instance-manager/pkg/types"
)

const (
	BackendStoreDriverTypeUndefined = ""
	BackendStoreDriverTypeLonghorn  = "longhorn"
	BackendStoreDriverTypeSpdkAio   = "spdk-aio"
)

type Server struct {
	logsDir                      string
	processManagerServiceAddress string
	diskServiceAddress           string
	spdkEnabled                  bool
	shutdownCh                   chan error
	HealthChecker                HealthChecker
}

func NewServer(logsDir, processManagerServiceAddress, diskServiceAddress string, spdkEnabled bool, shutdownCh chan error) (*Server, error) {
	s := &Server{
		logsDir:                      logsDir,
		processManagerServiceAddress: processManagerServiceAddress,
		diskServiceAddress:           diskServiceAddress,
		spdkEnabled:                  spdkEnabled,
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

func (s *Server) InstanceCreate(ctx context.Context, req *rpc.InstanceCreateRequest) (*rpc.InstanceResponse, error) {
	logrus.WithFields(logrus.Fields{
		"name":               req.Spec.Name,
		"type":               req.Spec.Type,
		"backendStoreDriver": req.Spec.BackendStoreDriver,
	}).Info("Creating instance")

	switch req.Spec.BackendStoreDriver {
	case BackendStoreDriverTypeLonghorn:
		pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
		if err != nil {
			return nil, err
		}
		defer pmClient.Close()

		if req.Spec.ProcessSpecific == nil {
			return nil, fmt.Errorf("process is required for longhorn backend store driver")
		}
		process, err := pmClient.ProcessCreate(req.Spec.Name, req.Spec.ProcessSpecific.Binary, int(req.Spec.PortCount), req.Spec.ProcessSpecific.Args, req.Spec.PortArgs)
		if err != nil {
			return nil, err
		}
		return processResponseToInstanceResponse(process), nil
	case BackendStoreDriverTypeSpdkAio:
		diskClient, err := client.NewDiskServiceClient("tcp://"+s.diskServiceAddress, nil)
		if err != nil {
			return nil, err
		}
		defer diskClient.Close()

		switch req.Spec.Type {
		case types.InstanceTypeEngine:
			engine, err := diskClient.EngineCreate(req.Spec.Name, req.Spec.VolumeName, req.Spec.SpdkSpecific.Frontend, req.Spec.SpdkSpecific.DiskUuid, req.Spec.SpdkSpecific.ReplicaAddressMap)
			if err != nil {
				return nil, err
			}
			return engineResponseToInstanceResponse(engine), nil
		case types.InstanceTypeReplica:
			replica, err := diskClient.ReplicaCreate(req.Spec.Name, req.Spec.SpdkSpecific.DiskUuid, req.Spec.Size, req.Spec.SpdkSpecific.Address)
			if err != nil {
				return nil, err
			}
			return replicaResponseToInstanceResponse(replica), nil
		default:
			return nil, fmt.Errorf("unknown instance type %v", req.Spec.Type)
		}
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.Spec.BackendStoreDriver)
	}
}

func (s *Server) InstanceDelete(ctx context.Context, req *rpc.InstanceDeleteRequest) (*rpc.InstanceResponse, error) {
	logrus.WithFields(logrus.Fields{
		"name":               req.Name,
		"type":               req.Type,
		"backendStoreDriver": req.BackendStoreDriver,
		"diskUuid":           req.DiskUuid,
	}).Info("Deleting instance")

	switch req.BackendStoreDriver {
	case BackendStoreDriverTypeLonghorn:
		pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
		if err != nil {
			return nil, err
		}
		defer pmClient.Close()

		process, err := pmClient.ProcessDelete(req.Name)
		if err != nil {
			return nil, err
		}
		return processResponseToInstanceResponse(process), nil
	case BackendStoreDriverTypeSpdkAio:
		diskClient, err := client.NewDiskServiceClient("tcp://"+s.diskServiceAddress, nil)
		if err != nil {
			return nil, err
		}
		defer diskClient.Close()

		err = diskClient.ReplicaDelete(req.Name, req.DiskUuid)
		if err != nil {
			return nil, err
		}
		return &rpc.InstanceResponse{
			Spec: &rpc.InstanceSpec{
				Name: req.Name,
			},
			Status: &rpc.InstanceStatus{
				State: types.ProcessStateStopped,
			},
			Deleted: true,
		}, nil
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

	switch req.BackendStoreDriver {
	case BackendStoreDriverTypeLonghorn:
		pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
		if err != nil {
			return nil, err
		}
		defer pmClient.Close()

		process, err := pmClient.ProcessGet(req.Name)
		if err != nil {
			return nil, err
		}
		return processResponseToInstanceResponse(process), nil
	case BackendStoreDriverTypeSpdkAio:
		diskClient, err := client.NewDiskServiceClient("tcp://"+s.diskServiceAddress, nil)
		if err != nil {
			return nil, err
		}
		defer diskClient.Close()
		switch req.Type {
		case types.InstanceTypeEngine:
			engine, err := diskClient.EngineGet(req.Name)
			if err != nil {
				return nil, err
			}
			return engineResponseToInstanceResponse(engine), nil
		case types.InstanceTypeReplica:
			replica, err := diskClient.ReplicaGet(req.Name, req.DiskUuid)
			if err != nil {
				return nil, err
			}
			return replicaResponseToInstanceResponse(replica), nil
		default:
			return nil, fmt.Errorf("unknown instance type %v", req.Type)
		}
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (s *Server) InstanceList(ctx context.Context, req *empty.Empty) (*rpc.InstanceListResponse, error) {
	logrus.WithFields(logrus.Fields{}).Info("Listing instances")

	instances := map[string]*rpc.InstanceResponse{}

	// Collect engine and replica processes as instances
	pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
	if err != nil {
		return nil, err
	}
	defer pmClient.Close()

	processes, err := pmClient.ProcessList()
	if err != nil {
		return nil, err
	}
	for _, process := range processes {
		instances[process.Spec.Name] = processResponseToInstanceResponse(process)
	}

	// Collect spdk instances as instances
	if s.spdkEnabled {
		diskClient, err := client.NewDiskServiceClient("tcp://"+s.diskServiceAddress, nil)
		if err != nil {
			return nil, err
		}
		defer diskClient.Close()

		replicas, err := diskClient.ReplicaList()
		if err != nil {
			return nil, err
		}
		for _, replica := range replicas {
			instances[replica.Name] = replicaResponseToInstanceResponse(replica)
		}

		engines, err := diskClient.EngineList()
		if err != nil {
			return nil, err
		}
		for _, engine := range engines {
			instances[engine.Name] = engineResponseToInstanceResponse(engine)
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

	switch req.Spec.BackendStoreDriver {
	case BackendStoreDriverTypeLonghorn:
		if req.Spec.ProcessSpecific == nil {
			return nil, fmt.Errorf("process is required for longhorn backend store driver")
		}

		pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
		if err != nil {
			return nil, err
		}
		defer pmClient.Close()

		process, err := pmClient.ProcessReplace(req.Spec.Name,
			req.Spec.ProcessSpecific.Binary, int(req.Spec.PortCount), req.Spec.ProcessSpecific.Args, req.Spec.PortArgs, req.TerminateSignal)
		if err != nil {
			return nil, err
		}

		return processResponseToInstanceResponse(process), nil
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

	switch req.BackendStoreDriver {
	case BackendStoreDriverTypeLonghorn:
		pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
		if err != nil {
			return err
		}
		defer pmClient.Close()

		stream, err := pmClient.ProcessLog(context.Background(), req.Name)
		if err != nil {
			return err
		}
		for {
			line, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				logrus.WithError(err).Error("Failed to receive log")
				return err
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
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.watchProcess(ctx, req, srv)
	})

	g.Go(func() error {
		return s.watchSpdkReplica(ctx, req, srv)
	})

	g.Go(func() error {
		return s.watchSpdkEngine(ctx, req, srv)
	})

	if err := g.Wait(); err != nil {
		return errors.Wrap(err, "failed to watch instances")
	}

	return nil
}

func (s *Server) watchSpdkReplica(ctx context.Context, req *emptypb.Empty, srv rpc.InstanceService_InstanceWatchServer) error {
	logrus.Info("Start watching SPDK replica")

	diskClient, err := client.NewDiskServiceClient("tcp://"+s.diskServiceAddress, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create DiskServiceClient")
	}
	defer diskClient.Close()

	notifier, err := diskClient.ReplicaWatch(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to create SPDK replica watch notifier")
	}

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Stop watching SPDK replica")
			return nil
		default:
			v, err := notifier.Recv()
			if err != nil {
				logrus.WithError(err).Error("Failed to receive next item in SPDK replica watch")
			} else {
				instance := replicaResponseToInstanceResponse(v)
				if err := srv.Send(instance); err != nil {
					return errors.Wrap(err, "failed to send instance response")
				}
			}
		}
	}
}

func (s *Server) watchSpdkEngine(ctx context.Context, req *emptypb.Empty, srv rpc.InstanceService_InstanceWatchServer) error {
	logrus.Info("Start watching SPDK engine")

	diskClient, err := client.NewDiskServiceClient("tcp://"+s.diskServiceAddress, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create DiskServiceClient")
	}
	defer diskClient.Close()

	notifier, err := diskClient.EngineWatch(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to create SPDK engine watch notifier")
	}

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Stop watching SPDK engine")
			return nil
		default:
			v, err := notifier.Recv()
			if err != nil {
				logrus.WithError(err).Error("Failed to receive next item in SPDK engine watch")
			} else {
				instance := engineResponseToInstanceResponse(v)
				if err := srv.Send(instance); err != nil {
					return errors.Wrap(err, "failed to send instance response")
				}
			}
		}
	}
}

func (s *Server) watchProcess(ctx context.Context, req *emptypb.Empty, srv rpc.InstanceService_InstanceWatchServer) error {
	logrus.Info("Start watching processes")

	pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create ProcessManagerClient")
	}
	defer pmClient.Close()

	notifier, err := pmClient.ProcessWatch(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to create process watch notifier")
	}

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Stop watching processes")
			return nil
		default:
			process, err := notifier.Recv()
			if err != nil {
				logrus.WithError(err).Error("Failed to receive next item in process watch")
			} else {
				instance := processResponseToInstanceResponse(process)
				if err := srv.Send(instance); err != nil {
					return errors.Wrap(err, "failed to send instance response")
				}
			}
		}
	}
}

func processResponseToInstanceResponse(p *rpc.ProcessResponse) *rpc.InstanceResponse {
	return &rpc.InstanceResponse{
		Spec: &rpc.InstanceSpec{
			Name: p.Spec.Name,
			// Leave Type empty. It will be determined in longhorn manager.
			Type:               BackendStoreDriverTypeUndefined,
			BackendStoreDriver: BackendStoreDriverTypeLonghorn,
			ProcessSpecific: &rpc.ProcessSpecific{
				Binary: p.Spec.Binary,
				Args:   p.Spec.Args,
			},
			PortCount: int32(p.Spec.PortCount),
			PortArgs:  p.Spec.PortArgs,
		},
		Status: &rpc.InstanceStatus{
			State:     p.Status.State,
			PortStart: p.Status.PortStart,
			PortEnd:   p.Status.PortEnd,
			ErrorMsg:  p.Status.ErrorMsg,
		},
		Deleted: p.Deleted,
	}
}

func replicaResponseToInstanceResponse(r *rpc.Replica) *rpc.InstanceResponse {
	return &rpc.InstanceResponse{
		Spec: &rpc.InstanceSpec{
			Name:               r.Name,
			Type:               types.InstanceTypeReplica,
			BackendStoreDriver: BackendStoreDriverTypeSpdkAio,
		},
		Status: &rpc.InstanceStatus{
			State:     types.ProcessStateRunning,
			PortStart: r.Port,
			PortEnd:   r.Port,
		},
	}
}

func engineResponseToInstanceResponse(e *rpc.Engine) *rpc.InstanceResponse {
	return &rpc.InstanceResponse{
		Spec: &rpc.InstanceSpec{
			Name:               e.Name,
			Type:               types.InstanceTypeEngine,
			BackendStoreDriver: BackendStoreDriverTypeSpdkAio,
		},
		Status: &rpc.InstanceStatus{
			State:     types.ProcessStateRunning,
			ErrorMsg:  "",
			PortStart: 0,
			PortEnd:   0,
		},
	}
}
