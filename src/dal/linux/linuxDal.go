package linux

import (
	"strings"
	"time"

	"strconv"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
)

// AssetCollection related constants
const (
	cAssetCreatedBy string = "/continuum/agent/plugin/asset"
	cAssetDataType  string = "assetCollection"
)

const (
	cSysProductCmd string = `lshw -c system | grep product | cut -d ":" -f2`
	cCPUArcCmd     string = `lscpu | grep Architecture | cut -d ":" -f2`
	cSysTz         string = "date +%z"
	cSysTzd        string = "date +%Z"
	cSysSerialNo   string = "dmidecode -s system-serial-number"
	cSysHostname   string = "hostname"
)

// Memory Proc related constants
const (
	cMemProcPath                   string = "/proc/meminfo"
	cMemProcPhysicalTotalBytes     string = "MemTotal"
	cMemProcPhysicalAvailableBytes string = "MemAvailable"
	cMemProcPageAvailableBytes     string = "SwapFree"
	cMemProcPageTotalBytes         string = "SwapTotal"
)

// AssetDalImpl ...
type AssetDalImpl struct {
	Factory model.AssetDalDependencies
	Logger  logging.Logger
}

//GetAssetData ...
func (a AssetDalImpl) GetAssetData() (*asset.AssetCollection, error) {
	o, err := a.GetOSInfo()
	if err != nil {
		return nil, err
	}
	a.Logger.Log(logging.INFO, "")
	s, err := a.GetSystemInfo()
	if err != nil {
		return nil, err
	}
	n, err := a.GetNetworkInfo()
	if err != nil {
		return nil, err
	}
	m, err := a.GetMemoryInfo()
	if err != nil {
		return nil, err
	}
	p, err := a.GetProcessorInfo()
	if err != nil {
		return nil, err
	}
	return &asset.AssetCollection{
		CreatedBy:     cAssetCreatedBy,
		CreateTimeUTC: time.Now().UTC(),
		Type:          cAssetDataType,
		Os:            *o,
		BaseBoard:     *(getBaseBoardInfo()),
		Bios:          *(getBiosInfo()),
		Memory:        *m,
		System:        *s,
		Networks:      n,
		Processors:    p,
	}, nil
}

// GetOSInfo returns the OS info
func (a AssetDalImpl) GetOSInfo() (*asset.AssetOs, error) {
	parser := a.Factory.GetParser()
	cfg := procParser.Config{
		ParserMode: procParser.ModeSeparator,
		Separator:  ":",
	}

	util := dalUtil{
		envDep: a.Factory,
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

// GetSystemInfo returns system info
func (a AssetDalImpl) GetSystemInfo() (*asset.AssetSystem, error) {
	//TODO - Below repetitive code needs to be refactored
	//This command require sudo access to execute
	product, err := a.Factory.GetEnv().ExecuteBash(cSysProductCmd)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	//time zone
	tz, err := a.Factory.GetEnv().ExecuteBash(cSysTz)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	//time zone description
	tzd, err := a.Factory.GetEnv().ExecuteBash(cSysTzd)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	//serial number
	srno, err := a.Factory.GetEnv().ExecuteBash(cSysSerialNo)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	//hostname
	hostname, err := a.Factory.GetEnv().ExecuteBash(cSysHostname)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	return &asset.AssetSystem{
		Product:             product,
		TimeZone:            tz,
		TimeZoneDescription: tzd,
		//Category - to be added
		//Model - to be added
		SerialNumber: srno,
		SystemName:   hostname,
	}, nil
}

// GetNetworkInfo returns network info
func (a AssetDalImpl) GetNetworkInfo() ([]asset.AssetNetwork, error) {
	parser := a.Factory.GetParser()
	cfg := procParser.Config{
		ParserMode: procParser.ModeSeparator,
		Separator:  ":",
	}
	util := dalUtil{
		envDep: a.Factory,
	}
	dataCmd, err := util.getCommandData(parser, cfg, "lshw", "-c", "network")
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	mapArr := util.getProcData(dataCmd, "*-network")
	networks := make([]asset.AssetNetwork, len(mapArr))
	for i := 0; i < len(mapArr); i++ {
		networks[i].Product = mapArr[i]["product"][1]
		networks[i].Vendor = mapArr[i]["vendor"][1]
	}
	return networks, nil
}

// GetMemoryInfo returns memory info
func (a AssetDalImpl) GetMemoryInfo() (*asset.AssetMemory, error) {
	parser := a.Factory.GetParser()
	cfg := procParser.Config{
		ParserMode:    procParser.ModeKeyValue,
		IgnoreNewLine: true,
	}
	util := dalUtil{
		envDep: a.Factory,
	}
	data, err := util.getFileData(parser, cfg, cMemProcPath)
	if err != nil {
		return nil, exception.New(model.ErrFileReadFailed, err)
	}

	memTotal := util.getDataFromMap(cMemProcPhysicalTotalBytes, data)
	memAvail := util.getDataFromMap(cMemProcPhysicalAvailableBytes, data)
	swapTotal := util.getDataFromMap(cMemProcPageTotalBytes, data)
	swapAvail := util.getDataFromMap(cMemProcPageAvailableBytes, data)

	return &asset.AssetMemory{
		TotalPhysicalMemoryBytes:     memTotal,
		AvailablePhysicalMemoryBytes: memAvail,
		TotalPageFileSpaceBytes:      swapTotal,
		AvailablePageFileSpaceBytes:  swapAvail,
		TotalVirtualMemoryBytes:      (memTotal + swapTotal),
		AvailableVirtualMemoryBytes:  (memAvail + swapAvail),
	}, nil
}

// GetProcessorInfo returns processor info
func (a AssetDalImpl) GetProcessorInfo() ([]asset.AssetProcessor, error) {
	parser := a.Factory.GetParser()
	cfg := procParser.Config{
		ParserMode: procParser.ModeSeparator,
		Separator:  ":",
	}
	util := dalUtil{
		envDep: a.Factory,
	}
	cpuType, err := a.Factory.GetEnv().ExecuteBash(cCPUArcCmd)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	dataFile, err := util.getFileData(parser, cfg, "/proc/cpuinfo")
	if err != nil {
		return nil, exception.New(model.ErrFileReadFailed, err)
	}
	mapArr := util.getProcData(dataFile, "processor")
	processors := make([]asset.AssetProcessor, len(mapArr))
	for i := 0; i < len(mapArr); i++ {
		processors[i].ClockSpeedMhz, _ = strconv.ParseFloat(mapArr[i]["cpu MHz"][1], 64)
		processors[i].Family, _ = strconv.Atoi(mapArr[i]["cpu family"][1])
		processors[i].Manufacturer = mapArr[i]["vendor_id"][1]
		processors[i].NumberOfCores, _ = strconv.Atoi(mapArr[i]["cpu cores"][1])
		processors[i].Product = mapArr[i]["model name"][1]
		processors[i].ProcessorType = cpuType
		//processors[i].SerialNumber  ... to be added
	}
	return processors, nil
}
