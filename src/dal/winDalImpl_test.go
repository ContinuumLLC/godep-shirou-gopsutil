// +build windows

package dal

import (
	"testing"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/golang/mock/gomock"
)

// TODO: Unit testing strategy for gopsutil/win32 system call
func TestGetAssetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAssetDalDep := mock.NewMockAssetDalDependencies(ctrl)

	_, err := assetDalImpl{
		Factory: mockAssetDalDep,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetAssetData()

	if err != nil {
		t.Errorf("Could not get AssetData %v", err)
	}
}

func TestConvertInstallDateToTime(t *testing.T) {
	var objInstallSoft assetDalImpl
	installDate := "20160725"

	_, err := objInstallSoft.convertInstallDateToTime(installDate)
	if nil != err {
		t.Errorf("Expected time object, but received error %v", err)
	}
}

func TestConvertInstallDateToTimeErr(t *testing.T) {
	var objInstallSoft assetDalImpl
	installDate := "20"

	_, err := objInstallSoft.convertInstallDateToTime(installDate)
	if nil == err {
		t.Error("Expected error")
	}
}
