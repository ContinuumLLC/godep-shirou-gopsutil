// +build windows

package dal

import "github.com/ContinuumLLC/platform-asset-plugin/src/model"

// AssetCollectionDalFactoryImpl return AssetCollectionDal
type AssetCollectionDalFactoryImpl struct {
}

// GetAssetCollectionDal returns Dal
func (AssetCollectionDalFactoryImpl) GetAssetCollectionDal(deps model.AssetCollectionDalDependencies) model.AssetCollectionDal {
	//TODO - to be implemented for Windows
	return nil
}
