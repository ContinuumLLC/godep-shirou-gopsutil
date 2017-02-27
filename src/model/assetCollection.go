package model

import (
	apiModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-common-lib/src/env"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
)

// AssetCollectionService captures and returns asset collection data
type AssetCollectionService interface {
	Process() (*apiModel.AssetCollection, error)
}

// AssetCollectionDal captures  asset collection metrics from underlying system
type AssetCollectionDal interface {
	GetAssetCollection() (*apiModel.AssetCollection, error)
}

// AssetCollectionServiceFactory returns AssetCollectionService
type AssetCollectionServiceFactory interface {
	GetAssetCollectionService(deps AssetCollectionServiceDependencies) AssetCollectionService
}

// AssetCollectionDalFactory returns instance of AssetCollectionDal
type AssetCollectionDalFactory interface {
	GetAssetCollectionDal(deps AssetCollectionDalDependencies) AssetCollectionDal
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
