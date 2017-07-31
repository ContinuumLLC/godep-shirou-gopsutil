package mem

import (
	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/wmi"
)

const memQuery = "SELECT Name, Manufacturer, SerialNumber, Capacity FROM Win32_PhysicalMemory"

// WMI by name itself is the WMI implementation to collect Asset Physical Memory information
type WMI struct {
	dep wmi.Wrapper
}

//GetByWMI returns the WMI implementation of Asset Physical Memory information
func GetByWMI() WMI {
	return WMI{
		dep: wmi.GetWrapper(),
	}
}

type win32PhysicalMemory struct {
	Manufacturer string
	SerialNumber string
	Capacity     uint64
}

// Info returns the Asset Physical Memory information
func (w WMI) Info() ([]asset.PhysicalMemory, error) {
	var dst []win32PhysicalMemory
	err := w.dep.Query(memQuery, &dst)
	if err != nil {
		return nil, err
	}
	return mapToMemModel(dst), nil
}

func mapToMemModel(dst []win32PhysicalMemory) (ret []asset.PhysicalMemory) {
	if dst != nil {
		l := len(dst)
		ret = make([]asset.PhysicalMemory, l)
		for i := 0; i < l; i++ {
			ret[i].Manufacturer = dst[i].Manufacturer
			ret[i].SerialNumber = dst[i].SerialNumber
			ret[i].SizeBytes = dst[i].Capacity
		}
	}
	return
}
