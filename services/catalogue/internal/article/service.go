package article

import "context"

// Radimo samo uslugu za dohvat svih artikala, jer je to jedina funkcionalnost koja je potrebna za US-01.
type Service interface {
	GetAll(ctx context.Context) ([]Article, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context) ([]Article, error) {
	return s.repo.FindAll(ctx)
}
