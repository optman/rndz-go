package udp

import (
	"context"
	"net"

	"github.com/optman/rndz-go/ctl"
)

//raw bind (unconnected udp)
//
//Connect() return a connnected udp, some 3rd package want a unconnected udp.
func Bind(localAddr string) (net.PacketConn, error) {
	cfg := net.ListenConfig{
		Control: ctl.ControlUDP,
	}

	return cfg.ListenPacket(context.Background(), "udp", localAddr)
}
