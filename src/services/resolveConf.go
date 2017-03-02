package services

import "github.com/ContinuumLLC/platform-asset-plugin/src/model"

// TODO - Add logging for setting up default config

//Default config values if not found from the config file
const (
	cAssetCollectionPluginPath string = "/asset"
)

const (
	cAgentServiceURLAssetCollection string = "/asset"
)

type resolveConf struct{}

func (r resolveConf) resolveValues(pcfg *model.AssetPluginConfig) {
	r.resolvePluginPathValues(&pcfg.PluginPath)
	r.resolveURLSuffix(pcfg)
}

func (r resolveConf) resolvePluginPathValues(path *model.AssetPluginPath) {
	if path.AssetCollection == "" {
		path.AssetCollection = cAssetCollectionPluginPath
	}
}

func (r resolveConf) resolveURLSuffix(cfg *model.AssetPluginConfig) {
	if cfg.URLSuffix[model.ConstURLSuffixAssetCollection] == "" {
		cfg.URLSuffix[model.ConstURLSuffixAssetCollection] = cAgentServiceURLAssetCollection
	}
}
