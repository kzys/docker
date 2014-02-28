// +build !linux

package system

import "syscall"

// to avoid "imported and not used" error
var _ = syscall.Open

func LUtimesNano(path string, ts []syscall.Timespec) error {
	return ErrNotImplemented
}

func UtimesNano(path string, ts []syscall.Timespec) error {
	return ErrNotImplemented
}
