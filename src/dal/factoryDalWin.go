// +build windows

package dal

import "github.com/ContinuumLLC/platform-asset-plugin/src/model"

// AssetCollectionDalFactoryImpl return AssetDal
type AssetCollectionDalFactoryImpl struct {
}

// GetAssetDal returns Dal
func (AssetCollectionDalFactoryImpl) GetAssetDal(deps model.AssetDalDependencies) model.AssetDal {
	//TODO - to be implemented for Windows
	return nil
}
