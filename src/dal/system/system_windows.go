package system

import (
	"fmt"
	"time"

	"math"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/wmi"
)

// WMI struct is to represent that System information is to be collected using WMI calls
type WMI struct {
	dep wmi.Wrapper
}

const (
	hour int = 60
)

//GetByWMI returns the WMI implementation System information
func GetByWMI() WMI {
	return WMI{
		dep: wmi.GetWrapper(),
	}
}

type win32ComputerSystem struct {
	Manufacturer    string
	Model           string
	Name            string
	CurrentTimeZone int
}

// Info returns the Asset System information
func (w WMI) Info() (*asset.AssetSystem, error) {
	var dst []win32ComputerSystem
	err := w.dep.Query("SELECT Manufacturer,Model,Name,CurrentTimeZone FROM Win32_ComputerSystem", &dst)
	if err != nil {
		return nil, err
	}
	tz, _ := time.Now().In(time.Local).Zone()
	return mapping(dst[0], tz), nil

}

func mapping(dst win32ComputerSystem, tz string) *asset.AssetSystem {
	return &asset.AssetSystem{
		Model:               dst.Model,
		Product:             dst.Manufacturer,
		SystemName:          dst.Name,
		TimeZone:            timeZoneMinuteToHourStr(dst.CurrentTimeZone),
		TimeZoneDescription: tz,
	}
}

//This function will convert time zone minute value to string hour
//For ex if tz = 330 then function should return +0530
func timeZoneMinuteToHourStr(tz int) string {
	sign := "+"
	if tz < 0 {
		sign = "-"
		tz = int(math.Abs(float64(tz)))
	}
	q := tz / hour
	r := tz % hour
	return fmt.Sprintf("%s%02d%02d", sign, q, r)
}
