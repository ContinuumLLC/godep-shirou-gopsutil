package linux

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	envMock "github.com/ContinuumLLC/platform-common-lib/src/env/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	pp "github.com/ContinuumLLC/platform-common-lib/src/procParser"
	procMock "github.com/ContinuumLLC/platform-common-lib/src/procParser/mock"
	"github.com/golang/mock/gomock"
)

func TestGetAssetCollectionParseError(t *testing.T) {
	parseError := "Parse Error"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assetDal, _ := setupConfigMocks(ctrl, errors.New(parseError), nil)
	_, err := assetDal.GetAssetData()
	if err == nil || err.Error() != parseError {
		t.Error("Error expected but not returned")
	}
}

func setupConfigMocks(ctrl *gomock.Controller, parseError error, parseData *pp.Data) (*AssetCollectionDalLinux, *procMock.MockParser) {
	mockParser := procMock.NewMockParser(ctrl)
	mockDep := mock.NewMockAssetCollectionDalDependencies(ctrl)
	assetDal := new(AssetCollectionDalLinux)
	assetDal.Factory = mockDep
	assetDal.Logger = logging.GetLoggerFactory().New("")
	assetDal.Logger.SetLogLevel(logging.OFF)
	err := parseError
	mockEnv := envMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte(""))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, err)
	mockDep.EXPECT().GetEnv().Return(mockEnv)
	mockDep.EXPECT().GetParser().Return(mockParser).AnyTimes()
	mockParser.EXPECT().Parse(gomock.Any(), gomock.Any()).Return(parseData, err).AnyTimes()
	return assetDal, mockParser
}

// func TestGetAssetCollectionDal(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	deps := mock.NewMockAssetCollectionDalDependencies(ctrl)
// 	dal := AssetCollectionDalFactoryImpl{}.GetAssetCollectionDal(deps)

// 	if dal == nil {
// 		t.Error("Dal not initialized")
// 	}
// }

func TestGetDataFromMap(t *testing.T) {
	data := new(pp.Data)

	data.Map = make(map[string]pp.Line, 1)
	data.Map["MemTotal"] = pp.Line{Values: []string{"MemTotal", "InvalidNumber", "KB"}}

	tm := translateAssetCollection{
		logger: logging.GetLoggerFactory().New(""),
	}
	tm.logger.SetLogLevel(logging.OFF)
	val := tm.getDataFromMap("MemTotal", data)

	if val != 0 {
		t.Errorf("Expected 0, returned %d", val)
	}
}
func TestGetAssetCollectionParseDataGetBytesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := new(pp.Data)

	data.Map = make(map[string]pp.Line, 1)
	data.Map["MemTotal"] = pp.Line{Values: []string{"physicalTotalBytes", "1", "KB"}}

	assetDal, _ := setupConfigMocks(ctrl, nil, data)
	_, err := assetDal.GetAssetData()
	if err != nil {
		t.Error("Unexpected Error")
	}
}
