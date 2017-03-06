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
		mockParser.EXPECT().Parse(gomock.Any(), gomock.Any()).Return(&procParser.Data{}, parseErr)
	}
	mockAssetDalD.EXPECT().GetParser().Return(mockParser)

	return ctrl, mockAssetDalD
}

func TestGetOS(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader(t, nil, errors.New(model.ErrExecuteCommandFailed))
	defer ctrl.Finish()

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, err := AssetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetOS()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}
