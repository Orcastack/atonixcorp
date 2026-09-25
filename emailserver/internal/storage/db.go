package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// OpenPostgresDB opens a PostgreSQL database connection using the DSN
func OpenPostgresDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
