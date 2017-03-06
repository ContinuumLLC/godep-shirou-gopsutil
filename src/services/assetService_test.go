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

	mockPerfSrvcD := mock.NewMockAssetServiceDependencies(ctrl)
	assetSrvc := AssetServiceFactoryImpl{}.GetAssetService(mockPerfSrvcD)
	typ, ok := assetSrvc.(assetServiceImpl)
	if !ok {
		t.Errorf("Expected type is assetServiceImpl but got %v", reflect.TypeOf(typ))
	}
}

func TestProcessGetAssetCollectionErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPerfSrvcD := mock.NewMockAssetServiceDependencies(ctrl)
	mockPerfMemDal := mock.NewMockAssetDal(ctrl)
	mockPerfSrvcD.EXPECT().GetAssetDal(gomock.Any()).Return(mockPerfMemDal)

	setupErr := errors.New("GetAssetData error")
	mockPerfMemDal.EXPECT().GetAssetData().Return(&amodel.AssetCollection{}, setupErr)
	srvc := &assetServiceImpl{
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
