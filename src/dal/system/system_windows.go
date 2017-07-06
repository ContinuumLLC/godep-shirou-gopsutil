package system

import (
	"fmt"
	"time"

	"math"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/StackExchange/wmi"
)

// WMI struct is to represent that System information is to be collected using WMI calls
type WMI struct{}

const (
	hour int = 60
)

type win32ComputerSystem struct {
	Manufacturer    string
	Model           string
	Name            string
	CurrentTimeZone int
}

// Info returns the Asset System information
func (WMI) Info() (*asset.AssetSystem, error) {
	var dst []win32ComputerSystem
	err := wmi.Query("SELECT Manufacturer,Model,Name,CurrentTimeZone FROM Win32_ComputerSystem", &dst)
	if err != nil {
		return nil, err
	}
	tz, _ := time.Now().In(time.Local).Zone()
	return &asset.AssetSystem{
		Model:               dst[0].Model,
		Product:             dst[0].Manufacturer,
		SystemName:          dst[0].Name,
		TimeZone:            timeZoneMinuteToHourStr(dst[0].CurrentTimeZone),
		TimeZoneDescription: tz,
	}, nil
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
