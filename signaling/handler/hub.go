package handler

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"signaling/proto"
	"sync"
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}
var hub *Hub

func init() {
	hub = New()
}

type Hub struct {
	// Registered clients.
	clients map[string]*websocket.Conn
	mu      sync.Mutex
}

func (h *Hub) Delete(id string) {
	h.mu.Lock()
	if c := h.clients[id]; c != nil {
		delete(h.clients, id)
	}
	h.mu.Unlock()
	return
}

func New() *Hub {
	h := new(Hub)
	h.clients = make(map[string]*websocket.Conn)
	return h
}

func (h *Hub) Register(id string, c *websocket.Conn) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conn := h.clients[id]; conn != nil {
		return false
	}
	h.clients[id] = c
	return true
}

func (h *Hub) SafeTransport(req *proto.Request, cmd int64) (err error) {
	var data []byte
	h.mu.Lock()
	defer h.mu.Unlock()
	c := h.clients[req.Target]
	if c == nil {
		return errors.New("remote user code is not valid")
	}
	switch cmd {
	case proto.CmdOffer:
		data = proto.NewResponse(converseCmd(cmd), req.ToOfferString())
	default:
		data = proto.NewResponse(converseCmd(cmd), req.ToString())
	}
	logrus.Printf("cmd: %d source: %s target: %s\n", req.Cmd, req.Source, req.Target)
	return c.WriteMessage(websocket.TextMessage, data)
}

func converseCmd(cmd int64) int64 {
	switch cmd {
	case proto.CmdOffer:
		return proto.OfferResponse
	case proto.CmdAnswer:
		return proto.AnswerResponse
	case proto.CmdCandidate:
		return proto.CandidateResponse
	}
	return 0
}
