package udp

import (
	"context"
	"net"

	"github.com/optman/rndz-go/ctl"
)

func Bind(localAddr string) (net.PacketConn, error) {
	cfg := net.ListenConfig{
		Control: ctl.ControlUDP,
	}

	return cfg.ListenPacket(context.Background(), "udp", localAddr)
}
