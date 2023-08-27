package service

import (
	"avito/pkg/repository"
	"avito/pkg/structures"
)

type UserSegmentsService struct {
	repo repository.UserSegments
}

func NewUserSegmentsService(repo repository.UserSegments) *UserSegmentsService {
	return &UserSegmentsService{repo: repo}
}

func (s *UserSegmentsService) Patch(userSegments structures.UserSegments) (int, error) {
	return s.repo.Patch(userSegments)
}

func (s *UserSegmentsService) GetUsersInSegment(user structures.User) ([]string, error) {
	return s.repo.GetUsersInSegment(user)
}
