package dal

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"strings"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	eMock "github.com/ContinuumLLC/platform-common-lib/src/env/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
	pMock "github.com/ContinuumLLC/platform-common-lib/src/procParser/mock"
	"github.com/golang/mock/gomock"
)

const (
	hwXML  string = "<list></list>"
	hwXML1 string = `<?xml version="1.0" standalone="yes" ?>
<list>
<node id="milinda-virtualbox" claimed="true" class="system" handle="DMI:0001">
 <description>Computer</description>
 <product>VirtualBox</product>
 <vendor>innotek GmbH</vendor>
 <version>1.2</version>
 <serial>0</serial>
 <width units="bits">64</width>
 <configuration>
  <setting id="family" value="Virtual Machine" />
  <setting id="uuid" value="E8F8DC61-D6D2-413D-86FE-A0E3E9807FC2" />
 </configuration>
 <capabilities>
  <capability id="smbios-2.5" >SMBIOS version 2.5</capability>
  <capability id="dmi-2.5" >DMI version 2.5</capability>
  <capability id="vsyscall32" >32-bit processes</capability>
 </capabilities>
  <node id="core" claimed="true" class="bus" handle="DMI:0008">
   <description>Motherboard</description>
   <product>VirtualBox</product>
   <vendor>Oracle Corporation</vendor>
   <physid>0</physid>
   <version>1.2</version>
   <serial>0</serial>
    <node id="firmware" claimed="true" class="memory" handle="">
     <description>BIOS</description>
     <vendor>innotek GmbH</vendor>
     <physid>0</physid>
     <version>VirtualBox</version>
     <date>12/01/2006</date>
     <size units="bytes">131072</size>
    </node>
  <node id="cdrom" claimed="true" class="disk" handle="SCSI:01:00:00:00">
   <description>DVD reader</description>
   <logicalname>/dev/cdrom</logicalname>
   <logicalname>/dev/dvd</logicalname>
   <logicalname>/dev/sr0</logicalname>
  </node>
  <node id="disk" claimed="true" class="disk" handle="SCSI:02:00:00:00">
   <description>ATA Disk</description>
   <product>VBOX HARDDISK</product>
   <logicalname>/dev/sda</logicalname>
   <version>1.0</version>
   <serial>VB7f4a1ba4-6ef7655d</serial>
   <size units="bytes">64424509440</size>
    <node id="volume:0" claimed="true" class="volume" handle="">
     <description>EXT4 volume</description>
     <vendor>Linux</vendor>
     <logicalname>/dev/sda1</logicalname>
     <logicalname>/</logicalname>
     <version>1.0</version>
     <serial>ecc36b30-d5d2-40a0-9962-88661930be29</serial>
     <size units="bytes">55833526272</size>
     <capacity>55833526272</capacity>
    </node>
    <node id="volume:1" claimed="true" class="volume" handle="">
     <description>Extended partition</description>
     <logicalname>/dev/sda2</logicalname>
     <size units="bytes">8587838464</size>
     <capacity>8587838464</capacity>
      <node id="logicalvolume" claimed="true" class="volume" handle="">
       <description>Linux swap / Solaris partition</description>
       <logicalname>/dev/sda5</logicalname>
       <capacity>8587837440</capacity>
      </node>
    </node>
  </node>
  </node>
</node>
</list>`
)

func setupGetCommandReader(t *testing.T, parseErr error, commandReaderErr error) (*gomock.Controller, *mock.MockAssetDalDependencies) {
	ctrl := gomock.NewController(t)
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any()).Return(reader, commandReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if commandReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(&procParser.Data{}, parseErr)
	}
	mockAssetDalD.EXPECT().GetParser().Return(mockParser)

	return ctrl, mockAssetDalD
}

//TODO - Duplicate function as setupGetCommandReader. Need to relook at it.
func setupGetCommandReader2(t *testing.T, parseErr error, commandReaderErr error, data *procParser.Data) (*gomock.Controller, *mock.MockAssetDalDependencies) {
	ctrl := gomock.NewController(t)
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any()).Return(reader, commandReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if commandReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(data, parseErr)
	}
	mockAssetDalD.EXPECT().GetParser().Return(mockParser)

	return ctrl, mockAssetDalD
}

