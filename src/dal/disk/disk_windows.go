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
	var dst []win32DiskDrive
	err := w.dep.Query(diskQ, &dst)
	if err != nil {
		return nil, err
	}
	var drives []asset.AssetDrive
	iDiskLen := len(dst)
	for i := 0; i < iDiskLen; i++ {
		// get partition data corresponding to a disk
		parts, err := w.diskToPartitions(dst[i].DeviceID)
		if err != nil {
			return nil, err
		}
		lDisks, err := w.logicalDiskInfo(parts)
		if err != nil {
			return nil, err
		}

		var srNum string
		if dst[i].SerialNumber != nil {
			srNum = strings.TrimSpace(*dst[i].SerialNumber)
		}
		disk := asset.AssetDrive{
			Product:            dst[i].Caption,
			Manufacturer:       dst[i].Manufacturer,
			MediaType:          dst[i].MediaType,
			LogicalName:        strings.Replace(dst[i].Name, `\\\\.\\`, `\\.\`, -1),
			NumberOfPartitions: int(dst[i].Partitions),
			SerialNumber:       srNum,
			SizeBytes:          int64(dst[i].Size),
			PartitionData:      lDisks,
		}
		drives = append(drives, disk)
	}
	return drives, nil
}

func (w WMI) diskToPartitions(deviceID string) ([]win32DiskPartition, error) {
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
		for _, v := range dstLD {
			part := asset.AssetDrivePartition{
				Name:        v.Name,
				Label:       v.VolumeName,
				Description: v.Description,
				FileSystem:  v.FileSystem,
				SizeBytes:   int64(v.Size),
			}
			partitions = append(partitions, part)
		}
	}
	return partitions, nil
}
