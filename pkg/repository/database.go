package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	usersTable               = "users"
	segmentsTable            = "segments"
	userSegmentsTable        = "user_segments"
	userSegmentsHistoryTable = "user_segments_history"
)

type Config struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

func NewDB(cfg Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := sql.Open(cfg.Driver, connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
