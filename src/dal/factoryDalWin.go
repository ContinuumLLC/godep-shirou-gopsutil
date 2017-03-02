// +build windows

package dal

import "github.com/ContinuumLLC/platform-asset-plugin/src/model"

// AssetDalFactoryImpl return AssetDal
type AssetDalFactoryImpl struct {
}

// GetAssetDal returns Dal
func (AssetDalFactoryImpl) GetAssetDal(deps model.AssetDalDependencies) model.AssetDal {
	//TODO - to be implemented for Windows
	return nil
}
