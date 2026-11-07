package catalogue

import (
	"context"
	"fmt"
)

// ArticleService fetches catalogue articles.
type ArticleService struct {
	repo *ArticleRepository
}

func NewArticleService(repo *ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) fetchAll(ctx context.Context, limit, offset int) ([]Article, error) {
	articles, err := s.repo.fetchAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("fetch all articles: %w", err)
	}
	if articles == nil {
		articles = []Article{} // JSON is [] instead of null when there are no articles
	}
	return articles, nil
}
