package dal

import (
	"errors"
	"testing"

	"strings"

	mockClar "github.com/ContinuumLLC/platform-common-lib/src/clar/mock"
	mockJson "github.com/ContinuumLLC/platform-common-lib/src/json/mock"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	"github.com/golang/mock/gomock"
)

func TestGetConfigDal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cDalDepMock := mock.NewMockConfigDalDependencies(ctrl)
	d := ConfigDalFactoryImpl{}.GetConfigDal(cDalDepMock)

	_, ok := d.(configDalImpl)
	if !ok {
		t.Errorf("Expected was configDalImpl")
	}
}

func setupPerPluginConf(t *testing.T, err error) (*gomock.Controller, error) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cDalDepMock := mock.NewMockConfigDalDependencies(ctrl)
	jsonMockObj := mockJson.NewMockDeserializerJSON(ctrl)
	jsonMockObj.EXPECT().ReadFile(gomock.Any(), "").Return(err)
	cDalDepMock.EXPECT().GetDeserializerJSON().Return(jsonMockObj)

	serverInitMock := mockClar.NewMockServiceInit(ctrl)
	serverInitMock.EXPECT().GetConfigPath().Return("")
	cDalDepMock.EXPECT().GetServiceInit().Return(serverInitMock)

	_, e := configDalImpl{
		factory: cDalDepMock,
	}.GetAssetPluginConf()
	return ctrl, e
}

func TestGetPerfPluginConfErr(t *testing.T) {
	errMsg := "ReadFile Error"
	ctrl, err := setupPerPluginConf(t, errors.New(errMsg))
	defer ctrl.Finish()
	if err == nil || !strings.HasPrefix(err.Error(), errMsg) {
		t.Errorf("Expected error is %s but got %v", errMsg, err)
	}
}

func TestGetPerfPluginConfNoErr(t *testing.T) {
	ctrl, err := setupPerPluginConf(t, nil)
	defer ctrl.Finish()
	if err != nil {
		t.Errorf("Unexpected error : %v", err)
	}
}
