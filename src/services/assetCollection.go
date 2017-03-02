package services

import (
	apiModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
)

// AssetCollectionServiceFactoryImpl factory implementation
type AssetCollectionServiceFactoryImpl struct{}

// GetAssetService returns Asset Service
func (AssetCollectionServiceFactoryImpl) GetAssetService(deps model.AssetServiceDependencies) model.AssetService {
	return assetCollectionServiceImpl{
		factory: deps,
	}
}

type assetCollectionServiceImpl struct {
	factory model.AssetServiceDependencies
}

// Process function processes.
func (srv assetCollectionServiceImpl) Process() (*apiModel.AssetCollection, error) {
	return srv.factory.GetAssetDal(srv.factory).GetAssetData()
}
