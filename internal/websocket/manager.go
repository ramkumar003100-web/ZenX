package websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type AuthzFunc func(*http.Request) (clientID, role, room string, allowed bool)

type Manager struct {
	Hub      *Hub
	Upgrader websocket.Upgrader
	Authz    AuthzFunc
}

func NewManager(hub *Hub, authz AuthzFunc) *Manager {
	return &Manager{
		Hub:   hub,
		Authz: authz,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(*http.Request) bool { return true },
		},
	}
}

func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, role, room, ok := m.Authz(r)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	conn, err := m.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "upgrade failed", http.StatusBadRequest)
		return
	}
	client := &Client{ID: id, Role: role, Room: room, Send: make(chan []byte, 64)}
	m.Hub.Register(client)
	defer m.Hub.Unregister(client.ID)

	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case msg, ok := <-client.Send:
				if !ok {
					_ = conn.Close()
					return
				}
				_ = conn.WriteMessage(websocket.TextMessage, msg)
			case <-ticker.C:
				_ = conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(5*time.Second))
			}
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		m.Hub.BroadcastRoom(r.Context(), room, msg)
	}
}
