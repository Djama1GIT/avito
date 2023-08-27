package repository

import (
	"avito/pkg/structures"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserSegmentsDB struct {
	db *sqlx.DB
}

func NewUserSegmentsDB(db *sqlx.DB) *UserSegmentsDB {
	return &UserSegmentsDB{db: db}
}

func (r *UserSegmentsDB) Patch(userSegments structures.UserSegments) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	for _, segment := range userSegments.SegmentsToAdd {
		createSegmentQuery := fmt.Sprintf("INSERT INTO %s (user_id, segment) VALUES ($1, $2)", userSegmentsTable)
		_, err = tx.Exec(createSegmentQuery, userSegments.UserId, segment)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	for _, segment := range userSegments.SegmentsToDelete {
		deleteSegmentQuery := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND segment = $2", userSegmentsTable)
		_, err = tx.Exec(deleteSegmentQuery, userSegments.UserId, segment)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	return userSegments.UserId, tx.Commit()
}

func (r *UserSegmentsDB) GetUsersInSegment(user structures.User) ([]string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var slugs []string
	createSegmentQuery := fmt.Sprintf("SELECT segment FROM %s WHERE user_id = $1", userSegmentsTable)
	rows, err := tx.Query(createSegmentQuery, user.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			tx.Rollback()
			return nil, err
		}
		slugs = append(slugs, slug)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return slugs, nil
}
