package types

import "sync"

type Mode string

const (
	ModeWO  = Mode("WO")
	ModeRW  = Mode("RW")
	ModeERR = Mode("ERR")
)

const (
	FrontendSPDKTCPNvmf     = "spdk-tcp-nvmf"
	FrontendSPDKTCPBlockdev = "spdk-tcp-blockdev"
	FrontendEmpty           = ""

	DefaultEngineReservedPortCount  = 1
	DefaultReplicaReservedPortCount = 5
)

type InstanceState string

const (
	InstanceStatePending = "pending"
	InstanceStateStopped = "stopped"
	InstanceStateRunning = "running"
	InstanceStateError   = "error"
)

type InstanceType string

const (
	InstanceTypeReplica = InstanceType("replica")
	InstanceTypeEngine  = InstanceType("engine")
)

const SPDKServicePort = 8504

var spdkLock sync.Mutex
