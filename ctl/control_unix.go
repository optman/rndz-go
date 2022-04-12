//go:build !windows

package ctl

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func ControlTCP(network, address string, c syscall.RawConn) error {
	var err error
	c.Control(func(fd uintptr) {
		if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
			return
		}
		if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
			return
		}
	})

	return err
}

func ControlUDP(network, address string, c syscall.RawConn) error {
	var err error
	c.Control(func(fd uintptr) {
		if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
			return
		}
	})

	return err
}
