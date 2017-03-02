// +build linux

package dal

import "github.com/ContinuumLLC/platform-asset-plugin/src/model"
import "github.com/ContinuumLLC/platform-common-lib/src/logging"

// AssetCollectionDalFactoryImpl return AssetDal
type AssetCollectionDalFactoryImpl struct {
}

// GetAssetCollectionDal returns Dal
func (AssetCollectionDalFactoryImpl) GetAssetCollectionDal(deps model.AssetCollectionDalDependencies) model.AssetDal {
	return &assetCollectionDalLinux{
		factory: deps,
		logger:  logging.GetLoggerFactory().New("AssetDal"),
	}
}
