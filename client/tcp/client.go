package tcp

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
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
	svrConn    net.Conn
	listener   net.Listener
	closeOnce  sync.Once
}

func New(rndzServer, id string, localAddr netip.AddrPort) *Client {
	return &Client{
		rndzServer: rndzServer,
		id:         id,
		localAddr:  localAddr,
	}
}

func (c *Client) Connect(ctx context.Context, targetId string) (net.Conn, error) {
	svrConn, err := dial(ctx, net.TCPAddrFromAddrPort(c.localAddr), c.rndzServer, time.Duration(0))
	if err != nil {
		return nil, err
	}
	defer svrConn.Close()

	isync := &pb.Request_Isync{
		Isync: &pb.Isync{
			Id: targetId,
		},
	}

	if err := c.writeReq(svrConn, isync); err != nil {
		return nil, err
	}

	resp, err := readResp(svrConn)
	if err != nil {
		return nil, err
	}

	rdr := resp.GetRedirect()
	if rdr == nil {
		return nil, fmt.Errorf("invalid response")
	}

	remoteAddr := rdr.GetAddr()
	if remoteAddr == "" {
		return nil, fmt.Errorf("target not found")
	}

	log.Println("remote addr ", remoteAddr)

	return dial(ctx, svrConn.LocalAddr(), remoteAddr, time.Duration(0))
}

func (c *Client) Listen(ctx context.Context) (net.Listener, error) {

	localAddr := c.localAddr
	if !localAddr.IsValid() {
		serverAddr, err := net.ResolveTCPAddr("tcp", c.rndzServer)
		if err != nil {
			return nil, err
		}

		if len(serverAddr.IP) == net.IPv4len {
			localAddr = netip.MustParseAddrPort("0.0.0.0:0")
		} else {
			localAddr = netip.MustParseAddrPort("[::]:0")
		}
	}

	cfg := net.ListenConfig{
		Control: ctl.ControlTCP,
	}

	listener, err := cfg.Listen(ctx, "tcp", localAddr.String())
	if err != nil {
		return nil, err
	}

	svrConn, err := dial(ctx, listener.Addr(), c.rndzServer, time.Duration(0))
	if err != nil {
		return nil, err
	}

	cmdChannel := make(chan pb.IsRequestCmd, 10)

	sendPing := func() {
		ping := &pb.Request_Ping{
			Ping: &pb.Ping{},
		}
		cmdChannel <- ping
	}

	sendPing()

	cancel := func() {
		c.Close()
	}

	go func() {
		defer cancel()

		dur := 10 * time.Second

		t := time.NewTimer(dur)
		defer t.Stop()

		for {

			select {
			case <-ctx.Done():
				return
			case <-t.C:
				sendPing()
			case cmd, ok := <-cmdChannel:
				if !ok {
					return
				}

				if err := c.writeReq(svrConn, cmd); err != nil {
					return
				}
			}

			t.Reset(dur)
		}
	}()

	go func() {
		defer close(cmdChannel)
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

				dial(context.Background(), svrConn.LocalAddr(), dstAddr, time.Millisecond)

				rsync := &pb.Request_Rsync{
					Rsync: &pb.Rsync{
						Id: dstId,
					},
				}

				cmdChannel <- rsync

			case *pb.Response_Redirect:
				rdr, _ := resp.Cmd.(*pb.Response_Redirect)
				log.Printf("Redirect?  %s %s", rdr.Redirect.GetId(), rdr.Redirect.GetAddr())

			default:
				log.Println("ignore unknown response cmd")
			}
		}

	}()

	c.svrConn = svrConn
	c.listener = listener

	return listener, nil
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		if c.svrConn != nil {
			c.svrConn.Close()
		}

		if c.listener != nil {
			c.listener.Close()
		}
	})
}

func dial(ctx context.Context, localAddr net.Addr, remoteAddr string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control:   ctl.ControlTCP,
		LocalAddr: localAddr,
		Timeout:   timeout,
	}

	return d.DialContext(ctx, "tcp", remoteAddr)
}

func (c *Client) writeReq(w io.Writer, cmd pb.IsRequestCmd) (err error) {
	req := &pb.Request{
		Id:  c.id,
		Cmd: cmd,
	}
	buf, _ := proto.Marshal(req)
	if err = binary.Write(w, binary.BigEndian, uint16(len(buf))); err != nil {
		return
	}

	_, err = w.Write(buf)
	return
}

func readResp(r io.Reader) (resp pb.Response, err error) {
	var n uint16
	if err = binary.Read(r, binary.BigEndian, &n); err != nil {
		return
	}

	buf := make([]byte, n)
	if _, err = r.Read(buf); err != nil {
		return
	}

	if err = proto.Unmarshal(buf, &resp); err != nil {
		return
	}

	return
}
