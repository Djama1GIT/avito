package service

import (
	"avito/repository"
)

type Segment interface {
}

type Service struct {
	Segment
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
