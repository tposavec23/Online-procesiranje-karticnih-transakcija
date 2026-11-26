package catalogue

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// CartHandler translates HTTP requests into cart service calls.
type CartHandler struct {
	service *CartService
}

func NewCartHandler(service *CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/cart", h.getCart)
	mux.HandleFunc("POST /api/cart/items", h.addItem)
}

type addItemRequest struct {
	ArticleID int `json:"articleId"`
	Quantity  int `json:"quantity"`
}

func (h *CartHandler) getCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	cart, err := h.service.fetchCart(r.Context(), userID)
	if err != nil {
		log.Printf("get cart: %v", err)
		http.Error(w, "failed to fetch cart", http.StatusInternalServerError)
		return
	}
	writeJSON(w, cart)
}

func (h *CartHandler) addItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(w, r)
	if !ok {
		return
	}

	var req addItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ArticleID < 1 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Quantity == 0 {
		req.Quantity = 1 // adding without a quantity adds one piece
	}

	cart, err := h.service.addItem(r.Context(), userID, req.ArticleID, req.Quantity)
	switch {
	case errors.Is(err, errInvalidQuantity):
		http.Error(w, "invalid quantity", http.StatusBadRequest)
	case errors.Is(err, errArticleNotFound):
		http.Error(w, "article not found", http.StatusNotFound)
	case errors.Is(err, errArticleUnavailable):
		http.Error(w, "article is not available", http.StatusConflict)
	case errors.Is(err, errStockExceeded):
		http.Error(w, "requested quantity exceeds available stock", http.StatusConflict)
	case err != nil:
		log.Printf("add item to cart: %v", err)
		http.Error(w, "failed to add item to cart", http.StatusInternalServerError)
	default:
		writeJSON(w, cart)
	}
}

// The login service returns no token, so for now the user id comes from a header.
func userIDFromRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.Header.Get("X-User-Id"))
	if err != nil || id < 1 {
		http.Error(w, "missing or invalid X-User-Id header", http.StatusUnauthorized)
		return 0, false
	}
	return id, true
}
