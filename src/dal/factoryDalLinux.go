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

// GetAssetDal returns Dal
func (AssetCollectionDalFactoryImpl) GetAssetDal(deps model.AssetDalDependencies) model.AssetDal {
	return &linux.AssetCollectionDalLinux{
		Factory: deps,
		Logger:  logging.GetLoggerFactory().New("AssetDal"),
	}
}
