package msgl

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	apiModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	mockJson "github.com/ContinuumLLC/platform-common-lib/src/json/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
	"github.com/golang/mock/gomock"
)

func createMock(ctrl *gomock.Controller, processError error, serializeError error, respBodyStr string) *mock.MockAssetServiceDependencies {
	mockServiceDep := mock.NewMockAssetServiceDependencies(ctrl)

	mockAssetCollectionService := mock.NewMockAssetCollectionService(ctrl)

	jsonMock := mockJson.NewMockSerializerJSON(ctrl)
	jsonMock.EXPECT().WriteByteStream(gomock.Any()).Return([]byte(respBodyStr), serializeError)
	mockServiceDep.EXPECT().GetSerializerJSON().Return(jsonMock)

	assetCollectionData := apiModel.AssetCollection{}
	mockServiceDep.EXPECT().GetAssetCollectionService(gomock.Any()).Return(mockAssetCollectionService)

	mockAssetCollectionService.EXPECT().Process().Return(&assetCollectionData, processError)

	mockServiceDep.EXPECT().GetAssetCollectionServiceDependencies().Return(mockServiceDep)

	confMock := mock.NewMockConfigService(ctrl)
	confMock.EXPECT().GetAssetPluginConfig().Return(&model.AssetPluginConfig{}, nil)
	mockServiceDep.EXPECT().GetConfigService(gomock.Any()).Return(confMock).AnyTimes()

	return mockServiceDep
}

func GetProcessAssetTest(t *testing.T) {
	assetProcFact := ProcessAssetFactoryImpl{}

	asset := assetProcFact.GetHandler(nil, &model.AssetPluginConfig{})
	if asset == nil {
		t.Error("New Handler struct expected, not returned")
	}
}

