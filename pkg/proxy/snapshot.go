package proxy

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/types"

	eclient "github.com/longhorn/longhorn-engine/pkg/controller/client"
	esync "github.com/longhorn/longhorn-engine/pkg/sync"
	eptypes "github.com/longhorn/longhorn-engine/proto/ptypes"
)

func (p *Proxy) VolumeSnapshot(ctx context.Context, req *rpc.EngineVolumeSnapshotRequest) (resp *rpc.EngineVolumeSnapshotProxyResponse, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.ProxyEngineRequest.Address,
		"engineName":         req.ProxyEngineRequest.EngineName,
		"backendStoreDriver": req.ProxyEngineRequest.BackendStoreDriver,
	})
	log.Infof("Snapshotting volume: snapshot %v", req.SnapshotVolume.Name)

	switch req.ProxyEngineRequest.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.volumeSnapshotFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.volumeSnapshotFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.ProxyEngineRequest.BackendStoreDriver)
	}
}

func (p *Proxy) volumeSnapshotFromEngine(ctx context.Context, req *rpc.EngineVolumeSnapshotRequest) (resp *rpc.EngineVolumeSnapshotProxyResponse, err error) {
	c, err := eclient.NewControllerClient(req.ProxyEngineRequest.Address)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	recv, err := c.VolumeSnapshot(req.SnapshotVolume.Name, req.SnapshotVolume.Labels)
	if err != nil {
		return nil, err
	}

	return &rpc.EngineVolumeSnapshotProxyResponse{
		Snapshot: &eptypes.VolumeSnapshotReply{
			Name: recv,
		},
	}, nil
}

func (p *Proxy) volumeSnapshotFromSpdkService(ctx context.Context, req *rpc.EngineVolumeSnapshotRequest) (resp *rpc.EngineVolumeSnapshotProxyResponse, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *Proxy) SnapshotList(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineSnapshotListProxyResponse, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.Address,
		"engineName":         req.EngineName,
		"backendStoreDriver": req.BackendStoreDriver,
	})
	log.Trace("Listing snapshots")

	switch req.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.snapshotListFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.snapshotListFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (p *Proxy) snapshotListFromEngine(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineSnapshotListProxyResponse, err error) {
	c, err := eclient.NewControllerClient(req.Address)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	recv, err := c.ReplicaList()
	if err != nil {
		return nil, err
	}

	snapshotsDiskInfo, err := esync.GetSnapshotsInfo(recv)
	if err != nil {
		return nil, err
	}

	resp = &rpc.EngineSnapshotListProxyResponse{
		Disks: map[string]*rpc.EngineSnapshotDiskInfo{},
	}
	for k, v := range snapshotsDiskInfo {
		resp.Disks[k] = &rpc.EngineSnapshotDiskInfo{
			Name:        v.Name,
			Parent:      v.Parent,
			Children:    v.Children,
			Removed:     v.Removed,
			UserCreated: v.UserCreated,
			Created:     v.Created,
			Size:        v.Size,
			Labels:      v.Labels,
		}
	}

	return resp, nil
}

func (p *Proxy) snapshotListFromSpdkService(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineSnapshotListProxyResponse, err error) {
	return &rpc.EngineSnapshotListProxyResponse{
		Disks: map[string]*rpc.EngineSnapshotDiskInfo{},
	}, nil
}

func (p *Proxy) SnapshotClone(ctx context.Context, req *rpc.EngineSnapshotCloneRequest) (resp *empty.Empty, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.ProxyEngineRequest.Address,
		"engineName":         req.ProxyEngineRequest.EngineName,
		"backendStoreDriver": req.ProxyEngineRequest.BackendStoreDriver,
	})
	log.Infof("Cloning snapshot from %v to %v", req.FromController, req.ProxyEngineRequest.Address)

	switch req.ProxyEngineRequest.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.snapshotCloneFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.snapshotCloneFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.ProxyEngineRequest.BackendStoreDriver)
	}
}

func (p *Proxy) snapshotCloneFromEngine(ctx context.Context, req *rpc.EngineSnapshotCloneRequest) (resp *empty.Empty, err error) {
	cFrom, err := eclient.NewControllerClient(req.FromController)
	if err != nil {
		return nil, err
	}
	defer cFrom.Close()

	cTo, err := eclient.NewControllerClient(req.ProxyEngineRequest.Address)
	if err != nil {
		return nil, err
	}
	defer cTo.Close()

	err = esync.CloneSnapshot(cTo, cFrom, req.SnapshotName, req.ExportBackingImageIfExist, int(req.FileSyncHttpClientTimeout))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *Proxy) snapshotCloneFromSpdkService(ctx context.Context, req *rpc.EngineSnapshotCloneRequest) (resp *empty.Empty, err error) {
	return &empty.Empty{}, nil
}

