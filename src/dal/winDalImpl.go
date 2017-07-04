// +build windows

package dal

import (
	"errors"
	"runtime"
	"time"

	"golang.org/x/sys/windows/registry"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/shirou/gopsutil/baseboard"
	"github.com/shirou/gopsutil/host"
)

const (
	baseRegString = "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
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
		Model:        bb.Model,
		InstallDate:  bb.InstallDate,
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
	var svcPack string
	if os.CSDVersion != nil {
		svcPack = *os.CSDVersion
	}
	return &asset.AssetOs{
		Product:      os.Caption,
		Manufacturer: os.Manufacturer,
		Version:      os.Version,
		ServicePack:  svcPack,
		SerialNumber: os.SerialNumber,
		InstallDate:  os.InstallDate,
		Type:         runtime.GOOS,
		Arch:         os.OSArchitecture,
	}, nil
}

// GetSystemInfo returns system info
func (a assetDalImpl) GetSystemInfo() (*asset.AssetSystem, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}
	return &asset.AssetSystem{
		SystemName: info.Hostname,
	}, nil
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

func (a assetDalImpl) GetInstalledSoftwareInfo() ([]asset.AssetInstalledSoftware, error) {

	var objAssetInstalledSoftware []asset.AssetInstalledSoftware
	objAsset32BitInstalledSoftware, err := a.getInstalledSoftInfo(registry.WOW64_32KEY)
	if nil != err {
		return nil, err
	}

	if nil != objAsset32BitInstalledSoftware {
		objAssetInstalledSoftware = append(objAssetInstalledSoftware, objAsset32BitInstalledSoftware...)
	}

	objAsset64BitInstalledSoftware, err := a.getInstalledSoftInfo(registry.WOW64_64KEY)
	if nil != err {
		return nil, err
	}

	if nil != objAsset64BitInstalledSoftware {
		objAssetInstalledSoftware = append(objAssetInstalledSoftware, objAsset64BitInstalledSoftware...)
	}

	return objAssetInstalledSoftware, nil
}

func (a assetDalImpl) getInstalledSoftInfo(access32or64 uint32) ([]asset.AssetInstalledSoftware, error) {
	var objInstallSoft installSoftwareImpl
	var objAssetInstalledSoftware []asset.AssetInstalledSoftware

	subKeys, err := objInstallSoft.getSoftwareRegistrySubKeys(baseRegString, access32or64)
	if nil != err {
		return nil, err
	}

	for _, value := range subKeys {
		regSubKeys := baseRegString
		regSubKeys += "\\"
		regSubKeys += value

		softAttributes, err := objInstallSoft.getSoftwareRegistryProperties(regSubKeys, access32or64)
		if nil == err {
			if objInstallSoft.validatePropertiesForInstallSoftware(*softAttributes) {
				a.appendAttributesToAssetInstalledSoftware(*softAttributes, objAssetInstalledSoftware)
			}
		}
	}
	return objAssetInstalledSoftware, nil
}

func (a assetDalImpl) appendAttributesToAssetInstalledSoftware(softAttributes softwareAttributes, assetInstallSoft []asset.AssetInstalledSoftware) []asset.AssetInstalledSoftware {
	var objAssetInstalledSoftware asset.AssetInstalledSoftware

	objAssetInstalledSoftware.Name = softAttributes.displayName
	objAssetInstalledSoftware.Publisher = softAttributes.publisher
	objAssetInstalledSoftware.Version = softAttributes.displayVersion
	objAssetInstalledSoftware.InstallDate, _ = a.convertInstallDateToTime(softAttributes.installDate)

	assetInstallSoft = append(assetInstallSoft, objAssetInstalledSoftware)

	return assetInstallSoft
}

func (a assetDalImpl) convertInstallDateToTime(installDate string) (tm time.Time, err error) {

	if len(installDate) != 8 {
		return tm, errors.New(model.ErrAssetInstallDate)
	}
	tempInstallDate := installDate[0:4]
	tempInstallDate += "-"
	tempInstallDate += installDate[4:6]
	tempInstallDate += "-"
	tempInstallDate += installDate[6:8]
	tempInstallDate += "T00:00:00Z"

	tm, err = time.Parse(time.RFC3339, tempInstallDate)
	return tm, nil
}
