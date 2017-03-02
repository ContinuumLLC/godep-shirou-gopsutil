package linux

import (
	"os/exec"
	"strings"
	amodel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
)

func getSystemInfo ()  *amodel.AssetSystem {
    system := new(amodel.AssetSystem)
   
    //Product
    cmdName := "sudo dmidecode -s system-product-name"
    out, _ := exec.Command("bash", "-c", cmdName).Output()
    product := strings.Replace(string(out), "\n","",-1)
    system.Product = product

    //TimeZone
    cmdName = "date +%z"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    timeZone := strings.Replace(string(out), "\n","",-1)
    system.TimeZone = timeZone

    //TimeZoneDescription
    cmdName = "date +%Z"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    timeZoneDescription := strings.Replace(string(out), "\n","",-1)
    system.TimeZoneDescription = timeZoneDescription

    //SystemSerialNumber
    cmdName = "sudo dmidecode -s system-serial-number"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    serialNumber := strings.Replace(string(out), "\n","",-1)
    system.SerialNumber = serialNumber

    //SystemName
    cmdName = "hostname"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    systemName := strings.Replace(string(out), "\n","",-1)
    system.SystemName = systemName

    return system
}
