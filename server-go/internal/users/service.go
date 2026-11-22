package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserSummary struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	AvatarURL  string `json:"avatar_url"`
	Department string `json:"department"`
	Program    string `json:"program"`
	Level      string `json:"level"`
	IsOnline   bool   `json:"is_online"`
	IsVerified bool   `json:"is_verified"`
}

func Search(ctx context.Context, db *pgxpool.Pool, query, excludeUserID string) ([]UserSummary, error) {
	rows, err := db.Query(ctx, `
		SELECT u.id::text, u.email, r.name, p.first_name, p.last_name,
		       COALESCE(p.avatar_url, ''), COALESCE(d.code, ''), COALESCE(pr.code, ''), COALESCE(al.code, ''),
		       u.is_online, u.is_verified
		FROM users u
		JOIN roles r ON r.id = u.role_id
		JOIN profiles p ON p.user_id = u.id
		LEFT JOIN departments d ON d.id = p.department_id
		LEFT JOIN programs pr ON pr.id = p.program_id
		LEFT JOIN academic_levels al ON al.id = p.academic_level_id
		WHERE u.is_active = TRUE AND u.is_suspended = FALSE
		  AND u.id::text <> $1
		  AND ($2 = '' OR u.email ILIKE '%'||$2||'%' OR p.first_name ILIKE '%'||$2||'%' OR p.last_name ILIKE '%'||$2||'%')
		ORDER BY p.last_name ASC, p.first_name ASC
		LIMIT 50
	`, excludeUserID, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []UserSummary
	for rows.Next() {
		var u UserSummary
		if err := rows.Scan(&u.ID, &u.Email, &u.Role, &u.FirstName, &u.LastName, &u.AvatarURL, &u.Department, &u.Program, &u.Level, &u.IsOnline, &u.IsVerified); err == nil {
			list = append(list, u)
		}
	}
	return list, nil
}
