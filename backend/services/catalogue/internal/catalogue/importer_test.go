package catalogue

import (
	"fmt"
	"strings"
	"testing"
)

const testCSVHeader = "code;name;description;category;price;quantity"

// buildCSV joins the header and the given data rows into CSV text.
func buildCSV(rows []string) string {
	return testCSVHeader + "\n" + strings.Join(rows, "\n")
}

// validRows returns n valid CSV rows with new codes N-001, N-002, ...
func validRows(n int) []string {
	rows := make([]string, n)
	for i := range rows {
		rows[i] = fmt.Sprintf("N-%03d;Kolač %d;Opis kolača;Kolači;%d,50;5", i+1, i+1, i+1)
	}
	return rows
}

func TestParseArticlesCSV_AllValid_ReturnsAllArticles(t *testing.T) {
	articles, rowErrors, err := parseArticlesCSV(strings.NewReader(buildCSV(validRows(8))))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(articles) != 8 {
		t.Errorf("expected 8 articles, got %d", len(articles))
	}
	if len(rowErrors) != 0 {
		t.Errorf("expected no row errors, got %v", rowErrors)
	}
}

func TestParseArticlesCSV_NegativePriceInRow4_SkipsRow4(t *testing.T) {
	rows := validRows(8)
	rows[3] = "N-004;Kolač 4;Opis kolača;Kolači;-2,00;5" // row 4

	articles, rowErrors, err := parseArticlesCSV(strings.NewReader(buildCSV(rows)))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(articles) != 7 {
		t.Errorf("expected 7 articles, got %d", len(articles))
	}
	want := []RowError{{Row: 4, Reason: reasonInvalidPrice}}
	if len(rowErrors) != 1 || rowErrors[0] != want[0] {
		t.Errorf("expected %v, got %v", want, rowErrors)
	}
}

func TestParseArticlesCSV_ParsesAllFields(t *testing.T) {
	csv := buildCSV([]string{"N-001;Kremšnita;Domaća;Kolači;3,20;4"})

	articles, _, err := parseArticlesCSV(strings.NewReader(csv))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := Article{
		Code:              "N-001",
		Name:              "Kremšnita",
		Description:       "Domaća",
		Category:          "Kolači",
		Price:             320,
		AvailableQuantity: 4,
	}
	if len(articles) != 1 || articles[0] != want {
		t.Errorf("expected %+v, got %+v", want, articles)
	}
}

func TestParseArticlesCSV_InvalidHeader_ReturnsError(t *testing.T) {
	_, _, err := parseArticlesCSV(strings.NewReader("code;name;price\nN-001;Kava;1,50"))

	if err != errInvalidHeader {
		t.Errorf("expected errInvalidHeader, got %v", err)
	}
}

func TestParsePriceCents(t *testing.T) {
	tests := []struct {
		input   string
		want    int64
		wantErr bool
	}{
		{"4,50", 450, false},
		{"4.50", 450, false},
		{"3", 300, false},
		{"0,5", 50, false},
		{"-2,00", -200, false}, // parses, but parseArticleRecord rejects prices <= 0
		{"abc", 0, true},
		{"", 0, true},
		{"1,505", 0, true}, // more than two decimals
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parsePriceCents(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("parsePriceCents(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parsePriceCents(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseArticleRecord_SkipReasons(t *testing.T) {
	tests := []struct {
		name       string
		row        string
		wantReason string
	}{
		{"valid row", "N-001;Kava;;Napici;1,50;", ""},
		{"zero price", "N-001;Kava;;Napici;0;1", reasonInvalidPrice},
		{"negative price", "N-001;Kava;;Napici;-1;1", reasonInvalidPrice},
		{"text price", "N-001;Kava;;Napici;skupo;1", reasonInvalidPrice},
		{"missing code", ";Kava;;Napici;1,50;1", reasonMissingField},
		{"missing name", "N-001;;;Napici;1,50;1", reasonMissingField},
		{"missing category", "N-001;Kava;;;1,50;1", reasonMissingField},
		{"negative quantity", "N-001;Kava;;Napici;1,50;-3", reasonInvalidQuantity},
		{"too few columns", "N-001;Kava;1,50", reasonWrongColumnCount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, reason := parseArticleRecord(strings.Split(tt.row, ";"))

			if reason != tt.wantReason {
				t.Errorf("row %q: got reason %q, want %q", tt.row, reason, tt.wantReason)
			}
		})
	}
}
