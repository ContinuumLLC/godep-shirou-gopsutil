// +build darwin

package bios

import (
	"errors"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
)

// Info returns baseboard information for MAC OS
// TODO: Impementation
func Info() (*asset.AssetBios, error) {
	return nil, errors.New(model.ErrNotImplemented)
}
