package disk

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	spdkclient "github.com/longhorn/go-spdk-helper/pkg/spdk/client"
	spdktypes "github.com/longhorn/go-spdk-helper/pkg/spdk/types"
	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func (s *Server) blockTypeDiskCreate(ctx context.Context, req *rpc.DiskCreateRequest) (*rpc.Disk, error) {
	log := logrus.WithFields(logrus.Fields{
		"diskName":  req.DiskName,
		"diskPath":  req.DiskPath,
		"blockSize": req.BlockSize,
	})

	if err := validateDiskCreateRequest(req); err != nil {
		log.WithError(err).Error("Failed to validate disk create request")
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, errors.Wrap(err, "failed to validate disk create request").Error())
	}

	blockSize := uint64(defaultBlockSize)
	if req.BlockSize > 0 {
		log.Infof("Using custom block size %v", req.BlockSize)
		blockSize = uint64(req.BlockSize)
	}

	uuid, err := s.addBlockDevice(req.DiskName, req.DiskPath, blockSize)
	if err != nil {
		log.WithError(err).Error("Failed to add block device")
		return nil, grpcstatus.Error(grpccodes.Internal, errors.Wrap(err, "failed to add block device").Error())
	}

	return s.lvstoreToDisk(req.DiskPath, "", uuid)
}

func (s *Server) blockTypeDiskGet(ctx context.Context, req *rpc.DiskGetRequest) (*rpc.Disk, error) {
	log := logrus.WithFields(logrus.Fields{
		"diskName": req.DiskName,
		"diskPath": req.DiskPath,
	})

	// Check if the disk exists
	bdevInfos, err := s.spdkClient.BdevAioGet(req.DiskName, 3000)
	if err != nil {
		resp, parseErr := parseErrorMessage(err.Error())
		if parseErr != nil || !isNoSuchDevice(resp.Message) {
			log.WithError(err).Errorf("Failed to get AIO bdev with name %v", req.DiskName)
			return nil, grpcstatus.Errorf(grpccodes.Internal, errors.Wrapf(err, "failed to get AIO bdev with name %v", req.DiskName).Error())
		}
	}
	if len(bdevInfos) == 0 {
		log.WithError(err).Errorf("Cannot find AIO bdev with name %v", req.DiskName)
		return nil, grpcstatus.Errorf(grpccodes.NotFound, "cannot find AIO bdev with name %v", req.DiskName)
	}

	diskPath := getDiskPath(req.DiskPath)

	var bdevInfo *spdktypes.BdevInfo
	for i, info := range bdevInfos {
		if info.DriverSpecific != nil ||
			info.DriverSpecific.Aio != nil ||
			info.DriverSpecific.Aio.FileName == diskPath {
			bdevInfo = &bdevInfos[i]
			break
		}
	}
	if bdevInfo == nil {
		log.WithError(err).Errorf("Failed to get AIO bdev name for disk path %v", diskPath)
		return nil, grpcstatus.Errorf(grpccodes.NotFound, errors.Wrapf(err, "failed to get AIO bdev name for disk path %v", diskPath).Error())
	}

	// Get the disk information from the lvstore
	return s.lvstoreToDisk(req.DiskPath, req.DiskName, "")
}

func (s *Server) blockTypeDiskDelete(ctx context.Context, req *rpc.DiskDeleteRequest) error {
	log := logrus.WithFields(logrus.Fields{
		"diskName": req.DiskName,
	})

	_, err := s.spdkClient.BdevAioDelete(req.DiskName)
	if err != nil {
		log.WithError(err).Errorf("Failed to delete AIO bdev with name %v", req.DiskName)
		return grpcstatus.Errorf(grpccodes.Internal, errors.Wrapf(err, "failed to delete AIO bdev with name %v", req.DiskName).Error())
	}
	return nil
}

func (s *Server) bdevLvolGetLvstore(spdkClient *spdkclient.Client, lvsName, uuid string) (*spdktypes.LvstoreInfo, error) {
	lvstoreInfos, err := spdkClient.BdevLvolGetLvstore(lvsName, uuid)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get lvstore info")
	}

	if len(lvstoreInfos) == 0 {
		return nil, errors.New(syscall.ENOENT.Error())
	}

	if len(lvstoreInfos) != 1 {
		return nil, fmt.Errorf("number of lvstores with name %v is not 1", lvsName)
	}

	return &lvstoreInfos[0], nil
}

func getDiskPath(path string) string {
	return filepath.Join(hostPrefix, path)
}

func getDeviceNumber(filename string) (string, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return "", err
	}
	dev := fi.Sys().(*syscall.Stat_t).Dev
	major := int(dev >> 8)
	minor := int(dev & 0xff)
	return fmt.Sprintf("%d-%d", major, minor), nil
}

