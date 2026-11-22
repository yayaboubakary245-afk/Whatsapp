package messaging

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Message struct {
	ID               string `json:"id"`
	ConversationID   string `json:"conversation_id"`
	SenderID         string `json:"sender_id"`
	SenderName       string `json:"sender_name"`
	ReplyToMessageID string `json:"reply_to_message_id,omitempty"`
	MessageType      string `json:"message_type"`
	Content          string `json:"content"`
	Status           string `json:"status"`
	IsDeleted        bool   `json:"is_deleted"`
	CreatedAt        string `json:"created_at"`
}

func Send(ctx context.Context, db *pgxpool.Pool, convID, senderID, replyToID, msgType, content string) (*Message, error) {
	if msgType == "" {
		msgType = "text"
	}
	var m Message
	err := db.QueryRow(ctx, `
		INSERT INTO messages (conversation_id, sender_id, reply_to_message_id, message_type, content)
		VALUES ($1::uuid, $2::uuid, NULLIF($3, '')::uuid, $4, $5)
		RETURNING id::text, conversation_id::text, sender_id::text, COALESCE(reply_to_message_id::text, ''), message_type, content, is_deleted, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`, convID, senderID, replyToID, msgType, content).Scan(
		&m.ID, &m.ConversationID, &m.SenderID, &m.ReplyToMessageID, &m.MessageType, &m.Content, &m.IsDeleted, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	m.Status = "sent"

	_, _ = db.Exec(ctx, `UPDATE conversations SET last_message_at = NOW() WHERE id = $1::uuid`, convID)
	_, _ = db.Exec(ctx, `
		INSERT INTO message_status (message_id, user_id, status)
		SELECT $1::uuid, cm.user_id, CASE WHEN cm.user_id = $2::uuid THEN 'read' ELSE 'sent' END
		FROM conversation_members cm
		WHERE cm.conversation_id = $3::uuid
		ON CONFLICT DO NOTHING
	`, m.ID, senderID, convID)

	_ = db.QueryRow(ctx, `SELECT first_name || ' ' || last_name FROM profiles WHERE user_id = $1::uuid`, senderID).Scan(&m.SenderName)
	return &m, nil
}

func ListByConversation(ctx context.Context, db *pgxpool.Pool, convID, viewerID, searchQuery string) ([]Message, error) {
	rows, err := db.Query(ctx, `
		SELECT m.id::text, m.conversation_id::text, m.sender_id::text,
		       COALESCE(p.first_name || ' ' || p.last_name, 'Utilisateur'),
		       COALESCE(m.reply_to_message_id::text, ''),
		       m.message_type,
		       CASE WHEN m.is_deleted THEN 'Message supprimé' ELSE m.content END,
		       COALESCE((SELECT ms.status FROM message_status ms WHERE ms.message_id = m.id AND ms.user_id <> m.sender_id LIMIT 1), 'sent'),
		       m.is_deleted,
		       to_char(m.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM messages m
		LEFT JOIN profiles p ON p.user_id = m.sender_id
		WHERE m.conversation_id = $1::uuid
		  AND ($2 = '' OR m.content ILIKE '%'||$2||'%')
		ORDER BY m.created_at ASC
		LIMIT 200
	`, convID, searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Message
	for rows.Next() {
		var m Message
		if rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.SenderName, &m.ReplyToMessageID, &m.MessageType, &m.Content, &m.Status, &m.IsDeleted, &m.CreatedAt) == nil {
			list = append(list, m)
		}
	}

	_, _ = db.Exec(ctx, `
		UPDATE message_status ms
		SET status = 'read', updated_at = NOW()
		FROM messages m
		WHERE ms.message_id = m.id AND m.conversation_id = $1::uuid AND ms.user_id = $2::uuid AND ms.status <> 'read'
	`, convID, viewerID)
	return list, nil
}
