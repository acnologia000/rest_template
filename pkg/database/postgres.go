package database

import (
	"database/sql"
	"rest_template/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetPostgresConnectionString())
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
