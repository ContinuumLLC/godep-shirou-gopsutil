package services

import (
	apiModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	//"github.com/ContinuumLLC/platform-common-lib/src/exception"
)

// AssetCollectionServiceFactoryImpl factory implementation
type AssetCollectionServiceFactoryImpl struct{}

// GetAssetCollectionService returns Asset Service
func (AssetCollectionServiceFactoryImpl) GetAssetCollectionService(deps model.AssetCollectionServiceDependencies) model.AssetCollectionService {
	return assetCollectionServiceImpl{
		factory: deps,
	}
}

type assetCollectionServiceImpl struct {
	factory model.AssetCollectionServiceDependencies
}

// Process function processes.
func (srv assetCollectionServiceImpl) Process() (*apiModel.AssetCollection, error) {
	//return srv.factory.GetAssetCollectionDal(srv.factory).GetAssetCollection()
	assetCol, err := srv.factory.GetAssetCollectionDal(srv.factory).GetAssetCollection()
        if err != nil {
                //err = exception.New(model.ErrPluginService, err)
                return nil, err
        }
        return assetCol, nil
}
