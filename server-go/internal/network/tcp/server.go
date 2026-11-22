package tcp

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/groupe12-supptic/whatsapp-campus/internal/protocol"
)

// Server implémente les sockets TCP 100 % Golang sur le port 5000 (C-T-07).
type Server struct {
	Addr     string
	Registry *ClientRegistry
	listener net.Listener
}

func NewServer(port string, registry *ClientRegistry) *Server {
	return &Server{
		Addr:     ":" + port,
		Registry: registry,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.listener = ln
	log.Printf("[TCP] Serveur Go à l'écoute sur %s (framing length-prefixed 4 octets big-endian)", s.Addr)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Printf("[TCP] Erreur accept: %v", err)
				continue
			}
			go s.handleConn(conn)
		}
	}()
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var boundUserID string
	defer func() {
		if boundUserID != "" {
			s.Registry.Unregister(boundUserID)
		}
	}()

	for {
		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		env, err := ReadFrame(reader)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("[TCP] Déconnexion %s: %v", conn.RemoteAddr(), err)
			}
			return
		}

		switch env.Type {
		case protocol.TypeLogin, protocol.TypeRegister:
			if env.UserID != "" {
				boundUserID = env.UserID
				s.Registry.Register(boundUserID, conn, writer)
			}
			ack, _ := protocol.NewEnvelope(protocol.TypePresence, boundUserID, map[string]any{
				"status": "connected",
				"online": true,
			})
			_ = WriteFrame(writer, ack)

		case protocol.TypePing:
			pong, _ := protocol.NewEnvelope(protocol.TypePong, boundUserID, map[string]any{
				"message": "PONG",
			})
			_ = WriteFrame(writer, pong)

		case protocol.TypeSendMessage:
			var body struct {
				RecipientID string `json:"recipient_id"`
			}
			_ = json.Unmarshal(env.Payload, &body)
			if body.RecipientID != "" {
				_ = s.Registry.SendToUser(body.RecipientID, env)
			}

		case protocol.TypeLogout:
			return
		}
	}
}

func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
