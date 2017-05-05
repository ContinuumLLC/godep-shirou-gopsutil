package model

import (
	"github.com/ContinuumLLC/platform-common-lib/src/clar"
	"github.com/ContinuumLLC/platform-common-lib/src/env"
	"github.com/ContinuumLLC/platform-common-lib/src/json"
)

// AssetPluginConfig is the overall config that asset plugin will use.
type AssetPluginConfig struct {
	PluginPath         AssetPluginPath   `json:"pluginPath"`
	URLSuffix          map[string]string `json:"urlSuffix"`
	LogLevel           string            `json:"logLevel"`
	MaxLogFileSizeInMB int64
	OldLogFileToKeep   int
}

// AssetPluginPath is to have plugin paths.
type AssetPluginPath struct {
	AssetCollection string `json:"assetCollection"`
	Configuration   string `json:"configuration"`
}

// ConfigServiceFactory interface returns the ConfigService
type ConfigServiceFactory interface {
	GetConfigService(ConfigServiceDependencies) ConfigService
}

//ConfigService interface provides methods to configurations
type ConfigService interface {
	GetAssetPluginConfig() (*AssetPluginConfig, error)
	GetAssetPluginConfMap() (map[string]interface{}, error)
	SetAssetPluginMap(map[string]interface{}) error
}

// ConfigDal interface to handle all dal operations of Config
type ConfigDal interface {
	GetAssetPluginConf() (*AssetPluginConfig, error)
	GetAssetPluginConfMap() (map[string]interface{}, error)
	SetAssetPluginMap(map[string]interface{}) error
}

// ConfigDalDependencies interface contains all the Config Dal dependencies
type ConfigDalDependencies interface {
	clar.ServiceInitFactory
	env.FactoryEnv
	json.FactoryJSON
}

// ConfigDalFactory interface returns the ConfigDal
type ConfigDalFactory interface {
	GetConfigDal(ConfigDalDependencies) ConfigDal
}

// ConfigServiceDependencies interface contains all the Config Service dependencies
type ConfigServiceDependencies interface {
	ConfigDalFactory
	ConfigDalDependencies
}

const (
	//ErrAssetPluginConfig error code for error in reading plugin config
	ErrAssetPluginConfig = "ErrAssetPluginConfig"
)
