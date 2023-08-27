package repository

import (
	"avito/structures"

	"github.com/jmoiron/sqlx"
)

type Segment interface {
	Create(segment structures.Segment) (string, error)
}

type UserSegments interface {
	Create(segment structures.Segment, user structures.User) (string, error)
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
