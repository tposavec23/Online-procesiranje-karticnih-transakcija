package catalogue

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errUnknownUser        = errors.New("unknown user")
	errArticleUnavailable = errors.New("article is not available")
)

// CartRepository reads and writes carts in PostgreSQL.
type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) addItem(ctx context.Context, userID, articleID, quantity int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var stock int
	err = tx.QueryRow(ctx, `
		SELECT available_quantity
		FROM catalogue.article
		WHERE article_id = $1
		FOR UPDATE`, articleID).Scan(&stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return errArticleNotFound
	}
	if err != nil {
		return fmt.Errorf("check article %d: %w", articleID, err)
	}
	if stock < 1 {
		return errArticleUnavailable
	}

	var cartID int
	err = tx.QueryRow(ctx, `
		INSERT INTO catalogue.cart (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) WHERE status = 'ACTIVE'
		DO UPDATE SET updated_at = CURRENT_TIMESTAMP
		RETURNING cart_id`, userID).Scan(&cartID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign key violation
		return errUnknownUser
	}
	if err != nil {
		return fmt.Errorf("get cart: %w", err)
	}

	_, err = tx.Exec(ctx, `
    INSERT INTO catalogue.cart_item (cart_id, article_id, quantity, unit_price)
    SELECT $1::int, article_id, $3::int, price
    FROM catalogue.article
    WHERE article_id = $2
    ON CONFLICT (cart_id, article_id)
    DO UPDATE SET quantity = catalogue.cart_item.quantity + EXCLUDED.quantity`,
		cartID, articleID, quantity)
	if err != nil {
		return fmt.Errorf("save cart item: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// fetchActive returns the user's active cart, or an empty cart if there is none.
func (r *CartRepository) fetchActive(ctx context.Context, userID int) (Cart, error) {
	cart := Cart{Items: []CartItem{}}

	err := r.db.QueryRow(ctx, `
		SELECT cart_id
		FROM catalogue.cart
		WHERE user_id = $1 AND status = 'ACTIVE'`, userID).Scan(&cart.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		return cart, nil
	}
	if err != nil {
		return Cart{}, fmt.Errorf("query cart: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT a.article_id, a.code, a.name, ci.quantity, (ci.unit_price * 100)::bigint
		FROM catalogue.cart_item ci
		JOIN catalogue.article a ON a.article_id = ci.article_id
		WHERE ci.cart_id = $1
		ORDER BY ci.added_at, ci.cart_item_id`, cart.ID)
	if err != nil {
		return Cart{}, fmt.Errorf("query cart items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item CartItem
		if err := rows.Scan(&item.ArticleID, &item.Code, &item.Name, &item.Quantity, &item.UnitPrice); err != nil {
			return Cart{}, fmt.Errorf("scan cart item: %w", err)
		}
		cart.Items = append(cart.Items, item)
	}
	if err := rows.Err(); err != nil {
		return Cart{}, fmt.Errorf("iterate cart items: %w", err)
	}

	return cart, nil
}
