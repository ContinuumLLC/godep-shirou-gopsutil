package baseboard

import (
	"encoding/json"
	"time"

	"github.com/shirou/gopsutil/internal/common"
)

var invoke common.Invoker

func init() {
	invoke = common.Invoke{}
}

// InfoStat is struct definition of baseboard attributes
type InfoStat struct {
	Product      string    `json:"product"`
	Manufacturer string    `json:"manufacturer"`
	Model        string    `json:"model"`
	SerialNumber string    `json:"serialNumber"`
	Version      string    `json:"version"`
	Name         string    `json:"name"`
	InstallDate  time.Time `json:"installDate"`
}

// String returns json string of baseboard inInfoStat struct
func (b InfoStat) String() string {
	s, _ := json.Marshal(b)
	return string(s)
}
