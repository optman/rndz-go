package udp

import (
	"context"
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
		return err
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
			return err
		}

		if err = s.handlePacket(buf[:n], addr); err != nil {
			log.Println(err)
		}
	}

}

func (s *Server) Close() {
	s.closeOnce.Do(func() {
		if s.listener != nil {
			s.Close()
		}
	})
}

func (s *Server) gc(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
		}

		now := time.Now()
		s.peers.Range(func(k, v any) bool {
			ps := v.(*peerState)
			if now.Sub(ps.lastPing) > 60*time.Second {
				log.Printf("%s expired", k)
				s.peers.Delete(k)
			}

			return true
		})
	}
}

func (s *Server) handlePacket(buf []byte, addr net.Addr) error {
	var req pb.Request
	if err := proto.Unmarshal(buf, &req); err != nil {
		return err
	}

	id := req.GetId()
	switch req.Cmd.(type) {
	case *pb.Request_Ping:
		s.handlePing(req.Cmd.(*pb.Request_Ping), id, addr)
	case *pb.Request_Isync:
		s.handleIsync(req.Cmd.(*pb.Request_Isync), id, addr)
	case *pb.Request_Rsync:
		s.handleRsync(req.Cmd.(*pb.Request_Rsync), id, addr)
	case *pb.Request_Bye:
		s.handleBye(id)
	default:
		log.Println("unkown cmd")
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

func (s *Server) getPeer(id string) (addr net.Addr, ok bool) {
	if v, ok := s.peers.Load(id); ok {
		ps := v.(*peerState)
		ps.RLock()
		defer ps.RUnlock()

		return ps.addr, true
	} else {
		return nil, false
	}
}

func (s *Server) handlePing(ping *pb.Request_Ping, id string, addr net.Addr) {
	s.updatePeer(id, addr)

	pong := &pb.Response_Pong{
		Pong: &pb.Pong{},
	}

	s.sendResp(pong, id, addr)

}

func (s *Server) handleIsync(isync *pb.Request_Isync, id string, addr net.Addr) {
	targetId := isync.Isync.GetId()
	log.Printf("isync %s -> %s", id, targetId)

	if targetAddr, ok := s.getPeer(targetId); ok {

		s.updatePeer(id, addr)

		fsync := &pb.Response_Fsync{
			Fsync: &pb.Fsync{
				Id:   id,
				Addr: addr.String(),
			},
		}

		s.sendResp(fsync, targetId, targetAddr)

	} else {
		rdr := &pb.Response_Redirect{
			Redirect: &pb.Redirect{
				Id: targetId,
			},
		}

		s.sendResp(rdr, id, addr)

		log.Printf("target %s not found", targetId)
	}

}

func (s *Server) handleRsync(rsync *pb.Request_Rsync, id string, addr net.Addr) {
	targetId := rsync.Rsync.GetId()
	log.Printf("rsync %s -> %s", id, targetId)

	if targetAddr, ok := s.getPeer(targetId); ok {
		rdr := &pb.Response_Redirect{
			Redirect: &pb.Redirect{
				Id:   id,
				Addr: addr.String(),
			},
		}

		s.sendResp(rdr, targetId, targetAddr)

	} else {
		log.Printf("target %s not found", targetId)
	}
}

func (s *Server) handleBye(id string) {
	log.Printf("bye %s", id)

	s.peers.Delete(id)
}

func (s *Server) sendResp(cmd pb.IsResponseCmd, id string, addr net.Addr) {
	resp := &pb.Response{
		Id:  id,
		Cmd: cmd,
	}

	buf, _ := proto.Marshal(resp)

	s.listener.WriteTo(buf, addr)
}
