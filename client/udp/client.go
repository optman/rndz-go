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

const (
    // connection timeout for dialing to the server
    connectionTimeout = 10 * time.Second
    // reconciliation retry limit
    retryLimit = 3
    // ping interval duration
    pingInterval = 10 * time.Second
)

type Client struct {
    rndzServer string
    id         string
    localAddr  netip.AddrPort
    listener   net.PacketConn
    svrConn    *net.UDPConn
    closeOnce  sync.Once
}

// New initializes a new Client with the provided rendezvous server, ID, and local address.
// If no local address is set, it chooses one based on the server address type (IPv4 or IPv6).
func New(rndzServer, id string, localAddr netip.AddrPort) *Client {
    return &Client{
        rndzServer: rndzServer,
        id:         id,
        localAddr:  localAddr,
    }
}

// Connect sends a connection request to the rendezvous server to connect to the target peer.
// It returns a connected UDP connection.
func (c *Client) Connect(ctx context.Context, targetId string) (*net.UDPConn, error) {
    svrConn, err := dial(c.localAddr, c.rndzServer)
    if err != nil {
        return nil, fmt.Errorf("failed to dial server: %w", err)
    }
    defer func() {
        bye := &pb.Request_Bye{Bye: &pb.Bye{}}
        c.writeReq(svrConn, bye)
        svrConn.Close()
    }()

    for i := 0; i < retryLimit; i++ {
        remoteAddr, err := c.connectToTarget(svrConn, targetId)
        if err == nil {
            return dial(svrConn.LocalAddr().(*net.UDPAddr).AddrPort(), remoteAddr)
        }
        log.Printf("connect attempt %d failed: %v", i+1, err)
    }
    return nil, fmt.Errorf("failed to connect after %d attempts", retryLimit)
}

func (c *Client) connectToTarget(svrConn *net.UDPConn, targetId string) (string, error) {
    isync := &pb.Request_Isync{Isync: &pb.Isync{Id: targetId}}

    if err := c.writeReq(svrConn, isync); err != nil {
        return "", fmt.Errorf("failed to write isync request: %w", err)
    }

    svrConn.SetReadDeadline(time.Now().Add(connectionTimeout))
    resp, err := readResp(svrConn)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %w", err)
    }
    if resp.GetId() != c.id {
        return "", fmt.Errorf("unexpected response ID: %s", resp.GetId())
    }

    rdr := resp.GetRedirect()
    if rdr == nil || rdr.GetAddr() == "" {
        return "", fmt.Errorf("redirect address not found for target ID: %s", targetId)
    }
    log.Printf("remote addr %s", rdr.GetAddr())
    return rdr.GetAddr(), nil
}

// Listen starts listening for peer connection requests while pinging the rendezvous server.
func (c *Client) Listen(ctx context.Context) (net.PacketConn, error) {
    if !c.localAddr.IsValid() {
        serverAddr, err := net.ResolveUDPAddr("udp", c.rndzServer)
        if err != nil {
            return nil, fmt.Errorf("failed to resolve server address: %w", err)
        }
        if len(serverAddr.IP) == net.IPv4len {
            c.localAddr = netip.MustParseAddrPort("0.0.0.0:0")
        } else {
            c.localAddr = netip.MustParseAddrPort("[::]:0")
        }
    }

    svrConn, err := dial(c.localAddr, c.rndzServer)
    if err != nil {
        return nil, fmt.Errorf("failed to dial server: %w", err)
    }
    c.svrConn = svrConn

    go c.pingServer(ctx)
    go c.handleServerResponses()

    cfg := net.ListenConfig{Control: ctl.ControlUDP}
    listener, err := cfg.ListenPacket(ctx, "udp", svrConn.LocalAddr().String())
    if err != nil {
        return nil, fmt.Errorf("failed to start listener: %w", err)
    }
    c.listener = listener
    return listener, nil
}

func (c *Client) pingServer(ctx context.Context) {
    defer c.Close()
    ping := &pb.Request_Ping{Ping: &pb.Ping{}}

    for {
        select {
        case <-ctx.Done():
            return
        case <-time.After(pingInterval):
            if err := c.writeReq(c.svrConn, ping); err != nil {
                log.Printf("failed to send ping: %v", err)
            }
        }
    }
}

func (c *Client) handleServerResponses() {
    defer c.Close()

    for {
        resp, err := readResp(c.svrConn)
        if err != nil {
            log.Printf("error reading server response: %v", err)
            return
        }

        switch cmd := resp.Cmd.(type) {
        case *pb.Response_Pong:
            // Handle pong response (no action required)
        case *pb.Response_Fsync:
            c.handleFsync(cmd.Fsync)
        case *pb.Response_Redirect:
            log.Printf("Redirect? %s %s", cmd.Redirect.GetId(), cmd.Redirect.GetAddr())
        default:
            log.Printf("ignored unknown response command: %T", cmd)
        }
    }
}

func (c *Client) handleFsync(fsync *pb.Fsync) {
    dstId := fsync.GetId()
    dstAddr := fsync.GetAddr()
    log.Printf("received fsync %s %s", dstId, dstAddr)

    rsync := &pb.Request_Rsync{Rsync: &pb.Rsync{Id: dstId}}
    if conn, err := dial(c.svrConn.LocalAddr().(*net.UDPAddr).AddrPort(), dstAddr); err == nil {
        c.writeReq(conn, rsync)
        conn.Close()
    } else {
        log.Printf("failed to dial fsync address %s: %v", dstAddr, err)
    }
    c.writeReq(c.svrConn, rsync)
}

func dial(srcAddr netip.AddrPort, dstAddr string) (*net.UDPConn, error) {
    d := net.Dialer{
        Control:   ctl.ControlUDP,
        LocalAddr: net.UDPAddrFromAddrPort(srcAddr),
    }

    conn, err := d.Dial("udp", dstAddr)
    if err != nil {
        return nil, fmt.Errorf("failed to dial destination: %w", err)
    }
    return conn.(*net.UDPConn), nil
}

func (c *Client) writeReq(conn *net.UDPConn, cmd pb.IsRequestCmd) error {
    req := &pb.Request{Id: c.id, Cmd: cmd}
    buf, err := proto.Marshal(req)
    if err != nil {
        return fmt.Errorf("failed to marshal request: %w", err)
    }
    if _, err = conn.Write(buf); err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    return nil
}

func readResp(conn *net.UDPConn) (pb.Response, error) {
    buf := make([]byte, 1500)
    n, err := conn.Read(buf)
    if err != nil {
        return pb.Response{}, fmt.Errorf("failed to read response: %w", err)
    }

    var resp pb.Response
    if err = proto.Unmarshal(buf[:n], &resp); err != nil {
        return pb.Response{}, fmt.Errorf("failed to unmarshal response: %w", err)
    }
    return resp, nil
}

// Close stops internal goroutines and closes connections.
func (c *Client) Close() {
    c.closeOnce.Do(func() {
        if c.svrConn != nil {
            bye := &pb.Request_Bye{Bye: &pb.Bye{}}
            if err := c.writeReq(c.svrConn, bye); err != nil {
                log.Printf("failed to send bye request: %v", err)
            }
            c.svrConn.Close()
        }

        if c.listener != nil {
            c.listener.Close()
        }
    })
}