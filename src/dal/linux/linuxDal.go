package linux

import (
	"strings"
	"time"

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
	//getProcData will return the map of map[key]valuesArray
	//dataCmd is proc data, second param(*-network) is to create another set from this key
	//third param(logical name) is to name the map key as the value
	mapArr := util.getProcData(dataCmd, "*-network", "logical name")
	networks := make(map[string]asset.AssetNetwork)
	for k := range mapArr {
		n := asset.AssetNetwork{}
		n.Product = mapArr[k]["product"][1]
		n.Vendor = mapArr[k]["vendor"][1]
		n.LogicalName = mapArr[k]["logical name"][1]
		networks[k] = n
	}
	dataCmd, err = util.getCommandData(parser, cfg, "bash", "-c", "nmcli dev list")
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	mapArr = util.getProcData(dataCmd, "GENERAL.DEVICE", "GENERAL.DEVICE")
	for k := range mapArr {
		var (
			dsi  string
			ds   []string
			ipv4 string
		)
		dsiVal := mapArr[k]["DHCP4.OPTION[11]"][1]
		dsiArr := strings.Split(dsiVal, "=")
		if len(dsiArr) > 0 {
			dsi = strings.TrimSpace(dsiArr[1])
		}

		dsVal := mapArr[k]["DHCP4.OPTION[9]"][1]
		dsArr := strings.Split(dsVal, "=")
		if len(dsArr) > 0 {
			ds = strings.Split(strings.TrimSpace(dsArr[1]), " ")
		}

		ipv4Val := mapArr[k]["DHCP4.OPTION[6]"][1]
		ipv4Arr := strings.Split(ipv4Val, "=")
		if len(ipv4Arr) > 0 {
			ipv4 = strings.TrimSpace(ipv4Arr[1])
		}
		n := networks[k]
		n.DhcpServer = dsi
		n.DnsServers = ds
		n.PrivateIPv4 = ipv4
		networks[k] = n
	}
	return mapToArr(networks), nil
}

func mapToArr(m map[string]asset.AssetNetwork) []asset.AssetNetwork {
	networks := make([]asset.AssetNetwork, len(m))
	var i int
	for d := range m {
		networks[i] = m[d]
		i++
	}
	return networks
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
