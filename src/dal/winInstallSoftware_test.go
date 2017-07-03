// +build windows

package dal

import (
	"testing"
)

func TestValidatePropertiesForInstallSoftwareTrue(t *testing.T) {
	var objInstallSoft installSoftwareImpl
	objSoftAttributes := softwareAttributes{
		displayName:     "ABC Soft",
		uninstallString: "ABC Uninstall",
	}

	if !objInstallSoft.validatePropertiesForInstallSoftware(objSoftAttributes) {
		t.Error("Expected true, but returned false")
	}
}

func TestValidatePropertiesForInstallSoftwareFalse(t *testing.T) {
	var objInstallSoft installSoftwareImpl
	objSoftAttributes := softwareAttributes{
		displayName: "ABC Soft",
	}

	if objInstallSoft.validatePropertiesForInstallSoftware(objSoftAttributes) {
		t.Error("Expected false, but returned true")
	}
}

func TestConvertInstallDateToTime(t *testing.T) {
	var objInstallSoft installSoftwareImpl
	installDate := "20160725"

	_, err := objInstallSoft.convertInstallDateToTime(installDate)
	if nil != err {
		t.Errorf("Expected time object, but received error %v", err)
	}
}

func TestConvertInstallDateToTimeErr(t *testing.T) {
	var objInstallSoft installSoftwareImpl
	installDate := "20"

	_, err := objInstallSoft.convertInstallDateToTime(installDate)
	if nil == err {
		t.Error("Expected error")
	}
}
