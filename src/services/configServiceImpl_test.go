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

func setupGetPerfPluginConfig(t *testing.T, pConf *model.PerfPluginConfig, err error) (*gomock.Controller, configServiceImpl) {
	ctrl := gomock.NewController(t)
	confSrvcDepMock := mock.NewMockConfigServiceDependencies(ctrl)
	confDalMock := mock.NewMockConfigDal(ctrl)
	confDalMock.EXPECT().GetPerfPluginConf().Return(pConf, err)
	confSrvcDepMock.EXPECT().GetConfigDal(gomock.Any()).Return(confDalMock)
	csrvc := configServiceImpl{
		factory: confSrvcDepMock,
	}
	//_, e := csrvc.GetPerfPluginConfig()
	return ctrl, csrvc
}

func TestGetPerfPluginConfigNoErr(t *testing.T) {
	sConfig = nil // setting this static variable to nil to make it re-initialize
	pConf := model.PerfPluginConfig{
		PluginPath: model.PerfPluginPath{
			Memory: "/asset/memory",
		},
		URLSuffix: make(map[string]string),
	}
	ctrl, csrvc := setupGetPerfPluginConfig(t, &pConf, nil)
	_, err := csrvc.GetPerfPluginConfig()
	defer ctrl.Finish()
	if err != nil {
		t.Errorf("Unexpected error : %v", err)
	}
}

func TestGetPerfPluginConfigErr(t *testing.T) {
	sConfig = nil // setting this static variable to nil to make it re-initialize
	ctrl, csrvc := setupGetPerfPluginConfig(t, nil, errors.New(model.ErrPerfPluginConfig))
	_, err := csrvc.GetPerfPluginConfig()
	defer ctrl.Finish()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrPerfPluginConfig) {
		t.Errorf("Expected error is %s but got %v", model.ErrPerfPluginConfig, err)
	}
}
