// +build windows

package dal

import "github.com/ContinuumLLC/platform-asset-plugin/src/model"

// AssetCollectionDalFactoryImpl return AssetDal
type AssetCollectionDalFactoryImpl struct {
}

// GetAssetCollectionDal returns Dal
func (AssetCollectionDalFactoryImpl) GetAssetCollectionDal(deps model.AssetCollectionDalDependencies) model.AssetDal {
	//TODO - to be implemented for Windows
	return nil
}
