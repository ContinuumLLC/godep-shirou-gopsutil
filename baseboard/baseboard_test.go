package baseboard

import (
	"testing"
)

func TestBaseboardInfo(t *testing.T) {
	v, err := Info()
	if err != nil {
		t.Errorf("error %v", err)
	}
	empty := &InfoStat{}
	if v == empty {
		t.Errorf("Could not get baseboard info %v", v)
	}
}
