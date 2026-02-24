package websocket

import (
	"context"
	"sync"

	"zenx/pkg/metrics"
)

type Client struct {
	ID    string
	Role  string
	Room  string
	Send  chan []byte
	close func()
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client
	rooms   map[string]map[string]*Client
}

func NewHub() *Hub {
	return &Hub{clients: map[string]*Client{}, rooms: map[string]map[string]*Client{}}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.ID] = c
	if c.Room != "" {
		if h.rooms[c.Room] == nil {
			h.rooms[c.Room] = map[string]*Client{}
		}
		h.rooms[c.Room][c.ID] = c
	}
	metrics.WebsocketConnections.Inc()
}

func (h *Hub) Unregister(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if c, ok := h.clients[clientID]; ok {
		delete(h.clients, clientID)
		if c.Room != "" && h.rooms[c.Room] != nil {
			delete(h.rooms[c.Room], c.ID)
			if len(h.rooms[c.Room]) == 0 {
				delete(h.rooms, c.Room)
			}
		}
		close(c.Send)
		metrics.WebsocketConnections.Dec()
	}
}

func (h *Hub) Broadcast(_ context.Context, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.clients {
		select {
		case c.Send <- msg:
		default:
		}
	}
}

func (h *Hub) BroadcastRoom(_ context.Context, room string, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.rooms[room] {
		select {
		case c.Send <- msg:
		default:
		}
	}
}

func (h *Hub) SendToClient(clientID string, msg []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	c, ok := h.clients[clientID]
	if !ok {
		return false
	}
	select {
	case c.Send <- msg:
	default:
	}
	return true
}
