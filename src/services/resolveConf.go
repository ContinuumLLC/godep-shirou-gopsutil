package services

import (
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
)

// TODO - Add logging for setting up default config

//Default config values if not found from the config file
const (
	cAssetCollectionPluginPath    string = "/asset"
	cAssetConfigurationPluginPath string = "/asset/configuration"
)

const (
	cAgentServiceURLAssetCollection string = "/asset"
)

type resolveConf struct {
	logger logging.Logger
}

func (r resolveConf) resolveValues(pcfg *model.AssetPluginConfig) {
	r.resolvePluginPathValues(&pcfg.PluginPath)
	r.resolveURLSuffix(pcfg)
}

func (r resolveConf) resolvePluginPathValues(path *model.AssetPluginPath) {
	if path.AssetCollection == "" {
		r.logger.Logf(logging.ERROR, "Plugin path for assetCollection not found. Using default")
		path.AssetCollection = cAssetCollectionPluginPath
	}
	if path.Configuration == "" {
		r.logger.Logf(logging.ERROR, "Plugin path for Asset Configuration not found. Using default")
		path.Configuration = cAssetConfigurationPluginPath
	}
}

func (r resolveConf) resolveURLSuffix(cfg *model.AssetPluginConfig) {
	if cfg.URLSuffix[model.ConstURLSuffixAssetCollection] == "" {
		cfg.URLSuffix[model.ConstURLSuffixAssetCollection] = cAgentServiceURLAssetCollection
	}
}
