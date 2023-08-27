package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable        = "users"
	segmentsTable     = "segments"
	userSegmentsTable = "user_segments"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("dbname=%s host=%s port=%s user=%s password=%s sslmode=%s",
		cfg.Name, cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
