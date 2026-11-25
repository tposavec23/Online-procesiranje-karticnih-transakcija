package catalogue

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// fakeStore replaces the PostgreSQL repository in tests and records saved articles.
type fakeStore struct {
	saved     []Article
	insertErr error
}

func (f *fakeStore) fetchAll(_ context.Context, _, _ int) ([]Article, error) {
	return f.saved, nil
}

func (f *fakeStore) fetchByID(_ context.Context, id int) (Article, error) {
	for _, a := range f.saved {
		if a.ID == id {
			return a, nil
		}
	}
	return Article{}, errArticleNotFound
}

func (f *fakeStore) insertMany(_ context.Context, articles []Article) error {
	if f.insertErr != nil {
		return f.insertErr
	}
	f.saved = append(f.saved, articles...)
	return nil
}

// Acceptance criterion 1:
// Given a CSV with 8 articles with new codes, when I import it, then all 8 are saved.
func TestImportArticles_EightNewArticles_AllSaved(t *testing.T) {
	store := &fakeStore{}
	svc := NewArticleService(store)

	result, err := svc.importArticles(context.Background(), strings.NewReader(buildCSV(validRows(8))))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.saved) != 8 {
		t.Errorf("expected 8 saved articles, got %d", len(store.saved))
	}
	if result.Imported != 8 || len(result.Errors) != 0 {
		t.Errorf("expected 8 imported and no errors, got %+v", result)
	}
}

// Acceptance criterion 2:
// Given a CSV with 8 articles where row 4 has a negative price, when I import it,
// then 7 articles are saved and row 4 is reported as having an invalid price.
func TestImportArticles_NegativePriceInRow4_SevenSavedAndRow4Reported(t *testing.T) {
	store := &fakeStore{}
	svc := NewArticleService(store)
	rows := validRows(8)
	rows[3] = "N-004;Kolač 4;Opis kolača;Kolači;-2,00;5"

	result, err := svc.importArticles(context.Background(), strings.NewReader(buildCSV(rows)))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.saved) != 7 {
		t.Errorf("expected 7 saved articles, got %d", len(store.saved))
	}
	for _, a := range store.saved {
		if a.Code == "N-004" {
			t.Errorf("article N-004 with a negative price must not be saved")
		}
	}
	want := RowError{Row: 4, Reason: reasonInvalidPrice}
	if len(result.Errors) != 1 || result.Errors[0] != want {
		t.Errorf("expected errors [%+v], got %+v", want, result.Errors)
	}
}

func TestImportArticles_DatabaseError_ReturnsError(t *testing.T) {
	store := &fakeStore{insertErr: errors.New("connection lost")}
	svc := NewArticleService(store)

	_, err := svc.importArticles(context.Background(), strings.NewReader(buildCSV(validRows(2))))

	if err == nil {
		t.Fatal("expected an error when the database insert fails")
	}
}
