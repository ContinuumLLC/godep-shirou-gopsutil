// +build windows

package disk

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/internal/common"
	"golang.org/x/sys/windows"
)

var (
	procGetDiskFreeSpaceExW     = common.Modkernel32.NewProc("GetDiskFreeSpaceExW")
	procGetLogicalDriveStringsW = common.Modkernel32.NewProc("GetLogicalDriveStringsW")
	procGetDriveType            = common.Modkernel32.NewProc("GetDriveTypeW")
	procGetVolumeInformation    = common.Modkernel32.NewProc("GetVolumeInformationW")
)

var (
	FileFileCompression = uint32(16)     // 0x00000010
	FileReadOnlyVolume  = uint32(524288) // 0x00080000
)

// diskPerformance is an equivalent representation of DISK_PERFORMANCE in the Windows API.
// https://docs.microsoft.com/fr-fr/windows/win32/api/winioctl/ns-winioctl-disk_performance
type diskPerformance struct {
	BytesRead           int64
	BytesWritten        int64
	ReadTime            int64
	WriteTime           int64
	IdleTime            int64
	ReadCount           uint32
	WriteCount          uint32
	QueueDepth          uint32
	SplitCount          uint32
	QueryTime           int64
	StorageDeviceNumber uint32
	StorageManagerName  [8]uint16
	alignmentPadding    uint32 // necessary for 32bit support, see https://github.com/elastic/beats/pull/16553
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
	WaitMSec                        = 500
	ioctlDiskGetDriveGeometry       = 458752
	ioctlStorageGetDeviceNumber     = 2953344
	ioctlVolumeGetVolumeDiskExtents = 5636096
)

type diskExtent struct {
	DiskNumber     uint32
	StartingOffset uint64
	ExtentLength   uint64
}

type volumeDiskExtents []byte

func (v *volumeDiskExtents) Len() uint {
	return uint(binary.LittleEndian.Uint32([]byte(*v)))
}

func (v *volumeDiskExtents) Extent(n uint) diskExtent {
	ba := []byte(*v)
	//This calculates the next offset in the structure
	offset := 8 + 24*n
	return diskExtent{
		DiskNumber:     binary.LittleEndian.Uint32(ba[offset:]),
		StartingOffset: binary.LittleEndian.Uint64(ba[offset+8:]),
		ExtentLength:   binary.LittleEndian.Uint64(ba[offset+16:]),
	}
}

func UsageWithContext(ctx context.Context, path string) (*UsageStat, error) {
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
	ret := &UsageStat{
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

func PartitionsWithContext(ctx context.Context, all bool) ([]PartitionStat, error) {
	var ret []PartitionStat

	bufferSize := uint32(256)
	volnameBuffer := make([]byte, bufferSize)

	// find first volume
	handle, err := windows.FindFirstVolume((*uint16)(unsafe.Pointer(&volnameBuffer[0])), bufferSize)
	if nil != err {
		return ret, windows.GetLastError()
	}

	defer func() {
		_ = windows.FindVolumeClose(handle)
	}()

	volumeInfo, err := getVolumeInfo(string(cleanVolumePath(volnameBuffer)))
	if err != nil {
		return ret, err
	}
	if volumeInfo != nil {
		ret = append(ret, *volumeInfo)
	}

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
		if volumeInfo != nil {
			ret = append(ret, *volumeInfo)
		}
	}

	return ret, nil
}

func getVolumeInfo(volumeName string) (partition *PartitionStat, err error) {
	bufferSize := uint32(256)
	volumeNameBuffer := make([]byte, bufferSize)
	mountPointBuffer := make([]byte, bufferSize)
	lpFileSystemNameBuffer := make([]byte, bufferSize)
	volumeNameSerialNumber := uint32(0)
	maximumComponentLength := uint32(0)
	lpFileSystemFlags := uint32(0)

	volpath, _ := windows.UTF16PtrFromString(volumeName)
	typeret, _, _ := procGetDriveType.Call(uintptr(unsafe.Pointer(volpath)))
	// 0: DRIVE_UNKNOWN 1: DRIVE_NO_ROOT_DIR 2: DRIVE_REMOVABLE 3: DRIVE_FIXED 4: DRIVE_REMOTE 5: DRIVE_CDROM
	if typeret == 0 {
		return nil, windows.GetLastError()
	}

	// 1: DRIVE_NO_ROOT_DIR The root path is invalid; for example, there is no volume mounted at the specified path.
	// Including type 1 because we also want to get details about Unmounted Partitions
	if typeret == 1 || typeret == 2 || typeret == 3 || typeret == 4 {
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
			if typeret == 2 {
				//device is not ready will happen if there is no disk in the drive
				// Should ignore it ?
				return nil, nil
			}
		}

		opts := "rw"
		if lpFileSystemFlags&FileReadOnlyVolume != 0 {
			opts = "ro"
		}
		if lpFileSystemFlags&FileFileCompression != 0 {
			opts += ".compress"
		}

		// Ignore the error, some volumes may not have a Mount Point (Unmounted Partitions)
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
			return nil, nil
		}
		mediaType, err := getMediaType(hFile)
		diskNumber, err := getDiskNumber(hFile)
		windows.CloseHandle(hFile)

		return &PartitionStat{
			Mountpoint: string(cleanVolumePath(mountPointBuffer)),
			Device:     volumeName,
			Fstype:     string(bytes.Replace(lpFileSystemNameBuffer, []byte("\x00"), []byte(""), -1)),
			Opts:       opts,
			DriveType:  uint32(typeret),
			VolumeName: string(bytes.Replace(volumeNameBuffer, []byte("\x00"), []byte(""), -1)),
			MediaType:  mediaType,
			DiskNumber: diskNumber,
		}, nil
	}

	return nil, nil
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
		size := uint32(16 * 1024)
		vols := make(volumeDiskExtents, size)
		var bytesReturned uint32
		err = windows.DeviceIoControl(hFile, ioctlVolumeGetVolumeDiskExtents, nil, 0, &vols[0], size, &bytesReturned, nil)
		if err != nil {
			return
		}
		if vols.Len() > 0 {
			return vols.Extent(0).DiskNumber, nil
		}
		return 0, errors.New("Unable to get disk number")
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

