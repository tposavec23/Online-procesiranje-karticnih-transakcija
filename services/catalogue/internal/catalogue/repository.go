package catalogue

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errArticleNotFound = errors.New("article not found")

const articleQuery = `
	SELECT a.article_id,
	       a.code,
	       a.name,
	       COALESCE(a.description, ''),
	       c.name,
	       (a.price * 100)::bigint,
	       a.available_quantity
	FROM catalogue.article a
	JOIN catalogue.category c ON c.category_id = a.category_id`

// ArticleRepository reads articles from PostgreSQL.
type ArticleRepository struct {
	db *pgxpool.Pool
}

func NewArticleRepository(db *pgxpool.Pool) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) fetchAll(ctx context.Context, limit, offset int) ([]Article, error) {
	rows, err := r.db.Query(ctx, articleQuery+` ORDER BY a.code LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query articles: %w", err)
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		a, err := scanArticle(rows)
		if err != nil {
			return nil, fmt.Errorf("scan article: %w", err)
		}
		articles = append(articles, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate articles: %w", err)
	}

	return articles, nil
}

func (r *ArticleRepository) fetchByID(ctx context.Context, id int) (Article, error) {
	a, err := scanArticle(r.db.QueryRow(ctx, articleQuery+` WHERE a.article_id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return Article{}, errArticleNotFound
	}
	if err != nil {
		return Article{}, fmt.Errorf("query article %d: %w", id, err)
	}
	return a, nil
}

func scanArticle(row pgx.Row) (Article, error) {
	var a Article
	err := row.Scan(&a.ID, &a.Code, &a.Name, &a.Description, &a.Category, &a.Price, &a.AvailableQuantity)
	return a, err
}
