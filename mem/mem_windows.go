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
	PagesOutputPerSec          uint32
}

// Win32_OperatingSystem struct to provide virtual memory values
type Win32_OperatingSystem struct {
	//Virtual memory total and free in KBs
	TotalVirtualMemorySize *uint64
	FreeVirtualMemory      *uint64
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

	// GlobalMemoryStatusEx WinAPI retrieves virtual memory information
	// but does not match with the one that is displayed by the system information application run on the same system.
	// (Start->Program->Accessories->System Tools->System Information).
	// https://groups.google.com/forum/#!topic/microsoft.public.vc.mfc/i7UzUJOYziE
	var dst []Win32_OperatingSystem
	var totalVirtualMemorySize, freeVirtualMemory uint64
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		return ret, err
	}
	if dst[0].TotalVirtualMemorySize != nil {
		totalVirtualMemorySize = *(dst[0].TotalVirtualMemorySize)
	}
	if dst[0].FreeVirtualMemory != nil {
		freeVirtualMemory = *(dst[0].FreeVirtualMemory)
	}

	ret.TotalVirtual = totalVirtualMemorySize * 1024 // in bytes
	ret.AvailableVirtual = freeVirtualMemory * 1024  // in bytes
	ret.UsedVirtual = ret.TotalVirtual - ret.AvailableVirtual
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