func setupGetFileReader(t *testing.T, parseErr error, fileReaderErr error, parseData *procParser.Data) (*gomock.Controller, *mock.MockAssetDalDependencies) {
	ctrl := gomock.NewController(t)
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, fileReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if fileReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(parseData, parseErr)
	}
	mockAssetDalD.EXPECT().GetParser().Return(mockParser)

	return ctrl, mockAssetDalD
}

func setupAddGetFileReader(ctrl *gomock.Controller, mockAssetDalD *mock.MockAssetDalDependencies, parseErr error, fileReaderErr error) {
	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, fileReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if fileReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(&procParser.Data{}, parseErr)
	}
}

func TestGetOSCommandErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader(t, nil, errors.New(model.ErrExecuteCommandFailed))
	defer ctrl.Finish()

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetOSInfo()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}

func TestGetOSFileErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader(t, nil, nil)
	defer ctrl.Finish()

	setupAddGetFileReader(ctrl, mockAssetDalD, nil, errors.New(model.ErrFileReadFailed))

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetOSInfo()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrFileReadFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrFileReadFailed, err)
	}
}

// TODO - fix error
// func TestGetOSNoErr(t *testing.T) {
// 	ctrl, mockAssetDalD := setupGetCommandReader(t, nil, nil)
// 	defer ctrl.Finish()

// 	setupAddGetFileReader(ctrl, mockAssetDalD, nil, nil)

// 	log := logging.GetLoggerFactory().New("")
// 	log.SetLogLevel(logging.OFF)
// 	_, err := assetDalImpl{
// 		Factory: mockAssetDalD,
// 		Logger:  log,
// 	}.GetOS()
// 	if err != nil {
// 		t.Errorf("Unexpected error : %v", err)
// 	}
// }

func setupGetSystemInfo(t *testing.T, times int, err error) (*gomock.Controller, error) {
	ctrl := gomock.NewController(t)

	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	var str string
	switch times {
	case 1:
		str = cSysProductCmd
	case 2:
		mockEnv.EXPECT().ExecuteBash(cSysProductCmd).Return("", nil)
		str = cSysTz
	case 3:
		mockEnv.EXPECT().ExecuteBash(cSysProductCmd).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTz).Return("", nil)
		str = cSysTzd
	case 4:
		mockEnv.EXPECT().ExecuteBash(cSysProductCmd).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTz).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTzd).Return("", nil)
		str = cSysSerialNo
	case 5:
		mockEnv.EXPECT().ExecuteBash(cSysProductCmd).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTz).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTzd).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysSerialNo).Return("", nil)
		str = cSysHostname

	}
	mockEnv.EXPECT().ExecuteBash(str).Return("", err)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv).Times(times)

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetSystemInfo()
	return ctrl, e
}

func TestGetSystemInfoErr(t *testing.T) {
	cmdExeArr := []int{1, 2, 3, 4, 5}
	for _, i := range cmdExeArr {
		ctrl, err := setupGetSystemInfo(t, i, errors.New(model.ErrExecuteCommandFailed))
		defer ctrl.Finish()
		if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
			t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
		}
	}
}

func TestGetSystemNoErr(t *testing.T) {
	ctrl, err := setupGetSystemInfo(t, 5, nil)
	defer ctrl.Finish()
	if err != nil {
		t.Errorf("Unexpected error received  : %v", err)
	}
}

func TestGetMemoryInfoErr(t *testing.T) {
	parseError := model.ErrFileReadFailed
	_, mockAssetDalD := setupGetFileReader(t, errors.New(parseError), nil, nil)
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetMemoryInfo()

	if err == nil || err.Error() != parseError {
		t.Error("Error expected but not returned")
	}
}

func TestGetDataFromMap(t *testing.T) {
	data := procParser.Data{
		Map: make(map[string]procParser.Line, 1),
	}
	data.Map["MemTotal"] = procParser.Line{Values: []string{"MemTotal", "InvalidNumber", "KB"}}

	util := dalUtil{}
	val := util.getDataFromMap("MemTotal", &data)

	if val != 0 {
		t.Errorf("Expected 0, returned %d", val)
	}
}

func TestGetMemoryInfoNoErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := procParser.Data{
		Map: make(map[string]procParser.Line, 1),
	}
	data.Map["MemTotal"] = procParser.Line{Values: []string{"physicalTotalBytes", "1", "KB"}}

	_, mockAssetDalD := setupGetFileReader(t, nil, nil, &data)
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetMemoryInfo()

	if err != nil {
		t.Errorf("Unexpected error received  : %v", err)
	}
}

func TestGetProcessorInfoErr(t *testing.T) {
	parseError := model.ErrFileReadFailed
	ctrl, mockAssetDalD := setupGetFileReader(t, errors.New(parseError), nil, nil)
	defer ctrl.Finish()

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetProcessorInfo()

	if err == nil || err.Error() != parseError {
		t.Error("Error expected but not returned")
	}
}

func TestGetProcessorInfoBashErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetFileReader(t, nil, nil, nil)
	defer ctrl.Finish()
	envMock := eMock.NewMockEnv(ctrl)
	envMock.EXPECT().ExecuteBash(cCPUArcCmd).Return("", errors.New(model.ErrExecuteCommandFailed))
	mockAssetDalD.EXPECT().GetEnv().Return(envMock)

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetProcessorInfo()

	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}

func TestGetProcessorInfoNoErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetFileReader(t, nil, nil, &procParser.Data{})
	defer ctrl.Finish()
	envMock := eMock.NewMockEnv(ctrl)
	envMock.EXPECT().ExecuteBash(cCPUArcCmd).Return("", nil)
	mockAssetDalD.EXPECT().GetEnv().Return(envMock)

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetProcessorInfo()

	if err != nil {
		t.Errorf("Unexpected error returned : %v", err)
	}
}

func setupEnv(t *testing.T) (*gomock.Controller, *mock.MockAssetDalDependencies, *eMock.MockEnv) {
	ctrl := gomock.NewController(t)
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv).Times(1)
	v = nil

	return ctrl, mockAssetDalD, mockEnv
}

func TestReadHwList(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return("<list></list>", nil)

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.readHwList()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestReadHwListError(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return("<list></list>", errors.New("readHwListErr"))

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	v = nil
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.readHwList()
	if e == nil || e.Error() != model.ErrExecuteCommandFailed {
		t.Errorf("Expecting model.ErrExecuteCommandFailed , Unexpected error %v", e)
	}
}

func TestReadHwListErr(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return("nu$756ll", nil)

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.readHwList()
	if e == nil {
		t.Error("Expecting EOF error ")
	}
}

func TestGetBiosInfo(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML, nil)
	// v := List{}
	// xml.Unmarshal([]byte(hwXML), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetBiosInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetBiosInfo2(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, nil)
	// v := List{}
	// xml.Unmarshal([]byte(hwXML1), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetBiosInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetBiosInfoError(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, errors.New("XMLReadError"))
	// v := List{}
	// xml.Unmarshal([]byte(hwXML1), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetBiosInfo()
	if e == nil || e.Error() != model.ErrExecuteCommandFailed {
		t.Errorf("Unexpected error %v, was expecting model.ErrExecuteCommandFailed", e)
	}
}

func TestGetBaseBoardInfo(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML, nil)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetBaseBoardInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetBaseBoardInfo2(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, nil)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetBaseBoardInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetBaseBoardInfoError(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, errors.New("XMLReadErr"))
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetBaseBoardInfo()
	if e == nil || e.Error() != model.ErrExecuteCommandFailed {
		t.Errorf("Unexpected error %v, was expecting model.ErrExecuteCommandFailed ", e)
	}
}

func TestGetDrivesInfo(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML, nil)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetDrivesInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}
func TestGetDrivesInfo2(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, nil)
	// v := List{}
	// xml.Unmarshal([]byte(hwXML1), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetDrivesInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetDrivesInfoError(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, errors.New("XMLReadErr"))
	// v := List{}
	// xml.Unmarshal([]byte(hwXML1), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetDrivesInfo()
	if e == nil || e.Error() != model.ErrExecuteCommandFailed {
		t.Errorf("Unexpected error %v, was expecting model.ErrExecuteCommandFailed", e)
	}
}

