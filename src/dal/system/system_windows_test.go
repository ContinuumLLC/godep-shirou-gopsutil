package system

import "testing"

func TestTimeZoneMinuteToHourStr(t *testing.T) {
	testCases := []struct {
		testName    string
		tzInMinute  int
		expectedVal string
	}{
		{"Test1", 330, "+0530"},
		{"Test2", -330, "-0530"},
		{"Test3", 1001, "+1641"},
	}
	for _, v := range testCases {
		rVal := timeZoneMinuteToHourStr(v.tzInMinute)
		if rVal != v.expectedVal {
			t.Errorf("%s : Expected value is %s but returned %s", v.testName, v.expectedVal, rVal)
			break
		}
	}
}
