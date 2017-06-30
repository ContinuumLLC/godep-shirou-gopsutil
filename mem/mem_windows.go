// +build windows

package mem

import (
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/internal/common"
	"golang.org/x/sys/windows"
)

var (
	procGlobalMemoryStatusEx = common.Modkernel32.NewProc("GlobalMemoryStatusEx")
)

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64 // in bytes
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

// Win32_PerfFormattedData_PerfOS_Memory struct to provide performance memory metrics for windows
type Win32_PerfFormattedData_PerfOS_Memory struct {
	AvailableBytes             uint64
	CommittedBytes             uint64
	PercentCommittedBytesInUse uint32
	FreeSystemPageTableEntries uint32
	PagesPerSec                uint32
	PoolNonpagedBytes          uint64
}

func VirtualMemory() (*VirtualMemoryStat, error) {
	var memInfo memoryStatusEx
	memInfo.cbSize = uint32(unsafe.Sizeof(memInfo))
	mem, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if mem == 0 {
		return nil, windows.GetLastError()
	}

	ret := &VirtualMemoryStat{
		Total:       memInfo.ullTotalPhys,
		Available:   memInfo.ullAvailPhys,
		UsedPercent: float64(memInfo.dwMemoryLoad),
	}

	ret.Used = ret.Total - ret.Available
	return ret, nil
}

func SwapMemory() (*SwapMemoryStat, error) {
	ret := &SwapMemoryStat{}

	return ret, nil
}

// PerfInfo returns the performance data from performance counters of memory object.
func PerfInfo() ([]Win32_PerfFormattedData_PerfOS_Memory, error) {
	var ret []Win32_PerfFormattedData_PerfOS_Memory
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}
