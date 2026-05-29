package sqlite

import (
	"database/sql"
	"note-buddy-main/internal/domain"
)

func GetLastSession(db *sql.DB) (*domain.Session, error){
	query := `
		SELECT
			s.id,
			s.subject,
			s.title,
			s.started_at,
			s.status
		FROM sessions s
		ORDER BY created_at DESC
		LIMIT 1
	`
	var s domain.Session
	err := db.QueryRow(query).Scan(&s.ID, &s.Subject, &s.Title, &s.StartedAt, &s.Status)

	return &s, err
}