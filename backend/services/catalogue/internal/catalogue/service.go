package catalogue

import (
	"context"
	"fmt"
)

type ArticleStore interface {
	fetchAll(ctx context.Context, limit, offset int) ([]Article, error)
	fetchByID(ctx context.Context, id int) (Article, error)
	fetchAllForExport(ctx context.Context) ([]Article, error)
	insertMany(ctx context.Context, articles []Article) error
}

// ArticleService fetches catalogue articles.
type ArticleService struct {
	repo ArticleStore
}

func NewArticleService(repo ArticleStore) *ArticleService {
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

func (s *ArticleService) fetchByID(ctx context.Context, id int) (Article, error) {
	article, err := s.repo.fetchByID(ctx, id)
	if err != nil {
		return Article{}, fmt.Errorf("fetch article: %w", err)
	}
	return article, nil
}

func (s *ArticleService) fetchAllForExport(ctx context.Context) ([]Article, error) {
	articles, err := s.repo.fetchAllForExport(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch articles for export: %w", err)
	}
	return articles, nil
}
