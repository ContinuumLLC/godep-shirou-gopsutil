package services

import (
	"errors"
	"testing"

	"strings"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	"github.com/golang/mock/gomock"
)

func TestGetConfigService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	confSrvcDepMock := mock.NewMockConfigServiceDependencies(ctrl)
	cs := ConfigServiceFactoryImpl{}.GetConfigService(confSrvcDepMock)
	if cs == nil {
		t.Error("Concrete implementation of ConfigService is expected")
	}
}

func setupGetAssetPluginConfig(t *testing.T, pConf *model.AssetPluginConfig, err error) (*gomock.Controller, configServiceImpl) {
	ctrl := gomock.NewController(t)
	confSrvcDepMock := mock.NewMockConfigServiceDependencies(ctrl)
	confDalMock := mock.NewMockConfigDal(ctrl)
	confDalMock.EXPECT().GetAssetPluginConf().Return(pConf, err)
	confSrvcDepMock.EXPECT().GetConfigDal(gomock.Any()).Return(confDalMock)
	csrvc := configServiceImpl{
		factory: confSrvcDepMock,
	}
	//_, e := csrvc.GetAssetPluginConfig()
	return ctrl, csrvc
}

func TestGetAssetPluginConfigNoErr(t *testing.T) {
	sConfig = nil // setting this static variable to nil to make it re-initialize
	pConf := model.AssetPluginConfig{
		PluginPath: model.AssetPluginPath{
			AssetCollection: "/asset",
		},
		URLSuffix: make(map[string]string),
	}
	ctrl, csrvc := setupGetAssetPluginConfig(t, &pConf, nil)
	_, err := csrvc.GetAssetPluginConfig()
	defer ctrl.Finish()
	if err != nil {
		t.Errorf("Unexpected error : %v", err)
	}
}

func TestGetAssetPluginConfMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	confSrvcDepMock := mock.NewMockConfigServiceDependencies(ctrl)
	confDalMock := mock.NewMockConfigDal(ctrl)
	confDalMock.EXPECT().GetAssetPluginConfMap().Return(nil, nil)
	confSrvcDepMock.EXPECT().GetConfigDal(gomock.Any()).Return(confDalMock)
	csrvc := &configServiceImpl{
		factory: confSrvcDepMock,
	}
	_, err := csrvc.GetAssetPluginConfMap()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

}

func TestSetAssetPluginMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	confSrvcDepMock := mock.NewMockConfigServiceDependencies(ctrl)
	confDalMock := mock.NewMockConfigDal(ctrl)
	var a map[string]interface{}
	confDalMock.EXPECT().SetAssetPluginMap(a).Return(nil)
	confSrvcDepMock.EXPECT().GetConfigDal(gomock.Any()).Return(confDalMock)

	csrvc := &configServiceImpl{
		factory: confSrvcDepMock,
	}
	err := csrvc.SetAssetPluginMap(a)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

}

func TestGetAssetPluginConfigErr(t *testing.T) {
	sConfig = nil // setting this static variable to nil to make it re-initialize
	ctrl, csrvc := setupGetAssetPluginConfig(t, nil, errors.New(model.ErrAssetPluginConfig))
	_, err := csrvc.GetAssetPluginConfig()
	defer ctrl.Finish()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrAssetPluginConfig) {
		t.Errorf("Expected error is %s but got %v", model.ErrAssetPluginConfig, err)
	}
}
