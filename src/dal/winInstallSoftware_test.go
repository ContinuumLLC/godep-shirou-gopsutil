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
