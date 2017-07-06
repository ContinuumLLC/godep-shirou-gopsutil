// +build linux

package net

import (
	"errors"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
)

// TODO: Info returns network information for Linux
func Info() ([]asset.AssetNetwork, error) {
	return nil, errors.New("not implemented yet")
}
