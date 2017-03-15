package dal

import (
	"errors"
	"testing"

	"strings"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	mockClar "github.com/ContinuumLLC/platform-common-lib/src/clar/mock"
	mockEnv "github.com/ContinuumLLC/platform-common-lib/src/env/mock"
	mockJson "github.com/ContinuumLLC/platform-common-lib/src/json/mock"
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

func setupAssetPluginConfMap(t *testing.T, err error) (*gomock.Controller, *mock.MockConfigDalDependencies) {
	ctrl := gomock.NewController(t)
	confDalMock := mock.NewMockConfigDalDependencies(ctrl)
	envMock := mockEnv.NewMockEnv(ctrl)
	envMock.EXPECT().GetExeDir().Return("", err)
	confDalMock.EXPECT().GetEnv().Return(envMock)
	return ctrl, confDalMock
}

func TestGetAssetPluginConfMapGetExeErr(t *testing.T) {
	errMsg := "GetExeDirErr"
	ctrl, confDalMock := setupAssetPluginConfMap(t, errors.New(errMsg))
	defer ctrl.Finish()
	_, err := configDalImpl{
		factory: confDalMock,
	}.GetAssetPluginConfMap()
	if err == nil || !strings.HasPrefix(err.Error(), errMsg) {
		t.Errorf("Expected error is :%v but got : %v", errMsg, err)
	}
}

//GetAssetPluginConfMap returning error doing DeserializerJSON
func TestGetAssetPluginConfMapReadFileErr(t *testing.T) {
	errMsg := "ReadFileErr"
	ctrl, confDalMock := setupAssetPluginConfMap(t, nil)
	defer ctrl.Finish()

	clarMock := mockClar.NewMockServiceInit(ctrl)
	clarMock.EXPECT().GetConfigPath().Return("") //returns any filename
	confDalMock.EXPECT().GetServiceInit().Return(clarMock)

	jsonMock := mockJson.NewMockDeserializerJSON(ctrl)
	jsonMock.EXPECT().ReadFile(gomock.Any(), gomock.Any()).Return(errors.New(errMsg))
	confDalMock.EXPECT().GetDeserializerJSON().Return(jsonMock)
	_, err := configDalImpl{
		factory: confDalMock,
	}.GetAssetPluginConfMap()
	if err == nil || !strings.HasPrefix(err.Error(), errMsg) {
		t.Errorf("Expected error is :%v but got : %v", errMsg, err)
	}
}

func TestAssetPluginMapGood(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	confDalMock := mock.NewMockConfigDalDependencies(ctrl)

	clarMock := mockClar.NewMockServiceInit(ctrl)
	clarMock.EXPECT().GetConfigPath().Return("") //returns any filename
	confDalMock.EXPECT().GetServiceInit().Return(clarMock)

	jsonMock := mockJson.NewMockSerializerJSON(ctrl)
	jsonMock.EXPECT().WriteFile(gomock.Any(), gomock.Any()).Return(nil)
	confDalMock.EXPECT().GetSerializerJSON().Return(jsonMock)

	var tt map[string]interface{}
	err := configDalImpl{
		factory: confDalMock,
	}.SetAssetPluginMap(tt)
	if err != nil {
		t.Errorf("Unexpected error :%v ", err)
	}
}
