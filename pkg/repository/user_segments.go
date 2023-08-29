package repository

import (
	"avito/pkg/structures"
	"database/sql"
	"fmt"
)

type UserSegmentsDB struct {
	db *sql.DB
}

func NewUserSegmentsDB(db *sql.DB) *UserSegmentsDB {
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
			return -1, fmt.Errorf("error occurred while processing segment to add '%s': %v", segment, err)
		}
		_, err = historyUpdate(tx, segment, userSegments.UserId, true)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	for _, segment := range userSegments.SegmentsToDelete {
		existsQuery := fmt.Sprintf("SELECT count(*) FROM %s WHERE user_id = $1 AND segment = $2", userSegmentsTable)
		var count int
		err := r.db.QueryRow(existsQuery, userSegments.UserId, segment).Scan(&count)
		if err != nil {
			tx.Rollback()
			return -1, fmt.Errorf("error occurred while checking segment to delete existence '%s': %v", segment, err)
		}

		if count < 1 {
			tx.Rollback()
			return -1, fmt.Errorf("error occurred while checking segment to delete existence '%s': user(%d) is not in this segment", segment, userSegments.UserId)
		}

		deleteSegmentQuery := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND segment = $2", userSegmentsTable)
		_, err = tx.Exec(deleteSegmentQuery, userSegments.UserId, segment)
		if err != nil {
			tx.Rollback()
			return -1, fmt.Errorf("error occurred while processing segment to delete '%s': %v", segment, err)
		}

		_, err = historyUpdate(tx, segment, userSegments.UserId, false)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	return userSegments.UserId, tx.Commit()
}

func (r *UserSegmentsDB) GetUserSegments(user structures.User) ([]string, error) {
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

func (r *UserSegmentsDB) GetSegmentUsers(segment structures.Segment) ([]int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var users []int
	createSegmentQuery := fmt.Sprintf("SELECT user_id FROM %s WHERE segment = $1", userSegmentsTable)
	rows, err := tx.Query(createSegmentQuery, segment.Slug)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user int
		if err := rows.Scan(&user); err != nil {
			tx.Rollback()
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return users, nil
}
