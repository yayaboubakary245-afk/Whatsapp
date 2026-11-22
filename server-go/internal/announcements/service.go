package announcements

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Announcement struct {
	ID          string `json:"id"`
	AuthorID    string `json:"author_id"`
	AuthorName  string `json:"author_name"`
	AuthorRole  string `json:"author_role"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	IsPinned    bool   `json:"is_pinned"`
	PublishedAt string `json:"published_at"`
}

func Publish(ctx context.Context, db *pgxpool.Pool, authorID, authorRole, title, content, category string, isPinned bool) (*Announcement, error) {
	if authorRole != "admin" && authorRole != "enseignant" {
		return nil, errors.New("seuls les enseignants et administrateurs peuvent publier des annonces")
	}
	var a Announcement
	err := db.QueryRow(ctx, `
		INSERT INTO announcements (author_id, title, content, category, is_pinned)
		VALUES ($1::uuid, $2, $3, $4, $5)
		RETURNING id::text, author_id::text, title, content, category, is_pinned, to_char(published_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`, authorID, title, content, category, isPinned).Scan(
		&a.ID, &a.AuthorID, &a.Title, &a.Content, &a.Category, &a.IsPinned, &a.PublishedAt,
	)
	a.AuthorRole = authorRole
	return &a, err
}

func List(ctx context.Context, db *pgxpool.Pool) ([]Announcement, error) {
	rows, err := db.Query(ctx, `
		SELECT a.id::text, a.author_id::text, p.first_name || ' ' || p.last_name, r.name,
		       a.title, a.content, a.category, a.is_pinned,
		       to_char(a.published_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM announcements a
		JOIN users u ON u.id = a.author_id
		JOIN roles r ON r.id = u.role_id
		JOIN profiles p ON p.user_id = u.id
		ORDER BY a.is_pinned DESC, a.published_at DESC
		LIMIT 100
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Announcement
	for rows.Next() {
		var a Announcement
		if rows.Scan(&a.ID, &a.AuthorID, &a.AuthorName, &a.AuthorRole, &a.Title, &a.Content, &a.Category, &a.IsPinned, &a.PublishedAt) == nil {
			list = append(list, a)
		}
	}
	return list, nil
}
