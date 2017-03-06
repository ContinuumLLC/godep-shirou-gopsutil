package linux

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetOS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	//mockAssetDalD := mock.NewMockAssetCollectionDalDependencies

	// log := logging.GetLoggerFactory().New("")
	// log.SetLogLevel(logging.OFF)
	// _, err := AssetDalImpl{
	// 	Factory: mockAssetDalD,
	// 	Logger:  log,
	// }.GetOS()
}
