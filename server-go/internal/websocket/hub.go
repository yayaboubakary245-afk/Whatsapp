package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub gère les connexions WebSocket pour les push temps réel et la signalisation WebRTC (Ndue).
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client
	OnEvent func(senderID string, msg json.RawMessage)
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

func (h *Hub) Register(userID string, conn *websocket.Conn) *Client {
	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 64),
	}
	h.mu.Lock()
	h.clients[userID] = client
	h.mu.Unlock()

	go h.writePump(client)
	go h.readPump(client)
	return client
}

func (h *Hub) Unregister(userID string) {
	h.mu.Lock()
	if c, ok := h.clients[userID]; ok {
		close(c.Send)
		delete(h.clients, userID)
	}
	h.mu.Unlock()
}

func (h *Hub) SendToUser(userID string, payload any) bool {
	b, err := json.Marshal(payload)
	if err != nil {
		return false
	}
	h.mu.RLock()
	c, ok := h.clients[userID]
	h.mu.RUnlock()
	if !ok {
		return false
	}
	select {
	case c.Send <- b:
		return true
	default:
		return false
	}
}

func (h *Hub) Broadcast(payload any) {
	b, err := json.Marshal(payload)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.clients {
		select {
		case c.Send <- b:
		default:
		}
	}
}

func (h *Hub) ConnectedCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

func (h *Hub) readPump(c *Client) {
	defer func() {
		h.Unregister(c.UserID)
		_ = c.Conn.Close()
	}()
	c.Conn.SetReadLimit(65536)
	_ = c.Conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		return nil
	})
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("[WS] Déconnexion utilisateur %s: %v", c.UserID, err)
			return
		}
		if h.OnEvent != nil {
			h.OnEvent(c.UserID, json.RawMessage(msg))
		}
	}
}

func (h *Hub) writePump(c *Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
