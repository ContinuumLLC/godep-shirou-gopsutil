package dal

import (
	//"time"

	amodel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
)

// AssetCollection Proc related constants
const (
	cMemCreatedBy                           string = "/continuum/agent/plugin/asset"
	cMemDataType                            string = "assetCollection"
	cMemProcPath                            string = "/proc/meminfo"
	cMemProcPhysicalTotalBytes              string = "MemTotal"
	cMemProcPhysicalAvailableBytes          string = "MemAvailable"
	cMemProcVirtualAvailableBytes           string = "SwapFree"
	cMemProcVirtualTotalBytes               string = "SwapTotal"
	cMemProcCommittedBytes                  string = "CommitLimit"
	cMemProcFreeSystemPageTableEntriesBytes string = "PageTables"
)

//Error Codes
const (
	 // INVALIDAssetCollectionMEASURE = "Invalid measure :"
)

type assetCollectionDalLinux struct {
	factory model.AssetCollectionDalDependencies
	logger  logging.Logger
}

//Gets AssetCollection data
func (dal *assetCollectionDalLinux) GetAssetCollection() (*amodel.AssetCollection, error) {
	reader, err := dal.factory.GetEnv().GetFileReader(cMemProcPath)
	if err != nil {
		dal.logger.Logf(logging.DEBUG, "Error in reading file %v", err)
		return nil, err
	}
	defer reader.Close()
	parser := dal.factory.GetParser()
	cfg := procParser.Config{
		ParserMode:    procParser.ModeKeyValue,
		IgnoreNewLine: true,
	}
	data, err := parser.Parse(cfg, reader)
	if err != nil {
		dal.logger.Logf(logging.DEBUG, "Error in parsing config %v", err)
		return nil, err
	}
	return translateAssetCollection{logger: dal.logger}.translateAssetCollectionProcToModel(data), nil
}

type translateAssetCollection struct {
	logger logging.Logger
}

func (t translateAssetCollection) translateAssetCollectionProcToModel(data *procParser.Data) *amodel.AssetCollection {
	assetCollection := new(amodel.AssetCollection)
	//assetCollection.Type = cMemDataType
	//assetCollection.CreatedBy = cMemCreatedBy
	//assetCollection.CreateTimeUTC = time.Now().UTC()
	//assetCollection.PhysicalTotalBytes = t.getDataFromMap(cMemProcPhysicalTotalBytes, data)
	//assetCollection.PhysicalAvailableBytes = t.getDataFromMap(cMemProcPhysicalAvailableBytes, data)
	//assetCollection.PhysicalInUseBytes = assetCollection.PhysicalTotalBytes - assetCollection.PhysicalAvailableBytes
	//assetCollection.VirtualAvailableBytes = t.getDataFromMap(cMemProcVirtualAvailableBytes, data)
	//assetCollection.VirtualTotalBytes = t.getDataFromMap(cMemProcVirtualTotalBytes, data)
	//assetCollection.VirtualInUseBytes = assetCollection.VirtualTotalBytes - assetCollection.VirtualAvailableBytes
	//assetCollection.CommittedBytes = t.getDataFromMap(cMemProcCommittedBytes, data)
	//assetCollection.FreeSystemPageTableEntriesBytes = t.getDataFromMap(cMemProcFreeSystemPageTableEntriesBytes, data)
	return assetCollection
}

func (t translateAssetCollection) getDataFromMap(key string, data *procParser.Data) int64 {
	if _, exists := data.Map[key]; !exists {
		return int64(0)
	}
	val, err := procParser.GetInt64(data.Map[key].Values[1])
	if err != nil {
		t.logger.Logf(logging.DEBUG, "Error in GetInt64 for key %s, Error : %v", key, err)
		return 0
	}
	return procParser.GetBytes(val, data.Map[key].Values[2])
}
