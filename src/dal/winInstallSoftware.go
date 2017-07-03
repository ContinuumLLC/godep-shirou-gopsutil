// +build windows

package dal

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
)

var releaseTypes = [...]string{"Hotfix", "Security Update", "Update Rollup", "Update", "Critical Update", "Security"}

type softwareAttributes struct {
	displayName      string
	displayIcon      string
	displayVersion   string
	installDate      string
	publisher        string
	uninstallString  string
	kbNumber         string
	parentKeyName    string
	releaseType      string
	systemComponent  uint64
	windowsInstaller uint64
}

type installSoftwareImpl struct {
}

func (installSoftwareImpl) validatePropertiesForInstallSoftware(objSoftAttributes softwareAttributes) bool {
	if objSoftAttributes.displayName == "" ||
		objSoftAttributes.uninstallString == "" ||
		objSoftAttributes.kbNumber != "" ||
		objSoftAttributes.parentKeyName != "" ||
		objSoftAttributes.systemComponent == 1 {
		return false
	}

	for _, releaseValue := range releaseTypes {
		if strings.EqualFold(objSoftAttributes.releaseType, releaseValue) {
			return false
		}
	}
	return true
}

func (installSoftwareImpl) convertInstallDateToTime(installDate string) (tm time.Time, err error) {

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

func (installSoftwareImpl) getSoftwareRegistryProperties(regPath string, access32or64 uint32) (*softwareAttributes, error) {
	var objSoftAttributes softwareAttributes

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.QUERY_VALUE|registry.READ|access32or64)
	if err != nil {
		return nil, err
	}
	defer key.Close()

	//Get the required properties
	objSoftAttributes.displayName, _, _ = key.GetStringValue("DisplayName")
	objSoftAttributes.displayIcon, _, _ = key.GetStringValue("DisplayIcon")
	objSoftAttributes.displayVersion, _, _ = key.GetStringValue("DisplayVersion")
	objSoftAttributes.installDate, _, _ = key.GetStringValue("InstallDate")
	objSoftAttributes.publisher, _, _ = key.GetStringValue("Publisher")
	objSoftAttributes.uninstallString, _, _ = key.GetStringValue("UninstallString")
	objSoftAttributes.kbNumber, _, _ = key.GetStringValue("KBNumber")
	objSoftAttributes.parentKeyName, _, _ = key.GetStringValue("ParentKeyName")
	objSoftAttributes.releaseType, _, _ = key.GetStringValue("ReleaseType")
	objSoftAttributes.systemComponent, _, _ = key.GetIntegerValue("SystemComponent")
	objSoftAttributes.windowsInstaller, _, _ = key.GetIntegerValue("WindowsInstaller")

	return &objSoftAttributes, nil
}

func (installSoftwareImpl) getSoftwareRegistrySubKeys(regPath string, access32or64 uint32) ([]string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.QUERY_VALUE|registry.READ|registry.ENUMERATE_SUB_KEYS|access32or64)
	if err != nil {
		return nil, err
	}
	defer key.Close()

	regStat, err := key.Stat()
	if err != nil {
		return nil, err
	}

	subKeys, err := key.ReadSubKeyNames(int(regStat.SubKeyCount))
	if err != nil {
		return nil, err
	}

	return subKeys, nil
}
