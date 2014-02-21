// +build !linux

package system

import "syscall"

func Lgetxattr(path string, attr string) ([]byte, error) {
	return nil, ErrNotImplemented
}

func Lsetxattr(path string, attr string, data []byte, flags int) error {
	return ErrNotImplemented
}
