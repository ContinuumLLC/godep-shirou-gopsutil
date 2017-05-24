// +build windows

package cpu

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"

	"github.com/shirou/gopsutil/internal/common"
)

type Win32_Processor struct {
	LoadPercentage            *uint16
	Family                    uint16
	Manufacturer              string
	Name                      string
	NumberOfLogicalProcessors uint32
	ProcessorID               *string
	Stepping                  *string
	MaxClockSpeed             uint32
}

// Win32_PerfFormattedData_Counters_ProcessorInformation stores instance value of the perf counters
type Win32_PerfFormattedData_Counters_ProcessorInformation struct {
	Name                        string
	PercentDPCTime              uint64
	PercentIdleTime             uint64
	PercentUserTime             uint64
	PercentProcessorTime        uint64
	PercentInterruptTime        uint64
	PercentProcessorUtility     uint64
	PercentPriorityTime         uint64
	PercentPrivilegedTime       uint64
	PercentProcessorPerformance uint64
	InterruptsPerSec            uint32
	ProcessorFrequency          uint32
	DPCRate                     uint32
}

// TODO: Get percpu
func Times(percpu bool) ([]TimesStat, error) {
	var ret []TimesStat

	var lpIdleTime common.FILETIME
	var lpKernelTime common.FILETIME
	var lpUserTime common.FILETIME
	r, _, _ := common.ProcGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&lpIdleTime)),
		uintptr(unsafe.Pointer(&lpKernelTime)),
		uintptr(unsafe.Pointer(&lpUserTime)))
	if r == 0 {
		return ret, syscall.GetLastError()
	}

	LOT := float64(0.0000001)
	HIT := (LOT * 4294967296.0)
	idle := ((HIT * float64(lpIdleTime.DwHighDateTime)) + (LOT * float64(lpIdleTime.DwLowDateTime)))
	user := ((HIT * float64(lpUserTime.DwHighDateTime)) + (LOT * float64(lpUserTime.DwLowDateTime)))
	kernel := ((HIT * float64(lpKernelTime.DwHighDateTime)) + (LOT * float64(lpKernelTime.DwLowDateTime)))
	system := (kernel - idle)

	ret = append(ret, TimesStat{
		Idle:   float64(idle),
		User:   float64(user),
		System: float64(system),
	})
	return ret, nil
}

func Info() ([]InfoStat, error) {
	var ret []InfoStat
	var dst []Win32_Processor
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		return ret, err
	}

	var procID string
	for i, l := range dst {
		procID = ""
		if l.ProcessorID != nil {
			procID = *l.ProcessorID
		}

		cpu := InfoStat{
			CPU:        int32(i),
			Family:     fmt.Sprintf("%d", l.Family),
			VendorID:   l.Manufacturer,
			ModelName:  l.Name,
			Cores:      int32(l.NumberOfLogicalProcessors),
			PhysicalID: procID,
			Mhz:        float64(l.MaxClockSpeed),
			Flags:      []string{},
		}
		ret = append(ret, cpu)
	}

	return ret, nil
}

// PerfInfo returns the performance counter's instance value for ProcessorInformation
func PerfInfo() ([]Win32_PerfFormattedData_Counters_ProcessorInformation, error) {
	var ret []Win32_PerfFormattedData_Counters_ProcessorInformation
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}
