package proxy

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/longhorn/longhorn-instance-manager/pkg/client"
	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/types"

	eclient "github.com/longhorn/longhorn-engine/pkg/controller/client"
	eptypes "github.com/longhorn/longhorn-engine/proto/ptypes"
)

func (p *Proxy) VolumeGet(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineVolumeGetProxyResponse, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.Address,
		"engineName":         req.EngineName,
		"backendStoreDriver": req.BackendStoreDriver,
	})
	log.Trace("Getting volume")

	switch req.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.volumeGetFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.volumeGetFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (p *Proxy) volumeGetFromEngine(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineVolumeGetProxyResponse, err error) {
	c, err := eclient.NewControllerClient(req.Address)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	recv, err := c.VolumeGet()
	if err != nil {
		return nil, err
	}

	return &rpc.EngineVolumeGetProxyResponse{
		Volume: &eptypes.Volume{
			Name:                      recv.Name,
			Size:                      recv.Size,
			ReplicaCount:              int32(recv.ReplicaCount),
			Endpoint:                  recv.Endpoint,
			Frontend:                  recv.Frontend,
			FrontendState:             recv.FrontendState,
			IsExpanding:               recv.IsExpanding,
			LastExpansionError:        recv.LastExpansionError,
			LastExpansionFailedAt:     recv.LastExpansionFailedAt,
			UnmapMarkSnapChainRemoved: recv.UnmapMarkSnapChainRemoved,
		},
	}, nil
}

func (p *Proxy) volumeGetFromSpdkService(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineVolumeGetProxyResponse, err error) {
	// TODO: Should connect to SPDK service
	c, err := client.NewDiskServiceClient("tcp://"+p.diskServiceAddress, nil)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	recv, err := c.EngineGet(req.EngineName)
	if err != nil {
		return nil, err
	}

	return &rpc.EngineVolumeGetProxyResponse{
		Volume: &eptypes.Volume{
			Name:                      recv.Name,
			Size:                      int64(recv.Size),
			ReplicaCount:              int32(len(recv.ReplicaAddressMap)),
			Endpoint:                  recv.Endpoint,
			Frontend:                  recv.Frontend,
			FrontendState:             recv.FrontendState,
			IsExpanding:               false,
			LastExpansionError:        "",
			LastExpansionFailedAt:     "",
			UnmapMarkSnapChainRemoved: false,
		},
	}, nil
}

func (p *Proxy) VolumeExpand(ctx context.Context, req *rpc.EngineVolumeExpandRequest) (resp *empty.Empty, err error) {
	log := logrus.WithFields(logrus.Fields{"serviceURL": req.ProxyEngineRequest.Address})
	log.Infof("Expanding volume to size %v", req.Expand.Size)

	c, err := eclient.NewControllerClient(req.ProxyEngineRequest.Address)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	err = c.VolumeExpand(req.Expand.Size)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *Proxy) VolumeFrontendStart(ctx context.Context, req *rpc.EngineVolumeFrontendStartRequest) (resp *empty.Empty, err error) {
	log := logrus.WithFields(logrus.Fields{"serviceURL": req.ProxyEngineRequest.Address})
	log.Infof("Starting volume frontend %v", req.FrontendStart.Frontend)

	c, err := eclient.NewControllerClient(req.ProxyEngineRequest.Address)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	err = c.VolumeFrontendStart(req.FrontendStart.Frontend)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *Proxy) VolumeFrontendShutdown(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *empty.Empty, err error) {
	log := logrus.WithFields(logrus.Fields{"serviceURL": req.Address})
	log.Info("Shutting down volume frontend")

	c, err := eclient.NewControllerClient(req.Address)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	err = c.VolumeFrontendShutdown()
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *Proxy) VolumeUnmapMarkSnapChainRemovedSet(ctx context.Context, req *rpc.EngineVolumeUnmapMarkSnapChainRemovedSetRequest) (resp *empty.Empty, err error) {
	log := logrus.WithFields(logrus.Fields{"serviceURL": req.ProxyEngineRequest.Address})
	log.Infof("Setting volume flag UnmapMarkSnapChainRemoved to %v", req.UnmapMarkSnap.Enabled)

	c, err := eclient.NewControllerClient(req.ProxyEngineRequest.Address)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	err = c.VolumeUnmapMarkSnapChainRemovedSet(req.UnmapMarkSnap.Enabled)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
