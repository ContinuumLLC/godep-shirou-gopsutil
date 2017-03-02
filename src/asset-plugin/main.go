package main

import "github.com/ContinuumLLC/platform-common-lib/src/logging"
import "os"

const (
	configFile  string = "ctm_asset_agent_plugin_cfg.json"
	logFile     string = "ctm_asset_agent_plugin.log"
	configIndex int    = 1
	logIndex    int    = 2
)

var logger logging.Logger

func main() {
	var logger logging.Logger
	factories := factory{}
	service := factories.GetServiceInit()
	service.SetupOsArgs(configFile, logFile, os.Args, configIndex, logIndex)
	logger = logging.GetLoggerFactory().New("Main ")
	err := factories.GetAssetListener(factories).Process()
	if err != nil {
		logger.Logf(logging.ERROR, "Error retrieving Asset data %+v", err)
	}
}
