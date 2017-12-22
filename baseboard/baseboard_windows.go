// +build windows

package baseboard

import (
	"time"

	"github.com/StackExchange/wmi"
)

// Win32_BaseBoard struct represents a baseboard
type Win32_BaseBoard struct {
	Name         string
	Product      string
	Manufacturer string
	SerialNumber string
	Version      *string
	Model        *string
	InstallDate  *time.Time
}

// Info returns baseboard information
func Info() (*InfoStat, error) {
	var dst []Win32_BaseBoard
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		return nil, err
	}
	var model, version string
	var installDate time.Time

	if dst[0].Model != nil {
		model = *(dst[0].Model)
	}
	if dst[0].Version != nil {
		version = *(dst[0].Version)
	}
	if dst[0].InstallDate != nil {
		installDate = *(dst[0].InstallDate)
	}
	return &InfoStat{
		Product:      dst[0].Product,
		Manufacturer: dst[0].Manufacturer,
		SerialNumber: dst[0].SerialNumber,
		Version:      version,
		Name:         dst[0].Name,
		Model:        model,
		InstallDate:  installDate,
	}, nil
}
