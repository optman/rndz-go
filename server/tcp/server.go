package tcp

import (
	"context"
	"encoding/binary"
	"errors"
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

type Server struct {
	listenAddr string
	listener   net.Listener
	closeOnce  sync.Once
	peers      sync.Map
	nextId     uint64
}

type peerState struct {
	id      uint64
	addr    net.Addr
	reqChan chan pb.Request
}

func (s *Server) Run(ctx context.Context) error {
	l, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.listenAddr, err)
	}
	defer l.Close()
	log.Printf("Server is listening on %s", s.listenAddr)

	s.listener = l
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				return fmt.Errorf("failed to accept connection: %w", err)
			}
		}

		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	h := &peerHandler{
		connId:  s.nextId,
		conn:    conn,
		peers:   &s.peers,
		reqChan: make(chan pb.Request, 10),
	}
	s.nextId++

	go h.readRequests()

	h.handleRequests(ctx)

	if h.id != "" {
		log.Printf("%s disconnected", h.id)
		h.peers.Delete(h.id)
	}
}

func (s *Server) Close() {
	s.closeOnce.Do(func() {
		if s.listener != nil {
			s.listener.Close()
		}
	})
}

type peerHandler struct {
	connId  uint64
	id      string
	conn    net.Conn
	peers   *sync.Map
	reqChan chan pb.Request
}

func (h *peerHandler) readRequests() {
	defer close(h.reqChan)

	for {
		h.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		req, err := h.readReq()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Printf("read request error: %v", err)
			}
			return
		}

		h.reqChan <- req
	}
}

func (h *peerHandler) handleRequests(ctx context.Context) {
	defer h.conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-h.reqChan:
			if !ok {
				return
			}
			if err := h.handleReq(&req); err != nil {
				log.Printf("handle request error: %v", err)
				return
			}
		}
	}
}

func (h *peerHandler) addPeer(id string) {
	if h.id == "" {
		h.id = id
	}

	ps := &peerState{
		id:      h.connId,
		addr:    h.conn.RemoteAddr(),
		reqChan: h.reqChan,
	}

	h.peers.LoadOrStore(id, ps)
}

func (h *peerHandler) getPeer(id string) (*peerState, bool) {
	v, ok := h.peers.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*peerState), true
}

func (h *peerHandler) handleReq(req *pb.Request) error {
	id := req.GetId()

	switch cmd := req.Cmd.(type) {
	case *pb.Request_Ping:
		return h.handlePing(id)
	case *pb.Request_Isync:
		return h.handleIsync(cmd.Isync, id)
	case *pb.Request_Fsync:
		return h.handleFsync(cmd.Fsync, id)
	case *pb.Request_Rsync:
		return h.handleRsync(cmd.Rsync, id)
	case *pb.Request_Bye:
		return fmt.Errorf("bye")
	default:
		return fmt.Errorf("unknown command")
	}
}

func (h *peerHandler) handlePing(id string) error {
	ps, ok := h.getPeer(id)
	if ok && ps.id != h.connId {
		log.Printf("peer %s updated", ps.id)
		h.peers.Delete(id)
	}

	h.addPeer(id)

	pong := &pb.Response_Pong{
		Pong: &pb.Pong{},
	}

	return h.writeResp(id, pong)
}

func (h *peerHandler) handleIsync(isync *pb.Isync, id string) error {
	targetId := isync.GetId()
	log.Printf("isync from %s to %s", id, targetId)

	ps, ok := h.getPeer(targetId)
	if !ok {
		log.Printf("target peer %s not found", targetId)
		rdr := &pb.Response_Redirect{
			Redirect: &pb.Redirect{
				Id: targetId,
			},
		}
		return h.writeResp(id, rdr)
	}

	h.addPeer(id)
	fsync := pb.Request{
		Id: targetId,
		Cmd: &pb.Request_Fsync{
			Fsync: &pb.Fsync{
				Id:   id,
				Addr: h.conn.RemoteAddr().String(),
			},
		},
	}
	ps.reqChan <- fsync

	return nil
}

func (h *peerHandler) handleFsync(fsync *pb.Fsync, id string) error {
	targetId := fsync.GetId()
	log.Printf("fsync from %s to %s", targetId, id)

	return h.writeResp(id, &pb.Response_Fsync{
		Fsync: fsync,
	})
}

func (h *peerHandler) handleRsync(rsync *pb.Rsync, id string) error {
	targetId := rsync.GetId()
	log.Printf("rsync from %s to %s", id, targetId)

	if targetId == h.id {
		ps, ok := h.getPeer(id)
		if !ok {
			log.Printf("target %s not found", id)
			return nil
		}

		rdr := &pb.Response_Redirect{
			Redirect: &pb.Redirect{
				Id:   id,
				Addr: ps.addr.String(),
			},
		}

		return h.writeResp(h.id, rdr)
	}

	ps, ok := h.getPeer(targetId)
	if !ok {
		log.Printf("target %s not found", targetId)
		return nil
	}

	rsyncReq := pb.Request{
		Id: id,
		Cmd: &pb.Request_Rsync{
			Rsync: rsync,
		},
	}
	ps.reqChan <- rsyncReq

	return nil
}

func (h *peerHandler) readReq() (pb.Request, error) {
	var n uint16
	if err := binary.Read(h.conn, binary.BigEndian, &n); err != nil {
		if errors.Is(err, net.ErrClosed) {
			return pb.Request{}, err
		}
		return pb.Request{}, fmt.Errorf("read error: %w", err)
	}

	if n > 1500 {
		return pb.Request{}, fmt.Errorf("invalid request size: %d", n)
	}

	buf := make([]byte, n)
	if _, err := h.conn.Read(buf); err != nil {
		return pb.Request{}, fmt.Errorf("read error: %w", err)
	}

	var req pb.Request
	if err := proto.Unmarshal(buf, &req); err != nil {
		return pb.Request{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return req, nil
}

func (h *peerHandler) writeResp(id string, cmd pb.IsResponseCmd) error {
	resp := &pb.Response{
		Id:  id,
		Cmd: cmd,
	}

	buf, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := binary.Write(h.conn, binary.BigEndian, uint16(len(buf))); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	if _, err := h.conn.Write(buf); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}
