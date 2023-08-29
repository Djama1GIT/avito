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

type Repository struct {
	Segment      Segment
	UserSegments UserSegments
}

func NewRepository(db *sql.DB) *Repository {
	segmentDB := NewSegmentDB(db)
	userSegmentsDB := NewUserSegmentsDB(db)

	return &Repository{
		Segment:      segmentDB,
		UserSegments: userSegmentsDB,
	}
}
