// +build windows

package dal

import (
	"runtime"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/shirou/gopsutil/baseboard"
	"github.com/shirou/gopsutil/host"
)

//GetBiosInfo ...
func (a assetDalImpl) GetBiosInfo() (*asset.AssetBios, error) {
	return &asset.AssetBios{}, nil
}

//GetBaseBoardInfo ...
func (a assetDalImpl) GetBaseBoardInfo() (*asset.AssetBaseBoard, error) {
	bb, err := baseboard.Info()
	if err != nil {
		return nil, err
	}
	return &asset.AssetBaseBoard{
		Name:         bb.Name,
		Product:      bb.Product,
		Manufacturer: bb.Manufacturer,
		SerialNumber: bb.SerialNumber,
		Version:      bb.Version,
	}, nil

}

//GetDrivesInfo ...
func (a assetDalImpl) GetDrivesInfo() ([]asset.AssetDrive, error) {
	return nil, nil
}

// GetOSInfo returns the OS info
func (a assetDalImpl) GetOSInfo() (*asset.AssetOs, error) {
	os, err := host.GetOSInfo()
	if err != nil {
		return nil, err
	}
	return &asset.AssetOs{
		Product:      os.Caption,
		Manufacturer: os.Manufacturer,
		Version:      os.Version,
		ServicePack:  os.CSDVersion,
		SerialNumber: os.SerialNumber,
		InstallDate:  os.InstallDate,
		Type:         runtime.GOOS,
		Arch:         os.OSArchitecture,
	}, nil
}

// GetSystemInfo returns system info
func (a assetDalImpl) GetSystemInfo() (*asset.AssetSystem, error) {
	return &asset.AssetSystem{}, nil
}

// GetNetworkInfo returns network info
func (a assetDalImpl) GetNetworkInfo() ([]asset.AssetNetwork, error) {
	return nil, nil
}

// GetMemoryInfo returns memory info
func (a assetDalImpl) GetMemoryInfo() (*asset.AssetMemory, error) {
	return &asset.AssetMemory{}, nil
}

// GetProcessorInfo returns processor info
func (a assetDalImpl) GetProcessorInfo() ([]asset.AssetProcessor, error) {
	return nil, nil
}
