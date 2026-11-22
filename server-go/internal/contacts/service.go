package contacts

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Contact struct {
	ID            string `json:"id"`
	ContactUserID string `json:"contact_user_id"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Role          string `json:"role"`
	IsBlocked     bool   `json:"is_blocked"`
	IsOnline      bool   `json:"is_online"`
}

func List(ctx context.Context, db *pgxpool.Pool, ownerID string) ([]Contact, error) {
	rows, err := db.Query(ctx, `
		SELECT c.id::text, u.id::text, u.email, p.first_name, p.last_name, r.name, c.is_blocked, u.is_online
		FROM contacts c
		JOIN users u ON u.id = c.contact_user_id
		JOIN roles r ON r.id = u.role_id
		JOIN profiles p ON p.user_id = u.id
		WHERE c.owner_id = $1::uuid
		ORDER BY p.first_name ASC, p.last_name ASC
	`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Contact
	for rows.Next() {
		var c Contact
		if rows.Scan(&c.ID, &c.ContactUserID, &c.Email, &c.FirstName, &c.LastName, &c.Role, &c.IsBlocked, &c.IsOnline) == nil {
			list = append(list, c)
		}
	}
	return list, nil
}

func Add(ctx context.Context, db *pgxpool.Pool, ownerID, contactUserID, nickname string) error {
	_, err := db.Exec(ctx, `
		INSERT INTO contacts (owner_id, contact_user_id, nickname)
		VALUES ($1::uuid, $2::uuid, $3)
		ON CONFLICT (owner_id, contact_user_id) DO UPDATE SET nickname = EXCLUDED.nickname
	`, ownerID, contactUserID, nickname)
	return err
}

func SetBlocked(ctx context.Context, db *pgxpool.Pool, ownerID, contactUserID string, blocked bool) error {
	_, err := db.Exec(ctx, `
		UPDATE contacts SET is_blocked = $3 WHERE owner_id = $1::uuid AND contact_user_id = $2::uuid
	`, ownerID, contactUserID, blocked)
	return err
}

func Remove(ctx context.Context, db *pgxpool.Pool, ownerID, contactUserID string) error {
	_, err := db.Exec(ctx, `
		DELETE FROM contacts WHERE owner_id = $1::uuid AND contact_user_id = $2::uuid
	`, ownerID, contactUserID)
	return err
}
