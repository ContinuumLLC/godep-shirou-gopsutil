package system

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/wmi/wmiMock"
	"github.com/golang/mock/gomock"
)

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

func mockSetup(ctrl *gomock.Controller, err error) *wmiMock.MockWrapper {
	mockObj := wmiMock.NewMockWrapper(ctrl)
	mockObj.EXPECT().Query(gomock.Any(), gomock.Any()).Return(err)
	return mockObj
}

func TestInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		testName    string
		mock        *wmiMock.MockWrapper
		expectedErr error
	}{
		{"Test1", mockSetup(ctrl, errors.New("ErrorQuery")), errors.New("ErrorQuery")},
	}

	for _, v := range testCases {
		_, err := WMI{
			dep: v.mock,
		}.Info()
		if v.expectedErr == nil && err != nil {
			t.Errorf("%s : Expected : %v, but Returned : %v", v.testName, v.expectedErr, err)
			break
		}
		if v.expectedErr != nil && err == nil {
			t.Errorf("%s : Expected : %v, but Returned : %v", v.testName, v.expectedErr, err)
			break
		}
		if v.expectedErr.Error() != err.Error() {
			t.Errorf("%s : Expected : %v, but Returned : %v", v.testName, v.expectedErr, err)
			break
		}
	}
}

func TestMapping(t *testing.T) {
	expectedObj := &asset.AssetSystem{
		Model: "Model", Product: "Manufacturer", SystemName: "Name", TimeZone: "+0530", TimeZoneDescription: "IST",
	}

	returnedObj := mapping(win32ComputerSystem{
		CurrentTimeZone: 330, Manufacturer: "Manufacturer", Model: "Model", Name: "Name",
	}, "IST")

	if !reflect.DeepEqual(returnedObj, expectedObj) {
		t.Errorf("Returned object is not equal to expected object")
	}
}
