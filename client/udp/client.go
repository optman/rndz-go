package udp

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/optman/rndz-go/ctl"
	pb "github.com/optman/rndz-go/proto"

	"google.golang.org/protobuf/proto"
)

type Client struct {
	rndzServer string
	id         string
	localAddr  netip.AddrPort
	listener   net.PacketConn
	svrConn    *net.UDPConn
	closeOnce  sync.Once
}

func NewClient(rndzServer, id string, localAddr netip.AddrPort) *Client {
	return &Client{
		rndzServer: rndzServer,
		id:         id,
		localAddr:  localAddr,
	}
}

func (c *Client) Connect(ctx context.Context, targetId string) (udpConn *net.UDPConn, err error) {

	svrConn, err := dial(c.localAddr, c.rndzServer)
	if err != nil {
		return nil, err
	}

	connect := func() (remoteAddr string, err error) {

		isync := &pb.Request_Isync{
			Isync: &pb.Isync{
				Id: targetId,
			},
		}

		if err = c.writeReq(svrConn, isync); err != nil {
			return
		}

		svrConn.SetReadDeadline(time.Now().Add(10 * time.Second))
		var resp pb.Response
		resp, err = readResp(svrConn)
		if err != nil {
			return
		}
		if resp.GetId() != c.id {
			err = fmt.Errorf("wrong response?")
			return
		}

		rdr := resp.GetRedirect()
		if rdr == nil {
			err = fmt.Errorf("invalid response")
			return

		}

		remoteAddr = rdr.GetAddr()
		if remoteAddr == "" {
			err = fmt.Errorf("target id not found")
			return
		}

		log.Printf("remote addr %s", remoteAddr)

		return
	}

	for i := 0; i < 3; i++ {
		var remoteAddr string
		remoteAddr, err = connect()
		if err == nil {
			return dial(netip.MustParseAddrPort(svrConn.LocalAddr().String()), remoteAddr)
		}
	}

	return
}

func (c *Client) Listen(ctx context.Context) (net.PacketConn, error) {

	localAddr := c.localAddr
	if !localAddr.IsValid() {
		serverAddr, err := net.ResolveUDPAddr("udp", c.rndzServer)
		if err != nil {
			return nil, err
		}

		if len(serverAddr.IP) == net.IPv4len {
			localAddr = netip.MustParseAddrPort("0.0.0.0:0")
		} else {
			localAddr = netip.MustParseAddrPort("[::]:0")
		}
	}

	svrConn, err := dial(localAddr, c.rndzServer)
	if err != nil {
		return nil, err
	}

	cancel := func() {
		c.Close()
	}

	go func() {

		defer cancel()

		ping := &pb.Request_Ping{
			Ping: &pb.Ping{},
		}

		c.writeReq(svrConn, ping)

		dur := 10 * time.Second
		t := time.NewTimer(dur)
		defer t.Stop()

		for {

			select {
			case <-ctx.Done():
				return
			case <-t.C:
				c.writeReq(svrConn, ping)
			}

			t.Reset(dur)
		}
	}()

	go func() {
		defer cancel()

		for {
			resp, err := readResp(svrConn)
			if err != nil {
				return
			}

			switch resp.Cmd.(type) {
			case *pb.Response_Pong:
			case *pb.Response_Fsync:
				fsync, _ := resp.Cmd.(*pb.Response_Fsync)

				dstId := fsync.Fsync.GetId()
				dstAddr := fsync.Fsync.GetAddr()
				log.Printf("fsync %s %s", dstId, dstAddr)

				rsync := &pb.Request_Rsync{
					Rsync: &pb.Rsync{
						Id: dstId,
					},
				}

				if conn, err := dial(netip.MustParseAddrPort(svrConn.LocalAddr().String()), dstAddr); err == nil {
					c.writeReq(conn, rsync)
					conn.Close()
				} else {
					log.Println(err)
				}

				c.writeReq(svrConn, rsync)

			case *pb.Response_Redirect:
				rdr, _ := resp.Cmd.(*pb.Response_Redirect)
				log.Printf("Redirect?  %s %s", rdr.Redirect.GetId(), rdr.Redirect.GetAddr())

			default:
				log.Println("ignore unknown response cmd")
			}

		}
	}()

	cfg := net.ListenConfig{
		Control: ctl.ControlUDP,
	}

	listener, err := cfg.ListenPacket(ctx, "udp", svrConn.LocalAddr().String())
	if err != nil {
		return nil, err
	}

	c.svrConn = svrConn
	c.listener = listener

	return listener, nil
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		if c.svrConn != nil {
			bye := &pb.Request_Bye{
				Bye: &pb.Bye{},
			}
			c.writeReq(c.svrConn, bye)
			c.svrConn.Close()
		}

		if c.listener != nil {
			c.listener.Close()
		}
	})
}

func dial(srcAddr netip.AddrPort, dstAddr string) (udpConn *net.UDPConn, err error) {

	d := net.Dialer{
		Control:   ctl.ControlUDP,
		LocalAddr: net.UDPAddrFromAddrPort(srcAddr),
	}

	var conn net.Conn
	conn, err = d.Dial("udp", dstAddr)
	if err != nil {
		return
	}

	return conn.(*net.UDPConn), nil
}

func (c *Client) writeReq(conn *net.UDPConn, cmd pb.IsRequestCmd) error {
	req := &pb.Request{
		Id:  c.id,
		Cmd: cmd,
	}
	buf, _ := proto.Marshal(req)

	_, err := conn.Write(buf)
	return err
}

func readResp(conn *net.UDPConn) (resp pb.Response, err error) {

	buf := make([]byte, 1500)
	var n int
	n, err = conn.Read(buf)
	if err != nil {
		return
	}

	if err = proto.Unmarshal(buf[:n], &resp); err != nil {
		return
	}

	return
}
