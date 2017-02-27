package services

import (
	"testing"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
)

func TestResolvePluginPathValues(t *testing.T) {
	path := model.AssetPluginPath{}
	resolveConf{}.resolvePluginPathValues(&path)
	if path.AssetCollection != cAssetCollectionPluginPath {
		t.Error("Mismatch default plugin path values")
	}
}

func TestResolveURLSuffix(t *testing.T) {
	cfg := model.AssetPluginConfig{URLSuffix: make(map[string]string)}
	resolveConf{}.resolveURLSuffix(&cfg)
	if cfg.URLSuffix[model.ConstURLSuffixAssetCollection] != cAgentServiceURLAssetCollection || cfg.URLSuffix[model.ConstURLSuffixProcessor] != cAgentServiceURLProcessor {
		t.Error("Mismatch default URL suffix value (assetCollection,processor)")
		return
	}
}
