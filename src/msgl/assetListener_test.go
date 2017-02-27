package msgl

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol/http"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	"github.com/golang/mock/gomock"
)

func TestGetAssetListener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assetListFact := AssetListenerFactory{}
	assetSrvcDepMock := mock.NewMockAssetServiceDependencies(ctrl)
	confMock := mock.NewMockConfigService(ctrl)
	confMock.EXPECT().GetAssetPluginConfig().Return(&model.AssetPluginConfig{}, nil)
	assetSrvcDepMock.EXPECT().GetConfigService(gomock.Any()).Return(confMock)
	assetList := assetListFact.GetAssetListener(assetSrvcDepMock)
	if assetList == nil {
		t.Error("Expected AssetListener returned nil")
	}
}

//TODO - Need to fix the test

// func TestProcess(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	reqr, reqw, _ := os.Pipe()
// 	resr, resw, _ := os.Pipe()
// 	respBodyString := "testResponse"

// 	servDep := createMock(ctrl, nil, nil, respBodyString)
// 	servDep.EXPECT().GetReader().Return(reqr)
// 	servDep.EXPECT().GetWriter().Return(resw)

// 	processPerf := ProcessAssetFactoryImpl{}.GetProcessAsset(servDep)

// 	servDep.EXPECT().GetProcessAsset(gomock.Any()).Return(processPerf).AnyTimes()

// 	assetListFact := AssetListenerFactory{}
// 	assetList := assetListFact.GetAssetListener(servDep)

// 	client := http.ClientHTTPFactory{}.GetClient(reqw, resr)
// 	request := createDummyRequest("/asset/processor")

// 	httpServer := http.ServerHTTPFactory{}
// 	server := httpServer.GetServer(reqr, resw)
// 	servDep.EXPECT().GetServer(gomock.Any(), gomock.Any()).Return(server)

// 	client.SendRequest(request)
// 	reqw.Close()

// 	var resp *protocol.Response
// 	var err error
// 	assetList.Process()
// 	resw.Close()
// 	resp, err = client.ReceiveResponse()
// 	if err != nil {
// 		t.Errorf("Unexpected error returned, Error: %v", err)
// 		return
// 	}
// 	if resp.Status != protocol.Ok {
// 		t.Errorf("Unexpected response status returned, Expected: %v, Returned: %v", protocol.Ok, resp.Status)
// 		return
// 	}

// 	data, _ := ioutil.ReadAll(resp.Body)
// 	if string(data) != respBodyString {
// 		t.Errorf("Unexpected response body returned, Expected: %s, Received: %s", respBodyString, string(data))
// 	}

// }

func TestProcessReceiveRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)

	_, resw, _ := os.Pipe()

	respBodyString := "testResponse"

	servDep := createMock(ctrl, nil, nil, respBodyString)
	servDep.EXPECT().GetReader().Return(resw) //send writer in instead of reader to  trigger an error
	servDep.EXPECT().GetWriter().Return(resw)

	processPerf := ProcessAssetFactoryImpl{}.GetProcessAsset(servDep, &model.AssetPluginConfig{})

	servDep.EXPECT().GetProcessAsset(gomock.Any(), &model.AssetPluginConfig{}).Return(processPerf).AnyTimes()

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

	processPerf := ProcessAssetFactoryImpl{}.GetProcessAsset(servDep, &model.AssetPluginConfig{})

	servDep.EXPECT().GetProcessAsset(gomock.Any(), &model.AssetPluginConfig{}).Return(processPerf).AnyTimes()

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

//TODO - Need to fix the test

// func TestProcessHandleError(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	reqr, reqw, _ := os.Pipe()
// 	resr, resw, _ := os.Pipe()
// 	respBodyString := "testResponse"
// 	processErr := errors.New("ProcessServiceError")
// 	servDep := createMock(ctrl, processErr, nil, respBodyString)
// 	servDep.EXPECT().GetReader().Return(reqr)
// 	servDep.EXPECT().GetWriter().Return(resw)

// 	processPerf := ProcessAssetFactoryImpl{}.GetProcessAsset(servDep)

// 	servDep.EXPECT().GetProcessAsset(gomock.Any()).Return(processPerf).AnyTimes()

// 	assetListFact := AssetListenerFactory{}
// 	assetList := assetListFact.GetAssetListener(servDep)

// 	client := http.ClientHTTPFactory{}.GetClient(reqw, resr)
// 	request := createDummyRequest("/asset/processor")

// 	httpServer := http.ServerHTTPFactory{}
// 	server := httpServer.GetServer(reqr, resw)
// 	servDep.EXPECT().GetServer(gomock.Any(), gomock.Any()).Return(server)

// 	client.SendRequest(request)
// 	reqw.Close()

// 	err := assetList.Process()
// 	if err == nil {
// 		t.Errorf("Expected error not returned, Expected Error: %v", err)
// 		return
// 	}
// 	if err.Error() != processErr.Error() {
// 		t.Errorf("Unexpected error returned, Expected: %v, Returned: %v", processErr, err)
// 	}
// }

func createDummyRequest(path string) *protocol.Request {
	request := protocol.NewRequest()
	request.Path = path
	request.Headers.SetKeyValue(protocol.HdrUserAgent, "testAgentCore")
	request.Headers.SetKeyValue(protocol.HdrContentType, "text/json")
	return request
}
