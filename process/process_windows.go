// +build windows

package process

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/StackExchange/wmi"
	cpu "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/internal/common"
	net "github.com/shirou/gopsutil/net"
	"github.com/shirou/w32"
	"golang.org/x/sys/windows"
)

type LUID struct {
	LowPart  uint32
	HighPart int32
}

type LUID_AND_ATTRIBUTES struct {
	Luid       LUID
	Attributes uint32
}

type TOKEN_PRIVILEGES struct {
	PrivilegeCount uint32
	Privileges     [1]LUID_AND_ATTRIBUTES
}

const (
	NoMoreFiles                  = 0x12
	MaxPathLength                = 260
	errnoERROR_IO_PENDING        = 997
	sePrivilegeEnabled    uint32 = 0x00000002
)

var (
	modpsapi                        = windows.NewLazyDLL("psapi.dll")
	procGetProcessMemoryInfo        = modpsapi.NewProc("GetProcessMemoryInfo")
	modadvapi32                     = syscall.NewLazyDLL("advapi32.dll")
	procLookupPrivilegeValueW       = modadvapi32.NewProc("LookupPrivilegeValueW")
	procAdjustTokenPrivileges       = modadvapi32.NewProc("AdjustTokenPrivileges")
	errERRORIOPENDING         error = syscall.Errno(errnoERROR_IO_PENDING)
)

type SystemProcessInformation struct {
	NextEntryOffset   uint64
	NumberOfThreads   uint64
	Reserved1         [48]byte
	Reserved2         [3]byte
	UniqueProcessID   uintptr
	Reserved3         uintptr
	HandleCount       uint64
	Reserved4         [4]byte
	Reserved5         [11]byte
	PeakPagefileUsage uint64
	PrivatePageCount  uint64
	Reserved6         [6]uint64
}

// Memory_info_ex is different between OSes
type MemoryInfoExStat struct {
}

type MemoryMapsStat struct {
}

type Win32_Process struct {
	Name                string
	ExecutablePath      *string
	CommandLine         *string
	Priority            uint32
	CreationDate        *time.Time
	ProcessID           uint32
	ThreadCount         uint32
	Status              *string
	ReadOperationCount  uint64
	ReadTransferCount   uint64
	WriteOperationCount uint64
	WriteTransferCount  uint64

	/*
		CSCreationClassName   string
		CSName                string
		Caption               *string
		CreationClassName     string
		Description           *string
		ExecutionState        *uint16
		HandleCount           uint32
		KernelModeTime        uint64
		MaximumWorkingSetSize *uint32
		MinimumWorkingSetSize *uint32
		OSCreationClassName   string
		OSName                string
		OtherOperationCount   uint64
		OtherTransferCount    uint64
		PageFaults            uint32
		PageFileUsage         uint32
		ParentProcessID       uint32
		PeakPageFileUsage     uint32
		PeakVirtualSize       uint64
		PeakWorkingSetSize    uint32
		PrivatePageCount      uint64
		TerminationDate       *time.Time
		UserModeTime          uint64
		WorkingSetSize        uint64
	*/
}

// Win32_PerfFormattedData_PerfProc_Process struct to provide performance process metrics for windows
type Win32_PerfFormattedData_PerfProc_Process struct {
	IDProcess            uint32
	Name                 string
	HandleCount          uint32
	PercentProcessorTime uint64
	PrivateBytes         uint64
	ThreadCount          uint32
	VirtualBytes         uint64
	WorkingSet           uint64
	WorkingSetPrivate    uint64
}

func Pids() ([]int32, error) {
	var ret []int32

	procs, err := processes()
	if err != nil {
		return ret, nil
	}

	for _, proc := range procs {
		ret = append(ret, proc.Pid)
	}

	return ret, nil
}

