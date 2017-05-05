package main

import (
	"fmt"
	"os"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
)

const (
	configFile  string = "ctm_asset_agent_plugin_cfg.json"
	logFile     string = "ctm_asset_agent_plugin.log"
	configIndex int    = 1
	logIndex    int    = 2
)

func main() {
	factories := factory{}
	service := factories.GetServiceInit()
	service.SetupOsArgs(configFile, logFile, os.Args, configIndex, logIndex)
	logger, err := logging.GetLoggerFactory().Init(logging.Config{
		AllowedLogLevel: logging.INFO,
		LogFileName:     service.GetLogFilePath(),
		MaxFileSizeInMB: 10,
		OldFileToKeep:   5,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	err = factories.GetAssetListener(factories).Process()
	if err != nil {
		logger.Logf(logging.ERROR, "Error retrieving Asset data %+v", err)
	}
}
