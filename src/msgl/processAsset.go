package msgl

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

//ProcessAssetFactoryImpl returns a asset Processor
type ProcessAssetFactoryImpl struct{}

//GetHandler returns a Processor processor
func (ProcessAssetFactoryImpl) GetHandler(deps model.HandlerDependencies, config *model.AssetPluginConfig) model.Handler {
	return processAsset{
		dep:    deps,
		cfg:    config,
		logger: logging.GetLoggerFactory().New("Handler"),
	}
}

type processAsset struct {
	dep    model.HandlerDependencies
	cfg    *model.AssetPluginConfig
	logger logging.Logger
}

//HandleAsset processes incoming Asset Collection request
func (p processAsset) HandleAsset(*protocol.Request) (*protocol.Response, error) {
	data, err := p.dep.GetAssetCollectionService(p.dep.GetAssetCollectionServiceDependencies()).Process()
	if err != nil {
		p.logger.Logf(logging.ERROR, "Error in ProcessProcessor %v", err)
		return nil, err
	}
	outBytes, err := p.dep.GetSerializerJSON().WriteByteStream(data)

	if err != nil {
		p.logger.Logf(logging.ERROR, "Error in Process Asset Collection %v", err)
		return nil, err
	}

	resp := createResponseBody(outBytes, p.cfg.URLSuffix[model.ConstURLSuffixAssetCollection], protocol.HdrConstPluginDataPersist)

	return resp, nil
}

func (p processAsset) HandleConfig(request *protocol.Request) (*protocol.Response, error) {
	p.logger.Logf(logging.INFO, "Received config to update")
	configData, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	newConfig := make(map[string]interface{})
	err = json.Unmarshal(configData, &newConfig)
	if err != nil {
		return nil, err
	}
	service := p.dep.GetConfigService(p.dep)
	currentConfig, err := service.GetAssetPluginConfMap()
	if err != nil {
		return nil, err
	}

	newConfig = copyMapValue(currentConfig, newConfig)
	err = service.SetAssetPluginMap(newConfig)
	if err != nil {
		return nil, err
	}

	resp := createResponseBody([]byte(`{"result":"success"}`), "", "")
	p.logger.Logf(logging.INFO, "Received config - operation completed")
	return resp, nil
}

func copyMapValue(currentConfigFile map[string]interface{}, newConfigValues map[string]interface{}) map[string]interface{} {
	for key, val := range newConfigValues {
		if mapval, ok := currentConfigFile[key].(map[string]interface{}); ok {
			val = copyMapValue(mapval, val.(map[string]interface{}))
		}
		currentConfigFile[key] = val
	}
	return currentConfigFile
}

func createResponseBody(outBytes []byte, brokerPath string, hdrPersistData string) *protocol.Response {
	resp := protocol.NewResponse()
	resp.Body = bytes.NewReader(outBytes)
	resp.Headers.SetKeyValue(protocol.HdrContentType, "text/json")
	if brokerPath != "" {
		resp.Headers.SetKeyValue(protocol.HdrBrokerPath, brokerPath)
	}
	if hdrPersistData != "" {
		resp.Headers.SetKeyValue(protocol.HdrPluginDataPersist, hdrPersistData)
	}
	resp.Status = protocol.Ok
	return resp
}
