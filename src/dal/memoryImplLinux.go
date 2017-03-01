package dal

import (
	"os/exec"
	"strings"
	"strconv"
	
	amodel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
)

func getMemoryInfo()  *amodel.AssetMemory {
    memory := new(amodel.AssetMemory)
    
    //TotalPhysicalMemoryBytes
    cmdName := "cat /proc/meminfo | grep MemTotal | cut -d \":\" -f2 | awk '{print $1}'"
    out, _ := exec.Command("bash", "-c", cmdName).Output()
    totalPhysicalMemoryBytes:= strings.Replace(string(out), "\n","",-1)
    memBytes, err := strconv.ParseInt(totalPhysicalMemoryBytes, 10, 64)
    if err != nil {
        panic(err)
    }
    memory.TotalPhysicalMemoryBytes = memBytes 

    //TotalVirtualMemoryBytes
    cmdName = "free -t | grep Total | cut -d \":\" -f2 | awk '{print $1}'"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    totalVirtualMemoryBytes := strings.Replace(string(out), "\n","",-1)
    memBytes, err = strconv.ParseInt(totalVirtualMemoryBytes, 10, 64)
    if err != nil {
        panic(err)
    }
    memory.TotalVirtualMemoryBytes = memBytes 

    //AvailableVirtualMemoryBytes
    cmdName = "free -t | grep Total | cut -d \":\" -f2 | awk '{print $3}'"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    availableVirtualMemoryBytes := strings.Replace(string(out), "\n","",-1)
    memBytes, err = strconv.ParseInt(availableVirtualMemoryBytes, 10, 64)
    if err != nil {
         panic(err)
    }
    memory.AvailableVirtualMemoryBytes = memBytes

    //AvailablePhysicalMemoryBytes
    cmdName = " cat /proc/meminfo | grep MemAvailable | cut -d \":\" -f2 | awk '{print $1}'"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    availablePhysicalMemoryBytes := strings.Replace(string(out), "\n","",-1)
    memBytes, err = strconv.ParseInt(availablePhysicalMemoryBytes, 10, 64)
    if err != nil {
        panic(err)
    }
    memory.AvailablePhysicalMemoryBytes = memBytes

    //TotalPageFileSpaceBytes
    cmdName = "cat /proc/meminfo | grep SwapTotal | cut -d \":\" -f2 |  awk '{print $1}'"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    totalPageFileSpaceBytes := strings.Replace(string(out), "\n","",-1)
    memBytes, err = strconv.ParseInt(totalPageFileSpaceBytes, 10, 64)
    if err != nil {
        panic(err)
    }
    memory.TotalPageFileSpaceBytes = memBytes

    //AvailablePageFileSpaceBytes
    cmdName = "cat /proc/meminfo | grep SwapFree | cut -d \":\" -f2 |  awk '{print $1}'"
    out, _ = exec.Command("bash", "-c", cmdName).Output()
    availablePageFileSpaceBytes := strings.Replace(string(out), "\n","",-1)
    memBytes, err = strconv.ParseInt(availablePageFileSpaceBytes, 10, 64)
    if err != nil {
        panic(err)
    }
    memory.AvailablePageFileSpaceBytes = memBytes

    return memory
}