func IOCountersWithContext(ctx context.Context, names ...string) (map[string]IOCountersStat, error) {
	// https://github.com/giampaolo/psutil/blob/544e9daa4f66a9f80d7bf6c7886d693ee42f0a13/psutil/arch/windows/disk.c#L83
	drivemap := make(map[string]IOCountersStat, 0)
	var diskPerformance diskPerformance

	lpBuffer := make([]uint16, 254)
	lpBufferLen, err := windows.GetLogicalDriveStrings(uint32(len(lpBuffer)), &lpBuffer[0])
	if err != nil {
		return drivemap, err
	}
	for _, v := range lpBuffer[:lpBufferLen] {
		if 'A' <= v && v <= 'Z' {
			path := string(rune(v)) + ":"
			typepath, _ := windows.UTF16PtrFromString(path)
			typeret := windows.GetDriveType(typepath)
			if typeret == 0 {
				return drivemap, windows.GetLastError()
			}
			if typeret != windows.DRIVE_FIXED {
				continue
			}
			szDevice := fmt.Sprintf(`\\.\%s`, path)
			const IOCTL_DISK_PERFORMANCE = 0x70020
			h, err := windows.CreateFile(syscall.StringToUTF16Ptr(szDevice), 0, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, nil, windows.OPEN_EXISTING, 0, 0)
			if err != nil {
				if err == windows.ERROR_FILE_NOT_FOUND {
					continue
				}
				return drivemap, err
			}
			defer windows.CloseHandle(h)

			var diskPerformanceSize uint32
			err = windows.DeviceIoControl(h, IOCTL_DISK_PERFORMANCE, nil, 0, (*byte)(unsafe.Pointer(&diskPerformance)), uint32(unsafe.Sizeof(diskPerformance)), &diskPerformanceSize, nil)
			if err != nil {
				return drivemap, err
			}
			drivemap[path] = IOCountersStat{
				ReadBytes:  uint64(diskPerformance.BytesRead),
				WriteBytes: uint64(diskPerformance.BytesWritten),
				ReadCount:  uint64(diskPerformance.ReadCount),
				WriteCount: uint64(diskPerformance.WriteCount),
				ReadTime:   uint64(diskPerformance.ReadTime / 10000 / 1000), // convert to ms: https://github.com/giampaolo/psutil/issues/1012
				WriteTime:  uint64(diskPerformance.WriteTime / 10000 / 1000),
				Name:       path,
			}
		}
	}
	return drivemap, nil
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
