package sqlite

import (
	"database/sql"
	_ "embed"
	_ "modernc.org/sqlite"
)

//go:embed migrations/001_init.sql
var initSQL string

//go:embed migrations/002_seed_mock.sql
var seedMockSQL string

func Open(path string, seed bool) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil	, err
	}

	if _, err = db.Exec(initSQL); err != nil {
		db.Close() 
		return nil, err
	}

	if seed {
		if _, err = db.Exec(seedMockSQL); err != nil {
			db.Close()
			return nil, err
		}
	}
	return db, nil
}