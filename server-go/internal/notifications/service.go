package notifications

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Notification struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

func Create(ctx context.Context, db *pgxpool.Pool, userID, nType, title, body, refID string) error {
	_, err := db.Exec(ctx, `
		INSERT INTO notifications (user_id, type, title, body, reference_id)
		VALUES ($1::uuid, $2, $3, $4, NULLIF($5, '')::uuid)
	`, userID, nType, title, body, refID)
	return err
}

func ListForUser(ctx context.Context, db *pgxpool.Pool, userID string) ([]Notification, error) {
	rows, err := db.Query(ctx, `
		SELECT id::text, type, title, body, is_read, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM notifications
		WHERE user_id = $1::uuid
		ORDER BY created_at DESC
		LIMIT 50
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Notification
	for rows.Next() {
		var n Notification
		if rows.Scan(&n.ID, &n.Type, &n.Title, &n.Body, &n.IsRead, &n.CreatedAt) == nil {
			list = append(list, n)
		}
	}
	return list, nil
}

func MarkAllRead(ctx context.Context, db *pgxpool.Pool, userID string) error {
	_, err := db.Exec(ctx, `UPDATE notifications SET is_read = TRUE WHERE user_id = $1::uuid AND is_read = FALSE`, userID)
	return err
}
