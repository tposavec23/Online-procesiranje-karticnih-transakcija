package catalogue

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const (
	defaultLimit = 10
	maxLimit     = 100
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
	mux.HandleFunc("GET /api/articles/{id}", h.getArticle)
}

func (h *Handler) getArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := queryInt(r, "limit", defaultLimit)
	if err != nil || limit < 1 {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset, err := queryInt(r, "offset", 0)
	if err != nil || offset < 0 {
		http.Error(w, "invalid offset", http.StatusBadRequest)
		return
	}

	articles, err := h.service.fetchAll(r.Context(), limit, offset)
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

func (h *Handler) getArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.Error(w, "invalid article id", http.StatusBadRequest)
		return
	}

	article, err := h.service.fetchByID(r.Context(), id)
	if err != nil {
		log.Printf("get article %d: %v", id, err)
		http.Error(w, "failed to fetch article", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(article); err != nil {
		log.Printf("encode article: %v", err)
	}
}

// queryInt reads an integer query parameter, returning fallback when it is absent.
func queryInt(r *http.Request, name string, fallback int) (int, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return fallback, nil
	}
	return strconv.Atoi(value)
}
