package catalogue

import (
	"context"
	"errors"
	"fmt"
)

var errInvalidQuantity = errors.New("quantity must be at least 1")

// CartService holds the cart rules.
type CartService struct {
	repo *CartRepository
}

func NewCartService(repo *CartRepository) *CartService {
	return &CartService{repo: repo}
}

// addItem adds an article to the user's cart and returns the updated cart.
func (s *CartService) addItem(ctx context.Context, userID, articleID, quantity int) (Cart, error) {
	if quantity < 1 {
		return Cart{}, errInvalidQuantity
	}
	if err := s.repo.addItem(ctx, userID, articleID, quantity); err != nil {
		return Cart{}, fmt.Errorf("add item to cart: %w", err)
	}
	return s.fetchCart(ctx, userID)
}

// fetchCart returns the active cart with item and cart totals (in cents).
func (s *CartService) fetchCart(ctx context.Context, userID int) (Cart, error) {
	cart, err := s.repo.fetchByUser(ctx, userID)
	if err != nil {
		return Cart{}, fmt.Errorf("fetch cart: %w", err)
	}
	for i := range cart.Items {
		item := &cart.Items[i]
		item.TotalPrice = item.UnitPrice * int64(item.Quantity)
		cart.TotalPrice += item.TotalPrice
	}
	return cart, nil
}
