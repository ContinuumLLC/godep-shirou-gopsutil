// +build linux

package dal

import "github.com/ContinuumLLC/platform-asset-plugin/src/model"
import "github.com/ContinuumLLC/platform-common-lib/src/logging"

// AssetCollectionDalFactoryImpl return AssetCollectionDal
type AssetCollectionDalFactoryImpl struct {
}

// GetAssetCollectionDal returns Dal
func (AssetCollectionDalFactoryImpl) GetAssetCollectionDal(deps model.AssetCollectionDalDependencies) model.AssetCollectionDal {
	return &assetCollectionDalLinux{
		factory: deps,
		logger:  logging.GetLoggerFactory().New("AssetCollectionDal"),
	}
}
