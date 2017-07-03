// +build windows

package baseboard

import (
	"github.com/StackExchange/wmi"
)

// Win32_BaseBoard struct represents a baseboard
type Win32_BaseBoard struct {
	Name         string
	Product      string
	Manufacturer string
	SerialNumber string
	Version      string
	//Model        string
	//InstallDate  time.Time
}

// Info returns baseboard information
func Info() (*InfoStat, error) {
	var dst []Win32_BaseBoard
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		return nil, err
	}
	return &InfoStat{
		Product:      dst[0].Product,
		Manufacturer: dst[0].Manufacturer,
		SerialNumber: dst[0].SerialNumber,
		Version:      dst[0].Version,
		Name:         dst[0].Name,
		//Model:        dst[0].Manufacturer,
		//InstallDate:  dst[0].InstallDate,
	}, nil
}
