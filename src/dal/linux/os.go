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
	parser := o.dep.GetParser()
	cfg := procParser.Config{
		ParserMode: procParser.ModeSeparator,
		Separator:  ":",
	}

	util := dalUtil{
		envDep: o.dep,
	}
	dataCmd, err := util.getCommandData(parser, cfg, "lsb_release", "-a")
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}

	//changing the separator for next file/command to parse and get data
	cfg.Separator = "="
	dataFile, err := util.getFileData(parser, cfg, "/etc/default/locale")
	if err != nil {
		return nil, exception.New(model.ErrFileReadFailed, err)
	}
	return &asset.AssetOs{
		Product:      dataCmd.Map["Distributor ID"].Values[1],
		Manufacturer: dataCmd.Map["Description"].Values[1],
		Version:      dataCmd.Map["Release"].Values[1],
		OsLanguage:   strings.Trim(dataFile.Map["LANG"].Values[1], "\""),
		//os.InstallDate - To be added
		//os.SerialNumber - Presently not able to find it for ubuntu
	}, nil
}
