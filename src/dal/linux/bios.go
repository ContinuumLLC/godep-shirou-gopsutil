package linux

import (
	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
)

func getBiosInfo() *asset.AssetBios {
	//TODO : error handling
	bios := new(asset.AssetBios)
	util := dalUtil{}
	//BiosManufacturer
	cmdName := "sudo dmidecode -s bios-vendor"
	manufacturer, _ := util.execCommand(cmdName)
	bios.Manufacturer = manufacturer

	//BiosSerialNumber
	cmdName = "sudo dmidecode -s baseboard-serial-number"
	serialNumber, _ := util.execCommand(cmdName)
	bios.SerialNumber = serialNumber

	//BiosVersion
	cmdName = "sudo dmidecode -s bios-version"
	ver, _ := util.execCommand(cmdName)
	bios.Version = ver

	//SMBIOSVersion
	cmdName = "sudo dmidecode --type bios | grep SMBIOS |  awk '{print $2}'"
	smbiosVer, _ := util.execCommand(cmdName)
	bios.SmbiosVersion = smbiosVer

	return bios
}
