package repository

import "github.com/jmoiron/sqlx"

type Segment interface {
}

type Repository struct {
	Segment
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
