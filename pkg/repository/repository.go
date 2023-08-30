package repository

import (
	"avito/pkg/structures"
	"database/sql"
)

type Segment interface {
	Create(segment structures.Segment) (string, error)
	Delete(segment structures.Segment) (string, error)
}

type UserSegments interface {
	Patch(userSegments structures.UserSegments) (int, error)
	GetUserSegments(user structures.User) ([]string, error)
	GetSegmentUsers(segment structures.Segment) ([]int, error)
}

type User interface {
	GetUserHistory(userHistory structures.UserHistory) (string, error)
	DeleteExpiredSegments() error
}

type Repository struct {
	Segment      Segment
	UserSegments UserSegments
	User         User
}

func NewRepository(db *sql.DB) *Repository {
	segmentDB := NewSegmentDB(db)
	userSegmentsDB := NewUserSegmentsDB(db)
	userDB := NewUserDB(db)

	return &Repository{
		Segment:      segmentDB,
		UserSegments: userSegmentsDB,
		User:         userDB,
	}
}
