// +build windows

package disk

import (
	"bytes"
	"errors"
	"strings"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/internal/common"
	"golang.org/x/sys/windows"
)

var (
	procGetDiskFreeSpaceExW = common.Modkernel32.NewProc("GetDiskFreeSpaceExW")
	procGetDriveType        = common.Modkernel32.NewProc("GetDriveTypeW")
)

var (
	FileFileCompression = uint32(16)     // 0x00000010
	FileReadOnlyVolume  = uint32(524288) // 0x00080000
)

type Win32_PerfFormattedData struct {
	Name                    string
	AvgDiskBytesPerRead     uint64
	AvgDiskBytesPerWrite    uint64
	AvgDiskReadQueueLength  uint64
	AvgDiskWriteQueueLength uint64
	AvgDisksecPerRead       uint64
	AvgDisksecPerWrite      uint64
}

type Win32_PerfFormattedData_PerfDisk_PhysicalDisk struct {
	AvgDiskBytesPerRead     uint64
	AvgDiskBytesPerTransfer uint64
	AvgDiskBytesPerWrite    uint64
	AvgDiskQueueLength      uint64
	AvgDiskReadQueueLength  uint64
	AvgDiskSecPerRead       uint32
	AvgDiskSecPerTransfer   uint32
	AvgDiskSecPerWrite      uint32
	AvgDiskWriteQueueLength uint64
	CurrentDiskQueueLength  uint32
	DiskBytesPerSec         uint64
	DiskReadBytesPerSec     uint64
	DiskReadsPerSec         uint32
	DiskTransfersPerSec     uint32
	DiskWriteBytesPerSec    uint64
	DiskWritesPerSec        uint32
	Name                    string
	PercentDiskReadTime     uint64
	PercentDiskTime         uint64
	PercentDiskWriteTime    uint64
	PercentIdleTime         uint64
	SplitIOPerSec           uint32
}

type Win32_PerfFormattedData_PerfDisk_LogicalDisk struct {
	AvgDiskBytesPerRead     uint64
	AvgDiskBytesPerTransfer uint64
	AvgDiskBytesPerWrite    uint64
	AvgDiskQueueLength      uint64
	AvgDiskReadQueueLength  uint64
	AvgDiskSecPerRead       uint32
	AvgDiskSecPerTransfer   uint32
	AvgDiskSecPerWrite      uint32
	AvgDiskWriteQueueLength uint64
	CurrentDiskQueueLength  uint32
	DiskBytesPerSec         uint64
	DiskReadBytesPerSec     uint64
	DiskReadsPerSec         uint32
	DiskTransfersPerSec     uint32
	DiskWriteBytesPerSec    uint64
	DiskWritesPerSec        uint32
	Name                    string
	PercentDiskReadTime     uint64
	PercentDiskTime         uint64
	PercentDiskWriteTime    uint64
	PercentIdleTime         uint64
	SplitIOPerSec           uint32
	FreeMegabytes           uint32
	PercentFreeSpace        uint32
}

type Win32_LogicalDisk struct {
	Name      string
	Size      *uint64 // null if no disk in CD/DVD drive
	FreeSpace *uint64 // null if no disk in CD/DVD drive
	DriveType uint32
}

type Win32_DiskDrive struct {
	Name       string
	Model      string
	Index      uint32
	Partitions uint32
	Size       *uint64
}

type diskGeometry struct {
	Cylinders         int64
	MediaType         int32
	TracksPerCylinder uint32
	SectorsPerTrack   uint32
	BytesPerSector    uint32
}

type diskDeviceNumber struct {
	DeviceType      uint32
	DeviceNumber    uint32
	PartitionNumber uint32
}

const (
	WaitMSec                    = 500
	ioctlDiskGetDriveGeometry   = 458752
	ioctlStorageGetDeviceNumber = 2953344
)

