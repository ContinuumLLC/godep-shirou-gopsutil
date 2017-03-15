// +build linux

package dal

import (
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
)

// AssetDalFactoryImpl return AssetDal
type AssetDalFactoryImpl struct {
}

// GetAssetDal returns Dal
func (AssetDalFactoryImpl) GetAssetDal(deps model.AssetDalDependencies) model.AssetDal {
	return &assetDalImpl{
		Factory: deps,
		Logger:  logging.GetLoggerFactory().New("AssetDal"),
	}
}
