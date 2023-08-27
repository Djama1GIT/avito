package service

import (
	"avito/repository"
	"avito/structures"
)

type SegmentService struct {
	repo repository.Segment
}

func NewSegmentService(repo repository.Segment) *SegmentService {
	return &SegmentService{repo: repo}
}

func (s *SegmentService) Create(Segment structures.Segment) (string, error) {
	return s.repo.Create(Segment)
}

func (s *SegmentService) Delete(Segment structures.Segment) (string, error) {
	return s.repo.Delete(Segment)
}