func Usage(path string) (*UsageStat, error) {
	ret := &UsageStat{}

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	diskret, _, err := procGetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	if diskret == 0 {
		return nil, err
	}
	ret = &UsageStat{
		Path:        path,
		Total:       uint64(lpTotalNumberOfBytes),
		Free:        uint64(lpTotalNumberOfFreeBytes),
		Used:        uint64(lpTotalNumberOfBytes) - uint64(lpTotalNumberOfFreeBytes),
		UsedPercent: (float64(lpTotalNumberOfBytes) - float64(lpTotalNumberOfFreeBytes)) / float64(lpTotalNumberOfBytes) * 100,
		// InodesTotal: 0,
		// InodesFree: 0,
		// InodesUsed: 0,
		// InodesUsedPercent: 0,
	}
	return ret, nil
}

func Partitions(all bool) ([]PartitionStat, error) {
	var ret []PartitionStat

	bufferSize := uint32(256)
	volnameBuffer := make([]byte, bufferSize)

	// find first volume
	handle, err := windows.FindFirstVolume((*uint16)(unsafe.Pointer(&volnameBuffer[0])), bufferSize)
	if nil != err {
		return ret, windows.GetLastError()
	}

	volumeInfo, err := getVolumeInfo(string(cleanVolumePath(volnameBuffer)))
	if err != nil {
		return ret, err
	}
	ret = append(ret, *volumeInfo)

	// loop over all volumes, excluding the first volume
	for {
		// If no more partiotions, returns error, exit the loop
		err = windows.FindNextVolume(handle, (*uint16)(unsafe.Pointer(&volnameBuffer[0])), bufferSize)
		if err != nil {
			break
		}

		volumeInfo, err := getVolumeInfo(string(cleanVolumePath(volnameBuffer)))
		if err != nil {
			return ret, err
		}
		ret = append(ret, *volumeInfo)
	}

	// Close the handle
	_ = windows.FindVolumeClose(handle)

	return ret, nil
}

