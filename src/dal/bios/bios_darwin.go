// +build darwin

package bios

import (
	"errors"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
)

// TODO: Info returns baseboard information for MAC OS
func Info() (*asset.AssetBios, error) {
	return nil, errors.New("not implemented yet")
}
