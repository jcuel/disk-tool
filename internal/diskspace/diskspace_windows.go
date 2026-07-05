//go:build windows

package diskspace

import (
	"syscall"
	"unsafe"
)

var (
	modKernel32         = syscall.NewLazyDLL("kernel32.dll")
	procGetDiskFreeExW  = modKernel32.NewProc("GetDiskFreeSpaceExW")
)

func platformUsage(path string) (total, free uint64, err error) {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, 0, err
	}
	var avail, tot, totFree uint64
	r, _, e := procGetDiskFreeExW.Call(
		uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(&avail)),
		uintptr(unsafe.Pointer(&tot)),
		uintptr(unsafe.Pointer(&totFree)),
	)
	if r == 0 {
		if e != syscall.Errno(0) {
			return 0, 0, e
		}
		return 0, 0, syscall.EINVAL
	}
	return tot, avail, nil
}
