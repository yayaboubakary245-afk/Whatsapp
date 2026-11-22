package conversations

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConversationSummary struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Title         string `json:"title"`
	PeerUserID    string `json:"peer_user_id,omitempty"`
	PeerOnline    bool   `json:"peer_online"`
	LastMessage   string `json:"last_message"`
	LastMessageAt string `json:"last_message_at"`
	UnreadCount   int    `json:"unread_count"`
}

func GetOrCreatePrivate(ctx context.Context, db *pgxpool.Pool, userA, userB string) (string, error) {
	var convID string
	err := db.QueryRow(ctx, `
		SELECT c.id::text
		FROM conversations c
		JOIN conversation_members m1 ON m1.conversation_id = c.id AND m1.user_id = $1::uuid
		JOIN conversation_members m2 ON m2.conversation_id = c.id AND m2.user_id = $2::uuid
		WHERE c.type = 'private'
		LIMIT 1
	`, userA, userB).Scan(&convID)
	if err == nil && convID != "" {
		return convID, nil
	}

	err = db.QueryRow(ctx, `
		INSERT INTO conversations (type, created_by) VALUES ('private', $1::uuid) RETURNING id::text
	`, userA).Scan(&convID)
	if err != nil {
		return "", err
	}
	_, err = db.Exec(ctx, `
		INSERT INTO conversation_members (conversation_id, user_id)
		VALUES ($1::uuid, $2::uuid), ($1::uuid, $3::uuid)
	`, convID, userA, userB)
	return convID, err
}

func ListForUser(ctx context.Context, db *pgxpool.Pool, userID string) ([]ConversationSummary, error) {
	rows, err := db.Query(ctx, `
		SELECT c.id::text, c.type,
		       COALESCE(g.name, peer_p.first_name || ' ' || peer_p.last_name, 'Discussion'),
		       COALESCE(peer_u.id::text, ''),
		       COALESCE(peer_u.is_online, FALSE),
		       COALESCE((SELECT content FROM messages m WHERE m.conversation_id = c.id AND m.is_deleted = FALSE ORDER BY m.created_at DESC LIMIT 1), ''),
		       to_char(COALESCE(c.last_message_at, c.created_at), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		       (SELECT COUNT(*) FROM messages m
		        LEFT JOIN message_status ms ON ms.message_id = m.id AND ms.user_id = $1::uuid
		        WHERE m.conversation_id = c.id AND m.sender_id <> $1::uuid AND COALESCE(ms.status, 'sent') <> 'read')
		FROM conversations c
		JOIN conversation_members cm ON cm.conversation_id = c.id AND cm.user_id = $1::uuid
		LEFT JOIN groups g ON g.conversation_id = c.id
		LEFT JOIN conversation_members peer_cm ON peer_cm.conversation_id = c.id AND peer_cm.user_id <> $1::uuid AND c.type = 'private'
		LEFT JOIN users peer_u ON peer_u.id = peer_cm.user_id
		LEFT JOIN profiles peer_p ON peer_p.user_id = peer_u.id
		ORDER BY c.last_message_at DESC NULLS LAST
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []ConversationSummary
	for rows.Next() {
		var cs ConversationSummary
		if rows.Scan(&cs.ID, &cs.Type, &cs.Title, &cs.PeerUserID, &cs.PeerOnline, &cs.LastMessage, &cs.LastMessageAt, &cs.UnreadCount) == nil {
			res = append(res, cs)
		}
	}
	return res, nil
}
