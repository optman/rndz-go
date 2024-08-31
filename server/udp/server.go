package udp

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/optman/rndz-go/proto"
	"google.golang.org/protobuf/proto"
)

func New(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

type peerState struct {
	sync.RWMutex
	addr     net.Addr
	lastPing time.Time
}

type Server struct {
	listenAddr string
	listener   net.PacketConn
	closeOnce  sync.Once
	peers      sync.Map
}

func (s *Server) Run(ctx context.Context) error {
	l, err := net.ListenPacket("udp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.listenAddr, err)
	}

	defer l.Close()
	s.listener = l

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go s.gc(ctx)

	buf := make([]byte, 1500)
	for {
		n, addr, err := l.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("error reading packet: %w", err)
		}

		if err = s.handlePacket(ctx, buf[:n], addr); err != nil {
			log.Printf("error handling packet from %s: %v", addr, err)
		}
	}
}

func (s *Server) Close() {
	s.closeOnce.Do(func() {
		if s.listener != nil {
			s.listener.Close()
		}
	})
}

func (s *Server) gc(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			s.peers.Range(func(k, v any) bool {
				ps := v.(*peerState)
				if now.Sub(ps.lastPing) > 60*time.Second {
					log.Printf("peer %s expired", k)
					s.peers.Delete(k)
				}
				return true
			})
		}
	}
}

func (s *Server) handlePacket(ctx context.Context, buf []byte, addr net.Addr) error {
	var req pb.Request
	if err := proto.Unmarshal(buf, &req); err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}

	id := req.GetId()
	switch cmd := req.Cmd.(type) {
	case *pb.Request_Ping:
		s.handlePing(ctx, cmd, id, addr)
	case *pb.Request_Isync:
		s.handleIsync(ctx, cmd, id, addr)
	case *pb.Request_Rsync:
		s.handleRsync(ctx, cmd, id, addr)
	case *pb.Request_Bye:
		s.handleBye(ctx, id)
	default:
		log.Println("unknown command")
	}

	return nil
}

func (s *Server) updatePeer(id string, addr net.Addr) {
	ps := &peerState{
		addr:     addr,
		lastPing: time.Now(),
	}
	v, load := s.peers.LoadOrStore(id, ps)
	if load {
		ps := v.(*peerState)
		ps.Lock()
		ps.addr = addr
		ps.lastPing = time.Now()
		ps.Unlock()
	}
}

func (s *Server) getPeer(id string) (net.Addr, bool) {
	v, ok := s.peers.Load(id)
	if !ok {
		return nil, false
	}
	ps := v.(*peerState)
	ps.RLock()
	defer ps.RUnlock()
	return ps.addr, true
}

func (s *Server) handlePing(ctx context.Context, ping *pb.Request_Ping, id string, addr net.Addr) {
	s.updatePeer(id, addr)

	pong := &pb.Response{Id: id, Cmd: &pb.Response_Pong{Pong: &pb.Pong{}}}

	s.sendResp(ctx, pong, addr)
}

func (s *Server) handleIsync(ctx context.Context, cmds *pb.Request_Isync, id string, addr net.Addr) {
	targetId := cmds.Isync.GetId()
	log.Printf("isync %s -> %s", id, targetId)

	if targetAddr, ok := s.getPeer(targetId); ok {
		s.updatePeer(id, addr)

		fsync := &pb.Response{
			Id: targetId,
			Cmd: &pb.Response_Fsync{
				Fsync: &pb.Fsync{
					Id:   id,
					Addr: addr.String(),
				},
			},
		}

		s.sendResp(ctx, fsync, targetAddr)
	} else {
		rdr := &pb.Response{
			Id: id,
			Cmd: &pb.Response_Redirect{
				Redirect: &pb.Redirect{
					Id: targetId,
				},
			},
		}

		s.sendResp(ctx, rdr, addr)

		log.Printf("target %s not found", targetId)
	}
}

func (s *Server) handleRsync(ctx context.Context, cmds *pb.Request_Rsync, id string, addr net.Addr) {
	targetId := cmds.Rsync.GetId()
	log.Printf("rsync %s -> %s", id, targetId)

	if targetAddr, ok := s.getPeer(targetId); ok {
		rdr := &pb.Response{
			Id: targetId,
			Cmd: &pb.Response_Redirect{
				Redirect: &pb.Redirect{
					Id:   id,
					Addr: addr.String(),
				},
			},
		}

		s.sendResp(ctx, rdr, targetAddr)
	} else {
		log.Printf("target %s not found", targetId)
	}
}

func (s *Server) handleBye(ctx context.Context, id string) {
	log.Printf("bye %s", id)

	s.peers.Delete(id)
}

func (s *Server) sendResp(ctx context.Context, resp *pb.Response, addr net.Addr) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		return
	}

	if _, err := s.listener.WriteTo(buf, addr); err != nil {
		log.Printf("failed to send response to %s: %v", addr, err)
	}
}
