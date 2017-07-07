// +build windows

package disk

import (
	"strings"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/StackExchange/wmi"
)

const diskQuery = "SELECT Name, Caption, Manufacturer, MediaType, SerialNumber, Index, Partitions, Size FROM Win32_DiskDrive"

// ByWMI implements Disk Information using WMI
type ByWMI struct {
}

// Win32_DiskDrive WMI class representation
type win32DiskDrive struct {
	Name         string
	Caption      string
	Manufacturer string
	MediaType    string
	SerialNumber string
	Index        uint32
	Partitions   uint32
	Size         uint64
}

// Info returns Disk information for Windows using WMI
func (ByWMI) Info() ([]asset.AssetDrive, error) {
	var dst []win32DiskDrive
	err := wmi.Query(diskQuery, &dst)
	if err != nil {
		return nil, err
	}
	var listOfDrives []asset.AssetDrive
	iDiskLen := len(dst)
	for i := 0; i < iDiskLen; i++ {
		tmp := asset.AssetDrive{
			Product:            dst[i].Caption,
			Manufacturer:       dst[i].Manufacturer,
			MediaType:          dst[i].MediaType,
			LogicalName:        dst[i].Name,
			NumberOfPartitions: int(dst[i].Partitions),
			SerialNumber:       strings.TrimSpace(dst[i].SerialNumber),
			SizeBytes:          int64(dst[i].Size),
		}
		listOfDrives = append(listOfDrives, tmp)
	}
	return listOfDrives, nil
}
