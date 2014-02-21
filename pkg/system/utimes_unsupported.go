// +build !linux

package system

import "syscall"

func LUtimesNano(path string, ts []syscall.Timespec) error {
	return ErrNotImplemented
}

func UtimesNano(path string, ts []syscall.Timespec) error {
	return ErrNotImplemented
}
