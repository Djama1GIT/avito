package repository

import (
	"avito/structures"

	"github.com/jmoiron/sqlx"
)

type UserSegmentsDB struct {
	db *sqlx.DB
}

func NewUserSegmentsDB(db *sqlx.DB) *UserSegmentsDB {
	return &UserSegmentsDB{db: db}
}

func (r *UserSegmentsDB) Create(segment structures.Segment, user structures.User) (string, error) {
	return "", nil
}
