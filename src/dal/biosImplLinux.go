package dal

import (
	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
)

func getBiosInfo() *asset.AssetBios {
	bios := new(asset.AssetBios)
	util := dalUtil{}
	//BiosManufacturer
	cmdName := "sudo dmidecode -s bios-vendor"
	manufacturer, err := util.execCommand(cmdName)
	if err != nil {
		//TODO : Add logging
	}
	bios.Manufacturer = manufacturer

	//BiosSerialNumber
	cmdName = "sudo dmidecode -s baseboard-serial-number"
	serialNumber, err := util.execCommand(cmdName)
	bios.SerialNumber = serialNumber

	//BiosVersion
	cmdName = "sudo dmidecode -s bios-version"
	ver, err := util.execCommand(cmdName)
	bios.Version = ver

	//SMBIOSVersion
	cmdName = "sudo dmidecode --type bios | grep SMBIOS |  awk '{print $2}'"
	smbiosVer, err := util.execCommand(cmdName)
	bios.SmbiosVersion = smbiosVer

	return bios
}
