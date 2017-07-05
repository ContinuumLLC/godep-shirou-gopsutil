package dal

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

var releaseTypes = [...]string{"Hotfix", "Security Update", "Update Rollup", "Update", "Critical Update", "Security"}

const (
	sysComponent = 1
)

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

func (installSoftwareImpl) validatePropertiesForInstallSoftware(objAttr softwareAttributes) bool {
	if objAttr.displayName == "" ||
		objAttr.uninstallString == "" ||
		objAttr.kbNumber != "" ||
		objAttr.parentKeyName != "" ||
		objAttr.systemComponent == sysComponent {
		return false
	}

	for _, releaseValue := range releaseTypes {
		if strings.EqualFold(objAttr.releaseType, releaseValue) {
			return false
		}
	}
	return true
}

func (installSoftwareImpl) getSoftwareRegistryProperties(regPath string, access32or64 uint32) (*softwareAttributes, error) {
	var objAttr softwareAttributes

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.QUERY_VALUE|registry.READ|access32or64)
	if err != nil {
		return nil, err
	}
	defer key.Close()

	//Get the required properties
	objAttr.displayName, _, _ = key.GetStringValue("DisplayName")
	objAttr.displayIcon, _, _ = key.GetStringValue("DisplayIcon")
	objAttr.displayVersion, _, _ = key.GetStringValue("DisplayVersion")
	objAttr.installDate, _, _ = key.GetStringValue("InstallDate")
	objAttr.publisher, _, _ = key.GetStringValue("Publisher")
	objAttr.uninstallString, _, _ = key.GetStringValue("UninstallString")
	objAttr.kbNumber, _, _ = key.GetStringValue("KBNumber")
	objAttr.parentKeyName, _, _ = key.GetStringValue("ParentKeyName")
	objAttr.releaseType, _, _ = key.GetStringValue("ReleaseType")
	objAttr.systemComponent, _, _ = key.GetIntegerValue("SystemComponent")
	objAttr.windowsInstaller, _, _ = key.GetIntegerValue("WindowsInstaller")

	return &objAttr, nil
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

	return key.ReadSubKeyNames(int(regStat.SubKeyCount))
}
