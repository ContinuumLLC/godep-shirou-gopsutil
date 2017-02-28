package dal

import (
	"os/exec"
	"strings"
	amodel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
)

func getBiosInfo ()  *amodel.AssetBios {
    bios := new(amodel.AssetBios)

    //BiosManufacturer
    cmdName := "sudo dmidecode -s bios-vendor"
    out, _ := exec.Command("bash", "-c", cmdName).Output()
    manufacturer := strings.Replace(string(out), "\n","",-1)
    bios.Manufacturer = manufacturer

    //BiosSerialNumber
    cmdName = "sudo dmidecode -s baseboard-serial-number"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    serialNumber := strings.Replace(string(out), "\n","",-1)
    bios.SerialNumber = serialNumber

    //BiosVersion
    cmdName = "sudo dmidecode -s bios-version"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    version := strings.Replace(string(out), "\n","",-1)
    bios.Version = version

    //SMBIOSVersion
    cmdName = "sudo dmidecode --type bios | grep SMBIOS |  awk '{print $2}'"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    smbiosVersion := strings.Replace(string(out), "\n","",-1)
    bios.SmbiosVersion = smbiosVersion

    return bios
}
