package catalogue

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errArticleUnavailable = errors.New("article is not available")
	errStockExceeded      = errors.New("requested quantity exceeds available stock")
)

// CartRepository reads and writes carts in PostgreSQL.
type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

// addItem adds an article to the user's cart and creates the cart if the user
// has none. The article row is locked while the transaction runs, so the stock
// cannot change between the checks and the insert.
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
		ON CONFLICT (user_id)
		DO UPDATE SET updated_at = CURRENT_TIMESTAMP
		RETURNING cart_id`, userID).Scan(&cartID)
	if err != nil {
		return fmt.Errorf("get cart: %w", err)
	}

	var inCart int
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(SUM(quantity), 0)
		FROM catalogue.cart_item
		WHERE cart_id = $1 AND article_id = $2`, cartID, articleID).Scan(&inCart)
	if err != nil {
		return fmt.Errorf("check cart quantity: %w", err)
	}
	if inCart+quantity > stock {
		return errStockExceeded
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO catalogue.cart_item (cart_id, article_id, quantity)
		VALUES ($1, $2, $3)
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

// fetchByUser returns the user's cart, or an empty cart if there is none.
// Prices are read from the article table, so they are always the current ones.
func (r *CartRepository) fetchByUser(ctx context.Context, userID int) (Cart, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.article_id, a.code, a.name, ci.quantity, (a.price * 100)::bigint
		FROM catalogue.cart c
		JOIN catalogue.cart_item ci ON ci.cart_id = c.cart_id
		JOIN catalogue.article a ON a.article_id = ci.article_id
		WHERE c.user_id = $1
		ORDER BY ci.added_at, ci.cart_item_id`, userID)
	if err != nil {
		return Cart{}, fmt.Errorf("query cart items: %w", err)
	}
	defer rows.Close()

	cart := Cart{Items: []CartItem{}}
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
