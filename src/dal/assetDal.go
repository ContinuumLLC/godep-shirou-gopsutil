package dal

import "github.com/ContinuumLLC/platform-asset-plugin/src/model"

//AssetDalFactoryImpl is an implementation of AssetDalFactory interface
type AssetDalFactoryImpl struct {
}

//GetAssetDal is a method of AssetDalFactory interface and returns the AssetDal interface
func (AssetDalFactoryImpl) GetAssetDal(d model.AssetDalDependencies) model.AssetDal {
	return assetDalImpl{
		factory: d,
	}
}

//assetDalImpl is an implementation of interface AssetDal
type assetDalImpl struct {
	factory model.AssetDalDependencies
}

//SerializeObject serializes the object and returns the byte[]
func (c assetDalImpl) SerializeObject(v interface{}) ([]byte, error) {
	//TODO - code to be added
	return nil, nil
}
