package udp

import (
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Server implémente les sockets UDP 100 % Golang sur le port 5001 (net.ListenPacket, ReadFromUDP, WriteToUDP).
type Server struct {
	Addr          string
	conn          *net.UDPConn
	PacketsIn     atomic.Int64
	PacketsOut    atomic.Int64
	LastPingAt    atomic.Int64
	mu            sync.RWMutex
	activePeers   map[string]time.Time
}

func NewServer(port string) *Server {
	return &Server{
		Addr:        ":" + port,
		activePeers: make(map[string]time.Time),
	}
}

func (s *Server) Start() error {
	pc, err := net.ListenPacket("udp", s.Addr)
	if err != nil {
		return err
	}
	udpConn, ok := pc.(*net.UDPConn)
	if !ok {
		_ = pc.Close()
		return net.InvalidAddrError("expected *net.UDPConn")
	}
	s.conn = udpConn
	log.Printf("[UDP] Serveur Go à l'écoute sur %s (PING/PONG & présence)", s.Addr)

	go func() {
		buf := make([]byte, 2048)
		for {
			n, remoteAddr, err := s.conn.ReadFromUDP(buf)
			if err != nil {
				return
			}
			s.PacketsIn.Add(1)
			msg := strings.TrimSpace(string(buf[:n]))

			s.mu.Lock()
			s.activePeers[remoteAddr.String()] = time.Now()
			s.mu.Unlock()

			var reply string
			switch {
			case strings.EqualFold(msg, "PING"):
				s.LastPingAt.Store(time.Now().UnixMilli())
				reply = "PONG"
			case strings.HasPrefix(strings.ToUpper(msg), "PRESENCE:"):
				reply = "PRESENCE_ACK"
			default:
				reply = "OK:" + msg
			}

			if _, err := s.conn.WriteToUDP([]byte(reply+"\n"), remoteAddr); err == nil {
				s.PacketsOut.Add(1)
			}
		}
	}()
	return nil
}

func (s *Server) Stats() map[string]any {
	s.mu.RLock()
	peerCount := len(s.activePeers)
	s.mu.RUnlock()

	return map[string]any{
		"addr":         s.Addr,
		"packets_in":   s.PacketsIn.Load(),
		"packets_out":  s.PacketsOut.Load(),
		"last_ping_ms": s.LastPingAt.Load(),
		"active_peers": peerCount,
		"status":       "listening",
	}
}

func (s *Server) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
