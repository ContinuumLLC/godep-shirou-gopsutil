package linux

import (
	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/exception"
)

type sysInfo struct {
	dep model.AssetDalDependencies
}

func (s sysInfo) getSystemInfo() (*asset.AssetSystem, error) {
	//This command require sudo access to execute
	product, err := s.dep.GetParser().ExecuteBash(`lshw -c system | grep product | cut -d ":" -f2`)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	//time zone
	tz, err := s.dep.GetParser().ExecuteBash("date +%z")
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	//time zone description
	tzd, err := s.dep.GetParser().ExecuteBash("date +%Z")
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	return &asset.AssetSystem{
		Product:             product,
		TimeZone:            tz,
		TimeZoneDescription: tzd,
		
	}, nil
}