func (p *Proxy) SnapshotCloneStatus(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineSnapshotCloneStatusProxyResponse, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.Address,
		"engineName":         req.EngineName,
		"backendStoreDriver": req.BackendStoreDriver,
	})
	log.Trace("Getting snapshot clone status")

	switch req.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.snapshotCloneStatusFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.snapshotCloneStatusFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (p *Proxy) snapshotCloneStatusFromEngine(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineSnapshotCloneStatusProxyResponse, err error) {
	c, err := eclient.NewControllerClient(req.Address)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	recv, err := esync.CloneStatus(c)
	if err != nil {
		return nil, err
	}

	resp = &rpc.EngineSnapshotCloneStatusProxyResponse{
		Status: map[string]*eptypes.SnapshotCloneStatusResponse{},
	}
	for k, v := range recv {
		resp.Status[k] = &eptypes.SnapshotCloneStatusResponse{
			IsCloning:          v.IsCloning,
			Error:              v.Error,
			Progress:           int32(v.Progress),
			State:              v.State,
			FromReplicaAddress: v.FromReplicaAddress,
			SnapshotName:       v.SnapshotName,
		}
	}

	return resp, nil
}

func (p *Proxy) snapshotCloneStatusFromSpdkService(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineSnapshotCloneStatusProxyResponse, err error) {
	return &rpc.EngineSnapshotCloneStatusProxyResponse{
		Status: map[string]*eptypes.SnapshotCloneStatusResponse{},
	}, nil
}

func (p *Proxy) SnapshotRevert(ctx context.Context, req *rpc.EngineSnapshotRevertRequest) (resp *empty.Empty, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.ProxyEngineRequest.Address,
		"engineName":         req.ProxyEngineRequest.EngineName,
		"backendStoreDriver": req.ProxyEngineRequest.BackendStoreDriver,
	})
	log.Infof("Reverting snapshot %v", req.Name)

	switch req.ProxyEngineRequest.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.snapshotRevertFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.snapshotRevertFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.ProxyEngineRequest.BackendStoreDriver)
	}
}

func (p *Proxy) snapshotRevertFromEngine(ctx context.Context, req *rpc.EngineSnapshotRevertRequest) (resp *empty.Empty, err error) {
	c, err := eclient.NewControllerClient(req.ProxyEngineRequest.Address)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	if err := c.VolumeRevert(req.Name); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *Proxy) snapshotRevertFromSpdkService(ctx context.Context, req *rpc.EngineSnapshotRevertRequest) (resp *empty.Empty, err error) {
	return &empty.Empty{}, nil
}

func (p *Proxy) SnapshotPurge(ctx context.Context, req *rpc.EngineSnapshotPurgeRequest) (resp *empty.Empty, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.ProxyEngineRequest.Address,
		"engineName":         req.ProxyEngineRequest.EngineName,
		"backendStoreDriver": req.ProxyEngineRequest.BackendStoreDriver,
	})
	log.Info("Purging snapshots")

	switch req.ProxyEngineRequest.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.snapshotPurgeFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.snapshotPurgeFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.ProxyEngineRequest.BackendStoreDriver)
	}
}

func (p *Proxy) snapshotPurgeFromEngine(ctx context.Context, req *rpc.EngineSnapshotPurgeRequest) (resp *empty.Empty, err error) {
	task, err := esync.NewTask(ctx, req.ProxyEngineRequest.Address)
	if err != nil {
		return nil, err
	}

	if err := task.PurgeSnapshots(req.SkipIfInProgress); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *Proxy) snapshotPurgeFromSpdkService(ctx context.Context, req *rpc.EngineSnapshotPurgeRequest) (resp *empty.Empty, err error) {
	return &empty.Empty{}, nil
}

func (p *Proxy) SnapshotPurgeStatus(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineSnapshotPurgeStatusProxyResponse, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.Address,
		"engineName":         req.EngineName,
		"backendStoreDriver": req.BackendStoreDriver,
	})
	log.Trace("Getting snapshot purge status")

	switch req.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.snapshotPurgeStatusFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.snapshotPurgeStatusFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.BackendStoreDriver)
	}
}

func (p *Proxy) snapshotPurgeStatusFromEngine(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineSnapshotPurgeStatusProxyResponse, err error) {
	task, err := esync.NewTask(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	recv, err := task.PurgeSnapshotStatus()
	if err != nil {
		return nil, err
	}

	resp = &rpc.EngineSnapshotPurgeStatusProxyResponse{
		Status: map[string]*eptypes.SnapshotPurgeStatusResponse{},
	}
	for k, v := range recv {
		resp.Status[k] = &eptypes.SnapshotPurgeStatusResponse{
			IsPurging: v.IsPurging,
			Error:     v.Error,
			Progress:  int32(v.Progress),
			State:     v.State,
		}
	}

	return resp, nil
}

func (p *Proxy) snapshotPurgeStatusFromSpdkService(ctx context.Context, req *rpc.ProxyEngineRequest) (resp *rpc.EngineSnapshotPurgeStatusProxyResponse, err error) {
	return &rpc.EngineSnapshotPurgeStatusProxyResponse{
		Status: map[string]*eptypes.SnapshotPurgeStatusResponse{},
	}, nil
}

func (p *Proxy) SnapshotRemove(ctx context.Context, req *rpc.EngineSnapshotRemoveRequest) (resp *empty.Empty, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.ProxyEngineRequest.Address,
		"engineName":         req.ProxyEngineRequest.EngineName,
		"backendStoreDriver": req.ProxyEngineRequest.BackendStoreDriver,
	})
	log.Infof("Removing snapshots %v", req.Names)

	switch req.ProxyEngineRequest.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.snapshotRemoveFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.snapshotRemoveFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.ProxyEngineRequest.BackendStoreDriver)
	}
}

