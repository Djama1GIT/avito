package service

import (
	"avito/pkg/repository"
	"avito/pkg/structures"
)

type Segment interface {
	Create(segment structures.Segment) (string, error)
	Delete(segment structures.Segment) (string, error)
}

type UserSegments interface {
	Patch(userSegments structures.UserSegments) (int, error)
	GetUsersInSegment(user structures.User) ([]string, error)
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