package tests

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/groupe12-supptic/whatsapp-campus/internal/network/tcp"
	"github.com/groupe12-supptic/whatsapp-campus/internal/network/udp"
	"github.com/groupe12-supptic/whatsapp-campus/internal/protocol"
)

// TestSprint1TCPFramingAndPing vérifie le serveur TCP Go et l'encadrement 4 octets big-endian (Ndue).
func TestSprint1TCPFramingAndPing(t *testing.T) {
	reg := tcp.NewClientRegistry()
	srv := tcp.NewServer("15000", reg)
	if err := srv.Start(); err != nil {
		t.Fatalf("TCP start failed: %v", err)
	}
	defer srv.Stop()

	conn, err := net.DialTimeout("tcp", "127.0.0.1:15000", 2*time.Second)
	if err != nil {
		t.Fatalf("TCP dial failed: %v", err)
	}
	defer conn.Close()

	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	pingEnv, _ := protocol.NewEnvelope(protocol.TypePing, "u-test-01", nil)
	if err := tcp.WriteFrame(w, pingEnv); err != nil {
		t.Fatalf("WriteFrame failed: %v", err)
	}

	reply, err := tcp.ReadFrame(r)
	if err != nil {
		t.Fatalf("ReadFrame failed: %v", err)
	}
	if reply.Type != protocol.TypePong {
		t.Fatalf("expected PONG, got %s", reply.Type)
	}
}

// TestSprint1UDPPingPong vérifie le serveur UDP Go (net.ListenPacket / ReadFromUDP / WriteToUDP).
func TestSprint1UDPPingPong(t *testing.T) {
	srv := udp.NewServer("15001")
	if err := srv.Start(); err != nil {
		t.Fatalf("UDP start failed: %v", err)
	}
	defer srv.Stop()

	conn, err := net.DialTimeout("udp", "127.0.0.1:15001", 2*time.Second)
	if err != nil {
		t.Fatalf("UDP dial failed: %v", err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
	if _, err := conn.Write([]byte("PING\n")); err != nil {
		t.Fatalf("UDP write failed: %v", err)
	}

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("UDP read failed: %v", err)
	}
	got := strings.TrimSpace(string(buf[:n]))
	if got != "PONG" {
		t.Fatalf("expected PONG, got %q", got)
	}
}
