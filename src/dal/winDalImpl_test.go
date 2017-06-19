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
