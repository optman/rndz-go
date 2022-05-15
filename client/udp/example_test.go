package udp_test

import (
	"context"
	"net/netip"
)

func Example_client1() {
	c := udp.NewClient(rndzServer, "c1", netip.AddrPort{})
	defer c.Close()
	l, _ := c.Listen(context.Background())
	defer l.Close()
	buf := make([]byte, 5)
	for {
		l.ReadFrom(buf)
		//    ...
	}
}

func Example_client2() {
	c := udp.NewClient(rndzServer, "c2", netip.AddrPort{})
	defer c.Close()
	conn, _ := c.Connect(context.Background(), "c1")
	conn.Write([]byte("hello"))
	defer conn.Close()
}
