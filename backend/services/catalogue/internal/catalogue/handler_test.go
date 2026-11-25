package catalogue

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newImportRequest builds a multipart POST request with the CSV in the "file" field.
func newImportRequest(t *testing.T, csv string) *http.Request {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "artikli.csv")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	part.Write([]byte(csv))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/articles/import", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

// newTestServer wires the real handler and service to a fake store.
func newTestServer(store *fakeStore) *http.ServeMux {
	mux := http.NewServeMux()
	NewHandler(NewArticleService(store)).RegisterRoutes(mux)
	return mux
}

// Acceptance criterion 2 through the HTTP API: the client gets 7 imported and an error for row 4.
func TestImportHandler_NegativePriceInRow4_ReturnsSummary(t *testing.T) {
	store := &fakeStore{}
	rows := validRows(8)
	rows[3] = "N-004;Kolač 4;Opis kolača;Kolači;-2,00;5"
	rec := httptest.NewRecorder()

	newTestServer(store).ServeHTTP(rec, newImportRequest(t, buildCSV(rows)))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var result ImportResult
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.Imported != 7 {
		t.Errorf("expected imported=7, got %d", result.Imported)
	}
	want := RowError{Row: 4, Reason: reasonInvalidPrice}
	if len(result.Errors) != 1 || result.Errors[0] != want {
		t.Errorf("expected errors [%+v], got %+v", want, result.Errors)
	}
}

func TestImportHandler_MissingFile_Returns400(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/articles/import", nil)
	rec := httptest.NewRecorder()

	newTestServer(&fakeStore{}).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestImportHandler_InvalidHeader_Returns400(t *testing.T) {
	rec := httptest.NewRecorder()

	newTestServer(&fakeStore{}).ServeHTTP(rec, newImportRequest(t, "code;name;price\nN-001;Kava;1,50"))

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}
