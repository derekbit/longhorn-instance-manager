package instance

import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/longhorn/longhorn-instance-manager/pkg/client"
	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/meta"
	"github.com/longhorn/longhorn-instance-manager/pkg/types"
)

type Server struct {
	logsDir                      string
	processManagerServiceAddress string
	diskServiceAddress           string
	spdkServiceAddress           string
	spdkEnabled                  bool
	shutdownCh                   chan error
	HealthChecker                HealthChecker
}

func NewServer(logsDir, processManagerServiceAddress, diskServiceAddress, spdkServiceAddress string, spdkEnabled bool, shutdownCh chan error) (*Server, error) {
	s := &Server{
		logsDir:                      logsDir,
		processManagerServiceAddress: processManagerServiceAddress,
		diskServiceAddress:           diskServiceAddress,
		spdkServiceAddress:           spdkServiceAddress,
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
	case rpc.BackendStoreDriver_longhorn:
		return s.processInstanceCreate(req)
	case rpc.BackendStoreDriver_spdk:
		return s.spdkInstanceCreate(req)
	default:
		return nil, grpcstatus.Errorf(grpccodes.InvalidArgument, "unknown backend store driver %v", req.Spec.BackendStoreDriver)
	}
}

func (s *Server) processInstanceCreate(req *rpc.InstanceCreateRequest) (*rpc.InstanceResponse, error) {
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
}

func (s *Server) spdkInstanceCreate(req *rpc.InstanceCreateRequest) (*rpc.InstanceResponse, error) {
	c, err := client.NewSPDKServiceClient("tcp://"+s.spdkServiceAddress, nil)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	switch req.Spec.Type {
	case types.InstanceTypeEngine:
		engine, err := c.EngineCreate(req.Spec.Name, req.Spec.VolumeName, req.Spec.SpdkSpecific.Frontend, req.Spec.SpdkSpecific.ReplicaAddressMap)
		if err != nil {
			return nil, err
		}
		return engineResponseToInstanceResponse(engine), nil
	case types.InstanceTypeReplica:
		replica, err := c.ReplicaCreate(req.Spec.Name, req.Spec.SpdkSpecific.DiskUuid, req.Spec.Size, req.ExposeRequired)
		if err != nil {
			return nil, err
		}
		return replicaResponseToInstanceResponse(replica), nil
	default:
		return nil, fmt.Errorf("unknown instance type %v", req.Spec.Type)
	}
}

