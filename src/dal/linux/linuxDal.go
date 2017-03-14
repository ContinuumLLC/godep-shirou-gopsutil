package linux

import (
	"encoding/xml"
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
	//cListHwAsJSON  string = "lshw -json"
	cListHwAsXML string = "lshw -c system,memory,bus,disk,volume -xml"
)

// Memory Proc related constants
const (
	cMemProcPath                   string = "/proc/meminfo"
	cMemProcPhysicalTotalBytes     string = "MemTotal"
	cMemProcPhysicalAvailableBytes string = "MemAvailable"
	cMemProcPageAvailableBytes     string = "SwapFree"
	cMemProcPageTotalBytes         string = "SwapTotal"
)

//List denotes list of hardware assets returned by lshw command
type List struct {
	XMLName  xml.Name `xml:"list"`
	Nodelist Node     `xml:"node"`
}

//Node denotes a particular hardware asset node returned by lshw command
type Node struct {
	Class        string        `xml:"class,attr"`
	ID           string        `xml:"id,attr"`
	Desc         string        `xml:"description"`
	Vendor       string        `xml:"vendor"`
	Product      string        `xml:"product"`
	Version      string        `xml:"version"`
	Serial       string        `xml:"serial"`
	SizeInBytes  int64         `xml:"size"`
	LogName      []Logicalname `xml:"logicalname"`
	Capabilities []Capability  `xml:"capabilities>capability"`
	Nodelist     []Node        `xml:"node"`
}

//Logicalname denotes logical name of a node
type Logicalname struct {
	Text string `xml:",chardata"`
}

//Capability denotes capability of a node
type Capability struct {
	ID   string `xml:"id,attr"`
	Text string `xml:",chardata"`
}

// AssetDalImpl ...
type AssetDalImpl struct {
	Factory model.AssetDalDependencies
	Logger  logging.Logger
}

func (a AssetDalImpl) readHwList() (*List, error) {
	v := List{}
	hw, err := a.Factory.GetEnv().ExecuteBash(cListHwAsXML)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
	}
	err = xml.Unmarshal([]byte(hw), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

//GetAssetData ...
func (a AssetDalImpl) GetAssetData() (*asset.AssetCollection, error) {
	v, err := a.readHwList()
	if err != nil {
		return nil, err
	}

	pp, err := a.GetBiosInformation(&v.Nodelist)
	if err != nil {
		return nil, err
	}

	pp1, err := a.GetBaseBoardInformation(&v.Nodelist)
	if err != nil {
		return nil, err
	}

	dd, err := a.GetDrivesInformation(&v.Nodelist)
	if err != nil {
		return nil, err
	}

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
	p, err := a.GetProcessorInfo()
	if err != nil {
		return nil, err
	}
	return &asset.AssetCollection{
		CreatedBy:     cAssetCreatedBy,
		CreateTimeUTC: time.Now().UTC(),
		Type:          cAssetDataType,
		Bios:          *pp,
		BaseBoard:     *pp1,
		Os:            *o,
		Memory:        *m,
		System:        *s,
		Networks:      n,
		Drives:        dd,
		Processors:    p,
	}, nil
}

func (a AssetDalImpl) getRequiredNode(l *Node, id string, class string) *Node {
	if l.ID == id && l.Class == class {
		return l
	}
	if len(l.Nodelist) > 0 {
		for i := range l.Nodelist {
			return a.getRequiredNode(&l.Nodelist[i], id, class)
		}
	}
	return nil
}

func (a AssetDalImpl) getAllNodes(root *Node, id string, class string, listOfNodes []Node) []Node {
	if root.ID == id && root.Class == class {
		return append(listOfNodes, *root)
	}
	if len(root.Nodelist) > 0 {
		for i := range root.Nodelist {
			listOfNodes = a.getAllNodes(&root.Nodelist[i], id, class, listOfNodes)
		}
	}
	return listOfNodes
}

func (a AssetDalImpl) getAllPartitions(root *Node, listOfPart []string) []string {
	if len(root.Nodelist) > 0 {
		var volume string
		for i, v := range root.Nodelist {
			if len(v.LogName) > 0 {
				volume = v.LogName[0].Text
			}
			listOfPart = append(listOfPart, volume)
			listOfPart = a.getAllPartitions(&root.Nodelist[i], listOfPart)
		}
	}
	return listOfPart
}

//GetBiosInformation ...
func (a AssetDalImpl) GetBiosInformation(l *Node) (*asset.AssetBios, error) {
	var smbiosVer string
	for _, v := range l.Capabilities {
		if strings.Contains(v.ID, "smbios") {
			smbiosVer = v.Text
			break
		}
	}
	n1 := a.getRequiredNode(l, "firmware", "memory")
	if n1 == nil {
		return &asset.AssetBios{}, nil
	}

	return &asset.AssetBios{
		Manufacturer:  n1.Vendor,
		Version:       n1.Version,
		SmbiosVersion: smbiosVer,
	}, nil

}

//GetBaseBoardInformation ...
func (a AssetDalImpl) GetBaseBoardInformation(l *Node) (*asset.AssetBaseBoard, error) {
	n1 := a.getRequiredNode(l, "core", "bus")
	if n1 == nil {
		return &asset.AssetBaseBoard{}, nil
	}

	return &asset.AssetBaseBoard{
		Product:      n1.Product,
		Manufacturer: n1.Vendor,
		Version:      n1.Version,
		SerialNumber: n1.Serial,
	}, nil

}

//GetDrivesInformation ...
func (a AssetDalImpl) GetDrivesInformation(l *Node) ([]asset.AssetDrive, error) {
	var listOfNodes []Node
	var listOfDrives []asset.AssetDrive
	var tmp asset.AssetDrive
	diskList := a.getAllNodes(l, "disk", "disk", listOfNodes)
	for _, value := range diskList {
		//Need description and logical name, snumber, version in place of model and media type
		tmp.Manufacturer = value.Vendor
		tmp.Product = value.Product
		tmp.SizeBytes = value.SizeInBytes
		var listOfPart []string
		tmp.Partitions = a.getAllPartitions(&value, listOfPart)

		listOfDrives = append(listOfDrives, tmp)
	}

	optDriveList := a.getAllNodes(l, "cdrom", "disk", listOfNodes)
	for _, value := range optDriveList {
		//Need description and logical name, snumber, version in place of model and media type
		tmp.Manufacturer = value.ID
		tmp.Product = value.Desc
		tmp.SizeBytes = value.SizeInBytes
		var listOfPart []string
		tmp.Partitions = a.getAllPartitions(&value, listOfPart)

		listOfDrives = append(listOfDrives, tmp)
	}

	return listOfDrives, nil

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
	mapArr := util.getProcData(dataFile, "processor", "processor")
	processors := make([]asset.AssetProcessor, len(mapArr))
	var i int
	for k := range mapArr {
		processors[i].ClockSpeedMhz, _ = strconv.ParseFloat(mapArr[k]["cpu MHz"][1], 64)
		processors[i].Family, _ = strconv.Atoi(mapArr[k]["cpu family"][1])
		processors[i].Manufacturer = mapArr[k]["vendor_id"][1]
		processors[i].NumberOfCores, _ = strconv.Atoi(mapArr[k]["cpu cores"][1])
		processors[i].Product = mapArr[k]["model name"][1]
		processors[i].ProcessorType = cpuType
		//processors[i].SerialNumber  ... to be added
		i++
	}
	return processors, nil
}