func TestAssetCollection(t *testing.T) {
	ctrl := gomock.NewController(t)
	respBodyStr := "testoutput"

	mockServiceDep := createMock(ctrl, nil, nil, respBodyStr)
	processAssetFact := ProcessAssetFactoryImpl{}
	processAsset := processAssetFact.GetHandler(mockServiceDep, &model.AssetPluginConfig{})
	req := createRequest()
	req.Path = "/assetCollection"
	resp, err := processAsset.HandleAsset(req)
	if err != nil {
		t.Errorf("Unexpected error returned %v", err)
	}
	if resp.Status != protocol.Ok {
		t.Errorf("Unexpected response status returned %v", resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unexpected error returned %v", err)
	}

	if string(data) != respBodyStr {
		t.Errorf("Unexpected data returned, expected data: %s, returned data: %s", respBodyStr, data)
	}

}

func TestAssetProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	respBodyStr := "testoutput"

	mockServiceDep := createMock(ctrl, nil, nil, respBodyStr)
	processAssetFact := ProcessAssetFactoryImpl{}
	processAsset := processAssetFact.GetHandler(mockServiceDep, &model.AssetPluginConfig{})
	req := createRequest()
	req.Path = "/assetCollection"
	resp, err := processAsset.HandleAsset(req)
	if err != nil {
		t.Errorf("Unexpected error returned %v", err)
	}
	if resp.Status != protocol.Ok {
		t.Errorf("Unexpected response status returned %v", resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unexpected error returned %v", err)
	}

	if string(data) != respBodyStr {
		t.Errorf("Unexpected data returned, expected data: %s, returned data: %s", respBodyStr, data)
	}

}

func TestAssetConfigurationGetAssetPluginConfMapError(t *testing.T) {
	ctrl := gomock.NewController(t)

	//config := make(map[string]interface{})
	getAssetPluginConfMapError := errors.New("GetAssetPluginConfMapError")

	dep := mock.NewMockAssetServiceDependencies(ctrl)
	serv := mock.NewMockConfigService(ctrl)
	serv.EXPECT().GetAssetPluginConfMap().Return(nil, getAssetPluginConfMapError)
	dep.EXPECT().GetConfigService(gomock.Any()).Return(serv)

	processAssetFact := ProcessAssetFactoryImpl{}
	processAsset := processAssetFact.GetHandler(dep, &model.AssetPluginConfig{})

	req := createRequest()
	req.Path = "/asset/configuration"
	data := `{"AgentServiceURL": "http://localhost:8081",
        "CommunicationBufferChannelLimit":1,
        "CommunicationMaxDataToRetrieve":1,
        "EndPointID":"e1",
		"HeartBeatCounterBaseValue":  100,
        "HeartBeatCounterMaxValue":1000,
		"HeartbeatPluginPath" : "Path1",		
		"PluginsLocation":
			{"asset1":"./asset","asset1":"/asset", "asset":"./asset","c":"d"}
		
	}`
	reader := bytes.NewReader([]byte(data))
	req.Body = reader
	_, err := processAsset.HandleConfig(req)
	if err == nil {
		t.Errorf("Expected error not returned, Expected:%v", getAssetPluginConfMapError)
	}

	if err != getAssetPluginConfMapError {
		t.Errorf("Unexpected error returned, Expected:%v, Returned:%v", getAssetPluginConfMapError, err)
	}
}

func TestAssetConfigurationSetAssetPluginConfMapError(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := make(map[string]interface{})
	setAssetPluginConfMapError := errors.New("SetAssetPluginConfMapError")

	dep := mock.NewMockAssetServiceDependencies(ctrl)
	serv := mock.NewMockConfigService(ctrl)
	serv.EXPECT().GetAssetPluginConfMap().Return(config, nil)
	serv.EXPECT().SetAssetPluginMap(gomock.Any()).Return(setAssetPluginConfMapError)

	dep.EXPECT().GetConfigService(gomock.Any()).Return(serv)

	processAssetFact := ProcessAssetFactoryImpl{}
	processAsset := processAssetFact.GetHandler(dep, &model.AssetPluginConfig{})

	req := createRequest()
	req.Path = "/asset/configuration"
	data := `{"AgentServiceURL": "http://localhost:8081",
        "CommunicationBufferChannelLimit":1,
        "CommunicationMaxDataToRetrieve":1,
        "EndPointID":"e1",
		"HeartBeatCounterBaseValue":  100,
        "HeartBeatCounterMaxValue":1000,
		"HeartbeatPluginPath" : "Path1",		
		"PluginsLocation":
			{"asset1":"./asset","asset1":"/asset", "asset":"./asset","c":"d"}
		
	}`
	reader := bytes.NewReader([]byte(data))
	req.Body = reader
	_, err := processAsset.HandleConfig(req)
	if err == nil {
		t.Errorf("Expected error not returned, Expected:%v", setAssetPluginConfMapError)
	}

	if err != setAssetPluginConfMapError {
		t.Errorf("Unexpected error returned, Expected:%v, Returned:%v", setAssetPluginConfMapError, err)
	}
}
func TestAssetConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := make(map[string]interface{})

	dep := mock.NewMockAssetServiceDependencies(ctrl)
	serv := mock.NewMockConfigService(ctrl)
	serv.EXPECT().GetAssetPluginConfMap().Return(config, nil)
	serv.EXPECT().SetAssetPluginMap(gomock.Any()).Return(nil)

	dep.EXPECT().GetConfigService(gomock.Any()).Return(serv)

	//processAssetFact := ProcessAssetFactoryImpl{}
	//processAsset := processAssetFact.GetHandler(dep, &model.AssetPluginConfig{})
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	ps := processAsset{
		dep:    dep,
		cfg:    &model.AssetPluginConfig{},
		logger: log,
	}

	req := createRequest()
	req.Path = "/asset/configuration"
	data := `{"AgentServiceURL": "http://localhost:8081",
        "CommunicationBufferChannelLimit":1,
        "CommunicationMaxDataToRetrieve":1,
        "EndPointID":"e1",
		"HeartBeatCounterBaseValue":  100,
        "HeartBeatCounterMaxValue":1000,
		"HeartbeatPluginPath" : "Path1",		
		"PluginsLocation":
			{"asset1":"./asset","asset1":"/asset", "asset":"./asset","c":"d"}
		
	}`
	reader := bytes.NewReader([]byte(data))
	req.Body = reader
	resp, err := ps.HandleConfig(req)
	if err != nil {
		t.Errorf("Unexpected error not returned, Expected:%v", err)
	}

	if resp.Status != protocol.Ok {
		t.Errorf("Unexpected status returned, Expected:%v, Returned:%v", protocol.Ok, resp.Status)
	}
}

func createRequest() *protocol.Request {
	request := protocol.NewRequest()
	request.Headers.SetKeyValue(protocol.HdrUserAgent, "AgentCore")
	request.Headers.SetKeyValue(protocol.HdrContentType, "text/json")
	return request
}
