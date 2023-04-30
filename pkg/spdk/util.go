package spdk

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	issiutil "github.com/longhorn/go-iscsi-helper/util"
	spdkclient "github.com/longhorn/go-spdk-helper/pkg/spdk/client"
	spdktypes "github.com/longhorn/go-spdk-helper/pkg/spdk/types"
	spdkutil "github.com/longhorn/go-spdk-helper/pkg/types"

	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
)

func bdevLvolGetLvstore(c *spdkclient.Client, log logrus.FieldLogger, LvsUuid string) (*spdktypes.LvstoreInfo, error) {
	lvstoreInfos, err := c.BdevLvolGetLvstore("", LvsUuid)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get lvstore info")
	}

	if len(lvstoreInfos) != 1 {
		return nil, fmt.Errorf("number of lvstore info is not 1")
	}

	return &lvstoreInfos[0], nil
}

func bdevLvolGet(c *spdkclient.Client, log logrus.FieldLogger, lvName string) (*spdktypes.BdevInfo, error) {
	lvolInfos, err := c.BdevLvolGet(lvName, 3000)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get lvol info")
	}
	if len(lvolInfos) != 1 {
		return nil, errors.Wrapf(err, "number of lvol info is not 1")
	}

	return &lvolInfos[0], nil
}

func getNameFromAlias(alias []string) string {
	if len(alias) != 1 {
		return ""
	}

	splitName := strings.Split(alias[0], "/")
	return splitName[1]
}

func getVolumeName(engineName string) string {
	parts := strings.Split(engineName, "-e-")
	return parts[0]
}

func getEngine(client *spdkclient.Client, name string, log logrus.FieldLogger) (*rpc.Engine, error) {
	bdevRaidInfos, err := client.BdevRaidGetBdevs(spdktypes.BdevRaidCategoryAll)
	if err != nil {
		log.WithError(err).Error("Failed to get bdev raid infos")
		return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
	}

	var engine *rpc.Engine
	for _, info := range bdevRaidInfos {
		if info.Name != name {
			continue
		}

		raidNQN := spdkutil.GetNQN(info.Name)
		listenerList, err := client.NvmfSubsystemGetListeners(raidNQN, "")
		if err != nil {
			log.WithError(err).Error("Failed to get NVMe subsystem listeners")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}

		if len(listenerList) != 1 {
			log.WithError(err).Error("Invalid NVMe subsystem listeners")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}

		port, err := strconv.Atoi(listenerList[0].Address.Trsvcid)
		if err != nil {
			log.WithError(err).Error("Failed to convert port")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}

		localIP, err := issiutil.GetIPToHost()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get local IP")
		}

		replicaModeMap := map[string]rpc.ReplicaMode{}
		replicaAddressMap := map[string]string{}
		for _, baseBdev := range info.BaseBdevsList {
			parts := strings.Split(baseBdev.Name, "/")
			name := parts[1]
			address := localIP + ":" + "4420"
			replicaModeMap[address] = rpc.ReplicaMode_RW
			replicaAddressMap[name] = localIP + ":" + strconv.Itoa(port)
		}

		endpoint := "/dev/longhorn/" + getVolumeName(name)
		size, err := getDiskDeviceSize(endpoint)
		if err != nil {
			log.WithError(err).Error("Failed to get disk device size")
			return nil, grpcstatus.Error(grpccodes.Internal, err.Error())
		}

		engine = &rpc.Engine{
			Name:              name,
			Uuid:              "",
			Size:              uint64(size),
			Address:           "",
			ReplicaAddressMap: replicaAddressMap,
			ReplicaModeMap:    replicaModeMap,
			Endpoint:          endpoint,
			Frontend:          "spdk-tcp-blockdev",
			FrontendState:     getFrontendState(spdktypes.BdevRaidCategoryOffline),
			Ip:                listenerList[0].Address.Traddr,
			Port:              int32(port),
		}
	}

	if engine == nil {
		return nil, grpcstatus.Error(grpccodes.NotFound, "")
	}

	return engine, nil
}

func getFrontendState(category spdktypes.BdevRaidCategory) string {
	if category == spdktypes.BdevRaidCategoryOffline {
		return "down"
	}
	return "up"
}

func createLonghornDevice(devicePath, name string) error {
	logrus.Infof("Creating longhorn device: devicePath=%s, name=%s", devicePath, name)

	if _, err := os.Stat(devPath); os.IsNotExist(err) {
		if err := os.MkdirAll(devPath, 0755); err != nil {
			logrus.Fatalf("device %v: Cannot create directory %v", name, devPath)
		}
	}

	// Get the major and minor numbers of the NVMe device.
	major, minor, err := getDeviceNumbers(devicePath)
	if err != nil {
		return err
	}

	longhornDevPath := filepath.Join(devPath, name)

	return duplicateDevice(major, minor, longhornDevPath)
}

func deleteLonghornDevice(devicePath string) error {
	if _, err := os.Stat(devicePath); err == nil {
		if err := remove(devicePath); err != nil {
			return fmt.Errorf("failed to removing device %s, %v", devicePath, err)
		}
	}
	return nil
}

func removeAsync(path string, done chan<- error) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		logrus.Errorf("Unable to remove: %v", path)
		done <- err
	}
	done <- nil
}

func remove(path string) error {
	done := make(chan error)
	go removeAsync(path, done)
	select {
	case err := <-done:
		return err
	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout trying to delete %s", path)
	}
}

func getDeviceNumbers(devicePath string) (major, minor uint32, err error) {
	fileInfo, err := os.Stat(devicePath)
	if err != nil {
		return 0, 0, err
	}

	statT := fileInfo.Sys().(*syscall.Stat_t)
	major = uint32(int(statT.Rdev) >> 8)
	minor = uint32(int(statT.Rdev) & 0xFF)
	return major, minor, nil
}

func duplicateDevice(major, minor uint32, dest string) error {
	if err := mknod(dest, major, minor); err != nil {
		return fmt.Errorf("couldn't create device %s: %w", dest, err)
	}
	if err := os.Chmod(dest, 0660); err != nil {
		return fmt.Errorf("couldn't change permission of the device %s: %w", dest, err)
	}
	return nil
}

func mknod(device string, major, minor uint32) error {
	var fileMode os.FileMode = 0660
	fileMode |= unix.S_IFBLK
	dev := int(unix.Mkdev(uint32(major), uint32(minor)))

	logrus.Infof("Creating device %s %d:%d", device, major, minor)
	return unix.Mknod(device, uint32(fileMode), dev)
}

func (s *Server) replicaBroadcastConnector() (chan interface{}, error) {
	return s.replicaBroadcastCh, nil
}

func (s *Server) engineBroadcastConnector() (chan interface{}, error) {
	return s.engineBroadcastCh, nil
}

func getDiskDeviceSize(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to open %s", path)
	}
	defer file.Close()

	pos, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to seek %s", path)
	}
	return pos, nil
}

func splitHostPort(addr string) (string, string, error) {
	parts := strings.Split(addr, ":")

	if len(parts) != 2 {
		return "", "", grpcstatus.Errorf(grpccodes.InvalidArgument, "Invalid address %v", addr)
	}

	return parts[0], parts[1], nil
}

// wait for device file comes up or timeout
func waitForDeviceReady(devPath string, seconds int) error {
	for i := 0; i <= seconds; i++ {
		time.Sleep(time.Second)
		_, err := os.Stat(devPath)
		if err == nil {
			return nil
		}
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return fmt.Errorf("device %s not found", devPath)
}