func (p *Process) Ppid() (int32, error) {
	ret, _, _, err := p.getFromSnapProcess(p.Pid)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func GetWin32Proc(pid int32) ([]Win32_Process, error) {
	var dst []Win32_Process
	query := fmt.Sprintf("WHERE ProcessId = %d", pid)
	q := wmi.CreateQuery(&dst, query)
	err := wmi.Query(q, &dst)
	if err != nil {
		return []Win32_Process{}, fmt.Errorf("could not get win32Proc: %s", err)
	}
	if len(dst) != 1 {
		return []Win32_Process{}, fmt.Errorf("could not get win32Proc: empty")
	}
	return dst, nil
}

// PerfProcessStats returns the performance data from performance counters of process object.
func PerfProcessStats() ([]Win32_PerfFormattedData_PerfProc_Process, error) {
	var ret []Win32_PerfFormattedData_PerfProc_Process
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}

func (p *Process) Name() (string, error) {
	dst, err := GetWin32Proc(p.Pid)
	if err != nil {
		return "", fmt.Errorf("could not get Name: %s", err)
	}
	return dst[0].Name, nil
}
func (p *Process) Exe() (string, error) {
	dst, err := GetWin32Proc(p.Pid)
	if err != nil {
		return "", fmt.Errorf("could not get ExecutablePath: %s", err)
	}
	return *dst[0].ExecutablePath, nil
}
func (p *Process) Cmdline() (string, error) {
	dst, err := GetWin32Proc(p.Pid)
	if err != nil {
		return "", fmt.Errorf("could not get CommandLine: %s", err)
	}
	return *dst[0].CommandLine, nil
}

// CmdlineSlice returns the command line arguments of the process as a slice with each
// element being an argument. This merely returns the CommandLine informations passed
// to the process split on the 0x20 ASCII character.
func (p *Process) CmdlineSlice() ([]string, error) {
	cmdline, err := p.Cmdline()
	if err != nil {
		return nil, err
	}
	return strings.Split(cmdline, " "), nil
}

func (p *Process) CreateTime() (int64, error) {
	ru, err := getRusage(p.Pid)
	if err != nil {
		return 0, fmt.Errorf("could not get CreationDate: %s", err)
	}

	return ru.CreationTime.Nanoseconds() / 1000000, nil
}

func (p *Process) Cwd() (string, error) {
	return "", common.ErrNotImplementedError
}
func (p *Process) Parent() (*Process, error) {
	return p, common.ErrNotImplementedError
}
func (p *Process) Status() (string, error) {
	return "", common.ErrNotImplementedError
}

func (p *Process) EnableCurrentProcessPrivilege(string strPrivilegeName) error {
	hCurrHandle, err := win.GetCurrentProcess()

	var tCurr win.Token
	err = windows.OpenProcessToken(hCurrHandle, windows.TOKEN_ADJUST_PRIVILEGES, &tCurr)
	if nil != err {
		return err
	}

	var tokPrev TOKEN_PRIVILEGES
	tokPrev.PrivilegeCount = 1

	uiSeDebugName, err := win.UTF16FromString(strPrivilegeName)

	_ = lookupPrivilegeValue(nil, &uiSeDebugName[0], &tokPrev.Privileges[0].Luid)

	tokPrev.PrivilegeCount = 1
	tokPrev.Privileges[0].Attributes = sePrivilegeEnabled

	cb := unsafe.Sizeof(tokPrev)

	_, err := adjustTokenPrivileges(tCurr, false, &tokPrev, uint32(cb), nil, nil)

	return err
}
func (p *Process) Username() (string, error) {
	hProcess, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, uint32(p.Pid))
	if nil != err {
		return "", err
	}
	defer windows.CloseHandle(hProcess)

	var t windows.Token
	err = windows.OpenProcessToken(hProcess, windows.TOKEN_QUERY, &t)
	if nil != err {
		return "", err
	}
	defer windows.CloseHandle(windows.Handle(t))

	var n uint32
	err = windows.GetTokenInformation(t, windows.TokenUser, nil, 0, &n)
	if n <= 0 {
		return "", err
	}

	b := make([]byte, n)

	err = windows.GetTokenInformation(t, windows.TokenUser, &b[0], uint32(len(b)), &n)
	if nil != err {
		return "", err
	}

	tokenUser, err := t.GetTokenUser()

	if nil != err {
		return "", err
	}

	name := uint32(250)
	domain := uint32(250)
	var accType uint32
	strname := make([]uint16, name)
	strdomain := make([]uint16, domain)

	err = windows.LookupAccountSid(nil, tokenUser.User.Sid, &strname[0], &name, &strdomain[0], &domain, &accType)
	if nil != err {
		return "", err
	}
	return windows.UTF16ToString(strname), nil
}
func (p *Process) Uids() ([]int32, error) {
	var uids []int32

	return uids, common.ErrNotImplementedError
}
func (p *Process) Gids() ([]int32, error) {
	var gids []int32
	return gids, common.ErrNotImplementedError
}
func (p *Process) Terminal() (string, error) {
	return "", common.ErrNotImplementedError
}

