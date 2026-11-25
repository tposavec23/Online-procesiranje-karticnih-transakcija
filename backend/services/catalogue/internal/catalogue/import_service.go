package catalogue

import (
	"context"
	"fmt"
	"io"
)

// ImportResult is what the user gets back after an import. Imported counts
// the articles that were saved (new or updated); Errors lists skipped rows.
type ImportResult struct {
	Imported int        `json:"imported"`
	Errors   []RowError `json:"errors"`
}

// importArticles reads a CSV file and saves every valid row. Invalid rows
// are skipped and reported, they do not stop the import.
func (s *ArticleService) importArticles(ctx context.Context, r io.Reader) (ImportResult, error) {
	articles, rowErrors, err := parseArticlesCSV(r)
	if err != nil {
		return ImportResult{}, err
	}

	if len(articles) > 0 {
		if err := s.repo.insertMany(ctx, articles); err != nil {
			return ImportResult{}, fmt.Errorf("insert articles: %w", err)
		}
	}

	return ImportResult{Imported: len(articles), Errors: rowErrors}, nil
}

// previewImport runs all checks of an import but saves nothing.
func (s *ArticleService) previewImport(r io.Reader) (ImportResult, error) {
	articles, rowErrors, err := parseArticlesCSV(r)
	if err != nil {
		return ImportResult{}, err
	}
	return ImportResult{Imported: len(articles), Errors: rowErrors}, nil
}
