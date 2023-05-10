package spdk

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	issiutil "github.com/longhorn/go-iscsi-helper/util"
	spdkclient "github.com/longhorn/go-spdk-helper/pkg/spdk/client"
	spdktypes "github.com/longhorn/go-spdk-helper/pkg/spdk/types"
	spdkutil "github.com/longhorn/go-spdk-helper/pkg/types"

	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
)

func bdevLvolGetLvstore(c *spdkclient.Client, log logrus.FieldLogger, LvsUUID string) (*spdktypes.LvstoreInfo, error) {
	lvstoreInfos, err := c.BdevLvolGetLvstore("", LvsUUID)
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
	bdevRaidInfos, err := client.BdevRaidGetInfoByCategory(spdktypes.BdevRaidCategoryAll)
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

			replicaName := ""
			ip := ""
			port := ""

			if len(parts) == 1 {
				controllerInfos, err := client.BdevNvmeGetControllers(replicaName)
				if err != nil {
					log.WithError(err).Error("Failed to get NVMe controllers")
					continue
				}
				controller := controllerInfos[0]

				replicaName = stripTail(baseBdev.Name, "n1")
				ip = controller.Ctrlrs[0].Trid.Traddr
				port = controller.Ctrlrs[0].Trid.Trsvcid
			} else {
				replicaName = parts[1]
				ip = localIP
				port = "0"
			}

			address := fmt.Sprintf("%s:%s", ip, port)
			replicaModeMap[address] = rpc.ReplicaMode_RW
			replicaAddressMap[replicaName] = address
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
			SpecSize:          uint64(size),
			ActualSize:        uint64(size),
			Ip:                listenerList[0].Address.Traddr,
			Port:              int32(port),
			ReplicaAddressMap: replicaAddressMap,
			ReplicaModeMap:    replicaModeMap,
			Endpoint:          endpoint,
		}
	}

	if engine == nil {
		return nil, grpcstatus.Error(grpccodes.NotFound, "")
	}

	return engine, nil
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

func parsePortRange(portRange string) (int32, int32, error) {
	if portRange == "" {
		return 0, 0, fmt.Errorf("empty port range")
	}
	parts := strings.Split(portRange, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format for range: %s", portRange)
	}
	portStart, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, errors.Wrap(err, "invalid start port for range")
	}
	portEnd, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, errors.Wrap(err, "invalid end port for range")
	}
	return int32(portStart), int32(portEnd), nil
}

func (s *Server) allocatePorts(portCount int32) (int32, int32, error) {
	if portCount < 0 {
		return 0, 0, fmt.Errorf("invalid port count %v", portCount)
	}
	if portCount == 0 {
		return 0, 0, nil
	}
	start, end, err := s.availablePorts.AllocateRange(portCount)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to allocate %v ports", portCount)
	}
	return int32(start), int32(end), nil
}

func (s *Server) allocateInstancePorts(name string, portCount int32) (int32, int32, error) {
	portStart, portEnd, err := s.allocatePorts(portCount)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "cannot allocate %v ports for %v", portCount, name)
	}
	return portStart, portEnd, nil
}

func (s *Server) releasePorts(start, end int32) error {
	if start < 0 || end < 0 {
		return fmt.Errorf("invalid start/end port %v %v", start, end)
	}
	return s.availablePorts.ReleaseRange(start, end)
}

func (s *Server) releaseInstancePorts(name string, portStart, portEnd int32) {
	if err := s.releasePorts(portStart, portEnd); err != nil {
		logrus.WithError(err).Errorf("cannot deallocate (%v-%v) ports for %v",
			portStart, portEnd, name)
	}
}

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func parseErrorMessage(errStr string) (*ErrorMessage, error) {
	r := regexp.MustCompile(`"code": (-?\d+),\n\t"message": "([^"]+)"`)
	match := r.FindStringSubmatch(errStr)
	if len(match) == 0 {
		return nil, fmt.Errorf("failed to parse error message")
	}

	code := 0
	if _, err := fmt.Sscanf(match[1], "%d", &code); err != nil {
		return nil, fmt.Errorf("failed to parse error code: %w", err)
	}

	em := &ErrorMessage{
		Code:    code,
		Message: match[2],
	}

	return em, nil
}

func isFileExists(message string) bool {
	return strings.EqualFold(message, syscall.Errno(syscall.EEXIST).Error())
}

func isNoSuchDevice(message string) bool {
	return strings.EqualFold(message, syscall.Errno(syscall.ENODEV).Error())
}

func stripTail(s string, tail string) string {
	if strings.HasSuffix(s, tail) {
		return s[:len(s)-len(tail)]
	}
	return s
}

const (
	EngineRandomIDLenth = 8
	EngineSuffix        = "-e"
)

func GetVolumeNameFromEngineName(engineName string) string {
	reg := regexp.MustCompile(fmt.Sprintf(`([^"]*)%s-[A-Za-z0-9]{%d,%d}$`, EngineSuffix, EngineRandomIDLenth, EngineRandomIDLenth))
	return reg.ReplaceAllString(engineName, "${1}")
}
