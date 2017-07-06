package processor

import (
	"strconv"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/StackExchange/wmi"
)

//WMI ...
type WMI struct{}

type win32Processor struct {
	Name                      string
	NumberOfCores             int
	CurrentClockSpeed         int
	Family                    int
	Manufacturer              string
	NumberOfLogicalProcessors int
	DataWidth                 int
}

//Info to return the Asset Processor information
func (WMI) Info() ([]asset.AssetProcessor, error) {
	var dst []win32Processor
	err := wmi.Query("SELECT Name,NumberOfCores,CurrentClockSpeed,Family,Manufacturer,NumberOfLogicalProcessors,DataWidth FROM Win32_Processor", &dst)
	if err != nil {
		return nil, err
	}
	l := len(dst)
	data := make([]asset.AssetProcessor, l)
	for i := 0; i < l; i++ {
		data[i].Product = dst[i].Name
		data[i].ClockSpeedMhz = float64(dst[i].CurrentClockSpeed)
		data[i].Family = dst[i].Family
		data[i].Manufacturer = dst[i].Manufacturer
		data[i].NumberOfCores = dst[i].NumberOfCores
		data[i].ProcessorType = strconv.Itoa(dst[i].DataWidth)
		//data[i].SerialNumber = dst[i].SerialNumber
	}
	return data, nil
}
