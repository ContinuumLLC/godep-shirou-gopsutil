package model

import (
	apiModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-common-lib/src/env"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
)

// AssetService captures and returns asset collection data
type AssetService interface {
	Process() (*apiModel.AssetCollection, error)
}

// AssetDal captures  asset collection metrics from underlying system
type AssetDal interface {
	GetAssetData() (*apiModel.AssetCollection, error)
}

// AssetServiceFactory returns AssetService
type AssetServiceFactory interface {
	GetAssetService(deps AssetServiceDependencies) AssetService
}

// AssetDalFactory returns instance of AssetDal
type AssetDalFactory interface {
	GetAssetDal(deps AssetDalDependencies) AssetDal
}

// AssetServiceDependencies are service dependencies
type AssetServiceDependencies interface {
	AssetDalDependencies
	AssetDalFactory
}

// AssetDalDependencies gathers dependencies of Dal
type AssetDalDependencies interface {
	procParser.ParserFactory
	env.FactoryEnv
}
