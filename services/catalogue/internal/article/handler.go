package article

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler je HTTP handler za artikle. On koristi Service za dohvat podataka i vraća ih kao JSON.
type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registrira rute za artikle na dani ServeMux. Trenutno registrira samo GET /api/articles.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/articles", h.getArticles)
}

func (h *Handler) getArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.service.GetAll(r.Context())
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
