package catalogue

import (
	"net/http/httptest"
	"testing"
)

func TestServeExport_headers(t *testing.T) {
	rec := httptest.NewRecorder()

	serveExport(rec, nil)

	if got := rec.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", got, "text/csv; charset=utf-8")
	}
	if got := rec.Header().Get("Content-Disposition"); got != `attachment; filename="artikli.csv"` {
		t.Errorf("Content-Disposition = %q, want %q", got, `attachment; filename="artikli.csv"`)
	}
}

func TestServeExport_body(t *testing.T) {
	rec := httptest.NewRecorder()
	articles := []Article{
		{Code: "A-1", Name: "Kolač", Category: "Kolači", Price: 400, AvailableQuantity: 10},
		{Code: "A-2", Name: "Torta", Category: "Torte", Price: 2500, AvailableQuantity: 2},
	}

	serveExport(rec, articles)

	got := rec.Body.String()
	want := exportHeader +
		"A-1,Kolač,,Kolači,4.00,10\n" +
		"A-2,Torta,,Torte,25.00,2\n"
	if got != want {
		t.Errorf("body =\n%s\nwant\n%s", got, want)
	}
}
