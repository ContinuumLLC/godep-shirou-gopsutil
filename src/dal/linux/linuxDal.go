package linux

import (
	"strings"
	"time"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
)

// AssetCollection related constants
const (
	cAssetCreatedBy string = "/continuum/agent/plugin/asset"
	cAssetDataType  string = "assetCollection"
)

const (
	cSysProductCmd string = `lshw -c system | grep product | cut -d ":" -f2`
	cSysTz         string = "date +%z"
	cSysTzd        string = "date +%Z"
)

// AssetDalImpl ...
type AssetDalImpl struct {
	Factory model.AssetDalDependencies
	Logger  logging.Logger
}

//GetAssetData ...
func (a AssetDalImpl) GetAssetData() (*asset.AssetCollection, error) {
	o, err := a.GetOSInfo()
	if err != nil {
		return nil, err
	}
	s, err := a.GetSystemInfo()
	if err != nil {
		return nil, err
	}
	return &asset.AssetCollection{
		CreatedBy:     cAssetCreatedBy,
		CreateTimeUTC: time.Now().UTC(),
		Type:          cAssetDataType,
		Os:            *o,
		BaseBoard:     *(getBaseBoardInfo()),
		Bios:          *(getBiosInfo()),
		Memory:        *(getMemoryInfo()),
		System:        *s,
	}, nil
}

// GetOSInfo returns the OS info
func (a AssetDalImpl) GetOSInfo() (*asset.AssetOs, error) {
	parser := a.Factory.GetParser()
	cfg := procParser.Config{
		ParserMode: procParser.ModeSeparator,
		Separator:  ":",
	}

	util := dalUtil{
		envDep: a.Factory,
	}
	dataCmd, err := util.getCommandData(parser, cfg, "lsb_release", "-a")
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	//changing the separator for next file/command to parse and get data
	cfg.Separator = "="
	dataFile, err := util.getFileData(parser, cfg, "/etc/default/locale")
	if err != nil {
		return nil, exception.New(model.ErrFileReadFailed, err)
	}
	return &asset.AssetOs{
		Product:      dataCmd.Map["Distributor ID"].Values[1],
		Manufacturer: dataCmd.Map["Description"].Values[1],
		Version:      dataCmd.Map["Release"].Values[1],
		OsLanguage:   strings.Trim(dataFile.Map["LANG"].Values[1], "\""),
		//os.InstallDate - To be added
		//os.SerialNumber - Presently not able to find it for ubuntu
	}, nil
}

// GetSystemInfo returns system info
func (a AssetDalImpl) GetSystemInfo() (*asset.AssetSystem, error) {
	//This command require sudo access to execute
	product, err := a.Factory.GetEnv().ExecuteBash(cSysProductCmd)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	//time zone
	tz, err := a.Factory.GetEnv().ExecuteBash(cSysTz)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	//time zone description
	tzd, err := a.Factory.GetEnv().ExecuteBash(cSysTzd)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	return &asset.AssetSystem{
		Product:             product,
		TimeZone:            tz,
		TimeZoneDescription: tzd,
		//Category - to be added
		//Model - to be added
		//SerialNumber - to be added
		//SystemName - to be added
	}, nil
}