// Nice returnes priority in Windows
func (p *Process) Nice() (int32, error) {
	dst, err := GetWin32Proc(p.Pid)
	if err != nil {
		return 0, fmt.Errorf("could not get Priority: %s", err)
	}
	return int32(dst[0].Priority), nil
}
func (p *Process) IOnice() (int32, error) {
	return 0, common.ErrNotImplementedError
}
func (p *Process) Rlimit() ([]RlimitStat, error) {
	var rlimit []RlimitStat

	return rlimit, common.ErrNotImplementedError
}

func (p *Process) IOCounters() (*IOCountersStat, error) {
	dst, err := GetWin32Proc(p.Pid)
	if err != nil || len(dst) == 0 {
		return nil, fmt.Errorf("could not get Win32Proc: %s", err)
	}
	ret := &IOCountersStat{
		ReadCount:  uint64(dst[0].ReadOperationCount),
		ReadBytes:  uint64(dst[0].ReadTransferCount),
		WriteCount: uint64(dst[0].WriteOperationCount),
		WriteBytes: uint64(dst[0].WriteTransferCount),
	}

	return ret, nil
}
func (p *Process) NumCtxSwitches() (*NumCtxSwitchesStat, error) {
	return nil, common.ErrNotImplementedError
}
func (p *Process) NumFDs() (int32, error) {
	return 0, common.ErrNotImplementedError
}
func (p *Process) NumThreads() (int32, error) {
	dst, err := GetWin32Proc(p.Pid)
	if err != nil {
		return 0, fmt.Errorf("could not get ThreadCount: %s", err)
	}
	return int32(dst[0].ThreadCount), nil
}
func (p *Process) Threads() (map[string]string, error) {
	ret := make(map[string]string, 0)
	return ret, common.ErrNotImplementedError
}
func (p *Process) Times() (*cpu.TimesStat, error) {
	return nil, common.ErrNotImplementedError
}
func (p *Process) CPUAffinity() ([]int32, error) {
	return nil, common.ErrNotImplementedError
}
func (p *Process) MemoryInfo() (*MemoryInfoStat, error) {
	mem, err := getMemoryInfo(p.Pid)
	if err != nil {
		return nil, err
	}

	ret := &MemoryInfoStat{
		RSS: uint64(mem.WorkingSetSize),
		VMS: uint64(mem.PagefileUsage),
	}

	return ret, nil
}
func (p *Process) MemoryInfoEx() (*MemoryInfoExStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) Children() ([]*Process, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) OpenFiles() ([]OpenFilesStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) Connections() ([]net.ConnectionStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) NetIOCounters(pernic bool) ([]net.IOCountersStat, error) {
	return nil, common.ErrNotImplementedError
}

func (p *Process) IsRunning() (bool, error) {
	return true, common.ErrNotImplementedError
}

func (p *Process) MemoryMaps(grouped bool) (*[]MemoryMapsStat, error) {
	var ret []MemoryMapsStat
	return &ret, common.ErrNotImplementedError
}

func NewProcess(pid int32) (*Process, error) {
	p := &Process{Pid: pid}

	return p, nil
}

func (p *Process) SendSignal(sig windows.Signal) error {
	return common.ErrNotImplementedError
}

func (p *Process) Suspend() error {
	return common.ErrNotImplementedError
}
func (p *Process) Resume() error {
	return common.ErrNotImplementedError
}

func (p *Process) Terminate() error {
	// PROCESS_TERMINATE = 0x0001
	proc := w32.OpenProcess(0x0001, false, uint32(p.Pid))
	ret := w32.TerminateProcess(proc, 0)
	w32.CloseHandle(proc)

	if ret == false {
		return windows.GetLastError()
	} else {
		return nil
	}
}

func (p *Process) Kill() error {
	return common.ErrNotImplementedError
}

