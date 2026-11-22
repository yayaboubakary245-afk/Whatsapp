package administration

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditEntry struct {
	ID         string `json:"id"`
	ActorID    string `json:"actor_id"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	IPAddress  string `json:"ip_address"`
	CreatedAt  string `json:"created_at"`
}

func SetUserSuspension(ctx context.Context, db *pgxpool.Pool, userID string, suspended bool) error {
	_, err := db.Exec(ctx, `UPDATE users SET is_suspended = $2, updated_at = NOW() WHERE id = $1::uuid`, userID, suspended)
	return err
}

func UpdateUserRole(ctx context.Context, db *pgxpool.Pool, userID, roleName string) error {
	_, err := db.Exec(ctx, `
		UPDATE users SET role_id = (SELECT id FROM roles WHERE name = $2), updated_at = NOW()
		WHERE id = $1::uuid
	`, userID, roleName)
	return err
}

func ListAuditLogs(ctx context.Context, db *pgxpool.Pool) ([]AuditEntry, error) {
	rows, err := db.Query(ctx, `
		SELECT id::text, COALESCE(actor_id::text, ''), action, target_type, COALESCE(target_id, ''), COALESCE(ip_address, ''),
		       to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT 100
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if rows.Scan(&e.ID, &e.ActorID, &e.Action, &e.TargetType, &e.TargetID, &e.IPAddress, &e.CreatedAt) == nil {
			list = append(list, e)
		}
	}
	return list, nil
}
