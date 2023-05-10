package nvme

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/longhorn/nsfilelock"

	"github.com/longhorn/go-spdk-helper/pkg/util"
)

const (
	LockFile    = "/var/run/longhorn-spdk.lock"
	LockTimeout = 120 * time.Second

	RetryCounts   = 5
	RetryInterval = 3 * time.Second

	HostProc = "/host/proc"
)

type Initiator struct {
	Name               string
	SubsystemNQN       string
	TransportAddress   string
	TransportServiceID string

	Endpoint       string
	ControllerName string
	NamespaceName  string
	dev            *util.KernelDevice
	isUp           bool

	// ControllerLossTimeout int64
	// FastIOFailTimeout     int64

	hostProc string
	executor util.Executor

	logger logrus.FieldLogger
}

func NewInitiator(name, subsystemNQN, hostProc string) (*Initiator, error) {
	if name == "" || subsystemNQN == "" {
		return nil, fmt.Errorf("empty name or subsystem for initiator creation")
	}

	// If transportAddress or transportServiceID is empty, the initiator is still valid for stopping

	var executor util.Executor
	if hostProc != "" {
		ne, err := util.NewNamespaceExecutor(util.GetHostNamespacePath(hostProc))
		if err != nil {
			return nil, err
		}
		executor = ne
	} else {
		executor = util.NewTimeoutExecutor(util.CmdTimeout)
	}

	if err := CheckForNVMeCliExistence(executor); err != nil {
		return nil, err
	}

	return &Initiator{
		Name:         name,
		SubsystemNQN: subsystemNQN,

		hostProc: hostProc,
		executor: executor,

		logger: logrus.WithFields(logrus.Fields{
			"name":         name,
			"subsystemNQN": subsystemNQN,
		}),
	}, nil
}

func (i *Initiator) Start(transportAddress, transportServiceID string) (err error) {
	if transportAddress == "" || transportServiceID == "" {
		return fmt.Errorf("empty TransportAddress %s and TransportServiceID %s for initiator %s start", transportAddress, transportServiceID, i.Name)
	}

	if i.hostProc != "" {
		lock := nsfilelock.NewLockWithTimeout(util.GetHostNamespacePath(i.hostProc), LockFile, LockTimeout)
		if err := lock.Lock(); err != nil {
			return errors.Wrapf(err, "failed to get file lock for initiator %s", i.Name)
		}
		defer lock.Unlock()
	}

	// Check if the initiator/NVMe device is already launched and matches the params
	if err := i.loadStartedDeviceInfoWithoutLock(); err == nil {
		if i.TransportAddress == transportAddress && i.TransportServiceID == transportServiceID {
			i.logger.Info("NVMe initiator is already launched with correct params")
			return nil
		}
		i.logger.Warnf("NVMe initiator is launched but with incorrect address, the required one is %s:%s, will try to stop then relaunch it",
			transportAddress, transportServiceID)
	}

	i.logger.Infof("Stopping NVMe initiator blindly  before starting")
	if err := i.stopWithoutLock(); err != nil {
		return errors.Wrapf(err, "failed to stop the mismatching NVMe initiator %s before starting", i.Name)
	}

	i.logger.Infof("Launching NVMe initiator")

	// Setup initiator
	for counter := 0; counter < RetryCounts; counter++ {
		// Rerun this API for a discovered target should be fine
		if i.SubsystemNQN, err = DiscoverTarget(transportAddress, transportServiceID, i.executor); err != nil {
			i.logger.WithError(err).Warnf("Failed to discover")
			time.Sleep(RetryInterval)
			continue
		}

		if i.ControllerName, err = ConnectTarget(transportAddress, transportServiceID, i.SubsystemNQN, i.executor); err != nil {
			i.logger.WithError(err).Warnf("Failed to connect target")
			time.Sleep(RetryInterval)
			continue
		}
		break
	}

	if i.ControllerName == "" {
		return fmt.Errorf("failed to start NVMe initiator %s within %d * %vsec retrys", i.Name, RetryCounts, RetryInterval.Seconds())
	}

	time.Sleep(1 * time.Second)

	if err := i.loadStartedDeviceInfoWithoutLock(); err != nil {
		return errors.Wrapf(err, "failed to load device info after starting NVMe initiator %s", i.Name)
	}

	if err := i.makeEndpoint(); err != nil {
		return err
	}

	i.isUp = true
	i.logger.Infof("Launched NVMe initiator")

	return nil
}

