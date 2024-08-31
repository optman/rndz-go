//Tcp connection builder
//
//Bind to a port, connect to a rendezvous server. Wait for peer connection request or initiate a
//peer connection request.
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

// Tcp connection builder

type Client struct {
	rndzServer string
	id         string
	localAddr  netip.AddrPort
	listener   net.Listener
	closeOnce  sync.Once
	closeChan  chan any
}

// New creates a new Client instance and sets the rendezvous server, peer identity, and local bind address.
// If no local address is set, it chooses according to the server address type (IPv4 or IPv6).
func New(rndzServer, id string, localAddr netip.AddrPort) *Client {
	return &Client{
		rndzServer: rndzServer,
		id:         id,
		localAddr:  localAddr,
		closeChan:  make(chan any),
	}
}

// Connect to the rendezvous server and request a connection to the target node.
// Returns a net.Conn with the remote peer.
// The connection with the rendezvous server will be dropped after return.
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

// Listen puts the socket in listen mode, creates a connection with the rendezvous server, and waits for peer
// connection requests. If the connection with the server is broken, it will auto-reconnect.
// When receiving a `Fsync` request from the server, it attempts to connect to the remote peer with a very short timeout,
// opening the firewall and NAT rule for the immediate peer connection. When the peer connection finally arrives,
// the listening socket accepts it as normal.
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

	{
		localAddr := listener.Addr()
		stopChan, err := c.connectServer(ctx, localAddr)
		if err != nil {
			return nil, err
		}

		go func() {
			var wait <-chan time.Time
			for {
				select {
				case <-ctx.Done():
					return
				case <-c.closeChan:
					return
				case <-stopChan:
					log.Printf("connection with rndz server is broken, try to reconnect.")
				case <-wait:
				}

				if stopChan, err = c.connectServer(ctx, localAddr); err != nil {
					log.Printf("connect rndz server fail, retry later. %s", err)
					wait = time.After(120 * time.Second)
				} else {
					log.Println("connect rndz server success")
					wait = nil
				}
			}
		}()
	}

	c.listener = listener

	return listener, nil
}

// Close stops the internal goroutine and closes the listener if it exists.
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.closeChan)

		if c.listener != nil {
			c.listener.Close()
		}
	})
}

// connectServer establishes a connection to the rendezvous server and returns a channel that signals if the connection is lost.
func (c *Client) connectServer(ctx context.Context, addr net.Addr) (<-chan any, error) {
	stopChan := make(chan any)

	svrConn, err := dial(ctx, addr, c.rndzServer, time.Duration(0))
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

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer svrConn.Close()

		dur := 10 * time.Second

		t := time.NewTimer(dur)
		defer t.Stop()

		for {
			select {
			case <-c.closeChan:
				return
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(cmdChannel)

		for {
			resp, err := readResp(svrConn)
			if err != nil {
				svrConn.Close()
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

	go func() {
		wg.Wait()
		close(stopChan)
	}()

	return stopChan, nil
}

// dial creates a TCP connection to the remote address using the specified local address and timeout.
func dial(ctx context.Context, localAddr net.Addr, remoteAddr string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control:   ctl.ControlTCP,
		LocalAddr: localAddr,
		Timeout:   timeout,
	}

	return d.DialContext(ctx, "tcp", remoteAddr)
}

// writeReq encodes and writes a request to the given io.Writer.
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

// readResp reads and decodes a response from the given io.Reader.
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