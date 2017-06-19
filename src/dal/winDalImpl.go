// +build windows

package dal

import (
	"time"

	"runtime"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/shirou/gopsutil/baseboard"
	"github.com/shirou/gopsutil/host"
)

// AssetCollection related constants
const (
	cAssetCreatedBy string = "/continuum/agent/plugin/asset"
	cAssetDataType  string = "assetCollection"
	cAssetDataName  string = "asset"
)

type assetDalImpl struct {
	Factory model.AssetDalDependencies
	Logger  logging.Logger
}

//GetAssetData ...
func (a assetDalImpl) GetAssetData() (*asset.AssetCollection, error) {
	var (
		bbrd  asset.AssetBaseBoard
		bios  asset.AssetBios
		opSys asset.AssetOs
		memry asset.AssetMemory
		syst  asset.AssetSystem
	)
	b, err := a.GetBiosInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetBiosInfo() %v", err)
	} else {
		bios = *b
	}

	bb, err := a.GetBaseBoardInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetBaseBoardInfo() %v", err)
	} else {
		bbrd = *bb
	}

	dd, err := a.GetDrivesInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetDrivesInfo() %v", err)
	}

	o, err := a.GetOSInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetOSInfo() %v", err)
	} else {
		opSys = *o
	}

	s, err := a.GetSystemInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetSystemInfo() %v", err)
	} else {
		syst = *s
	}

	n, err := a.GetNetworkInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetNetworkInfo() %v", err)
	}
	m, err := a.GetMemoryInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetMemoryInfo() %v", err)
	} else {
		memry = *m
	}
	p, err := a.GetProcessorInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetProcessorInfo() %v", err)
	}

	return &asset.AssetCollection{
		CreatedBy:     cAssetCreatedBy,
		CreateTimeUTC: time.Now().UTC(),
		Type:          cAssetDataType,
		Name:          cAssetDataName,
		Bios:          bios,
		BaseBoard:     bbrd,
		Os:            opSys,
		Memory:        memry,
		System:        syst,
		Networks:      n,
		Drives:        dd,
		Processors:    p,
	}, nil
}

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
