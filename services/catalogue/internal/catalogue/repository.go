package catalogue

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ArticleRepository reads articles from PostgreSQL.
type ArticleRepository struct {
	db *pgxpool.Pool
}

func NewArticleRepository(db *pgxpool.Pool) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) fetchAll(ctx context.Context, limit, offset int) ([]Article, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.article_id,
		       a.code,
		       a.name,
		       COALESCE(a.description, ''),
		       c.name,
		       (a.price * 100)::bigint,
		       a.available_quantity
		FROM catalogue.article a
		JOIN catalogue.category c ON c.category_id = a.category_id
		ORDER BY a.code
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query articles: %w", err)
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Code, &a.Name, &a.Description, &a.Category, &a.Price, &a.AvailableQuantity); err != nil {
			return nil, fmt.Errorf("scan article: %w", err)
		}
		articles = append(articles, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate articles: %w", err)
	}

	return articles, nil
}
