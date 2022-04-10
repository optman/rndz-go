package main

import (
	"context"
	"fmt"
	"log"
	"net/netip"
	"os"

	"github.com/optman/rndz-go/client/tcp"
	"github.com/optman/rndz-go/client/udp"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Usage: "a simple rendezvous protocol implementation to help NAT traversal or hole punching",
		Name:  "rndz-go",
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runClient(c *cli.Context) error {
	rndzServer := c.String("server-addr")
	remotePeer := c.String("remote-peer")
	id := c.String("id")

	if c.Bool("tcp") {
		c := tcp.NewClient(rndzServer, id, netip.AddrPort{})
		defer c.Close()
		if remotePeer == "" {
			return tcpServer(c)
		} else {
			return tcpClient(c, remotePeer)
		}

	} else {
		c := udp.NewClient(rndzServer, id, netip.AddrPort{})
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
		return err
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go func() {
			defer conn.Close()
			buf := make([]byte, 5)
			n, err := conn.Read(buf)
			if err != nil {
				log.Println(err)
				return
			}

			log.Printf("read %d bytes", n)

			if _, err := conn.Write(buf[:n]); err != nil {
				log.Println(err)
				return
			}

		}()

	}
}

func tcpClient(c *tcp.Client, remotePeer string) error {

	conn, err := c.Connect(context.Background(), remotePeer)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("hello")); err != nil {
		return err
	}

	buf := make([]byte, 5)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	log.Printf("read %d bytes", n)

	return nil
}

func udpServer(c *udp.Client) error {
	l, err := c.Listen(context.Background())
	if err != nil {
		return err
	}
	defer l.Close()

	buf := make([]byte, 1500)
	for {
		n, addr, err := l.ReadFrom(buf)
		if err != nil {
			return err
		}

		log.Printf("read %d bytes from %s", n, addr)

		l.WriteTo(buf[:n], addr)
	}
}

func udpClient(c *udp.Client, remotePeer string) error {

	conn, err := c.Connect(context.Background(), remotePeer)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := make([]byte, 1500)
	if _, err = conn.Write([]byte("hello")); err != nil {
		return err
	}

	if n, err := conn.Read(buf); err == nil {
		log.Printf("read %d bytes from %s", n, conn.RemoteAddr())
	} else {
		return err
	}

	return err
}
