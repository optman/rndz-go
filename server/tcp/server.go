package tcp

import (
	"context"
	"encoding/binary"
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
		return err
	}
	defer l.Close()
	s.listener = l

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		h := &peerHandler{
			connId:  s.nextId,
			conn:    conn,
			peers:   &s.peers,
			reqChan: make(chan pb.Request, 10),
		}
		go h.handleClient(ctx)

		s.nextId++
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

func (h *peerHandler) handleClient(ctx context.Context) {
	defer h.conn.Close()

	go func() {
		defer close(h.reqChan)
		for {
			h.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			if req, err := h.readReq(); err == nil {
				h.reqChan <- req
			} else {
				break
			}
		}
	}()

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP

		case req, ok := <-h.reqChan:
			if !ok || h.handleReq(&req) != nil {
				break LOOP
			}
		}
	}

	if h.id != "" {
		log.Printf("%s disconnected", h.id)
		h.peers.Delete(h.id)
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

	if v, ok := h.peers.Load(id); ok {
		return v.(*peerState), true
	} else {
		return nil, false
	}
}

func (h *peerHandler) handleReq(req *pb.Request) (err error) {
	id := req.GetId()

	switch req.Cmd.(type) {
	case *pb.Request_Ping:
		err = h.handlePing(id)
	case *pb.Request_Isync:
		err = h.handleIsync(req.Cmd.(*pb.Request_Isync).Isync, id)
	case *pb.Request_Fsync:
		err = h.handleFsync(req.Cmd.(*pb.Request_Fsync).Fsync, id)
	case *pb.Request_Rsync:
		err = h.handleRsync(req.Cmd.(*pb.Request_Rsync).Rsync, id)
	case *pb.Request_Bye:
		err = fmt.Errorf("bye")
	default:
		err = fmt.Errorf("unknown cmd")
	}

	if err != nil {
		log.Printf("%s fail, %s", h.id, err)
	}

	return
}

func (h *peerHandler) handlePing(id string) error {
	if ps, ok := h.getPeer(id); ok {
		if ps.id != h.connId {
			log.Printf("%s updated", ps.id)
			h.peers.Delete(id)
		}
	}

	h.addPeer(id)

	pong := &pb.Response_Pong{
		Pong: &pb.Pong{},
	}

	return h.writeResp(id, pong)
}

func (h *peerHandler) handleIsync(isync *pb.Isync, id string) error {
	targetId := isync.GetId()
	log.Printf("isync %s -> %s", id, targetId)

	if ps, ok := h.getPeer(targetId); ok {

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

	} else {
		log.Printf("target %s not found", targetId)

		rdr := &pb.Response_Redirect{
			Redirect: &pb.Redirect{
				Id: targetId,
			},
		}

		return h.writeResp(id, rdr)

	}
}

func (h *peerHandler) handleFsync(fsync *pb.Fsync, id string) error {
	targetId := fsync.GetId()
	log.Printf("fsync %s -> %s", targetId, id)

	return h.writeResp(id, &pb.Response_Fsync{
		Fsync: fsync})

}

func (h *peerHandler) handleRsync(rsync *pb.Rsync, id string) error {
	targetId := rsync.GetId()
	log.Printf("rsync %s -> %s@%s", id, targetId, h.id)

	if targetId == h.id {
		if ps, ok := h.getPeer(id); ok {
			rdr := &pb.Response_Redirect{Redirect: &pb.Redirect{
				Id:   id,
				Addr: ps.addr.String(),
			}}

			return h.writeResp(h.id, rdr)

		} else {
			log.Printf("target %s not found", id)
			return nil
		}

	}

	if ps, ok := h.getPeer(targetId); ok {
		rsync := pb.Request{
			Id: id,
			Cmd: &pb.Request_Rsync{
				Rsync: rsync,
			},
		}

		ps.reqChan <- rsync

	} else {
		log.Printf("target %s not found", targetId)
	}

	return nil
}

func (h *peerHandler) readReq() (req pb.Request, err error) {
	var n uint16
	if err = binary.Read(h.conn, binary.BigEndian, &n); err != nil {
		return
	}

	if n > 1500 {
		return pb.Request{}, fmt.Errorf("invalid request")
	}

	buf := make([]byte, n)
	if _, err = h.conn.Read(buf); err != nil {
		return
	}

	if err = proto.Unmarshal(buf, &req); err != nil {
		return
	}

	return
}

func (h *peerHandler) writeResp(id string, cmd pb.IsResponseCmd) (err error) {
	resp := &pb.Response{
		Id:  id,
		Cmd: cmd,
	}
	buf, _ := proto.Marshal(resp)
	if err = binary.Write(h.conn, binary.BigEndian, uint16(len(buf))); err != nil {
		return
	}

	_, err = h.conn.Write(buf)
	return
}