func TestMapToArr(t *testing.T) {
	m := map[string]asset.AssetNetwork{
		"eth0": asset.AssetNetwork{},
		"eth1": asset.AssetNetwork{},
	}
	nArr := mapToArr(m)
	if l := len(nArr); l != 2 {
		t.Errorf("Expected length is %d but received %d", 2, l)
	}
}

func TestSetValnmcli(t *testing.T) {
	networks := map[string]asset.AssetNetwork{
		"eth0": asset.AssetNetwork{},
		"eth1": asset.AssetNetwork{},
	}
	mapArr := map[string]map[string][]string{
		"eth0": {
			"DHCP4.OPTION[11]":         []string{"DHCP4.OPTION[11]", "dhcp_server_identifier = 10.0.3.2"},
			"DHCP4.OPTION[9]":          []string{"DHCP4.OPTION[9]", "domain_name_servers = 10.2.17.6 10.2.17.25 10.2.17.17"},
			"DHCP4.OPTION[6]":          []string{"DHCP4.OPTION[6]", "ip_address = 10.0.3.15"},
			"DHCP4.OPTION[5]":          []string{"DHCP4.OPTION[5]", "ip_address 10.0.3.15"},
			"DHCP4.OPTION[7]":          []string{"DHCP4.OPTION[7]", "subnet_mask = 255.255.255.0"},
			"GENERAL.HWADDR":           []string{"GENERAL.HWADDR", "08:00:27:09:C7:82"},
			"GENERAL.FIRMWARE-VERSION": []string{"GENERAL.FIRMWARE-VERSION"},
		},
	}
	setValnmcli(networks, mapArr)
	if d := networks["eth0"].DhcpServer; d != "10.0.3.2" {
		t.Errorf("Expected value is 10.0.3.2 but received %s", d)
	}
	if i := networks["eth0"].IPv4; i != "10.0.3.15" {
		t.Errorf("Expected value is 10.0.3.15 but received %s", i)
	}
}

func TestGetNetworkInfo(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader2(t, nil, errors.New(model.ErrExecuteCommandFailed), &procParser.Data{})
	defer ctrl.Finish()

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetNetworkInfo()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}

func TestGetNetworkInfoCommandDataErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader2(t, nil, nil, &procParser.Data{
		Lines: []procParser.Line{
			procParser.Line{
				Values: []string{"*-network", "0"},
			},
			procParser.Line{
				Values: []string{"*product", "82540EM Gigabit Ethernet Controller"},
			},
			procParser.Line{
				Values: []string{"*-network", "1"},
			},
		},
	})

	mockEnv := eMock.NewMockEnv(ctrl)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv).Times(1)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New(model.ErrExecuteCommandFailed))

	defer ctrl.Finish()

	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetNetworkInfo()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}

func TestGetNetworkInfoCommandDataErr1(t *testing.T) {
	data := &procParser.Data{
		Lines: []procParser.Line{
			procParser.Line{
				Values: []string{"*-network", "0"},
			},
			procParser.Line{
				Values: []string{"*product", "82540EM Gigabit Ethernet Controller"},
			},
			procParser.Line{
				Values: []string{"logical name", "enp0s3"},
			},
			procParser.Line{
				Values: []string{"*-network", "1"},
			},
			procParser.Line{
				Values: []string{"logical name", "enp0s4"},
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	mockParser := pMock.NewMockParser(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	gomock.InOrder(
		mockAssetDalD.EXPECT().GetParser().Return(mockParser),
		mockAssetDalD.EXPECT().GetEnv().Return(mockEnv),
		mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any()).Return(reader, nil),
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(data, nil),

		mockAssetDalD.EXPECT().GetEnv().Return(mockEnv),
		mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any()).Return(reader, errors.New("Err")),
	)
	log := logging.GetLoggerFactory().New("")
	log.SetLogLevel(logging.OFF)
	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  log,
	}.GetNetworkInfo()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}
