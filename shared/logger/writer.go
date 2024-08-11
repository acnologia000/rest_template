package logger

import (
	"database/sql"
	"time"
)

type DBWriter struct {
	DB *sql.DB
}

func (w *DBWriter) Write(p []byte) (n int, err error) {
	_, err = w.DB.Exec("INSERT INTO logs (message, created_at) VALUES ($1, $2)", string(p), time.Now())
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
