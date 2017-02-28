package dal

import (
	"os/exec"
	"strings"
	amodel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
)

func getBaseBoardInfo ()  *amodel.AssetBaseBoard {
    baseBoard := new(amodel.AssetBaseBoard)
    cmdName := "sudo dmidecode -s baseboard-product-name"
    out, _ := exec.Command("bash", "-c", cmdName).Output()
    product := strings.Replace(string(out), "\n","",-1)
    baseBoard.Product = product

    //BaseBoardManufacturer
    cmdName = "sudo dmidecode -s baseboard-manufacturer"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    manufacturer := strings.Replace(string(out), "\n","",-1)
    baseBoard.Manufacturer = manufacturer

    //BaseBoardSerialNumber
    cmdName = "sudo dmidecode -s baseboard-serial-number"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    serialNumber := strings.Replace(string(out), "\n","",-1)
    baseBoard.SerialNumber = serialNumber

    //BaseBoardVersion
    cmdName = "sudo dmidecode -s baseboard-version"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    version := strings.Replace(string(out), "\n","",-1)
    baseBoard.Version = version

    return baseBoard
}
