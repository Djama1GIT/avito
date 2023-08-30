package repository

import (
	"avito/pkg/structures"
	"database/sql"
	"fmt"
)

type SegmentDB struct {
	db *sql.DB
}

func NewSegmentDB(db *sql.DB) *SegmentDB {
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

func (r *SegmentDB) Delete(segment structures.Segment) (string, error) {
	existsQuery := fmt.Sprintf("SELECT COUNT(slug) FROM %s WHERE slug = $1", segmentsTable)
	var count int
	err := r.db.QueryRow(existsQuery, segment.Slug).Scan(&count)
	if err != nil {
		return "", err
	}

	if count == 0 {
		return "", fmt.Errorf("segment with slug %s does not exist", segment.Slug)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	deleteSegmentQuery := fmt.Sprintf("DELETE FROM %s WHERE slug = $1", segmentsTable)
	_, err = tx.Exec(deleteSegmentQuery, segment.Slug)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	repo := NewRepository(r.db)
	user_ids, err := repo.UserSegments.GetSegmentUsers(segment)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	for _, user_id := range user_ids {
		_, err := historyUpdate(tx, segment.Slug, user_id, false)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	return segment.Slug, tx.Commit()
}
