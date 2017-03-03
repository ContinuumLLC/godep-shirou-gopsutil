package linux

import (
	"strings"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
)

type osInfo struct {
	dep model.AssetDalDependencies
}

func (o osInfo) getOSInfo() (*asset.AssetOs, error) {
	os := new(asset.AssetOs)
	util := dalUtil{
		envDep: o.dep,
	}

	parser := o.dep.GetParser()
	cfg := procParser.Config{
		ParserMode: procParser.ModeSeparator,
		Separator:  ":",
	}

	dataCmd, err := util.getCommandData(parser, cfg, "lsb_release", "-a")
	if err != nil {
		return os, exception.New(model.ErrOSExecuteCommandFailed, err)
	}
	dataFile, err := util.getFileData(parser, procParser.Config{
		ParserMode: procParser.ModeSeparator,
		Separator:  "=",
	}, "/etc/default/locale")
	if err != nil {
		return os, exception.New(model.ErrOSExecuteCommandFailed, err)
	}
	os.Product = dataCmd.Map["Distributor ID"].Values[1]
	os.Manufacturer = dataCmd.Map["Description"].Values[1]
	os.Version = dataCmd.Map["Release"].Values[1]
	os.OsLanguage = strings.Trim(dataFile.Map["LANG"].Values[1], "\"")

	return os, nil
}
