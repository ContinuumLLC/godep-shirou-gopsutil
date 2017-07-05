// +build windows

package bios

import (
	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/StackExchange/wmi"
)

// ByWMI implements bios Information using WMI
type ByWMI struct {
}

// Info returns baseboard information for Windows using WMI
func (ByWMI) Info() (*asset.AssetBios, error) {
	// Win32_BIOS struct represents a bios
	type Win32_BIOS struct {
		Name              string
		Manufacturer      string
		SerialNumber      string
		Version           string
		SMBIOSBIOSVersion string
	}
	var dst []Win32_BIOS
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		return nil, err
	}
	return &asset.AssetBios{
		Product:       dst[0].Name,
		Manufacturer:  dst[0].Manufacturer,
		SerialNumber:  dst[0].SerialNumber,
		Version:       dst[0].Version,
		SmbiosVersion: dst[0].SMBIOSBIOSVersion,
	}, nil
}
