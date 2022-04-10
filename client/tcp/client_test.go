package tcp

import (
	"context"
	"log"
	"net"
	"net/netip"
	"sync"
	"testing"
	"time"
)

func TestTCPClient(t *testing.T) {
	rndzServer := "localhost:8888"

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		c1 := NewClient(rndzServer, "c1", netip.AddrPort{})
		defer c1.Close()
		l, err := c1.Listen(context.Background())
		if err != nil {
			panic(err)
		}

		defer l.Close()

		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		defer conn.Close()

		buf := make([]byte, 5)
		if _, err := conn.Read(buf); err != nil {
			panic(err)
		}

		if _, err := conn.Write(buf); err != nil {
			panic(err)
		}

	}()

	{
		c2 := NewClient(rndzServer, "c2", netip.AddrPort{})
		defer c2.Close()
		var conn net.Conn
		var err error
		for {
			conn, err = c2.Connect(context.Background(), "c1")
			if err != nil {
				log.Println(err)
				<-time.After(1 * time.Second)
			} else {
				break
			}
		}

		defer conn.Close()

		if _, err := conn.Write([]byte("hello")); err != nil {
			panic(err)
		}

		buf := make([]byte, 5)
		if _, err := conn.Read(buf); err != nil {
			panic(err)
		}
	}

	wg.Wait()
}
