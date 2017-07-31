package mem

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/wmi/wmiMock"
	"github.com/golang/mock/gomock"
)

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
		{"TestInfo1", mockSetup(ctrl, errors.New("ErrorQuery")), errors.New("ErrorQuery")},
		{"TestInfo2", mockSetup(ctrl, nil), nil},
	}

	for _, v := range testCases {
		_, err := WMI{
			dep: v.mock,
		}.Info()

		if v.expectedErr == nil && err == nil {
			//Test Pass, further check not required
			break
		}
		if v.expectedErr == nil && err != nil {
			t.Errorf("%s : Expected : %v, but Returned : %v", v.testName, v.expectedErr, err)
			break
		}
		if v.expectedErr != nil && err == nil {
			t.Errorf("%s : Expected error: %v, but Returned : nil", v.testName, v.expectedErr)
			break
		}
		if v.expectedErr.Error() != err.Error() {
			t.Errorf("%s : Expected : %v, but Returned : %v", v.testName, v.expectedErr, err)
			break
		}
	}
}

func TestMapping(t *testing.T) {
	expectedObj := []asset.PhysicalMemory{
		asset.PhysicalMemory{
			Manufacturer: "Samsung",
			SerialNumber: "348AE941",
			SizeBytes:    8589934592,
		},
	}

	dst := []win32PhysicalMemory{
		win32PhysicalMemory{
			Manufacturer: "Samsung",
			SerialNumber: "348AE941",
			Capacity:     8589934592,
		},
	}

	actualObj := mapToMemModel(dst)

	if !reflect.DeepEqual(actualObj, expectedObj) {
		t.Errorf("Returned object is not equal to expected object")
	}
}

func TestGetByWMI(t *testing.T) {
	wmi := GetByWMI()
	empty := WMI{}
	if wmi == empty {
		t.Errorf("Returned object is not equal to expected object")
	}
}
