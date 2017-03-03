package model

import (
	"github.com/ContinuumLLC/platform-common-lib/src/clar"
	"github.com/ContinuumLLC/platform-common-lib/src/env"
	cjson "github.com/ContinuumLLC/platform-common-lib/src/json"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
	"github.com/ContinuumLLC/platform-common-lib/src/pluginUtils"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
)

//Handler provides methods to Process incoming request
type Handler interface {
	HandleAsset(*protocol.Request) (*protocol.Response, error)
	HandleConfig(*protocol.Request) (*protocol.Response, error)
}

//HandlerFactory returns processAsset
type HandlerFactory interface {
	GetHandler(deps HandlerDependencies, cfg *AssetPluginConfig) Handler
}

//AssetListener interface provides methods to start processing incoming data
type AssetListener interface {
	Process() error
}

// HandlerDependencies is the dependency interface for AssetService
type HandlerDependencies interface {
	clar.ServiceInitFactory
	HandlerFactory
	AssetServiceFactory
	GetAssetCollectionServiceDependencies() AssetServiceDependencies
	AssetDalFactory
	procParser.ParserFactory
	ConfigDalFactory
	ConfigServiceFactory
	env.FactoryEnv
	protocol.ServerFactory
	cjson.FactoryJSON
	pluginUtils.PluginIOReader
	pluginUtils.PluginIOWriter
}
