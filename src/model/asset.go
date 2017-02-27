package model

import (
	"github.com/ContinuumLLC/platform-common-lib/src/clar"
	"github.com/ContinuumLLC/platform-common-lib/src/env"
	cjson "github.com/ContinuumLLC/platform-common-lib/src/json"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
	"github.com/ContinuumLLC/platform-common-lib/src/pluginUtils"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
)

//ProcessAsset provides methods to Process incoming request
type ProcessAsset interface {
	ProcessAssetCollection(*protocol.Request) (*protocol.Response, error)
	ProcessConfiguration(*protocol.Request) (*protocol.Response, error)
}

//ProcessAssetFactory returns processAsset
type ProcessAssetFactory interface {
	GetProcessAsset(deps AssetServiceDependencies, cfg *AssetPluginConfig) ProcessAsset
}

//AssetListener interface provides methods to start processing incoming data
type AssetListener interface {
	Process() error
}

// AssetService captures and returns memory data
type AssetService interface {
	Process() error
}

// AssetnServiceFactory returns AssetService
type AssetServiceFactory interface {
	GetAssetService(deps AssetServiceDependencies) AssetService
}

// AssetDalFactory returns instance of AssetDal
type AssetDalFactory interface {
	GetAssetDal(deps AssetDalDependencies) AssetDal
}

// AssetDal captures memory metrics from underlying system
type AssetDal interface {
	SerializeObject(v interface{}) ([]byte, error)
}

// AssetDalDependencies gathers dependencies of Dal
type AssetDalDependencies interface {
	pluginUtils.PluginIOReader
	pluginUtils.PluginIOWriter
	cjson.FactoryJSON
}

// AssetServiceDependencies is the dependency interface for AssetService
type AssetServiceDependencies interface {
	clar.ServiceInitFactory
	AssetDalFactory
	AssetDalDependencies
	ProcessAssetFactory
	AssetCollectionServiceFactory
	GetAssetCollectionServiceDependencies() AssetCollectionServiceDependencies
	AssetCollectionDalFactory
	procParser.ParserFactory
	ConfigDalFactory
	ConfigServiceFactory
	env.FactoryEnv
	protocol.ServerFactory
}
