package dal

import (
	"testing"

	"reflect"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	"github.com/golang/mock/gomock"
)

func TestGetAssetDal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assetDalDepMock := mock.NewMockAssetDalDependencies(ctrl)
	returnedT := AssetDalFactoryImpl{}.GetAssetDal(assetDalDepMock)
	if expectedT, ok := returnedT.(assetDalImpl); !ok {
		t.Errorf("Unexpected type returned : %v", reflect.TypeOf(expectedT))
	}
}

//func TestSerializeObject(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()

//	assetDalDepMock := mock.NewMockAssetDalDependencies(ctrl)
//	_, err := assetDalImpl{
//		factory: assetDalDepMock,
//	}.SerializeObject(gomock.Any())
//	//TODO - proper test to be added once SerializeObject will have the code
//	if err != nil {
//		t.Errorf("Exexpected error : %v", err)
//	}
//}
