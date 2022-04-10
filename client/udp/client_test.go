package udp

import (
	"context"
	"net"
	"net/netip"
	"sync"
	"testing"
	"time"
)

func TestUdpClient(t *testing.T) {
	rndzServer := "127.0.0.1:8888"

	var wg sync.WaitGroup

	ready := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		c1 := NewClient(rndzServer, "c1", netip.AddrPort{})
		defer c1.Close()
		l, err := c1.Listen(context.Background())
		if err != nil {
			t.Error(err)
		}
		defer l.Close()

		close(ready)

		buf := make([]byte, 5)
		n, _, err := l.ReadFrom(buf)
		if err != nil {
			t.Error(err)
			return
		}

		if string(buf[:n]) != "hello" {
			t.Errorf("invalid data")
			return
		}

	}()

	<-ready

	{
		c2 := NewClient(rndzServer, "c2", netip.AddrPort{})
		defer c2.Close()

		var err error
		var conn *net.UDPConn
		for {
			conn, err = c2.Connect(context.Background(), "c1")
			if err != nil {
				t.Error(err)

				<-time.After(1 * time.Second)
				continue
			}

			defer conn.Close()

			break
		}

		if _, err = conn.Write([]byte("hello")); err != nil {
			t.Error(err)
		}

	}

	wg.Wait()

}
