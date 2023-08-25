package repository

type Segment interface {
}

type Repository struct {
	Segment
}

func NewRepository() *Repository {
	return &Repository{}
}