func (p *Process) getFromSnapProcess(pid int32) (int32, int32, string, error) {
	snap := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPPROCESS, uint32(pid))
	if snap == 0 {
		return 0, 0, "", windows.GetLastError()
	}
	defer w32.CloseHandle(snap)
	var pe32 w32.PROCESSENTRY32
	pe32.DwSize = uint32(unsafe.Sizeof(pe32))
	if w32.Process32First(snap, &pe32) == false {
		return 0, 0, "", windows.GetLastError()
	}

	if pe32.Th32ProcessID == uint32(pid) {
		szexe := windows.UTF16ToString(pe32.SzExeFile[:])
		return int32(pe32.Th32ParentProcessID), int32(pe32.CntThreads), szexe, nil
	}

	for w32.Process32Next(snap, &pe32) {
		if pe32.Th32ProcessID == uint32(pid) {
			szexe := windows.UTF16ToString(pe32.SzExeFile[:])
			return int32(pe32.Th32ParentProcessID), int32(pe32.CntThreads), szexe, nil
		}
	}
	return 0, 0, "", errors.New("Couldn't find pid:" + string(pid))
}

// Get processes
func processes() ([]*Process, error) {

	var dst []Win32_Process
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		return []*Process{}, err
	}
	if len(dst) == 0 {
		return []*Process{}, fmt.Errorf("could not get Process")
	}
	results := make([]*Process, 0, len(dst))
	for _, proc := range dst {
		p, err := NewProcess(int32(proc.ProcessID))
		if err != nil {
			continue
		}
		results = append(results, p)
	}

	return results, nil
}

func getProcInfo(pid int32) (*SystemProcessInformation, error) {
	initialBufferSize := uint64(0x4000)
	bufferSize := initialBufferSize
	buffer := make([]byte, bufferSize)

	var sysProcInfo SystemProcessInformation
	ret, _, _ := common.ProcNtQuerySystemInformation.Call(
		uintptr(unsafe.Pointer(&sysProcInfo)),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&bufferSize)),
		uintptr(unsafe.Pointer(&bufferSize)))
	if ret != 0 {
		return nil, windows.GetLastError()
	}

	return &sysProcInfo, nil
}

func getRusage(pid int32) (*windows.Rusage, error) {
	var CPU windows.Rusage

	c, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(c)

	if err := windows.GetProcessTimes(c, &CPU.CreationTime, &CPU.ExitTime, &CPU.KernelTime, &CPU.UserTime); err != nil {
		return nil, err
	}

	return &CPU, nil
}

func getMemoryInfo(pid int32) (PROCESS_MEMORY_COUNTERS, error) {
	var mem PROCESS_MEMORY_COUNTERS
	c, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return mem, err
	}
	defer windows.CloseHandle(c)
	if err := getProcessMemoryInfo(c, &mem); err != nil {
		return mem, err
	}

	return mem, err
}

func getProcessMemoryInfo(h windows.Handle, mem *PROCESS_MEMORY_COUNTERS) (err error) {
	r1, _, e1 := syscall.Syscall(procGetProcessMemoryInfo.Addr(), 3, uintptr(h), uintptr(unsafe.Pointer(mem)), uintptr(unsafe.Sizeof(*mem)))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func errnoErr(e syscall.Errno) error {

	switch e {

	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERRORIOPENDING
	}
	return e
}

func lookupPrivilegeValue(systemname *uint16, name *uint16, luid *LUID) (err error) {
	r1, _, e1 := syscall.Syscall(procLookupPrivilegeValueW.Addr(), 3, uintptr(unsafe.Pointer(systemname)), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(luid)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func adjustTokenPrivileges(token win.Token, disableAllPrivileges bool, newstate *TOKEN_PRIVILEGES, buflen uint32, prevstate *TOKEN_PRIVILEGES, returnlen *uint32) (ret uint32, err error) {
	var _p0 uint32
	if disableAllPrivileges {
		_p0 = 1
	} else {
		_p0 = 0
	}

	r0, _, e1 := syscall.Syscall6(procAdjustTokenPrivileges.Addr(), 6, uintptr(token), uintptr(_p0), uintptr(unsafe.Pointer(newstate)), uintptr(buflen), uintptr(unsafe.Pointer(prevstate)), uintptr(unsafe.Pointer(returnlen)))
	ret = uint32(r0)

	if true {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}

	return
}
