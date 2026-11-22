package tcp

import (
	"bufio"
	"net"
	"sync"

	"github.com/groupe12-supptic/whatsapp-campus/internal/protocol"
)

type ClientConn struct {
	UserID string
	Conn   net.Conn
	Writer *bufio.Writer
	Mu     sync.Mutex
}

// ClientRegistry gère les connexions TCP actives indexées par UserID (Ndue).
type ClientRegistry struct {
	mu      sync.RWMutex
	clients map[string]*ClientConn
}

func NewClientRegistry() *ClientRegistry {
	return &ClientRegistry{
		clients: make(map[string]*ClientConn),
	}
}

func (r *ClientRegistry) Register(userID string, conn net.Conn, writer *bufio.Writer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[userID] = &ClientConn{
		UserID: userID,
		Conn:   conn,
		Writer: writer,
	}
}

func (r *ClientRegistry) Unregister(userID string) {
	if userID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, userID)
}

func (r *ClientRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

func (r *ClientRegistry) SendToUser(userID string, env *protocol.Envelope) error {
	r.mu.RLock()
	client, ok := r.clients[userID]
	r.mu.RUnlock()
	if !ok {
		return nil
	}
	client.Mu.Lock()
	defer client.Mu.Unlock()
	return WriteFrame(client.Writer, env)
}