func getVolumeInfo(volumeName string) (partition *PartitionStat, err error) {
	d := PartitionStat{}

	bufferSize := uint32(256)
	volumeNameBuffer := make([]byte, bufferSize)
	mountPointBuffer := make([]byte, bufferSize)
	lpFileSystemNameBuffer := make([]byte, bufferSize)
	volumeNameSerialNumber := uint32(0)
	maximumComponentLength := uint32(0)
	lpFileSystemFlags := uint32(0)

	volpath, _ := windows.UTF16PtrFromString(volumeName)
	typeret, _, _ := procGetDriveType.Call(uintptr(unsafe.Pointer(volpath)))
	if typeret == 0 {
		return nil, windows.GetLastError()
	}
	// 2: DRIVE_REMOVABLE 3: DRIVE_FIXED 4: DRIVE_REMOTE 5: DRIVE_CDROM
	if typeret == 2 || typeret == 3 || typeret == 4 || typeret == 5 {
		err = windows.GetVolumeInformation(
			(*uint16)(unsafe.Pointer(volpath)),
			(*uint16)(unsafe.Pointer(&volumeNameBuffer[0])),
			bufferSize,
			&volumeNameSerialNumber,
			&maximumComponentLength,
			&lpFileSystemFlags,
			(*uint16)(unsafe.Pointer(&lpFileSystemNameBuffer[0])),
			bufferSize)

		if err != nil {
			if typeret == 5 || typeret == 2 {
				//device is not ready will happen if there is no disk in the drive
				return nil, errors.New("Device is not ready!")
			}
			return nil, windows.GetLastError()
		}

		opts := "rw"
		if lpFileSystemFlags&FileReadOnlyVolume != 0 {
			opts = "ro"
		}
		if lpFileSystemFlags&FileFileCompression != 0 {
			opts += ".compress"
		}

		_ = windows.GetVolumePathNamesForVolumeName(
			(*uint16)(unsafe.Pointer(volpath)),
			(*uint16)(unsafe.Pointer(&mountPointBuffer[0])),
			bufferSize,
			&bufferSize)

		drivePath := strings.Replace(volumeName, "?", ".", 1)
		drivePath = strings.TrimSuffix(drivePath, `\`)

		hFile, err := windows.CreateFile(syscall.StringToUTF16Ptr(drivePath),
			0,
			windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
			nil,
			windows.OPEN_EXISTING,
			0,
			0)

		if nil != err {
			return nil, err
		}
		mediaType, err := getMediaType(hFile)
		diskNumber, err := getDiskNumber(hFile)
		windows.CloseHandle(hFile)

		d.Mountpoint = string(cleanVolumePath(mountPointBuffer))
		d.Device = volumeName
		d.Fstype = string(bytes.Replace(lpFileSystemNameBuffer, []byte("\x00"), []byte(""), -1))
		d.Opts = opts
		d.DriveType = uint32(typeret)
		d.VolumeName = string(bytes.Replace(volumeNameBuffer, []byte("\x00"), []byte(""), -1))
		d.MediaType = mediaType
		d.DiskNumber = diskNumber
	}

	return &d, nil
}

func getMediaType(hFile windows.Handle) (mediaType uint32, err error) {

	var diskGeo diskGeometry
	var dsikGeoSize uint32
	err = windows.DeviceIoControl(hFile,
		ioctlDiskGetDriveGeometry,
		nil,
		0,
		(*byte)(unsafe.Pointer(&diskGeo)),
		uint32(unsafe.Sizeof(diskGeo)),
		&dsikGeoSize,
		nil)
	if nil != err {
		return
	}

	return uint32(diskGeo.MediaType), nil
}

func getDiskNumber(hFile windows.Handle) (diskNumber uint32, err error) {

	var diskNum diskDeviceNumber
	var dsikNumSize uint32
	err = windows.DeviceIoControl(hFile,
		ioctlStorageGetDeviceNumber,
		nil,
		0,
		(*byte)(unsafe.Pointer(&diskNum)),
		uint32(unsafe.Sizeof(diskNum)),
		&dsikNumSize,
		nil)
	if nil != err {
		return
	}

	return uint32(diskNum.DeviceNumber), nil
}

func cleanVolumePath(data []byte) []byte {
	var res []byte
	for _, key := range data {
		if key != 0 {
			res = append(res, key)
		}
	}
	return res
}

func IOCounters(names ...string) (map[string]IOCountersStat, error) {
	ret := make(map[string]IOCountersStat, 0)
	var dst []Win32_PerfFormattedData

	err := wmi.Query("SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk ", &dst)
	if err != nil {
		return ret, err
	}
	for _, d := range dst {

		if len(d.Name) > 3 { // not get _Total or Harddrive
			continue
		}
		if len(names) > 0 && !common.StringsHas(names, d.Name) {
			continue
		}

		ret[d.Name] = IOCountersStat{
			Name:       d.Name,
			ReadCount:  uint64(d.AvgDiskReadQueueLength),
			WriteCount: d.AvgDiskWriteQueueLength,
			ReadBytes:  uint64(d.AvgDiskBytesPerRead),
			WriteBytes: uint64(d.AvgDiskBytesPerWrite),
			ReadTime:   d.AvgDisksecPerRead,
			WriteTime:  d.AvgDisksecPerWrite,
		}
	}
	return ret, nil
}

func PhysicalDisksStats() ([]Win32_PerfFormattedData_PerfDisk_PhysicalDisk, error) {
	var ret []Win32_PerfFormattedData_PerfDisk_PhysicalDisk
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}

func LogicalPartitionsStats() ([]Win32_PerfFormattedData_PerfDisk_LogicalDisk, error) {
	var ret []Win32_PerfFormattedData_PerfDisk_LogicalDisk
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}

func LogicalDiskSize() ([]Win32_LogicalDisk, error) {
	var ret []Win32_LogicalDisk
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}

func Info() ([]Win32_DiskDrive, error) {
	var ret []Win32_DiskDrive
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}
