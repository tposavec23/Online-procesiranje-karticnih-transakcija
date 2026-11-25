package catalogue

import (
	"errors"
	"log"
	"net/http"
)

const maxImportBytes = 1 << 20 // 1 MB is plenty for a catalogue CSV

// ImportHandler translates the import request into parser and repository calls.
type ImportHandler struct {
	repo *ImportRepository
}

func NewImportHandler(repo *ImportRepository) *ImportHandler {
	return &ImportHandler{repo: repo}
}

func (h *ImportHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/articles/import", h.importArticles)
}

type importResponse struct {
	Imported int        `json:"imported"`
	Updated  int        `json:"updated"`
	Errors   []RowError `json:"errors"`
}

// importArticles reads a CSV from the request body. With ?dryRun=true it only
// checks the file and saves nothing.
func (h *ImportHandler) importArticles(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, maxImportBytes)

	articles, rowErrors, err := parseArticlesCSV(body)
	if errors.Is(err, errInvalidHeader) {
		http.Error(w, "invalid CSV header", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("parse import: %v", err)
		http.Error(w, "failed to read file", http.StatusBadRequest)
		return
	}

	response := importResponse{Errors: rowErrors}

	if r.URL.Query().Get("dryRun") == "true" {
		response.Imported = len(articles) // how many would be saved
		writeJSON(w, response)
		return
	}

	result, err := h.repo.saveArticles(r.Context(), articles)
	if err != nil {
		log.Printf("save import: %v", err)
		http.Error(w, "failed to save articles", http.StatusInternalServerError)
		return
	}

	response.Imported = result.Imported
	response.Updated = result.Updated
	writeJSON(w, response)
}
