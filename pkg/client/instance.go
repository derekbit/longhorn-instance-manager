package client

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/longhorn/longhorn-instance-manager/pkg/api"
	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/meta"
	"github.com/longhorn/longhorn-instance-manager/pkg/types"
	"github.com/longhorn/longhorn-instance-manager/pkg/util"
)

type InstanceServiceContext struct {
	cc      *grpc.ClientConn
	service rpc.InstanceServiceClient
}

func (c InstanceServiceContext) Close() error {
	if c.cc == nil {
		return nil
	}
	return c.cc.Close()
}

func (c *InstanceServiceClient) getControllerServiceClient() rpc.InstanceServiceClient {
	return c.service
}

type InstanceServiceClient struct {
	serviceURL string
	tlsConfig  *tls.Config
	InstanceServiceContext
}

func NewInstanceServiceClient(serviceURL string, tlsConfig *tls.Config) (*InstanceServiceClient, error) {
	getInstanceServiceContext := func(serviceUrl string, tlsConfig *tls.Config) (InstanceServiceContext, error) {
		connection, err := util.Connect(serviceUrl, tlsConfig)
		if err != nil {
			return InstanceServiceContext{}, errors.Wrapf(err, "cannot connect to Instance Service %v", serviceUrl)
		}

		return InstanceServiceContext{
			cc:      connection,
			service: rpc.NewInstanceServiceClient(connection),
		}, nil
	}

	serviceContext, err := getInstanceServiceContext(serviceURL, tlsConfig)
	if err != nil {
		return nil, err
	}

	return &InstanceServiceClient{
		serviceURL:             serviceURL,
		tlsConfig:              tlsConfig,
		InstanceServiceContext: serviceContext,
	}, nil
}

func NewInstanceServiceClientWithTLS(serviceURL, caFile, certFile, keyFile, peerName string) (*InstanceServiceClient, error) {
	tlsConfig, err := util.LoadClientTLS(caFile, certFile, keyFile, peerName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load tls key pair from file")
	}

	return NewInstanceServiceClient(serviceURL, tlsConfig)
}

