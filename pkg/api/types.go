package api

import (
	"github.com/golang/protobuf/ptypes/empty"
	spdkrpc "github.com/longhorn/longhorn-spdk-engine/proto/spdkrpc"

	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
)

type Instance struct {
	Name               string   `json:"name"`
	Type               string   `json:"type"`
	BackendStoreDriver string   `json:"backendStoreDriver"`
	Binary             string   `json:"binary"`
	Args               []string `json:"args"`
	PortCount          int32    `json:"portCount"`
	PortArgs           []string `json:"portArgs"`

	InstanceStatus InstanceStatus `json:"instanceStatus"`

	Deleted bool `json:"deleted"`
}

func RPCToInstance(obj *rpc.InstanceResponse) *Instance {
	instance := &Instance{
		Name:               obj.Spec.Name,
		Type:               obj.Spec.Type,
		BackendStoreDriver: obj.Spec.BackendStoreDriver.String(),
		PortCount:          obj.Spec.PortCount,
		PortArgs:           obj.Spec.PortArgs,
		InstanceStatus:     RPCToInstanceStatus(obj.Status),
	}

	if obj.Spec.ProcessSpecific != nil {
		instance.Binary = obj.Spec.ProcessSpecific.Binary
		instance.Args = obj.Spec.ProcessSpecific.Args
	}

	return instance
}

func RPCToInstanceList(obj *rpc.InstanceListResponse) map[string]*Instance {
	ret := map[string]*Instance{}
	for name, p := range obj.Instances {
		ret[name] = RPCToInstance(p)
	}
	return ret
}

type InstanceStatus struct {
	State     string `json:"state"`
	ErrorMsg  string `json:"errorMsg"`
	PortStart int32  `json:"portStart"`
	PortEnd   int32  `json:"portEnd"`
}

func RPCToInstanceStatus(obj *rpc.InstanceStatus) InstanceStatus {
	return InstanceStatus{
		State:     obj.State,
		ErrorMsg:  obj.ErrorMsg,
		PortStart: obj.PortStart,
		PortEnd:   obj.PortEnd,
	}
}

type InstanceStream struct {
	stream rpc.InstanceService_InstanceWatchClient
}

func NewInstanceStream(stream rpc.InstanceService_InstanceWatchClient) *InstanceStream {
	return &InstanceStream{
		stream,
	}
}

func (s *InstanceStream) Recv() (*empty.Empty, error) {
	return s.stream.Recv()
}

type Process struct {
	Name      string   `json:"name"`
	Binary    string   `json:"binary"`
	Args      []string `json:"args"`
	PortCount int32    `json:"portCount"`
	PortArgs  []string `json:"portArgs"`

	ProcessStatus ProcessStatus `json:"processStatus"`

	Deleted bool `json:"deleted"`
}

func RPCToProcess(obj *rpc.ProcessResponse) *Process {
	return &Process{
		Name:          obj.Spec.Name,
		Binary:        obj.Spec.Binary,
		Args:          obj.Spec.Args,
		PortCount:     obj.Spec.PortCount,
		PortArgs:      obj.Spec.PortArgs,
		ProcessStatus: RPCToProcessStatus(obj.Status),
	}
}

func RPCToProcessList(obj *rpc.ProcessListResponse) map[string]*Process {
	ret := map[string]*Process{}
	for name, p := range obj.Processes {
		ret[name] = RPCToProcess(p)
	}
	return ret
}

type ProcessStatus struct {
	State     string `json:"state"`
	ErrorMsg  string `json:"errorMsg"`
	PortStart int32  `json:"portStart"`
	PortEnd   int32  `json:"portEnd"`
}

func RPCToProcessStatus(obj *rpc.ProcessStatus) ProcessStatus {
	return ProcessStatus{
		State:     obj.State,
		ErrorMsg:  obj.ErrorMsg,
		PortStart: obj.PortStart,
		PortEnd:   obj.PortEnd,
	}
}

type ProcessStream struct {
	stream rpc.ProcessManagerService_ProcessWatchClient
}

func NewProcessStream(stream rpc.ProcessManagerService_ProcessWatchClient) *ProcessStream {
	return &ProcessStream{
		stream,
	}
}

func (s *ProcessStream) Recv() (*rpc.ProcessResponse, error) {
	return s.stream.Recv()
}

type ReplicaStream struct {
	stream spdkrpc.SPDKService_ReplicaWatchClient
}

func NewReplicaStream(stream spdkrpc.SPDKService_ReplicaWatchClient) *ReplicaStream {
	return &ReplicaStream{
		stream,
	}
}

func (s *ReplicaStream) Recv() (*empty.Empty, error) {
	return s.stream.Recv()
}

type EngineStream struct {
	stream spdkrpc.SPDKService_EngineWatchClient
}

func NewEngineStream(stream spdkrpc.SPDKService_EngineWatchClient) *EngineStream {
	return &EngineStream{
		stream,
	}
}

func (s *EngineStream) Recv() (*empty.Empty, error) {
	return s.stream.Recv()
}

func NewLogStream(stream rpc.ProcessManagerService_ProcessLogClient) *LogStream {
	return &LogStream{
		stream,
	}
}

type LogStream struct {
	stream rpc.ProcessManagerService_ProcessLogClient
}

func (s *LogStream) Recv() (string, error) {
	resp, err := s.stream.Recv()
	if err != nil {
		return "", err
	}
	return resp.Line, nil
}

type DiskInfo struct {
	ID          string
	UUID        string
	Path        string
	Type        string
	TotalSize   int64
	FreeSize    int64
	TotalBlocks int64
	FreeBlocks  int64
	BlockSize   int64
	ClusterSize int64
}

type ReplicaInstance struct {
	UUID       string
	SpecSize   uint64
	ActualSize uint64
}
