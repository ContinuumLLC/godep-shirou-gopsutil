package dal

import (
	"time"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
)

// AssetCollection related constants
const (
	cAssetCreatedBy string = "/agent/plugin/asset"
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
		baseboard asset.AssetBaseBoard
		bios      asset.AssetBios
		os        asset.AssetOs
		sys       asset.AssetSystem
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
		baseboard = *bb
	}

	o, err := a.GetOSInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetOSInfo() %v", err)
	} else {
		os = *o
	}

	s, err := a.GetSystemInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetSystemInfo() %v", err)
	} else {
		sys = *s
	}

	mem, err := a.GetMemoryInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetMemoryInfo() %v", err)
	}

	storages, err := a.GetDrivesInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetDrivesInfo() %v", err)
	}

	net, err := a.GetNetworkInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetNetworkInfo() %v", err)
	}

	cpu, err := a.GetProcessorInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetProcessorInfo() %v", err)
	}

	installSoft, err := a.GetInstalledSoftwareInfo()
	if err != nil {
		a.Logger.Logf(logging.ERROR, "GetInstalledSoftwareInfo() %v", err)
	}

	return &asset.AssetCollection{
		CreatedBy:          cAssetCreatedBy,
		CreateTimeUTC:      time.Now().UTC(),
		Type:               cAssetDataType,
		Name:               cAssetDataName,
		Bios:               bios,
		BaseBoard:          baseboard,
		Os:                 os,
		Memory:             mem,
		System:             sys,
		Networks:           net,
		Drives:             storages,
		Processors:         cpu,
		InstalledSoftwares: installSoft,
	}, nil
}
