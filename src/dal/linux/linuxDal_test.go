package linux

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"strings"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	eMock "github.com/ContinuumLLC/platform-common-lib/src/env/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
	pMock "github.com/ContinuumLLC/platform-common-lib/src/procParser/mock"
	"github.com/golang/mock/gomock"
)

func setupGetCommandReader(t *testing.T, parseErr error, commandReaderErr error) (*gomock.Controller, *mock.MockAssetDalDependencies) {
	ctrl := gomock.NewController(t)
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any()).Return(reader, commandReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if commandReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(&procParser.Data{}, parseErr)
	}
	mockAssetDalD.EXPECT().GetParser().Return(mockParser)

	return ctrl, mockAssetDalD
}

func setupGetFileReader(t *testing.T, parseErr error, fileReaderErr error) (*gomock.Controller, *mock.MockAssetDalDependencies) {
	ctrl := gomock.NewController(t)
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, fileReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if fileReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(&procParser.Data{}, parseErr)
	}
	mockAssetDalD.EXPECT().GetParser().Return(mockParser)

	return ctrl, mockAssetDalD
}

func setupAddGetFileReader(ctrl *gomock.Controller, mockAssetDalD *mock.MockAssetDalDependencies, parseErr error, fileReaderErr error) {
	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, fileReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if fileReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(&procParser.Data{}, parseErr)
	}
}

func TestGetOSCommandErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader(t, nil, errors.New(model.ErrExecuteCommandFailed))
	defer ctrl.Finish()

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, err := AssetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetOSInfo()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}

func TestGetOSFileErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader(t, nil, nil)
	defer ctrl.Finish()

	setupAddGetFileReader(ctrl, mockAssetDalD, nil, errors.New(model.ErrFileReadFailed))

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, err := AssetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetOSInfo()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrFileReadFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrFileReadFailed, err)
	}
}

// TODO - fix error
// func TestGetOSNoErr(t *testing.T) {
// 	ctrl, mockAssetDalD := setupGetCommandReader(t, nil, nil)
// 	defer ctrl.Finish()

// 	setupAddGetFileReader(ctrl, mockAssetDalD, nil, nil)

// 	log := logging.GetLoggerFactory().New("")
// 	log.SetLogLevel(logging.OFF)
// 	_, err := AssetDalImpl{
// 		Factory: mockAssetDalD,
// 		Logger:  log,
// 	}.GetOS()
// 	if err != nil {
// 		t.Errorf("Unexpected error : %v", err)
// 	}
// }

func setupGetSystemInfo(t *testing.T, times int, err error) (*gomock.Controller, error) {
	ctrl := gomock.NewController(t)

	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	var str string
	switch times {
	case 1:
		str = cSysProductCmd
	case 2:
		mockEnv.EXPECT().ExecuteBash(cSysProductCmd).Return("", nil)
		str = cSysTz
	case 3:
		mockEnv.EXPECT().ExecuteBash(cSysProductCmd).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTz).Return("", nil)
		str = cSysTzd

	}
	mockEnv.EXPECT().ExecuteBash(str).Return("", err)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv).Times(times)

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, e := AssetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetSystemInfo()
	return ctrl, e
}

func TestGetSystemInfoErr(t *testing.T) {
	cmdExeArr := []int{1, 2, 3}
	for _, i := range cmdExeArr {
		ctrl, err := setupGetSystemInfo(t, i, errors.New(model.ErrExecuteCommandFailed))
		defer ctrl.Finish()
		if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
			t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
		}
	}
}

func TestGetSystemNoErr(t *testing.T) {
	ctrl, err := setupGetSystemInfo(t, 3, nil)
	defer ctrl.Finish()
	if err != nil {
		t.Errorf("Unexpected error received  : %v", err)
	}
}
