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
	GetAssetCollectionService(deps AssetCollectionServiceDependencies) AssetService
}

// AssetCollectionDalFactory returns instance of AssetDal
type AssetCollectionDalFactory interface {
	GetAssetCollectionDal(deps AssetCollectionDalDependencies) AssetDal
}

// AssetCollectionServiceDependencies are service dependencies
type AssetCollectionServiceDependencies interface {
	AssetCollectionDalDependencies
	AssetCollectionDalFactory
}

// AssetCollectionDalDependencies gathers dependencies of Dal
type AssetCollectionDalDependencies interface {
	procParser.ParserFactory
	env.FactoryEnv
}
