package dal

import (
	"testing"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	"github.com/golang/mock/gomock"
)

func TestGetAssetCollectionDal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deps := mock.NewMockAssetCollectionDalDependencies(ctrl)
	dal := AssetDalFactoryImpl{}.GetAssetDal(deps)

	if dal == nil {
		t.Error("Dal not initialized")
	}
}
