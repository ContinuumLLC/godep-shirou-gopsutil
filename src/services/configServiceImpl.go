package services

import (
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"

	exc "github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
)

//Singleton object
var (
	sConfig *model.AssetPluginConfig
)

// ConfigServiceFactoryImpl is the implementation of ConfigServiceFactory interface
type ConfigServiceFactoryImpl struct{}

// GetConfigService (method of ConfigService interface) returns ConfigService interface
func (ConfigServiceFactoryImpl) GetConfigService(csd model.ConfigServiceDependencies) model.ConfigService {
	return &configServiceImpl{
		factory: csd,
		logger:  logging.GetLoggerFactory().New("ConfigService "),
	}
}

type configServiceImpl struct {
	factory model.ConfigServiceDependencies
	logger  logging.Logger
}

//GetAssetPluginConfig reads the config and returns the AssetPluginConfig object
func (c *configServiceImpl) GetAssetPluginConfig() (*model.AssetPluginConfig, error) {
	var err error
	if sConfig == nil {
		sConfig, err = c.factory.GetConfigDal(c.factory).GetAssetPluginConf()
		if err != nil {
			err = exc.New(model.ErrAssetPluginConfig, err)
			//In case there is no config file, it will work with the default constants defined
			sConfig = new(model.AssetPluginConfig)
			sConfig.PluginPath = model.AssetPluginPath{}
			sConfig.URLSuffix = make(map[string]string)
		}
		rc := resolveConf{logger: logging.GetLoggerFactory().New("")}
		rc.resolveValues(sConfig)
		if sConfig.LogLevel != "" {
			c.logger.SetLogLevelConfig(sConfig.LogLevel)
		}
	}
	return sConfig, err
}

//GetAssetPluginMap reads the config and returns the AssetPluginConfig object
func (c *configServiceImpl) GetAssetPluginConfMap() (map[string]interface{}, error) {
	return c.factory.GetConfigDal(c.factory).GetAssetPluginConfMap()
}

func (c *configServiceImpl) SetAssetPluginMap(conf map[string]interface{}) error {
	return c.factory.GetConfigDal(c.factory).SetAssetPluginMap(conf)

}
