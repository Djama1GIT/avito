package service

import (
	"avito/pkg/repository"
	"avito/pkg/structures"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserHistory(userHistory structures.UserHistory) (string, error) {
	return s.repo.GetUserHistory(userHistory)
}
