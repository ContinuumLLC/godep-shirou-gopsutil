// +build linux

package dal

import (
	"github.com/ContinuumLLC/platform-asset-plugin/src/dal/linux"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
)

// AssetCollectionDalFactoryImpl return AssetDal
type AssetCollectionDalFactoryImpl struct {
}

// GetAssetCollectionDal returns Dal
func (AssetCollectionDalFactoryImpl) GetAssetCollectionDal(deps model.AssetCollectionDalDependencies) model.AssetDal {
	return &linux.AssetCollectionDalLinux{
		Factory: deps,
		Logger:  logging.GetLoggerFactory().New("AssetDal"),
	}
}
