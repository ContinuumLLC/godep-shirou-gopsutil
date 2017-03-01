package msgl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

//ProcessAssetFactoryImpl returns a asset Processor
type ProcessAssetFactoryImpl struct{}

//GetProcessAsset returns a Processor processor
func (ProcessAssetFactoryImpl) GetProcessAsset(deps model.AssetServiceDependencies, config *model.AssetPluginConfig) model.ProcessAsset {
	return processAsset{
		dep:    deps,
		cfg:    config,
		logger: logging.GetLoggerFactory().New("ProcessAsset"),
	}
}

type processAsset struct {
	dep    model.AssetServiceDependencies
	cfg    *model.AssetPluginConfig
	logger logging.Logger
}

//ProcessAssetCollection processes incoming Asset Collection request
func (p processAsset) ProcessAssetCollection(*protocol.Request) (*protocol.Response, error) {
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

func (p processAsset) ProcessConfiguration(request *protocol.Request) (*protocol.Response, error) {
	f, _ := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	log.SetOutput(f)

	log.Println("called")
	configData, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("read")
	newConfig := make(map[string]interface{})
	log.Println(configData)
	err = json.Unmarshal(configData, &newConfig)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	service := p.dep.GetConfigService(p.dep)

	currentConfig, err := service.GetAssetPluginConfMap()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	newConfig = copyMapValue(currentConfig, newConfig)
	log.Println("trying to save")
	err = service.SetAssetPluginMap(newConfig)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp := createResponseBody([]byte(`{"result":"success"}`), "", "")
	fmt.Println(resp)
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
