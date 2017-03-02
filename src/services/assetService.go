package services

import (
	apiModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
)

// AssetServiceFactoryImpl factory implementation
type AssetServiceFactoryImpl struct{}

// GetAssetService returns Asset Service
func (AssetServiceFactoryImpl) GetAssetService(deps model.AssetServiceDependencies) model.AssetService {
	return assetServiceImpl{
		factory: deps,
	}
}

type assetServiceImpl struct {
	factory model.AssetServiceDependencies
}

// Process function processes.
func (srv assetServiceImpl) Process() (*apiModel.AssetCollection, error) {
	return srv.factory.GetAssetDal(srv.factory).GetAssetData()
}
