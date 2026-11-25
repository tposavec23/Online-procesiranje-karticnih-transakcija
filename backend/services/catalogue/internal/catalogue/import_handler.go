package catalogue

import (
	"errors"
	"log"
	"net/http"
)

const maxImportBytes = 1 << 20 // 1 MB is plenty for a catalogue CSV

// ImportHandler translates the import request into article service calls.
type ImportHandler struct {
	service *ArticleService
}

func NewImportHandler(service *ArticleService) *ImportHandler {
	return &ImportHandler{service: service}
}

func (h *ImportHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/articles/import", h.importArticles)
}

// importArticles reads a CSV from the request body. With ?dryRun=true it only
// checks the file and saves nothing.
func (h *ImportHandler) importArticles(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxImportBytes)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing CSV file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var result ImportResult
	if r.URL.Query().Get("dryRun") == "true" {
		result, err = h.service.previewImport(file)
	} else {
		result, err = h.service.importArticles(r.Context(), file)
	}

	switch {
	case errors.Is(err, errInvalidHeader):
		http.Error(w, "invalid CSV header", http.StatusBadRequest)
	case err != nil:
		log.Printf("import articles: %v", err)
		http.Error(w, "failed to import articles", http.StatusInternalServerError)
	default:
		writeJSON(w, result)
	}
}
