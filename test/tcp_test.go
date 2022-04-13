package tcp

import (
	"context"
	"log"
	"net"
	"net/netip"
	"sync"
	"testing"
	"time"

	client "github.com/optman/rndz-go/client/tcp"
	server "github.com/optman/rndz-go/server/tcp"
)

func TestTCPClient(t *testing.T) {
	rndzServer := "127.0.0.1:8888"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.New(rndzServer).Run(ctx)

	var wg sync.WaitGroup

	ready := make(chan any)

	wg.Add(1)
	go func() {
		defer wg.Done()
		c1 := client.New(rndzServer, "c1", netip.AddrPort{})
		defer c1.Close()
		l, err := c1.Listen(context.Background())
		if err != nil {
			panic(err)
		}

		defer l.Close()

		close(ready)

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

	<-ready

	{
		c2 := client.New(rndzServer, "c2", netip.AddrPort{})
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
