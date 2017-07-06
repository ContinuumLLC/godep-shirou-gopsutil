// +build windows

package bios

import (
	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/StackExchange/wmi"
)

const biosQuery = "SELECT Name,Manufacturer,SerialNumber,Version,SMBIOSBIOSVersion FROM Win32_BIOS"

// ByWMI implements BIOS Information using WMI
type ByWMI struct {
}

// Win32_BIOS WMI struct representation
type win32BIOS struct {
	Name              string
	Manufacturer      string
	SerialNumber      string
	Version           string
	SMBIOSBIOSVersion string
}

// Info returns BIOS information for Windows using WMI
func (ByWMI) Info() (*asset.AssetBios, error) {
	var dst []win32BIOS
	err := wmi.Query(biosQuery, &dst)
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
