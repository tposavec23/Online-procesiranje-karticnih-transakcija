package catalogue

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler translates HTTP requests into service calls.
type Handler struct {
	service *ArticleService
}

func NewHandler(service *ArticleService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/articles", h.getArticles)
}

func (h *Handler) getArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.service.fetchAll(r.Context())
	if err != nil {
		log.Printf("get articles: %v", err)
		http.Error(w, "failed to fetch articles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		log.Printf("encode articles: %v", err)
	}
}
