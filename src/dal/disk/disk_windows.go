// +build windows

package disk

import (
	"fmt"
	"strings"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/wmi"
)

const (
	diskQ          = "SELECT Name, Caption, Manufacturer, MediaType, SerialNumber, Index, Partitions, Size, DeviceID FROM Win32_DiskDrive"
	diskPartitionQ = "ASSOCIATORS OF {Win32_DiskDrive.DeviceID='%s'} WHERE AssocClass = Win32_DiskDriveToDiskPartition"
	logicalDiskQ   = "ASSOCIATORS OF {Win32_DiskPartition.DeviceID='%s'} WHERE AssocClass = Win32_LogicalDiskToPartition"
)

// WMI implements Disk Information using WMI
type WMI struct {
	dep wmi.Wrapper
}

// GetByWMI returns the WMI implementation Disk information
func GetByWMI() WMI {
	return WMI{
		dep: wmi.GetWrapper(),
	}
}

// Win32_DiskDrive WMI class representation
type win32DiskDrive struct {
	Name         string
	Caption      string
	DeviceID     string
	Manufacturer string
	MediaType    string
	SerialNumber *string
	Index        uint32
	Partitions   uint32
	Size         uint64
}

// Win32_DiskPartition WMI class representation
type win32DiskPartition struct {
	DeviceID string
}

// Win32_LogicalDisk WMI class representation
type win32LogicalDisk struct {
	Name        string
	Description string
	FileSystem  string
	VolumeName  string
	Size        uint64
}

// Info returns Disk information for Windows using WMI
func (w WMI) Info() ([]asset.AssetDrive, error) {
	dst, err := w.disk()
	if err != nil {
		return nil, err
	}
	var drives []asset.AssetDrive
	l := len(dst)
	for i := 0; i < l; i++ {
		// get partition data corresponding to a disk
		parts, err := w.diskToPartition(dst[i].DeviceID)
		if err != nil {
			return nil, err
		}
		lDisks, err := w.logicalDiskInfo(parts)
		if err != nil {
			return nil, err
		}
		disk := mapToDriveModel(&dst[i], lDisks)
		drives = append(drives, *disk)
	}
	return drives, nil
}

func (w WMI) disk() ([]win32DiskDrive, error) {
	var dst []win32DiskDrive
	err := w.dep.Query(diskQ, &dst)
	return dst, err
}

func (w WMI) diskToPartition(deviceID string) ([]win32DiskPartition, error) {
	var dst []win32DiskPartition
	q := fmt.Sprintf(diskPartitionQ, deviceID)
	err := w.dep.Query(q, &dst)
	return dst, err
}

func (w WMI) partitionToLogicalDisk(deviceID string) ([]win32LogicalDisk, error) {
	var dst []win32LogicalDisk
	q := fmt.Sprintf(logicalDiskQ, deviceID)
	err := w.dep.Query(q, &dst)
	return dst, err
}

func (w WMI) logicalDiskInfo(dstDP []win32DiskPartition) ([]asset.AssetDrivePartition, error) {
	var partitions []asset.AssetDrivePartition
	iPartitions := len(dstDP)
	for i := 0; i < iPartitions; i++ {
		dstLD, err := w.partitionToLogicalDisk(dstDP[i].DeviceID)
		if err != nil {
			return nil, err
		}
		parts := mapToPartitionDataModel(dstLD)
		partitions = append(partitions, parts...)
	}
	return partitions, nil
}

func mapToPartitionDataModel(ld []win32LogicalDisk) []asset.AssetDrivePartition {
	var parts []asset.AssetDrivePartition
	for _, v := range ld {
		part := asset.AssetDrivePartition{
			Name:        v.Name,
			Label:       v.VolumeName,
			Description: v.Description,
			FileSystem:  v.FileSystem,
			SizeBytes:   int64(v.Size),
		}
		parts = append(parts, part)
	}
	return parts
}

func mapToDriveModel(disk *win32DiskDrive, lDisks []asset.AssetDrivePartition) *asset.AssetDrive {
	var srNum string
	if disk.SerialNumber != nil {
		srNum = strings.TrimSpace(*disk.SerialNumber)
	}
	return &asset.AssetDrive{
		Product:            disk.Caption,
		Manufacturer:       disk.Manufacturer,
		MediaType:          disk.MediaType,
		LogicalName:        strings.Replace(disk.Name, `\\.\`, ``, -1),
		NumberOfPartitions: int(disk.Partitions),
		SerialNumber:       srNum,
		SizeBytes:          int64(disk.Size),
		PartitionData:      lDisks,
	}
}
