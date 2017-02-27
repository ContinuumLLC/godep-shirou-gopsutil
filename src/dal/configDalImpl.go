package dal

import (
	"fmt"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
)

type configDalImpl struct {
	factory model.ConfigDalDependencies
}

const (
	configFilename = "ctm_asset_agent_plugin_cfg.json"
)

// GetAssetConf is the ConfigDal interface method which returns the Config
func (c configDalImpl) GetAssetPluginConf() (*model.AssetPluginConfig, error) {
	var obj model.AssetPluginConfig
	err := c.factory.GetDeserializerJSON().ReadFile(&obj, c.factory.GetServiceInit().GetConfigPath())
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func (c configDalImpl) GetAssetPluginConfMap() (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	dir, err := c.factory.GetEnv().GetExeDir()
	if err != nil {
		return nil, err
	}

	filename := c.factory.GetServiceInit().GetConfigPath()

	if filename == "" {
		filename = configFilename
	}
	filename = dir + "/" + filename
	fmt.Println(filename)
	err = c.factory.GetDeserializerJSON().ReadFile(&obj, filename)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (c configDalImpl) SetAssetPluginMap(conf map[string]interface{}) error {
	filename := c.factory.GetServiceInit().GetConfigPath()

	if filename == "" {
		filename = configFilename
	}
	return c.factory.GetSerializerJSON().WriteFile(filename, conf)
}

// ConfigDalFactoryImpl is the implementation of ConfigDalFactory interface
type ConfigDalFactoryImpl struct{}

// GetConfigDal (a method of ConfigDalFactory interface) returns ConfigDal
func (ConfigDalFactoryImpl) GetConfigDal(f model.ConfigDalDependencies) model.ConfigDal {
	return configDalImpl{
		factory: f,
	}
}
