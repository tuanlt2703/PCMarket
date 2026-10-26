package infrastructure

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"product-service/config"
)

func NewSQLite(cfg *config.Config) (*sql.DB, error) {

	db, err := sql.Open("sqlite", cfg.DBName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// SQLite works best with a single writer connection.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return db, nil
}
