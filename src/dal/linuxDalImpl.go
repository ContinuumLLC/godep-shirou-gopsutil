// +build linux

package dal

import (
	"encoding/xml"
	"net"
	"strings"

	"strconv"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
	//net "github.com/shirou/gopsutil/net"
)

const (
	cSysProductCmd string = `lshw -c system | grep product | cut -d ":" -f2`
	cCPUArcCmd     string = `lscpu | grep Architecture | cut -d ":" -f2`
	cSysTz         string = "date +%z"
	cSysTzd        string = "date +%Z"
	cSysSerialNo   string = "dmidecode -s system-serial-number"
	cSysHostname   string = "hostname"
	cListHwAsXML   string = "lshw -c system,memory,bus,disk,volume,network -xml"
)

var (
	v              *List
	supportedDisks = [...]string{"hd", "sd", "xvd", "sr"}
)

// Memory Proc related constants
const (
	cMemProcPath                   string = "/proc/meminfo"
	cMemProcPhysicalTotalBytes     string = "MemTotal"
	cMemProcPhysicalAvailableBytes string = "MemFree"
	cMemProcPageAvailableBytes     string = "SwapFree"
	cMemProcPageTotalBytes         string = "SwapTotal"
	cPartitionProcPath             string = "/proc/partitions"
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

func (a assetDalImpl) readHwList() (*List, error) {

	if v == nil {
		hw, err := a.Factory.GetEnv().ExecuteBash(cListHwAsXML)
		if err != nil {
			return nil, exception.New(model.ErrExecuteCommandFailed, err)
		}
		v = new(List)
		err = xml.Unmarshal([]byte(hw), v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func (a assetDalImpl) getRequiredNode(l *Node, id string, class string) *Node {
	//if l.ID == id && l.Class == class {
	if strings.Contains(l.ID, id) && l.Class == class {
		return l
	}
	if len(l.Nodelist) > 0 {
		for i := range l.Nodelist {
			tmp := a.getRequiredNode(&l.Nodelist[i], id, class)
			if tmp != nil {
				return tmp
			}

		}
	}
	return nil
}

func (a assetDalImpl) getAllNodes(root *Node, id string, class string, listOfNodes []Node) []Node {
	//if root.ID == id && root.Class == class {
	if strings.Contains(root.ID, id) && root.Class == class {
		return append(listOfNodes, *root)
	}
	if len(root.Nodelist) > 0 {
		for i := range root.Nodelist {
			listOfNodes = a.getAllNodes(&root.Nodelist[i], id, class, listOfNodes)
		}
	}
	return listOfNodes
}

func (a assetDalImpl) getAllPartitions(root *Node, listOfPart []string) []string {
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

//GetBiosInfo ...
func (a assetDalImpl) GetBiosInfo() (*asset.AssetBios, error) {
	hlist, err := a.readHwList()
	if err != nil {
		return nil, err
	}
	l := &hlist.Nodelist

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

//GetBaseBoardInfo ...
func (a assetDalImpl) GetBaseBoardInfo() (*asset.AssetBaseBoard, error) {
	hlist, err := a.readHwList()
	if err != nil {
		return nil, err
	}
	l := &hlist.Nodelist
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

//GetDrivesInfo ...
func (a assetDalImpl) GetDrivesInfo() ([]asset.AssetDrive, error) {
	hlist, err := a.readHwList()
	if err != nil {
		return nil, err
	}
	l := &hlist.Nodelist
	var listOfNodes []Node
	var listOfDrives []asset.AssetDrive
	var tmp asset.AssetDrive
	diskList := a.getAllNodes(l, "disk", "disk", listOfNodes)
	for _, value := range diskList {
		tmp.Manufacturer = value.Vendor
		tmp.Product = value.Product
		tmp.SizeBytes = value.SizeInBytes
		if len(value.LogName) > 0 {
			tmp.LogicalName = value.LogName[0].Text
		}
		tmp.SerialNumber = value.Serial
		var listOfPart []string
		tmp.Partitions = a.getAllPartitions(&value, listOfPart)

		listOfDrives = append(listOfDrives, tmp)
	}

	optDriveList := a.getAllNodes(l, "cdrom", "disk", listOfNodes)
	for _, value := range optDriveList {
		tmp.Manufacturer = value.Vendor
		tmp.Product = value.Desc
		tmp.SizeBytes = value.SizeInBytes
		if len(value.LogName) > 0 {
			tmp.LogicalName = value.LogName[0].Text
		}
		tmp.SerialNumber = value.Serial
		var listOfPart []string
		tmp.Partitions = a.getAllPartitions(&value, listOfPart)

		listOfDrives = append(listOfDrives, tmp)
	}
	//Disk information for Amazon cloud which are based on Xen hypervisors is not returned through lshw command
	//use other method like lsblk instead.

	if len(listOfDrives) == 0 {
		listD, err := a.getDiskInfo()
		if err != nil {
			return listOfDrives, err
		}
		return listD, nil
	}

	return listOfDrives, nil
}

func (a assetDalImpl) getDiskInfo() ([]asset.AssetDrive, error) {

	reader, err := a.Factory.GetEnv().GetFileReader(cPartitionProcPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	cfg := procParser.Config{
		ParserMode:    procParser.ModeTabular,
		IgnoreNewLine: true,
	}
	data, err := a.Factory.GetParser().Parse(cfg, reader)
	if err != nil {
		return nil, err
	}
	var listOfDrives []asset.AssetDrive
	for i := 0; i < len(data.Lines); i++ {

		if len(data.Lines[i].Values) > 3 {
			drive := data.Lines[i].Values[3]
			if !a.isValidDisk(drive) {
				continue
			}
			var tmp asset.AssetDrive

			driveSize, _ := procParser.GetInt64(data.Lines[i].Values[2])
			tmp.LogicalName = "/dev/" + drive
			tmp.SizeBytes = driveSize * 1024
			var k int
			for k = i + 1; k < len(data.Lines); k++ {
				if !strings.HasPrefix(data.Lines[k].Values[3], drive) {
					break
				}
				parition := data.Lines[k].Values[3]
				tmp.Partitions = append(tmp.Partitions, "/dev/"+parition)
			}
			i = k - 1

			listOfDrives = append(listOfDrives, tmp)
		}

	}

	return listOfDrives, nil
}

func (a assetDalImpl) isValidDisk(drive string) bool {
	for _, sDisk := range supportedDisks {
		if strings.HasPrefix(drive, sDisk) {
			return true
		}
	}
	return false
}

// GetOSInfo returns the OS info
func (a assetDalImpl) GetOSInfo() (*asset.AssetOs, error) {
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
func (a assetDalImpl) GetSystemInfo() (*asset.AssetSystem, error) {
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

func (a assetDalImpl) getVendorProduct(n1 []Node, s *asset.AssetNetwork) {
	//Try to get vendor and product for this network interface using lshw -c network command
	//Correlate lshw and golang's net package on interfce's logical name
	for _, v3 := range n1 {
		for _, v2 := range v3.LogName {
			if v2.Text == s.LogicalName {
				s.Vendor = v3.Vendor
				s.Product = v3.Product
			}
		}
	}

}

// GetNetworkInfo returns network info
func (a assetDalImpl) GetNetworkInfo() ([]asset.AssetNetwork, error) {
	var array []asset.AssetNetwork
	//var s asset.AssetNetwork
	s := asset.AssetNetwork{
		DhcpServer:       "0.0.0.0",
		IPv4:             "0.0.0.0",
		IPv6:             "::",
		SubnetMask:       "0.0.0.0",
		DefaultIPGateway: "0.0.0.0",
	}
	var n1 []Node
	var listOfNodes []Node
	//Get the result of  lshw -c network command
	hlist, err := a.readHwList()
	if err == nil {
		l := &hlist.Nodelist
		n1 = a.getAllNodes(l, "network", "network", listOfNodes)
	}

	//Get Network info using golang's net package
	interfStat, err := net.Interfaces()
	if err != nil {
		return array, err
	}
	for _, interf := range interfStat {
		s.LogicalName = interf.Name
		s.MacAddress = interf.HardwareAddr.String()

		a.getVendorProduct(n1, &s)
		//Continue getting remaining Network info using golang's net package
		s1, err := interf.Addrs()
		if err != nil {
			array = append(array, s)
			return array, err
		}
		for _, value := range s1 {
			ip, subnet, _ := net.ParseCIDR(value.String())
			if strings.Contains(ip.String(), ":") {
				s.IPv6 = ip.String()
			} else {
				s.IPv4 = ip.String()
			}
			if len(subnet.Mask) == 4 {
				s.SubnetMask = net.IPv4(subnet.Mask[0], subnet.Mask[1], subnet.Mask[2], subnet.Mask[3]).String()
			}
			//s.SubnetMask = subnet.Mask.String()
		}

		array = append(array, s)
	}
	return array, nil
}

// GetNetworkInfo1 returns network info
func (a assetDalImpl) GetNetworkInfo1() ([]asset.AssetNetwork, error) {
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
		n := asset.AssetNetwork{
			DhcpServer:       "0.0.0.0",
			IPv4:             "0.0.0.0",
			IPv6:             "::",
			SubnetMask:       "0.0.0.0",
			DefaultIPGateway: "0.0.0.0",
		}
		if val, ok := mapArr[k]["product"]; ok {
			n.Product = val[1]
		}
		if val, ok := mapArr[k]["vendor"]; ok {
			n.Vendor = val[1]
		}
		if val, ok := mapArr[k]["logical name"]; ok {
			n.LogicalName = val[1]
		}
		networks[k] = n
	}
	dataCmd, err = util.getCommandData(parser, cfg, "bash", "-c", "nmcli dev list")
	if err != nil {
		return mapToArr(networks), exception.New(model.ErrExecuteCommandFailed, err)
	}
	mapArr = util.getProcData(dataCmd, "GENERAL.DEVICE", "GENERAL.DEVICE")

	setValnmcli(networks, mapArr)

	return mapToArr(networks), nil
}

func setValnmcli(networks map[string]asset.AssetNetwork, mapArr map[string]map[string][]string) {
	for k := range mapArr {

		for mk := range mapArr[k] {
			if len(mapArr[k][mk]) == 1 {
				continue
			}
			n := networks[k]
			setValnmcli2(mapArr[k], mk, &n)
			networks[k] = n
		}
	}
}

func setValnmcli2(mapA map[string][]string, mk string, n *asset.AssetNetwork) {
	val := mapA[mk][1]
	if strings.HasPrefix(mk, "DHCP4.OPTION") {
		arr := strings.Split(val, "=")
		var key string
		if len(arr) == 0 {
			return
		}
		key = strings.TrimSpace(arr[0])
		switch key {
		case "dhcp_server_identifier":
			n.DhcpServer = strings.TrimSpace(arr[1])
		case "ip_address":
			n.IPv4 = strings.TrimSpace(arr[1])
		case "subnet_mask":
			n.SubnetMask = strings.TrimSpace(arr[1])
		case "domain_name_servers":
			n.DnsServers = strings.Split(strings.TrimSpace(arr[1]), " ")
		}
	}
	//Mac address
	if mk == "GENERAL.HWADDR" {
		n.MacAddress = strings.Join(mapA[mk][1:], ":")
	}
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
func (a assetDalImpl) GetMemoryInfo() (*asset.AssetMemory, error) {
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
func (a assetDalImpl) GetProcessorInfo() ([]asset.AssetProcessor, error) {
	parser := a.Factory.GetParser()
	cfg := procParser.Config{
		ParserMode: procParser.ModeSeparator,
		Separator:  ":",
	}
	util := dalUtil{
		envDep: a.Factory,
	}
	dataFile, err := util.getFileData(parser, cfg, "/proc/cpuinfo")
	if err != nil {
		return nil, exception.New(model.ErrFileReadFailed, err)
	}
	cpuType, err := a.Factory.GetEnv().ExecuteBash(cCPUArcCmd)
	if err != nil {
		return nil, exception.New(model.ErrExecuteCommandFailed, err)
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

func (a assetDalImpl) GetInstalledSoftwareInfo() ([]asset.AssetInstalledSoftware, error) {
	return nil, nil
}
