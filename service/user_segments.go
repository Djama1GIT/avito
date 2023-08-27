package service

import (
	"avito/repository"
	"avito/structures"
)

type UserSegmentsService struct {
	repo repository.UserSegments
}

func NewUserSegmentsService(repo repository.UserSegments) *UserSegmentsService {
	return &UserSegmentsService{repo: repo}
}

func (s *UserSegmentsService) Create(Segment structures.Segment, user structures.User) (string, error) {
	return s.repo.Create(Segment, user)
}
