package service

import (
	"avito/pkg/repository"
	"avito/pkg/structures"
)

type SegmentService struct {
	repo repository.Segment
}

func NewSegmentService(repo repository.Segment) *SegmentService {
	return &SegmentService{repo: repo}
}

func (s *SegmentService) Create(segment structures.Segment) (string, error) {
	return s.repo.Create(segment)
}

func (s *SegmentService) Delete(segment structures.Segment) (string, error) {
	return s.repo.Delete(segment)
}
