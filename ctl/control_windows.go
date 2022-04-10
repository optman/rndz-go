package ctl

import (
	"golang.org/x/sys/windows"
	"syscall"
)

func ControlUDP(network, address string, c syscall.RawConn) (err error) {
	c.Control(func(fd uintptr) {
		err = windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
	})

	return
}

var ControlTCP = ControlUDP
