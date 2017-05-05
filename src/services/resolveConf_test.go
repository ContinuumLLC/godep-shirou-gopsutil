package services

import (
	"testing"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
)

func TestResolvePluginPathValues(t *testing.T) {
	path := model.AssetPluginPath{}
	resolveConf{logger: logging.GetLoggerFactory().Get()}.resolvePluginPathValues(&path)
	if path.AssetCollection != cAssetCollectionPluginPath {
		t.Error("Mismatch default plugin path values")
	}
}

func TestResolveURLSuffix(t *testing.T) {
	cfg := model.AssetPluginConfig{URLSuffix: make(map[string]string)}
	resolveConf{logger: logging.GetLoggerFactory().Get()}.resolveURLSuffix(&cfg)
	if cfg.URLSuffix[model.ConstURLSuffixAssetCollection] != cAgentServiceURLAssetCollection {
		t.Error("Mismatch default URL suffix value (assetCollection,processor)")
		return
	}
}