func isBlockDevice(fullPath string) (bool, error) {
	var st unix.Stat_t
	err := unix.Stat(fullPath, &st)
	if err != nil {
		return false, err
	}

	return (st.Mode & unix.S_IFMT) == unix.S_IFBLK, nil
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

func validateDiskCreateRequest(req *rpc.DiskCreateRequest) error {
	ok, err := isBlockDevice(req.DiskPath)
	if err != nil {
		return errors.Wrap(err, "failed to check if disk is a block device")
	}
	if !ok {
		return errors.Wrapf(err, "disk %v is not a block device", req.DiskPath)
	}

	size, err := getDiskDeviceSize(req.DiskPath)
	if err != nil {
		return errors.Wrap(err, "failed to get disk size")
	}
	if size == 0 {
		return fmt.Errorf("disk %v size is 0", req.DiskPath)
	}

	return nil
}

func getLvstoresAdded(lvstoresBeforeAdd, lvstoresAfterAdd []spdktypes.LvstoreInfo) []spdktypes.LvstoreInfo {
	lvstoresAdded := []spdktypes.LvstoreInfo{}
	for _, lvstoreAfterAdd := range lvstoresAfterAdd {
		found := false
		for _, lvstoreBeforeAdd := range lvstoresBeforeAdd {
			if lvstoreAfterAdd.UUID == lvstoreBeforeAdd.UUID {
				found = true
				break
			}
		}
		if !found {
			lvstoresAdded = append(lvstoresAdded, lvstoreAfterAdd)
		}
	}
	return lvstoresAdded
}

func (s *Server) addBlockDevice(diskName, diskPath string, blockSize uint64) (string, error) {
	log := logrus.WithFields(logrus.Fields{
		"diskName":  diskName,
		"diskPath":  diskPath,
		"blockSize": blockSize,
	})

	log.Infof("Creating AIO bdev %v with block size %v", diskName, blockSize)
	bdevName, err := s.spdkClient.BdevAioCreate(getDiskPath(diskPath), diskName, blockSize)
	if err != nil {
		resp, parseErr := parseErrorMessage(err.Error())
		if parseErr != nil || !isFileExists(resp.Message) {
			return "", errors.Wrapf(err, "failed to create AIO bdev")
		}
	}

	log.Infof("Creating lvstore %v", diskName)

	// Name of the lvstore is the same as the name of the aio bdev
	lvstoreName := bdevName
	lvstores, err := s.spdkClient.BdevLvolGetLvstore("", "")
	if err != nil {
		return "", errors.Wrapf(err, "failed to get lvstores")
	}

	for _, lvstore := range lvstores {
		if lvstore.BaseBdev != bdevName {
			continue
		}

		log.Infof("Found an existing lvstore %v", lvstore.Name)
		if lvstore.Name == lvstoreName {
			return lvstore.UUID, fmt.Errorf("lvstore %v already exists", lvstoreName)
		}

		log.Infof("Renaming the existing lvstore %v to %v", lvstore.Name, lvstoreName)
		renamed, err := s.spdkClient.BdevLvolRenameLvstore(lvstore.Name, lvstoreName)
		if err != nil {
			return "", errors.Wrapf(err, "failed to rename lvstore from %v to %v", lvstore.Name, lvstoreName)
		}
		if !renamed {
			return "", fmt.Errorf("failed to rename lvstore from %v to %v", lvstore.Name, lvstoreName)
		}
		return lvstore.UUID, nil
	}

	return s.spdkClient.BdevLvolCreateLvstore(lvstoreName, diskName, defaultClusterSize)
}

func (s *Server) lvstoreToDisk(diskPath, lvstoreName, lvstoreUUID string) (*rpc.Disk, error) {
	lvstoreInfo, err := s.bdevLvolGetLvstore(s.spdkClient, lvstoreName, lvstoreUUID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get lvstore with name %v and UUID %v", lvstoreName, lvstoreUUID)
	}

	// A disk does not have a fsid, so we use the device number as the disk ID
	diskID, err := getDeviceNumber(getDiskPath(diskPath))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get disk ID")
	}

	return &rpc.Disk{
		Id:          diskID,
		Uuid:        lvstoreInfo.UUID,
		Path:        diskPath,
		Type:        DiskTypeBlock,
		TotalSize:   int64(lvstoreInfo.TotalDataClusters * lvstoreInfo.ClusterSize),
		FreeSize:    int64(lvstoreInfo.FreeClusters * lvstoreInfo.ClusterSize),
		TotalBlocks: int64(lvstoreInfo.TotalDataClusters * lvstoreInfo.ClusterSize / lvstoreInfo.BlockSize),
		FreeBlocks:  int64(lvstoreInfo.FreeClusters * lvstoreInfo.ClusterSize / lvstoreInfo.BlockSize),
		BlockSize:   int64(lvstoreInfo.BlockSize),
		ClusterSize: int64(lvstoreInfo.ClusterSize),
	}, nil
}
