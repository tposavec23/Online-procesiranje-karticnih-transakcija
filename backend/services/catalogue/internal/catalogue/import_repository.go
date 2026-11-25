package catalogue

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// importResult counts what a saved import did.
type importResult struct {
	Imported int // new articles
	Updated  int // existing articles (same code) that were changed
}

// ImportRepository saves imported articles to PostgreSQL.
type ImportRepository struct {
	db *pgxpool.Pool
}

func NewImportRepository(db *pgxpool.Pool) *ImportRepository {
	return &ImportRepository{db: db}
}

// saveArticles stores all articles in one transaction: either every article
// is saved or none is. An article whose code already exists is updated, and
// a category that does not exist yet is created.
func (r *ImportRepository) saveArticles(ctx context.Context, articles []Article) (importResult, error) {
	var result importResult

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return result, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	categoryIDs := map[string]int{} // category name -> id, so each is looked up once

	for _, a := range articles {
		categoryID, ok := categoryIDs[a.Category]
		if !ok {
			categoryID, err = findOrCreateCategory(ctx, tx, a.Category)
			if err != nil {
				return result, err
			}
			categoryIDs[a.Category] = categoryID
		}

		var inserted bool
		err = tx.QueryRow(ctx, `
			INSERT INTO catalogue.article (category_id, code, name, description, price, available_quantity)
			VALUES ($1, $2, $3, NULLIF($4, ''), ($5::bigint)::numeric / 100, $6)
			ON CONFLICT (code) DO UPDATE SET
				category_id        = EXCLUDED.category_id,
				name               = EXCLUDED.name,
				description        = EXCLUDED.description,
				price              = EXCLUDED.price,
				available_quantity = EXCLUDED.available_quantity,
				updated_at         = CURRENT_TIMESTAMP
			RETURNING (xmax = 0)`, // xmax = 0 only for a freshly inserted row
			categoryID, a.Code, a.Name, a.Description, a.Price, a.AvailableQuantity).Scan(&inserted)
		if err != nil {
			return result, fmt.Errorf("save article %s: %w", a.Code, err)
		}

		if inserted {
			result.Imported++
		} else {
			result.Updated++
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit: %w", err)
	}
	return result, nil
}

func findOrCreateCategory(ctx context.Context, tx pgx.Tx, name string) (int, error) {
	var id int
	err := tx.QueryRow(ctx,
		`SELECT category_id FROM catalogue.category WHERE name = $1 ORDER BY category_id LIMIT 1`,
		name).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("find category %q: %w", name, err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO catalogue.category (name) VALUES ($1) RETURNING category_id`,
		name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create category %q: %w", name, err)
	}
	return id, nil
}
