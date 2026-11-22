package call

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"

	ws "github.com/groupe12-supptic/whatsapp-campus/internal/websocket"
)

// SignalPayload transporte SDP (offer/answer), candidats ICE et événements d'appel (C-T-10).
type SignalPayload struct {
	Type         string          `json:"type"` // call_offer, call_answer, ice_candidate, call_accept, call_reject, call_end
	CallID       string          `json:"call_id"`
	FromUserID   string          `json:"from_user_id"`
	TargetUserID string          `json:"target_user_id"`
	CallType     string          `json:"call_type,omitempty"` // audio | video
	SDP          string          `json:"sdp,omitempty"`
	Candidate    json.RawMessage `json:"candidate,omitempty"`
	Timestamp    int64           `json:"timestamp"`
}

// SignalingService gère la signalisation WebRTC via gorilla/websocket et la persistance des appels (Ndue).
type SignalingService struct {
	DB  *pgxpool.Pool
	Hub *ws.Hub
}

func NewSignalingService(db *pgxpool.Pool, hub *ws.Hub) *SignalingService {
	svc := &SignalingService{DB: db, Hub: hub}
	hub.OnEvent = svc.HandleRawSignal
	return svc
}

// HandleWebSocket met à niveau la connexion HTTP en WebSocket (gorilla/websocket).
func (s *SignalingService) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[CALL] Upgrade WebSocket échoué: %v", err)
		return
	}
	// Garantit que gorilla/websocket est actif pour la signalisation
	_ = websocket.TextMessage
	s.Hub.Register(userID, conn)
	s.Hub.SendToUser(userID, map[string]any{
		"type":      "ws_ready",
		"user_id":   userID,
		"timestamp": time.Now().UnixMilli(),
	})
}

func (s *SignalingService) HandleRawSignal(senderID string, raw json.RawMessage) {
	var sig SignalPayload
	if err := json.Unmarshal(raw, &sig); err != nil {
		return
	}
	sig.FromUserID = senderID
	sig.Timestamp = time.Now().UnixMilli()

	if sig.TargetUserID != "" {
		s.Hub.SendToUser(sig.TargetUserID, sig)
	}
}

func (s *SignalingService) StartCall(ctx context.Context, callerID, targetUserID, conversationID, callType string) (string, error) {
	var callID string
	err := s.DB.QueryRow(ctx, `
		INSERT INTO calls (caller_id, conversation_id, call_type, status)
		VALUES ($1, NULLIF($2, '')::uuid, $3, 'ringing')
		RETURNING id::text
	`, callerID, conversationID, callType).Scan(&callID)
	if err != nil {
		return "", err
	}

	_, _ = s.DB.Exec(ctx, `
		INSERT INTO call_participants (call_id, user_id) VALUES ($1::uuid, $2::uuid), ($1::uuid, $3::uuid)
		ON CONFLICT DO NOTHING
	`, callID, callerID, targetUserID)

	s.Hub.SendToUser(targetUserID, SignalPayload{
		Type:         "incoming_call",
		CallID:       callID,
		FromUserID:   callerID,
		TargetUserID: targetUserID,
		CallType:     callType,
		Timestamp:    time.Now().UnixMilli(),
	})
	return callID, nil
}

func (s *SignalingService) AcceptCall(ctx context.Context, callID, userID string) error {
	_, err := s.DB.Exec(ctx, `
		UPDATE calls SET status = 'active', answered_at = NOW() WHERE id = $1::uuid
	`, callID)
	if err != nil {
		return err
	}
	var callerID string
	_ = s.DB.QueryRow(ctx, `SELECT caller_id::text FROM calls WHERE id = $1::uuid`, callID).Scan(&callerID)
	if callerID != "" {
		s.Hub.SendToUser(callerID, SignalPayload{
			Type:         "call_accept",
			CallID:       callID,
			FromUserID:   userID,
			TargetUserID: callerID,
			Timestamp:    time.Now().UnixMilli(),
		})
	}
	return nil
}

func (s *SignalingService) EndCall(ctx context.Context, callID, userID, endStatus string) error {
	if endStatus == "" {
		endStatus = "ended"
	}
	_, err := s.DB.Exec(ctx, `
		UPDATE calls
		SET status = $2,
			ended_at = NOW(),
			duration_seconds = COALESCE(EXTRACT(EPOCH FROM (NOW() - answered_at))::int, 0)
		WHERE id = $1::uuid
	`, callID, endStatus)
	if err != nil {
		return err
	}

	rows, err := s.DB.Query(ctx, `SELECT user_id::text FROM call_participants WHERE call_id = $1::uuid`, callID)
	if err == nil {
		defer rows.Close()
		var callerID string
		var duration int
		_ = s.DB.QueryRow(ctx, `SELECT caller_id::text, duration_seconds FROM calls WHERE id = $1::uuid`, callID).Scan(&callerID, &duration)
		for rows.Next() {
			var uid string
			if rows.Scan(&uid) == nil {
				dir := "incoming"
				if uid == callerID {
					dir = "outgoing"
				} else if endStatus == "missed" {
					dir = "missed"
				}
				_, _ = s.DB.Exec(ctx, `
					INSERT INTO call_history (call_id, user_id, direction, status, duration_seconds)
					VALUES ($1::uuid, $2::uuid, $3, $4, $5)
				`, callID, uid, dir, endStatus, duration)
				if uid != userID {
					s.Hub.SendToUser(uid, SignalPayload{
						Type:         "call_end",
						CallID:       callID,
						FromUserID:   userID,
						TargetUserID: uid,
						Timestamp:    time.Now().UnixMilli(),
					})
				}
			}
		}
	}
	return nil
}
