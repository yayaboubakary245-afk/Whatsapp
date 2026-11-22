package groups

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ValidGroupTypes = map[string]bool{
	"classe":      true,
	"filiere":     true,
	"niveau":      true,
	"departement": true,
	"td":          true,
	"projet":      true,
	"pedagogique": true,
}

type Group struct {
	ID                string `json:"id"`
	ConversationID    string `json:"conversation_id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	GroupType         string `json:"group_type"`
	OnlyAdminsCanPost bool   `json:"only_admins_can_post"`
	MemberCount       int    `json:"member_count"`
	CreatedAt         string `json:"created_at"`
}

func Create(ctx context.Context, db *pgxpool.Pool, creatorID, name, description, groupType string, onlyAdmins bool, memberIDs []string) (*Group, error) {
	if !ValidGroupTypes[groupType] {
		return nil, fmt.Errorf("type de groupe invalide (%s) : classe, filiere, niveau, departement, td, projet, pedagogique", groupType)
	}
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var convID string
	if err := tx.QueryRow(ctx, `INSERT INTO conversations (type, created_by) VALUES ('group', $1::uuid) RETURNING id::text`, creatorID).Scan(&convID); err != nil {
		return nil, err
	}

	var g Group
	err = tx.QueryRow(ctx, `
		INSERT INTO groups (conversation_id, name, description, group_type, only_admins_can_post, created_by)
		VALUES ($1::uuid, $2, $3, $4, $5, $6::uuid)
		RETURNING id::text, conversation_id::text, name, description, group_type, only_admins_can_post, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`, convID, name, description, groupType, onlyAdmins, creatorID).Scan(
		&g.ID, &g.ConversationID, &g.Name, &g.Description, &g.GroupType, &g.OnlyAdminsCanPost, &g.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	_, _ = tx.Exec(ctx, `INSERT INTO group_members (group_id, user_id, member_role) VALUES ($1::uuid, $2::uuid, 'admin')`, g.ID, creatorID)
	_, _ = tx.Exec(ctx, `INSERT INTO conversation_members (conversation_id, user_id) VALUES ($1::uuid, $2::uuid)`, convID, creatorID)

	for _, uid := range memberIDs {
		if uid == "" || uid == creatorID {
			continue
		}
		_, _ = tx.Exec(ctx, `INSERT INTO group_members (group_id, user_id, member_role) VALUES ($1::uuid, $2::uuid, 'member') ON CONFLICT DO NOTHING`, g.ID, uid)
		_, _ = tx.Exec(ctx, `INSERT INTO conversation_members (conversation_id, user_id) VALUES ($1::uuid, $2::uuid) ON CONFLICT DO NOTHING`, convID, uid)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	g.MemberCount = 1 + len(memberIDs)
	return &g, nil
}

func ListForUser(ctx context.Context, db *pgxpool.Pool, userID string) ([]Group, error) {
	rows, err := db.Query(ctx, `
		SELECT g.id::text, g.conversation_id::text, g.name, COALESCE(g.description, ''), g.group_type, g.only_admins_can_post,
		       (SELECT COUNT(*) FROM group_members gm2 WHERE gm2.group_id = g.id),
		       to_char(g.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM groups g
		JOIN group_members gm ON gm.group_id = g.id AND gm.user_id = $1::uuid
		ORDER BY g.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Group
	for rows.Next() {
		var g Group
		if rows.Scan(&g.ID, &g.ConversationID, &g.Name, &g.Description, &g.GroupType, &g.OnlyAdminsCanPost, &g.MemberCount, &g.CreatedAt) == nil {
			list = append(list, g)
		}
	}
	return list, nil
}
