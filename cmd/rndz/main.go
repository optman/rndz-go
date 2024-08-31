package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/netip"
	"os"

	"github.com/optman/rndz-go/client/tcp"
	"github.com/optman/rndz-go/client/udp"
	TcpServer "github.com/optman/rndz-go/server/tcp"
	UdpServer "github.com/optman/rndz-go/server/udp"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Usage: "a simple rendezvous protocol implementation to help NAT traversal or hole punching",
		Name:  "rndz-go",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "verbose", Aliases: []string{"debug"}},
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name: "client",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "tcp"},
					&cli.StringFlag{
						Name:     "server-addr",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "id",
						Required: true,
					},
					&cli.StringFlag{
						Name: "remote-peer",
					},
				},
				Action: runClient,
			},

			&cli.Command{
				Name: "server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "listen-addr",
						Value: ":8888",
					},
				},
				Action: runServer,
			},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("verbose") {
				log.SetOutput(os.Stdout)
			} else {
				log.SetOutput(io.Discard)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServer(c *cli.Context) error {
	listenAddr := c.String("listen-addr")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go TcpServer.New(listenAddr).Run(ctx)

	return UdpServer.New(listenAddr).Run(ctx)
}

func runClient(c *cli.Context) error {
	rndzServer := c.String("server-addr")
	remotePeer := c.String("remote-peer")
	id := c.String("id")

	if c.Bool("tcp") {
		c := tcp.New(rndzServer, id, netip.AddrPort{})
		defer c.Close()
		if remotePeer == "" {
			return tcpServer(c)
		} else {
			return tcpClient(c, remotePeer)
		}

	} else {
		c := udp.New(rndzServer, id, netip.AddrPort{})
		defer c.Close()
		if remotePeer == "" {
			return udpServer(c)
		} else {
			return udpClient(c, remotePeer)
		}
	}
}

func tcpServer(c *tcp.Client) error {
	l, err := c.Listen(context.Background())
	if err != nil {
		return fmt.Errorf("error starting TCP server: %w", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting TCP connection: %w", err)
		}

		go func() {
			defer conn.Close()

			buf := make([]byte, 5)
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("error reading from TCP connection:", err)
				return
			}

			log.Printf("read %d bytes", n)

			if _, err := conn.Write(buf[:n]); err != nil {
				log.Println("error writing to TCP connection:", err)
				return
			}
		}()
	}
}

func tcpClient(c *tcp.Client, remotePeer string) error {
	conn, err := c.Connect(context.Background(), remotePeer)
	if err != nil {
		return fmt.Errorf("error connecting to TCP remote peer: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("hello")); err != nil {
		return fmt.Errorf("error writing to remote TCP peer: %w", err)
	}

	buf := make([]byte, 5)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("error reading from remote TCP peer: %w", err)
	}

	log.Printf("read %d bytes from remote peer", n)

	return nil
}

func udpServer(c *udp.Client) error {
	l, err := c.Listen(context.Background())
	if err != nil {
		return fmt.Errorf("error starting UDP server: %w", err)
	}
	defer l.Close()

	buf := make([]byte, 1500)
	for {
		n, addr, err := l.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("error reading from UDP connection: %w", err)
		}

		log.Printf("read %d bytes from %s", n, addr)

		if _, err := l.WriteTo(buf[:n], addr); err != nil {
			log.Println("error writing to UDP connection:", err)
		}
	}
}

func udpClient(c *udp.Client, remotePeer string) error {
	conn, err := c.Connect(context.Background(), remotePeer)
	if err != nil {
		return fmt.Errorf("error connecting to UDP remote peer: %w", err)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte("hello")); err != nil {
		return fmt.Errorf("error writing to remote UDP peer: %w", err)
	}

	buf := make([]byte, 1500)
	if n, err := conn.Read(buf); err == nil {
		log.Printf("read %d bytes from %s", n, conn.RemoteAddr())
	} else {
		return fmt.Errorf("error reading from remote UDP peer: %w", err)
	}

	return nil
}
