package nvme

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	//"github.com/longhorn/nsfilelock"

	"github.com/longhorn/go-spdk-helper/pkg/types"
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
	SubsystemNQN       string
	TransportAddress   string
	TransportServiceID string

	ControllerName string

	// ControllerLossTimeout int64
	// FastIOFailTimeout     int64

	logger logrus.FieldLogger
}

func NewInitiator(subsystemNQN, transportServiceID string) (*Initiator, error) {
	// localIP, err := util.GetIPToHost()
	// if err != nil {
	// 	return err
	// }

	dev := &Initiator{
		SubsystemNQN:       subsystemNQN,
		TransportAddress:   types.LocalIP,
		TransportServiceID: transportServiceID,

		logger: logrus.WithFields(logrus.Fields{
			"subsystemNQN":       subsystemNQN,
			"transportAddress":   types.LocalIP,
			"transportServiceID": transportServiceID,
		}),
	}
	return dev, nil
}

func (i *Initiator) StartInitiator() error {
	/*
			lock := nsfilelock.NewLockWithTimeout(util.GetHostNamespacePath(HostProc), LockFile, LockTimeout)
			if err := lock.Lock(); err != nil {
				return errors.Wrap(err, "failed to lock")
			}
			defer lock.Unlock()

		ne, err := util.NewNamespaceExecutor(util.GetHostNamespacePath(HostProc))
		if err != nil {
			return err
		}
	*/
	if err := CheckForNVMeCliExistence(util.Execute); err != nil {
		return err
	}

	// Check if the initiator/NVMe device is already launched and matches the params
	if nvmeDevices, err := GetDevices(i.TransportAddress, i.TransportServiceID, i.SubsystemNQN, util.Execute); err == nil && len(nvmeDevices) == 1 {
		i.ControllerName = nvmeDevices[0].Controllers[0].Controller
		i.logger.WithField("controllerName", i.ControllerName)
		i.logger.Infof("the NVMe initiator is already launched with correct params")
		return nil
	}

	i.logger.Infof("Prepare to blindly do cleanup before starting")
	if err := DisconnectTarget(i.SubsystemNQN, util.Execute); err != nil {
		return errors.Wrapf(err, "failed to logout the mismatching target before starting")
	}

	i.logger.Infof("Prepare to launch NVMe initiator")

	var err error
	// Setup initiator
	for counter := 0; counter < RetryCounts; counter++ {
		// Rerun this API for a discovered target should be fine
		i.SubsystemNQN, err = DiscoverTarget(i.TransportAddress, i.TransportServiceID, util.Execute)
		if err != nil {
			logrus.WithError(err).Warnf("Failed to discover")
			time.Sleep(RetryInterval)
			continue
		}
		if i.ControllerName, err = ConnectTarget(i.TransportAddress, i.TransportServiceID, i.SubsystemNQN, util.Execute); err != nil {
			logrus.WithError(err).Warnf("Failed to connect target, ")
			time.Sleep(RetryInterval)
			continue
		}
		i.logger.WithField("controllerName", i.ControllerName)
		break
	}

	if i.ControllerName == "" {
		return fmt.Errorf("failed to start initiator within %d * %vsec retrys", RetryCounts, RetryInterval.Seconds())
	}

	return nil
}

func (i *Initiator) StopInitator() error {
	/*
		lock := nsfilelock.NewLockWithTimeout(util.GetHostNamespacePath(HostProc), LockFile, LockTimeout)
		if err := lock.Lock(); err != nil {
			return errors.Wrap(err, "failed to lock")
		}
		defer lock.Unlock()

		ne, err := util.NewNamespaceExecutor(util.GetHostNamespacePath(HostProc))
		if err != nil {
			return err
		}
	*/
	if err := DisconnectTarget(i.SubsystemNQN, util.Execute); err != nil {
		return errors.Wrapf(err, "failed to logout target")
	}
	return nil
}