func (c *InstanceServiceClient) InstanceCreate(backendStoreDriver, name, volumeName, instanceType, diskUUID string, size uint64,
	binary string, args []string, frontend, address string, replicaAddressMap map[string]string, portCount int, portArgs []string, exposeRequired bool) (*api.Instance, error) {
	if name == "" || instanceType == "" {
		return nil, fmt.Errorf("failed to create instance: missing required parameter")
	}

	driver, ok := rpc.BackendStoreDriver_value[backendStoreDriver]
	if !ok {
		return nil, fmt.Errorf("failed to delete instance: invalid backend store driver %v", backendStoreDriver)
	}

	client := c.getControllerServiceClient()
	ctx, cancel := context.WithTimeout(context.Background(), types.GRPCServiceTimeout)
	defer cancel()

	p, err := client.InstanceCreate(ctx, &rpc.InstanceCreateRequest{
		Spec: &rpc.InstanceSpec{
			Name:               name,
			Type:               instanceType,
			BackendStoreDriver: rpc.BackendStoreDriver(driver),
			VolumeName:         volumeName,
			Size:               size,

			PortCount: int32(portCount),
			PortArgs:  portArgs,

			ProcessSpecific: &rpc.ProcessSpecific{
				Binary: binary,
				Args:   args,
			},

			SpdkSpecific: &rpc.SpdkSpecific{
				Address:           address,
				Frontend:          frontend,
				ReplicaAddressMap: replicaAddressMap,
				DiskUuid:          diskUUID,
			},
		},
		ExposeRequired: exposeRequired,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create instance")
	}

	return api.RPCToInstance(p), nil
}

func (c *InstanceServiceClient) InstanceDelete(backendStoreDriver, name, instanceType, diskUUID string, cleanupRequired bool) (*api.Instance, error) {
	if name == "" {
		return nil, fmt.Errorf("failed to delete instance: missing required parameter name")
	}

	driver, ok := rpc.BackendStoreDriver_value[backendStoreDriver]
	if !ok {
		return nil, fmt.Errorf("failed to delete instance: invalid backend store driver %v", backendStoreDriver)
	}

	client := c.getControllerServiceClient()
	ctx, cancel := context.WithTimeout(context.Background(), types.GRPCServiceTimeout)
	defer cancel()

	p, err := client.InstanceDelete(ctx, &rpc.InstanceDeleteRequest{
		Name:               name,
		Type:               instanceType,
		BackendStoreDriver: rpc.BackendStoreDriver(driver),
		DiskUuid:           diskUUID,
		CleanupRequired:    cleanupRequired,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete instance %v", name)
	}
	return api.RPCToInstance(p), nil
}

func (c *InstanceServiceClient) InstanceGet(backendStoreDriver, name, instanceType, diskUUID string) (*api.Instance, error) {
	if name == "" {
		return nil, fmt.Errorf("failed to get instance: missing required parameter name")
	}

	driver, ok := rpc.BackendStoreDriver_value[backendStoreDriver]
	if !ok {
		return nil, fmt.Errorf("failed to get instance: invalid backend store driver %v", backendStoreDriver)
	}

	client := c.getControllerServiceClient()
	ctx, cancel := context.WithTimeout(context.Background(), types.GRPCServiceTimeout)
	defer cancel()

	p, err := client.InstanceGet(ctx, &rpc.InstanceGetRequest{
		Name:               name,
		Type:               instanceType,
		BackendStoreDriver: rpc.BackendStoreDriver(driver),
		DiskUuid:           diskUUID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get instance %v", name)
	}
	return api.RPCToInstance(p), nil
}

func (c *InstanceServiceClient) InstanceList() (map[string]*api.Instance, error) {
	client := c.getControllerServiceClient()
	ctx, cancel := context.WithTimeout(context.Background(), types.GRPCServiceTimeout)
	defer cancel()

	instances, err := client.InstanceList(ctx, &empty.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list instances")
	}
	return api.RPCToInstanceList(instances), nil
}

func (c *InstanceServiceClient) InstanceLog(ctx context.Context, backendStoreDriver, name, instanceType string) (*api.LogStream, error) {
	if name == "" {
		return nil, fmt.Errorf("failed to get instance: missing required parameter name")
	}

	driver, ok := rpc.BackendStoreDriver_value[backendStoreDriver]
	if !ok {
		return nil, fmt.Errorf("failed to log instance: invalid backend store driver %v", backendStoreDriver)
	}

	client := c.getControllerServiceClient()
	stream, err := client.InstanceLog(ctx, &rpc.InstanceLogRequest{
		Name:               name,
		Type:               instanceType,
		BackendStoreDriver: rpc.BackendStoreDriver(driver),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get instance log of %v", name)
	}
	return api.NewLogStream(stream), nil
}

func (c *InstanceServiceClient) InstanceWatch(ctx context.Context) (*api.InstanceStream, error) {
	client := c.getControllerServiceClient()
	stream, err := client.InstanceWatch(ctx, &empty.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open instance update stream")
	}

	return api.NewInstanceStream(stream), nil
}

func (c *InstanceServiceClient) InstanceReplace(backendStoreDriver, name, instanceType, binary string, portCount int, args, portArgs []string, terminateSignal string) (*api.Instance, error) {
	if name == "" || binary == "" {
		return nil, fmt.Errorf("failed to replace instance: missing required parameter")
	}

	driver, ok := rpc.BackendStoreDriver_value[backendStoreDriver]
	if !ok {
		return nil, fmt.Errorf("failed to replace instance: invalid backend store driver %v", backendStoreDriver)
	}

	if terminateSignal != "SIGHUP" {
		return nil, fmt.Errorf("unsupported terminate signal %v", terminateSignal)
	}

	client := c.getControllerServiceClient()
	ctx, cancel := context.WithTimeout(context.Background(), types.GRPCServiceTimeout)
	defer cancel()

	p, err := client.InstanceReplace(ctx, &rpc.InstanceReplaceRequest{
		Spec: &rpc.InstanceSpec{
			Name:               name,
			Type:               instanceType,
			BackendStoreDriver: rpc.BackendStoreDriver(driver),
			ProcessSpecific: &rpc.ProcessSpecific{
				Binary: binary,
				Args:   args,
			},
			PortCount: int32(portCount),
			PortArgs:  portArgs,
		},
		TerminateSignal: terminateSignal,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to replace instance")
	}
	return api.RPCToInstance(p), nil
}

func (c *InstanceServiceClient) VersionGet() (*meta.VersionOutput, error) {
	client := c.getControllerServiceClient()
	ctx, cancel := context.WithTimeout(context.Background(), types.GRPCServiceTimeout)
	defer cancel()

	resp, err := client.VersionGet(ctx, &empty.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get version")
	}

	return &meta.VersionOutput{
		Version:   resp.Version,
		GitCommit: resp.GitCommit,
		BuildDate: resp.BuildDate,

		InstanceManagerAPIVersion:    int(resp.InstanceManagerAPIVersion),
		InstanceManagerAPIMinVersion: int(resp.InstanceManagerAPIMinVersion),

		InstanceManagerProxyAPIVersion:    int(resp.InstanceManagerProxyAPIVersion),
		InstanceManagerProxyAPIMinVersion: int(resp.InstanceManagerProxyAPIMinVersion),

		InstanceManagerDiskServiceAPIVersion:    int(resp.InstanceManagerDiskServiceAPIVersion),
		InstanceManagerDiskServiceAPIMinVersion: int(resp.InstanceManagerDiskServiceAPIMinVersion),
	}, nil
}