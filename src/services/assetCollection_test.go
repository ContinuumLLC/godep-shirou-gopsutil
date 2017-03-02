package services

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	amodel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
        "github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	"github.com/golang/mock/gomock"
)

func TestGetAssetCollectionService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPerfSrvcD := mock.NewMockAssetCollectionServiceDependencies(ctrl)
	assetSrvc := AssetCollectionServiceFactoryImpl{}.GetAssetService(mockPerfSrvcD)
	typ, ok := assetSrvc.(assetCollectionServiceImpl)
	if !ok {
		t.Errorf("Expected type is assetCollectionServiceImpl but got %v", reflect.TypeOf(typ))
	}
}

func TestProcessGetAssetCollectionErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPerfSrvcD := mock.NewMockAssetCollectionServiceDependencies(ctrl)
	mockPerfMemDal := mock.NewMockAssetCollectionDal(ctrl)
	mockPerfSrvcD.EXPECT().GetAssetDal(gomock.Any()).Return(mockPerfMemDal)

	setupErr := errors.New("GetAssetData error")
	mockPerfMemDal.EXPECT().GetAssetData().Return(&amodel.AssetCollection{}, setupErr)
	srvc := &assetCollectionServiceImpl{
		factory: mockPerfSrvcD,
	}
	_, err := srvc.Process()
	if err == nil {
		t.Errorf("Expected error %v", setupErr)
		return
	}
	if !strings.HasPrefix(err.Error(), setupErr.Error()) {
		t.Errorf("Expected error is %v but got %v", setupErr, err)
	}
}