func (i *Initiator) Stop() error {
	if i.hostProc != "" {
		lock := nsfilelock.NewLockWithTimeout(util.GetHostNamespacePath(i.hostProc), LockFile, LockTimeout)
		if err := lock.Lock(); err != nil {
			return errors.Wrapf(err, "failed to get file lock for NVMe initiator %s", i.Name)
		}
		defer lock.Unlock()
	}

	return i.stopWithoutLock()
}

func (i *Initiator) stopWithoutLock() error {
	if err := i.removeEndpoint(); err != nil {
		return err
	}

	if err := DisconnectTarget(i.SubsystemNQN, i.executor); err != nil {
		return errors.Wrapf(err, "failed to logout target")
	}
	i.ControllerName = ""
	i.NamespaceName = ""
	i.TransportAddress = ""
	i.TransportServiceID = ""

	return nil
}

func (i *Initiator) GetControllerName() string {
	return i.ControllerName
}

func (i *Initiator) GetNamespaceName() string {
	return i.NamespaceName
}

func (i *Initiator) GetTransportAddress() string {
	return i.TransportAddress
}

func (i *Initiator) GetTransportServiceID() string {
	return i.TransportServiceID
}

func (i *Initiator) GetEndpoint() string {
	return i.Endpoint
}

func (i *Initiator) LoadStartedDeviceInfo() error {
	if i.hostProc != "" {
		lock := nsfilelock.NewLockWithTimeout(util.GetHostNamespacePath(i.hostProc), LockFile, LockTimeout)
		if err := lock.Lock(); err != nil {
			return errors.Wrapf(err, "failed to get file lock for NVMe initiator %s", i.Name)
		}
		defer lock.Unlock()
	}

	return i.loadStartedDeviceInfoWithoutLock()
}

func (i *Initiator) loadStartedDeviceInfoWithoutLock() error {
	nvmeDevices, err := GetDevices(i.TransportAddress, i.TransportServiceID, i.SubsystemNQN, i.executor)
	if err != nil {
		return err
	}
	if len(nvmeDevices) != 1 {
		return fmt.Errorf("found zero or multiple devices NVMe initiator %s", i.Name)
	}
	if len(nvmeDevices[0].Namespaces) != 1 {
		return fmt.Errorf("found zero or multiple devices for NVMe initiator %s", i.Name)
	}
	if i.ControllerName != "" && i.ControllerName != nvmeDevices[0].Controllers[0].Controller {
		return fmt.Errorf("found mismatching between the detected controller name %s and the recorded value %s for NVMe initiator %s", nvmeDevices[0].Controllers[0].Controller, i.ControllerName, i.Name)
	}
	i.ControllerName = nvmeDevices[0].Controllers[0].Controller
	i.NamespaceName = nvmeDevices[0].Namespaces[0].NameSpace
	i.TransportAddress, i.TransportServiceID = GetIPAndPortFromControllerAddress(nvmeDevices[0].Controllers[0].Address)
	i.logger.WithFields(logrus.Fields{
		"controllerName":     i.ControllerName,
		"namespaceName":      i.NamespaceName,
		"transportAddress":   i.TransportAddress,
		"transportServiceID": i.TransportServiceID,
	})

	devices, err := util.GetKnownDevices(i.executor)
	if err != nil {
		return err
	}
	i.dev = devices[i.NamespaceName]
	if i.dev == nil {
		return fmt.Errorf("cannot find the device for NVMe initiator %s with namespace name %s", i.Name, i.NamespaceName)
	}

	return nil
}

func (i *Initiator) makeEndpoint() error {
	endpoint := util.GetLonghornDevicePath(i.Name)
	if i.Endpoint != "" && i.Endpoint != endpoint {
		return fmt.Errorf("NVMe initiator %s already has an endpoint %s and it's different from the target endpoint %s", i.Name, i.Endpoint, endpoint)
	}
	if err := util.DuplicateDevice(i.dev, endpoint); err != nil {
		return err
	}
	i.Endpoint = endpoint

	return nil
}

func (i *Initiator) removeEndpoint() error {
	if err := util.RemoveDevice(i.Endpoint); err != nil {
		return err
	}
	i.Endpoint = ""
	i.dev = nil

	return nil
}
