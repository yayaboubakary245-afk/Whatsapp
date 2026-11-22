package tcp

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/groupe12-supptic/whatsapp-campus/internal/protocol"
)

const MaxFrameSize = 4 * 1024 * 1024 // 4 MiB max par trame TCP

// WriteFrame encode un message au format [4 octets big-endian : longueur N][N octets : payload JSON].
func WriteFrame(w *bufio.Writer, env *protocol.Envelope) error {
	payload, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal envelope: %w", err)
	}
	if len(payload) > MaxFrameSize {
		return fmt.Errorf("frame too large: %d bytes", len(payload))
	}

	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(payload)))

	if _, err := w.Write(lenBuf[:]); err != nil {
		return fmt.Errorf("write frame length: %w", err)
	}
	if _, err := w.Write(payload); err != nil {
		return fmt.Errorf("write frame payload: %w", err)
	}
	return w.Flush()
}

// ReadFrame lit une trame [4 octets big-endian : N][N octets JSON] depuis bufio.Reader.
func ReadFrame(r *bufio.Reader) (*protocol.Envelope, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf[:])
	if length == 0 || length > MaxFrameSize {
		return nil, fmt.Errorf("invalid frame length: %d", length)
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	var env protocol.Envelope
	if err := json.Unmarshal(payload, &env); err != nil {
		return nil, fmt.Errorf("invalid json frame: %w", err)
	}
	return &env, nil
}
