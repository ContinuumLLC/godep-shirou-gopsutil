package model

import (
	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-common-lib/src/env"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
)

const (
	//ErrExecuteCommandFailed error code for execute command failed
	ErrExecuteCommandFailed = "ErrExecuteCommandFailed"
	//ErrFileReadFailed error code for file read failed
	ErrFileReadFailed = "ErrFileReadFailed"
)

// AssetService captures and returns asset collection data
type AssetService interface {
	Process() (*asset.AssetCollection, error)
}

// AssetDal captures  asset collection metrics from underlying system
type AssetDal interface {
	GetAssetData() (*asset.AssetCollection, error)
	GetOSInfo() (*asset.AssetOs, error)
	GetSystemInfo() (*asset.AssetSystem, error)
	GetNetworkInfo() ([]asset.AssetNetwork, error)
}

// AssetServiceFactory returns AssetService
type AssetServiceFactory interface {
	GetAssetService(deps AssetServiceDependencies) AssetService
}

// AssetDalFactory returns instance of AssetDal
type AssetDalFactory interface {
	GetAssetDal(deps AssetDalDependencies) AssetDal
}

// AssetServiceDependencies are service dependencies
type AssetServiceDependencies interface {
	AssetDalDependencies
	AssetDalFactory
}

// AssetDalDependencies gathers dependencies of Dal
type AssetDalDependencies interface {
	procParser.ParserFactory
	env.FactoryEnv
}
