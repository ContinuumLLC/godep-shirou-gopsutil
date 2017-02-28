package dal

import (
	"os/exec"
	"strings"
	"time"

	amodel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
)

func getOsInfo ()  *amodel.AssetOs {
    os := new(amodel.AssetOs)

    //Product
    cmdName := "lsb_release -a | grep Description | cut -d \":\" -f2"
    out, _ := exec.Command("bash", "-c", cmdName).Output()
    product := strings.Replace(string(out), "\n","",-1)
    product = strings.Replace(product, "\t","",-1)
    os.Product = product

    //OsManufacturer
    cmdName = "lsb_release -a | grep \"Distributor ID\" | cut -d \":\" -f2"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    manufacturer := strings.Replace(string(out), "\n","",-1)
    manufacturer = strings.Replace(manufacturer, "\t","",-1)
    os.Manufacturer = manufacturer

    //OsLanguage
    cmdName = "cat /etc/default/locale | grep  \"\bLANG\b\" | cut -d \"=\" -f2"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    language := strings.Replace(string(out), "\n","",-1)
    language = strings.Replace(language, "\t","",-1)
    os.OsLanguage = language

    //OsVersion
    cmdName = "lsb_release -a | grep Release | cut -d \":\" -f2"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    version := strings.Replace(string(out), "\n","",-1)
    version = strings.Replace(version, "\t","",-1)
    os.Version = version

    //OsInstallDate
    cmdName = "ls -ld /var/log/installer | cut -d \" \" -f6,7,8"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    installDateString := strings.Replace(string(out), "\n","",-1)
    installDateString = strings.Replace(installDateString, "\t","",-1)
    t, _ := time.Parse("Jan 13, at 12:05", installDateString)
    os.InstallDate = t

    return os
}
