package service

import (
	"avito/repository"
	"avito/structures"
)

type Segment interface {
	Create(segment structures.Segment) (string, error)
}

type UserSegments interface {
	Create(segment structures.Segment, user structures.User) (string, error)
}

type Service struct {
	Segment
	UserSegments
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Segment:      NewSegmentService(repos.Segment),
		UserSegments: NewUserSegmentsService(repos.UserSegments),
	}
}