func (s *Server) InstanceDelete(ctx context.Context, req *rpc.InstanceDeleteRequest) (*rpc.InstanceResponse, error) {
	logrus.WithFields(logrus.Fields{
		"name":               req.Name,
		"type":               req.Type,
		"backendStoreDriver": req.BackendStoreDriver,
		"diskUuid":           req.DiskUuid,
		"cleanupRequired":    req.CleanupRequired,
	}).Info("Deleting instance")

	switch req.BackendStoreDriver {
	case rpc.BackendStoreDriver_longhorn:
		return s.processInstanceDelete(req)
	case rpc.BackendStoreDriver_spdk:
		return s.spdkInstanceDelete(req)
	default:
		return nil, grpcstatus.Errorf(grpccodes.InvalidArgument, "unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (s *Server) processInstanceDelete(req *rpc.InstanceDeleteRequest) (*rpc.InstanceResponse, error) {
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
}

func (s *Server) spdkInstanceDelete(req *rpc.InstanceDeleteRequest) (*rpc.InstanceResponse, error) {
	c, err := client.NewSPDKServiceClient("tcp://"+s.spdkServiceAddress, nil)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	switch req.Type {
	case types.InstanceTypeEngine:
		err = c.EngineDelete(req.Name)
	case types.InstanceTypeReplica:
		err = c.ReplicaDelete(req.DiskUuid+"/"+req.Name, req.CleanupRequired)
	default:
		err = fmt.Errorf("unknown instance type %v", req.Type)
	}
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
}

func (s *Server) InstanceGet(ctx context.Context, req *rpc.InstanceGetRequest) (*rpc.InstanceResponse, error) {
	logrus.WithFields(logrus.Fields{
		"name":               req.Name,
		"type":               req.Type,
		"backendStoreDriver": req.BackendStoreDriver,
	}).Info("Getting instance")

	switch req.BackendStoreDriver {
	case rpc.BackendStoreDriver_longhorn:
		return s.processInstanceGet(req)
	case rpc.BackendStoreDriver_spdk:
		return s.spdkInstanceGet(req)
	default:
		return nil, grpcstatus.Errorf(grpccodes.InvalidArgument, "unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (s *Server) processInstanceGet(req *rpc.InstanceGetRequest) (*rpc.InstanceResponse, error) {
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
}

func (s *Server) spdkInstanceGet(req *rpc.InstanceGetRequest) (*rpc.InstanceResponse, error) {
	c, err := client.NewSPDKServiceClient("tcp://"+s.spdkServiceAddress, nil)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	switch req.Type {
	case types.InstanceTypeEngine:
		engine, err := c.EngineGet(req.Name)
		if err != nil {
			return nil, err
		}
		return engineResponseToInstanceResponse(engine), nil
	case types.InstanceTypeReplica:
		replica, err := c.ReplicaGet(req.DiskUuid + "/" + req.Name)
		if err != nil {
			return nil, err
		}
		return replicaResponseToInstanceResponse(replica), nil
	default:
		return nil, fmt.Errorf("unknown instance type %v", req.Type)
	}
}

func (s *Server) InstanceList(ctx context.Context, req *empty.Empty) (*rpc.InstanceListResponse, error) {
	logrus.WithFields(logrus.Fields{}).Info("Listing instances")

	instances := map[string]*rpc.InstanceResponse{}

	err := s.processInstanceList(instances)
	if err != nil {
		return nil, err
	}

	if s.spdkEnabled {
		err := s.spdkInstanceList(instances)
		if err != nil {
			return nil, err
		}
	}

	return &rpc.InstanceListResponse{
		Instances: instances,
	}, nil
}

func (s *Server) processInstanceList(instances map[string]*rpc.InstanceResponse) error {
	pmClient, err := client.NewProcessManagerClient("tcp://"+s.processManagerServiceAddress, nil)
	if err != nil {
		return err
	}
	defer pmClient.Close()

	processes, err := pmClient.ProcessList()
	if err != nil {
		return err
	}
	for _, process := range processes {
		instances[process.Spec.Name] = processResponseToInstanceResponse(process)
	}
	return nil
}

func (s *Server) spdkInstanceList(instances map[string]*rpc.InstanceResponse) error {
	c, err := client.NewSPDKServiceClient("tcp://"+s.spdkServiceAddress, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	replicas, err := c.ReplicaList()
	if err != nil {
		return err
	}
	for _, replica := range replicas {
		instances[replica.Name] = replicaResponseToInstanceResponse(replica)
	}

	engines, err := c.EngineList()
	if err != nil {
		return err
	}
	for _, engine := range engines {
		instances[engine.Name] = engineResponseToInstanceResponse(engine)
	}
	return nil
}

func (s *Server) InstanceReplace(ctx context.Context, req *rpc.InstanceReplaceRequest) (*rpc.InstanceResponse, error) {
	logrus.WithFields(logrus.Fields{
		"name":               req.Spec.Name,
		"type":               req.Spec.Type,
		"backendStoreDriver": req.Spec.BackendStoreDriver,
	}).Info("Replacing instance")

	switch req.Spec.BackendStoreDriver {
	case rpc.BackendStoreDriver_longhorn:
		return s.processInstanceReplace(req)
	case rpc.BackendStoreDriver_spdk:
		return s.spdkInstanceReplace(req)
	default:
		return nil, grpcstatus.Errorf(grpccodes.InvalidArgument, "unknown backend store driver %v", req.Spec.BackendStoreDriver)
	}
}

func (s *Server) processInstanceReplace(req *rpc.InstanceReplaceRequest) (*rpc.InstanceResponse, error) {
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
}

func (s *Server) spdkInstanceReplace(req *rpc.InstanceReplaceRequest) (*rpc.InstanceResponse, error) {
	return nil, fmt.Errorf("spdk instance replace is not supported")
}

func (s *Server) InstanceLog(req *rpc.InstanceLogRequest, srv rpc.InstanceService_InstanceLogServer) error {
	logrus.WithFields(logrus.Fields{
		"name":               req.Name,
		"type":               req.Type,
		"backendStoreDriver": req.BackendStoreDriver,
	}).Info("Getting instance log")

	switch req.BackendStoreDriver {
	case rpc.BackendStoreDriver_longhorn:
		return s.processInstanceLog(req, srv)
	case rpc.BackendStoreDriver_spdk:
		return s.spdkInstanceLog(req, srv)
	default:
		return grpcstatus.Errorf(grpccodes.InvalidArgument, "unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (s *Server) processInstanceLog(req *rpc.InstanceLogRequest, srv rpc.InstanceService_InstanceLogServer) error {
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
}

func (s *Server) spdkInstanceLog(req *rpc.InstanceLogRequest, srv rpc.InstanceService_InstanceLogServer) error {
	return fmt.Errorf("spdk instance log is not supported")
}

func (s *Server) InstanceWatch(req *empty.Empty, srv rpc.InstanceService_InstanceWatchServer) error {
	logrus.Info("Start watching instances")

	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.watchProcess(ctx, req, srv)
	})

	if s.spdkEnabled {
		g.Go(func() error {
			return s.watchSpdkReplica(ctx, req, srv)
		})

		g.Go(func() error {
			return s.watchSpdkEngine(ctx, req, srv)
		})
	}

	if err := g.Wait(); err != nil {
		return errors.Wrap(err, "failed to watch instances")
	}

	return nil
}

func (s *Server) watchSpdkReplica(ctx context.Context, req *emptypb.Empty, srv rpc.InstanceService_InstanceWatchServer) error {
	logrus.Info("Start watching SPDK replica")

	c, err := client.NewSPDKServiceClient("tcp://"+s.spdkServiceAddress, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create DiskServiceClient")
	}
	defer c.Close()

	notifier, err := c.ReplicaWatch(context.Background())
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
	c, err := client.NewSPDKServiceClient("tcp://"+s.spdkServiceAddress, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create DiskServiceClient")
	}
	defer c.Close()

	notifier, err := c.EngineWatch(context.Background())
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
			Type:               "",
			BackendStoreDriver: rpc.BackendStoreDriver_longhorn,
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
	logrus.Infof("Debug ====> port %v", r.Port)
	return &rpc.InstanceResponse{
		Spec: &rpc.InstanceSpec{
			Name:               r.Name,
			Type:               types.InstanceTypeReplica,
			BackendStoreDriver: rpc.BackendStoreDriver_spdk,
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
			BackendStoreDriver: rpc.BackendStoreDriver_spdk,
		},
		Status: &rpc.InstanceStatus{
			State:     types.ProcessStateRunning,
			ErrorMsg:  "",
			PortStart: e.Port,
			PortEnd:   e.Port,
		},
	}
}
