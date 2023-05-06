package disk

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

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

	if err := s.validateDiskCreateRequest(req); err != nil {
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
	bdevs, err := s.spdkClient.BdevAioGet(req.DiskName, 3000)
	if err != nil {
		resp, parseErr := parseErrorMessage(err.Error())
		if parseErr != nil || !isNoSuchDevice(resp.Message) {
			log.WithError(err).Errorf("Failed to get AIO bdev with name %v", req.DiskName)
			return nil, grpcstatus.Errorf(grpccodes.Internal, errors.Wrapf(err, "failed to get AIO bdev with name %v", req.DiskName).Error())
		}
	}
	if len(bdevs) == 0 {
		log.WithError(err).Errorf("Cannot find AIO bdev with name %v", req.DiskName)
		return nil, grpcstatus.Errorf(grpccodes.NotFound, "cannot find AIO bdev with name %v", req.DiskName)
	}

	diskPath := getDiskPath(req.DiskPath)

	var targetBdev *spdktypes.BdevInfo
	for i, bdev := range bdevs {
		if bdev.DriverSpecific != nil ||
			bdev.DriverSpecific.Aio != nil ||
			bdev.DriverSpecific.Aio.FileName == diskPath {
			targetBdev = &bdevs[i]
			break
		}
	}
	if targetBdev == nil {
		log.WithError(err).Errorf("Failed to get AIO bdev name for disk path %v", diskPath)
		return nil, grpcstatus.Errorf(grpccodes.NotFound, errors.Wrapf(err, "failed to get AIO bdev name for disk path %v", diskPath).Error())
	}

	return s.lvstoreToDisk(req.DiskPath, req.DiskName, "")
}

func (s *Server) blockTypeDiskDelete(ctx context.Context, req *rpc.DiskDeleteRequest) error {
	log := logrus.WithFields(logrus.Fields{
		"diskName": req.DiskName,
		"diskUUID": req.DiskUuid,
	})

	lvstores, err := s.spdkClient.BdevLvolGetLvstore("", req.DiskUuid)
	if err != nil {
		resp, parseErr := parseErrorMessage(err.Error())
		if parseErr != nil || !isNoSuchDevice(resp.Message) {
			return errors.Wrapf(err, "failed to get lvstore with UUID %v", req.DiskUuid)
		}
		log.WithError(err).Errorf("Cannot find lvstore with UUID %v", req.DiskUuid)
		return nil
	}

	lvstore := &lvstores[0]

	if lvstore.Name != req.DiskName {
		log.Warnf("Disk name %v does not match lvstore name %v", req.DiskName, lvstore.Name)
		return nil
	}

	_, err = s.spdkClient.BdevAioDelete(lvstore.BaseBdev)
	if err != nil {
		log.WithError(err).Errorf("Failed to delete AIO bdev %v", lvstore.BaseBdev)
		return errors.Wrapf(err, "failed to delete AIO bdev %v", lvstore.BaseBdev)
	}
	return nil
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

func (s *Server) validateDiskCreateRequest(req *rpc.DiskCreateRequest) error {
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

	diskID, err := getDeviceNumber(getDiskPath(req.DiskPath))
	if err != nil {
		return errors.Wrap(err, "failed to get disk device number")
	}

	bdevs, err := s.spdkClient.BdevAioGet("", 0)
	if err != nil {
		return errors.Wrap(err, "failed to get AIO bdevs")
	}

	for _, bdev := range bdevs {
		id, err := getDeviceNumber(bdev.DriverSpecific.Aio.FileName)
		if err != nil {
			return errors.Wrap(err, "failed to get disk device number")
		}

		if id == diskID {
			return fmt.Errorf("disk %v is already used by AIO bdev %v", req.DiskPath, bdev.Name)
		}
	}

	return nil
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
			return lvstore.UUID, nil
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
	lvstores, err := s.spdkClient.BdevLvolGetLvstore(lvstoreName, lvstoreUUID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get lvstore with name %v and UUID %v", lvstoreName, lvstoreUUID)
	}

	lvstore := &lvstores[0]

	// A disk does not have a fsid, so we use the device number as the disk ID
	diskID, err := getDeviceNumber(getDiskPath(diskPath))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get disk ID")
	}

	return &rpc.Disk{
		Id:          diskID,
		Uuid:        lvstore.UUID,
		Path:        diskPath,
		Type:        DiskTypeBlock,
		TotalSize:   int64(lvstore.TotalDataClusters * lvstore.ClusterSize),
		FreeSize:    int64(lvstore.FreeClusters * lvstore.ClusterSize),
		TotalBlocks: int64(lvstore.TotalDataClusters * lvstore.ClusterSize / lvstore.BlockSize),
		FreeBlocks:  int64(lvstore.FreeClusters * lvstore.ClusterSize / lvstore.BlockSize),
		BlockSize:   int64(lvstore.BlockSize),
		ClusterSize: int64(lvstore.ClusterSize),
	}, nil
}
