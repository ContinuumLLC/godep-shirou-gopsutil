package linux

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
    memory.TotalPhysicalMemoryBytes = execCommand(cmdName) 

    //TotalVirtualMemoryBytes
    cmdName = "free -t | grep Total | cut -d \":\" -f2 | awk '{print $1}'"
    memory.TotalVirtualMemoryBytes = execCommand(cmdName) 

    //AvailableVirtualMemoryBytes
    cmdName = "free -t | grep Total | cut -d \":\" -f2 | awk '{print $3}'"
    memory.AvailableVirtualMemoryBytes = execCommand(cmdName) 

    //AvailablePhysicalMemoryBytes
    cmdName = " cat /proc/meminfo | grep MemAvailable | cut -d \":\" -f2 | awk '{print $1}'"
    memory.AvailablePhysicalMemoryBytes = execCommand(cmdName) 

    //TotalPageFileSpaceBytes
    cmdName = "cat /proc/meminfo | grep SwapTotal | cut -d \":\" -f2 |  awk '{print $1}'"
    memory.TotalPageFileSpaceBytes = execCommand(cmdName) 

    //AvailablePageFileSpaceBytes
    cmdName = "cat /proc/meminfo | grep SwapFree | cut -d \":\" -f2 |  awk '{print $1}'"
    memory.AvailablePageFileSpaceBytes = execCommand(cmdName) 

    return memory
}

func execCommand(cmdName string) int64 {
    out, _ := exec.Command("bash", "-c", cmdName).Output()
    memStr := strings.Replace(string(out), "\n","",-1)
    memBytes, _ := strconv.ParseInt(memStr, 10, 64)
    return memBytes
}
