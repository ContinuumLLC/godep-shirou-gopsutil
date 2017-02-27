package main

import (
	"github.com/ContinuumLLC/platform-common-lib/src/clar"
	"github.com/ContinuumLLC/platform-common-lib/src/env"
	cjson "github.com/ContinuumLLC/platform-common-lib/src/json"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol/http"
	"github.com/ContinuumLLC/platform-common-lib/src/pluginUtils"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
	"github.com/ContinuumLLC/platform-asset-plugin/src/dal"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-asset-plugin/src/msgl"
	"github.com/ContinuumLLC/platform-asset-plugin/src/services"
)

type factory struct {
	clar.ServiceInitFactoryImpl
	dal.AssetCollectionDalFactoryImpl
	dal.ConfigDalFactoryImpl

	dal.AssetDalFactoryImpl
	services.AssetCollectionServiceFactoryImpl
	services.ConfigServiceFactoryImpl

	pluginUtils.StandardIOReaderImpl
	pluginUtils.StandardIOWriterImpl
	cjson.FactoryJSONImpl
	env.FactoryEnvImpl
	procParser.ParserFactoryImpl

	msgl.AssetListenerFactory
	msgl.ProcessAssetFactoryImpl

	http.ServerHTTPFactory
}

func (f factory) GetAssetCollectionServiceDependencies() model.AssetCollectionServiceDependencies {
	var ff model.AssetCollectionServiceDependencies = f
	ret, _ := ff.(model.AssetCollectionServiceDependencies)
	return ret
}
