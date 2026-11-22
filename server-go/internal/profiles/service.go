package profiles

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
	UserID          string `json:"user_id"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	AvatarURL       string `json:"avatar_url"`
	Bio             string `json:"bio"`
	DepartmentID    *int   `json:"department_id,omitempty"`
	DepartmentName  string `json:"department_name,omitempty"`
	ProgramID       *int   `json:"program_id,omitempty"`
	ProgramName     string `json:"program_name,omitempty"`
	AcademicLevelID *int   `json:"academic_level_id,omitempty"`
	AcademicLevel   string `json:"academic_level,omitempty"`
	Matricule       string `json:"matricule,omitempty"`
}

func GetByUserID(ctx context.Context, db *pgxpool.Pool, userID string) (*Profile, error) {
	var p Profile
	err := db.QueryRow(ctx, `
		SELECT u.id::text, u.email, r.name, p.first_name, p.last_name,
		       COALESCE(p.avatar_url, ''), COALESCE(p.bio, ''),
		       p.department_id, COALESCE(d.name, ''),
		       p.program_id, COALESCE(pr.name, ''),
		       p.academic_level_id, COALESCE(al.label, ''),
		       COALESCE(p.matricule, '')
		FROM users u
		JOIN roles r ON r.id = u.role_id
		JOIN profiles p ON p.user_id = u.id
		LEFT JOIN departments d ON d.id = p.department_id
		LEFT JOIN programs pr ON pr.id = p.program_id
		LEFT JOIN academic_levels al ON al.id = p.academic_level_id
		WHERE u.id = $1::uuid
	`, userID).Scan(
		&p.UserID, &p.Email, &p.Role, &p.FirstName, &p.LastName,
		&p.AvatarURL, &p.Bio,
		&p.DepartmentID, &p.DepartmentName,
		&p.ProgramID, &p.ProgramName,
		&p.AcademicLevelID, &p.AcademicLevel,
		&p.Matricule,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func Update(ctx context.Context, db *pgxpool.Pool, userID string, firstName, lastName, bio, avatarURL string, deptID, progID, levelID *int) error {
	_, err := db.Exec(ctx, `
		UPDATE profiles
		SET first_name = COALESCE(NULLIF($2, ''), first_name),
		    last_name  = COALESCE(NULLIF($3, ''), last_name),
		    bio        = $4,
		    avatar_url = COALESCE(NULLIF($5, ''), avatar_url),
		    department_id = COALESCE($6, department_id),
		    program_id    = COALESCE($7, program_id),
		    academic_level_id = COALESCE($8, academic_level_id),
		    updated_at = NOW()
		WHERE user_id = $1::uuid
	`, userID, firstName, lastName, bio, avatarURL, deptID, progID, levelID)
	return err
}
