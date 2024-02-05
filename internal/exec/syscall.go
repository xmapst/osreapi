//go:build windows

package exec

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	MaxPath           = 260
	Th32CsSnapProcess = 0x00000002
)

type ProcessInfo struct {
	Name string
	Pid  uint32
	PPid uint32
}

type ProcessEnTry32 struct {
	DwSize              uint32
	CntUsage            uint32
	Th32ProcessID       uint32
	Th32DefaultHeapID   uintptr
	Th32ModuleID        uint32
	CntThreads          uint32
	Th32ParentProcessID uint32
	PcPriClassBase      int32
	DwFlags             uint32
	SzExeFile           [MaxPath]uint16
}

type HANDLE uintptr

var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procCreateToolHelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
)

func KillAll(pid int) error {
	pids := GetPPids(uint32(pid))
	Kill(pids)
	return nil
}

func Kill(pids []uint32) {
	for _, pid := range pids {
		pro, err := os.FindProcess(int(pid))
		if err != nil {
			continue
		}
		_ = pro.Kill()

	}
}

func GetPPids(pid uint32) []uint32 {
	infos, err := GetProCs()
	if err != nil {
		return []uint32{pid}
	}
	var pids = make([]uint32, 0, len(infos))
	var index = 0
	pids = append(pids, pid)

	var length = len(pids)
	for index < length {
		for _, info := range infos {
			if info.PPid == pids[index] {
				pids = append(pids, info.Pid)
			}
		}
		index += 1
		length = len(pids)
	}
	return pids
}

func GetProCs() (procs []ProcessInfo, err error) {
	snap := createToolHelp32Snapshot(Th32CsSnapProcess, uint32(0))
	if snap == 0 {
		err = syscall.GetLastError()
		return
	}

	defer closeHandle(snap)

	var pe32 ProcessEnTry32

	pe32.DwSize = uint32(unsafe.Sizeof(pe32))
	if process32First(snap, &pe32) == false {
		err = syscall.GetLastError()
		return
	}
	procs = append(procs, ProcessInfo{syscall.UTF16ToString(pe32.SzExeFile[:260]), pe32.Th32ProcessID, pe32.Th32ParentProcessID})
	for process32Next(snap, &pe32) {
		procs = append(procs, ProcessInfo{syscall.UTF16ToString(pe32.SzExeFile[:260]), pe32.Th32ProcessID, pe32.Th32ParentProcessID})
	}
	return
}

func createToolHelp32Snapshot(flags, processId uint32) HANDLE {
	ret, _, _ := procCreateToolHelp32Snapshot.Call(
		uintptr(flags),
		uintptr(processId))

	if ret <= 0 {
		return HANDLE(0)
	}
	return HANDLE(ret)
}

func process32First(snapshot HANDLE, pe *ProcessEnTry32) bool {
	ret, _, _ := procProcess32First.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(pe)))

	return ret != 0
}

func process32Next(snapshot HANDLE, pe *ProcessEnTry32) bool {
	ret, _, _ := procProcess32Next.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(pe)))

	return ret != 0
}

func closeHandle(object HANDLE) bool {
	ret, _, _ := procCloseHandle.Call(
		uintptr(object))
	return ret != 0
}