func (p *Proxy) snapshotRemoveFromEngine(ctx context.Context, req *rpc.EngineSnapshotRemoveRequest) (resp *empty.Empty, err error) {
	task, err := esync.NewTask(ctx, req.ProxyEngineRequest.Address)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for _, name := range req.Names {
		if err := task.DeleteSnapshot(name); err != nil {
			lastErr = err
			logrus.WithError(err).Warnf("Failed to delete %s", name)
		}
	}

	return &empty.Empty{}, lastErr
}

func (p *Proxy) snapshotRemoveFromSpdkService(ctx context.Context, req *rpc.EngineSnapshotRemoveRequest) (resp *empty.Empty, err error) {
	return &empty.Empty{}, nil
}

func (p *Proxy) SnapshotHash(ctx context.Context, req *rpc.EngineSnapshotHashRequest) (resp *empty.Empty, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.ProxyEngineRequest.Address,
		"engineName":         req.ProxyEngineRequest.EngineName,
		"backendStoreDriver": req.ProxyEngineRequest.BackendStoreDriver,
	})
	log.Infof("Hashing snapshot %v with rehash %v", req.SnapshotName, req.Rehash)

	switch req.ProxyEngineRequest.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.snapshotHashFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.snapshotHashFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.ProxyEngineRequest.BackendStoreDriver)
	}
}

func (p *Proxy) snapshotHashFromEngine(ctx context.Context, req *rpc.EngineSnapshotHashRequest) (resp *empty.Empty, err error) {
	task, err := esync.NewTask(ctx, req.ProxyEngineRequest.Address)
	if err != nil {
		return nil, err
	}

	if err := task.HashSnapshot(req.SnapshotName, req.Rehash); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (p *Proxy) snapshotHashFromSpdkService(ctx context.Context, req *rpc.EngineSnapshotHashRequest) (resp *empty.Empty, err error) {
	return &empty.Empty{}, nil
}

func (p *Proxy) SnapshotHashStatus(ctx context.Context, req *rpc.EngineSnapshotHashStatusRequest) (resp *rpc.EngineSnapshotHashStatusProxyResponse, err error) {
	log := logrus.WithFields(logrus.Fields{
		"serviceURL":         req.ProxyEngineRequest.Address,
		"engineName":         req.ProxyEngineRequest.EngineName,
		"backendStoreDriver": req.ProxyEngineRequest.BackendStoreDriver,
	})
	log.Trace("Getting snapshot hash status")

	switch req.ProxyEngineRequest.BackendStoreDriver {
	case types.BackendStoreDriverTypeLonghorn:
		return p.snapshotHashStatusFromEngine(ctx, req)
	case types.BackendStoreDriverTypeSpdkAio:
		return p.snapshotHashStatusFromSpdkService(ctx, req)
	default:
		return nil, fmt.Errorf("unknown backend store driver %v", req.ProxyEngineRequest.BackendStoreDriver)
	}
}

func (p *Proxy) snapshotHashStatusFromEngine(ctx context.Context, req *rpc.EngineSnapshotHashStatusRequest) (resp *rpc.EngineSnapshotHashStatusProxyResponse, err error) {
	task, err := esync.NewTask(ctx, req.ProxyEngineRequest.Address)
	if err != nil {
		return nil, err
	}

	recv, err := task.HashSnapshotStatus(req.SnapshotName)
	if err != nil {
		return nil, err
	}

	resp = &rpc.EngineSnapshotHashStatusProxyResponse{
		Status: map[string]*eptypes.SnapshotHashStatusResponse{},
	}
	for k, v := range recv {
		resp.Status[k] = &eptypes.SnapshotHashStatusResponse{
			State:             v.State,
			Checksum:          v.Checksum,
			Error:             v.Error,
			SilentlyCorrupted: v.SilentlyCorrupted,
		}
	}

	return resp, nil
}

func (p *Proxy) snapshotHashStatusFromSpdkService(ctx context.Context, req *rpc.EngineSnapshotHashStatusRequest) (resp *rpc.EngineSnapshotHashStatusProxyResponse, err error) {
	return &rpc.EngineSnapshotHashStatusProxyResponse{
		Status: map[string]*eptypes.SnapshotHashStatusResponse{},
	}, nil
}
