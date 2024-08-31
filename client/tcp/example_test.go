package tcp_test

import (
	"context"
	"net/netip"
)

func Example_client1(rndzServer String) {
	c := tcp.NewClient(rndzServer, "c1", netip.AddrPort{})
	defer c.Close()
	l, _ := c.Listen(context.Background())
	defer l.Close()
	for {
		conn, _ := l.Accept()
		defer conn.Close()
		//    ...
	}
}

func Example_client2(rndzServer String) {
	c := tcp.NewClient(rndzServer, "c2", netip.AddrPort{})
	defer c.Close()
	conn, _ := c.Connect(context.Background(), "c1")
	defer conn.Close()
}
