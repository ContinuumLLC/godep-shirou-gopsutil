// +build darwin

package baseboard

import "github.com/shirou/gopsutil/internal/common"

// Info returns baseboard information
func Info() (*InfoStat, error) {
	return nil, common.ErrNotImplementedError
}
