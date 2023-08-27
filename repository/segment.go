package repository

import (
	"avito/structures"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SegmentDB struct {
	db *sqlx.DB
}

func NewSegmentDB(db *sqlx.DB) *SegmentDB {
	return &SegmentDB{db: db}
}

func (r *SegmentDB) Create(segment structures.Segment) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	var slug string
	createSegmentQuery := fmt.Sprintf("INSERT INTO %s (slug) VALUES ($1) RETURNING slug", segmentsTable)
	row := tx.QueryRow(createSegmentQuery, segment.Slug)
	if err := row.Scan(&slug); err != nil {
		tx.Rollback()
		return "", err
	}

	return slug, tx.Commit()
}
