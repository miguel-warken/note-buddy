package sqlite

import (
	_ "modernc.org/sqlite"
	"database/sql"
	"path/filepath"
	"testing"
)

func TestOpenRunsMigrations(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "notes.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("An error occurred when opening database: %s", err)
	}
	t.Cleanup(func() {_ = db.Close()})
	
	assertTableIsCreated(t, db, "draft_sessions")
	assertTableIsCreated(t, db, "draft_entries")
	assertTableIsCreated(t, db, "sessions")
	assertTableIsCreated(t, db, "session_entries")
}

func assertTableIsCreated(t *testing.T, db *sql.DB, tableName string) {
	t.Helper()
	
	var exists int
	err := 
		db.QueryRow(`
			SELECT 1
			FROM sqlite_master
			WHERE type = 'table'
				AND name = ?
			LIMIT 1
		`, tableName).Scan(&exists)
	if err != nil {
		t.Fatalf("Table was not created: %s", tableName)
	}
}