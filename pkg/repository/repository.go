package repository

import (
	"avito/pkg/structures"

	"github.com/jmoiron/sqlx"
)

type Segment interface {
	Create(segment structures.Segment) (string, error)
	Delete(segment structures.Segment) (string, error)
}

type UserSegments interface {
	Patch(userSegments structures.UserSegments) (int, error)
	GetUsersInSegment(user structures.User) ([]string, error)
}

type Repository struct {
	Segment
	UserSegments
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Segment:      NewSegmentDB(db),
		UserSegments: NewUserSegmentsDB(db),
	}
}