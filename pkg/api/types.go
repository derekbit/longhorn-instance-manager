package api

import (
	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
)

type Instance struct {
	Name      string   `json:"name"`
	Binary    string   `json:"binary"`
	Args      []string `json:"args"`
	PortCount int32    `json:"portCount"`
	PortArgs  []string `json:"portArgs"`

	InstanceStatus InstanceStatus `json:"instanceStatus"`

	Deleted bool `json:"deleted"`
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

func RPCToInstance(obj *rpc.InstanceResponse) *Instance {
	return &Instance{
		Name:           obj.Spec.Name,
		Binary:         obj.Spec.Process.Binary,
		Args:           obj.Spec.Process.Args,
		PortCount:      obj.Spec.PortCount,
		PortArgs:       obj.Spec.PortArgs,
		InstanceStatus: RPCToInstanceStatus(obj.Status),
	}
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

func (s *InstanceStream) Recv() (*Instance, error) {
	resp, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return RPCToInstance(resp), nil
}

func NewInstanceLogStream(stream rpc.InstanceService_InstanceLogClient) *InstanceLogStream {
	return &InstanceLogStream{
		stream,
	}
}

type InstanceLogStream struct {
	stream rpc.InstanceService_InstanceLogClient
}

func (s *InstanceLogStream) Recv() (string, error) {
	resp, err := s.stream.Recv()
	if err != nil {
		return "", err
	}
	return resp.Line, nil
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

func (s *ProcessStream) Recv() (*Process, error) {
	resp, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return RPCToProcess(resp), nil
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
