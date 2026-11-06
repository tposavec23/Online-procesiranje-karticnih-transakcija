package article

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Article, error) //sto repozitorij mora moci raditi
}

type PostgresRepository struct {
	db *pgxpool.Pool //za konekciju sa bazom
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db} //stvaramo novi repozitorij sa konekcijom na bazu
}

func (r *PostgresRepository) FindAll(ctx context.Context) ([]Article, error) {
	rows, err := r.db.Query(ctx, `
		SELECT artikl_id,
		       sifra,
		       naziv,
		       COALESCE(opis, ''),
		       cijena::float8,
		       kolicina_na_zalihi
		FROM catalogue.artikl
		ORDER BY sifra`)

	if err != nil {
		return nil, fmt.Errorf("query articles: %w", err)
	}
	defer rows.Close()

	articles := []Article{} //Prazna lista artikala, JSON postaje [] umjesto null ako nema artikala
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Code, &a.Name, &a.Description, &a.Price, &a.AvailableQuantity); err != nil {
			return nil, fmt.Errorf("scan article: %w", err)
		}
		articles = append(articles, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate articles: %w", err)
	}
	return articles, nil
}
