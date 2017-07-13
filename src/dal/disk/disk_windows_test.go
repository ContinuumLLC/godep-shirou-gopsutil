// +build windows

package disk

import (
	"reflect"
	"testing"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
)

func TestMapToPartitionDataModel(t *testing.T) {
	expectedObj := []asset.AssetDrivePartition{
		{
			Name:        "C:",
			Label:       "Windows",
			FileSystem:  "NTFS",
			Description: "Local Fixed Disk",
			SizeBytes:   247962005504,
		},
		{
			Name:        "D:",
			Label:       "HP_TOOLS",
			FileSystem:  "FAT32",
			Description: "Local Fixed Disk",
			SizeBytes:   2139095040,
		},
	}

	wLogicalDisk := []win32LogicalDisk{
		{
			Name:        "C:",
			VolumeName:  "Windows",
			FileSystem:  "NTFS",
			Description: "Local Fixed Disk",
			Size:        247962005504,
		},
		{
			Name:        "D:",
			VolumeName:  "HP_TOOLS",
			FileSystem:  "FAT32",
			Description: "Local Fixed Disk",
			Size:        2139095040,
		},
	}

	actualObj := mapToPartitionDataModel(wLogicalDisk)

	if !reflect.DeepEqual(actualObj, expectedObj) {
		t.Errorf("Actual object is not equal to expected object")
	}
}

func TestMapToDriveModel(t *testing.T) {
	expectedObj := asset.AssetDrive{
		Product:            "TOSHIB  MQ01ACF050 SCSI Disk Device",
		Manufacturer:       "(Standard disk drives)",
		MediaType:          "Fixed hard disk media",
		LogicalName:        "PHYSICALDRIVE0",
		SerialNumber:       "Y6QZCSJ3T",
		SizeBytes:          500105249280,
		NumberOfPartitions: 1,
		PartitionData: []asset.AssetDrivePartition{
			{
				Name:        "C:",
				Label:       "Windows",
				FileSystem:  "NTFS",
				Description: "Local Fixed Disk",
				SizeBytes:   247962005504,
			},
		},
	}

	srNum := "       Y6QZCSJ3T"
	wPhysicalDisk := win32DiskDrive{
		Name:         "\\\\.\\PHYSICALDRIVE0",
		Caption:      "TOSHIB  MQ01ACF050 SCSI Disk Device",
		Manufacturer: "(Standard disk drives)",
		MediaType:    "Fixed hard disk media",
		SerialNumber: &srNum,
		Index:        0,
		Partitions:   1,
		Size:         500105249280,
	}

	logicalDisks := []asset.AssetDrivePartition{
		{
			Name:        "C:",
			Label:       "Windows",
			FileSystem:  "NTFS",
			Description: "Local Fixed Disk",
			SizeBytes:   247962005504,
		},
	}

	actualObj := mapToDriveModel(&wPhysicalDisk, logicalDisks)
	if actualObj == nil {
		t.Errorf("Could not get physical disk info.")
	}
	if !reflect.DeepEqual(*actualObj, expectedObj) {
		t.Errorf("Actual object is not equal to expected object")
	}
}
