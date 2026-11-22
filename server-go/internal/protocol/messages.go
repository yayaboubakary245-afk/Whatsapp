package protocol

import (
	"encoding/json"
	"time"
)

// Envelope représente un message JSON encadré par 4 octets big-endian sur TCP (Ndue).
type Envelope struct {
	Type      MessageType     `json:"type"`
	RequestID string          `json:"request_id,omitempty"`
	UserID    string          `json:"user_id,omitempty"`
	Timestamp int64           `json:"timestamp"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

func NewEnvelope(msgType MessageType, userID string, payload any) (*Envelope, error) {
	var raw json.RawMessage
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		raw = b
	}
	return &Envelope{
		Type:      msgType,
		UserID:    userID,
		Timestamp: time.Now().UnixMilli(),
		Payload:   raw,
	}, nil
}
