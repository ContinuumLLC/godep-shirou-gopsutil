package msgl

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol/http"
	"github.com/golang/mock/gomock"
)

func TestGetAssetListener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assetListFact := AssetListenerFactory{}
	assetSrvcDepMock := mock.NewMockHandlerDependencies(ctrl)
	confMock := mock.NewMockConfigService(ctrl)
	confMock.EXPECT().GetAssetPluginConfig().Return(&model.AssetPluginConfig{}, nil)
	assetSrvcDepMock.EXPECT().GetConfigService(gomock.Any()).Return(confMock)
	assetList := assetListFact.GetAssetListener(assetSrvcDepMock)
	if assetList == nil {
		t.Error("Expected AssetListener returned nil")
	}
}

func TestProcessReceiveRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)

	_, resw, _ := os.Pipe()

	respBodyString := "testResponse"

	servDep := createMock(ctrl, nil, nil, respBodyString)
	servDep.EXPECT().GetReader().Return(resw) //send writer in instead of reader to  trigger an error
	servDep.EXPECT().GetWriter().Return(resw)

	processPerf := ProcessAssetFactoryImpl{}.GetHandler(servDep, &model.AssetPluginConfig{})

	servDep.EXPECT().GetHandler(gomock.Any(), &model.AssetPluginConfig{}).Return(processPerf).AnyTimes()

	httpServer := http.ServerHTTPFactory{}
	server := httpServer.GetServer(resw, resw)
	servDep.EXPECT().GetServer(gomock.Any(), gomock.Any()).Return(server)

	assetListFact := AssetListenerFactory{}
	assetList := assetListFact.GetAssetListener(servDep)

	err := assetList.Process()
	//checking what the error is for badWrite, since this can be different for different platforms
	_, errBadWrite := ioutil.ReadAll(resw)

	if err == nil {
		t.Errorf("Expected error not returned, Expected Error: %v", err)
		return
	}
	if errBadWrite.Error() != err.Error() {
		t.Errorf("Unexpected error returned, Expected: %v, Returned: %v", errBadWrite, err)
	}
}

func TestProcessIncorrectRouteError(t *testing.T) {
	ctrl := gomock.NewController(t)

	reqr, reqw, _ := os.Pipe()
	resr, resw, _ := os.Pipe()
	respBodyString := "testResponse"

	servDep := createMock(ctrl, nil, nil, respBodyString)
	servDep.EXPECT().GetReader().Return(reqr)
	servDep.EXPECT().GetWriter().Return(resw)

	processPerf := ProcessAssetFactoryImpl{}.GetHandler(servDep, &model.AssetPluginConfig{})

	servDep.EXPECT().GetHandler(gomock.Any(), &model.AssetPluginConfig{}).Return(processPerf).AnyTimes()

	assetListFact := AssetListenerFactory{}
	assetList := assetListFact.GetAssetListener(servDep)

	client := http.ClientHTTPFactory{}.GetClient(reqw, resr)
	request := createDummyRequest("/testRouote/IncorrectRoute")

	httpServer := http.ServerHTTPFactory{}
	server := httpServer.GetServer(reqr, resw)
	servDep.EXPECT().GetServer(gomock.Any(), gomock.Any()).Return(server)

	client.SendRequest(request)
	reqw.Close()

	err := assetList.Process()
	if err == nil {
		t.Errorf("Expected error not returned, Expected Error: %v", err)
		return
	}
	if !strings.HasPrefix(err.Error(), ErrInvalidPluginPath) {
		t.Errorf("Unexpected error returned, Expected: %s, Returned: %v", ErrInvalidPluginPath, err)
	}
}

func createDummyRequest(path string) *protocol.Request {
	request := protocol.NewRequest()
	request.Path = path
	request.Headers.SetKeyValue(protocol.HdrUserAgent, "testAgentCore")
	request.Headers.SetKeyValue(protocol.HdrContentType, "text/json")
	return request
}
