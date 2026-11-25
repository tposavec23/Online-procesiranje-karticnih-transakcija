package catalogue

import (
	"fmt"
	"strings"
	"testing"
)

const exportHeader = "code,name,description,category,price,availableQuantity\n"

func TestWriteExport_rows(t *testing.T) {
	var articles []Article
	for i := 1; i <= 12; i++ {
		articles = append(articles, Article{Code: fmt.Sprintf("A-%d", i), Name: "Kolač", Category: "Kolači", Price: 400, AvailableQuantity: 10})
	}

	var b strings.Builder
	if err := writeExport(&b, articles); err != nil {
		t.Fatalf("writeExport() error = %v", err)
	}
	if lines := strings.Count(b.String(), "\n"); lines != 13 {
		t.Errorf("writeExport() wrote %d rows, want 13 (header + 12 articles)", lines)
	}
}

func TestWriteExport_noArticles(t *testing.T) {
	var b strings.Builder
	if err := writeExport(&b, nil); err != nil {
		t.Fatalf("writeExport() error = %v", err)
	}
	if got := b.String(); got != exportHeader {
		t.Errorf("writeExport() = %q, want only the header %q", got, exportHeader)
	}
}

func TestWriteExport_format(t *testing.T) {
	tests := []struct {
		name    string
		article Article
		wantRow string
	}{
		{
			name:    "simple article",
			article: Article{Code: "A-1", Name: "Kolač", Description: "Čokoladni", Category: "Kolači", Price: 400, AvailableQuantity: 10},
			wantRow: "A-1,Kolač,Čokoladni,Kolači,4.00,10\n",
		},
		{
			name:    "free article",
			article: Article{Code: "A-2", Name: "Uzorak", Category: "Kolači", Price: 0, AvailableQuantity: 0},
			wantRow: "A-2,Uzorak,,Kolači,0.00,0\n",
		},
		{
			name:    "price under ten cents",
			article: Article{Code: "A-3", Name: "Bombon", Category: "Slatkiši", Price: 5, AvailableQuantity: 100},
			wantRow: "A-3,Bombon,,Slatkiši,0.05,100\n",
		},
		{
			name:    "description with a comma",
			article: Article{Code: "A-6", Name: "Kolač", Description: "Opis, s zarezom", Category: "Kolači", Price: 400, AvailableQuantity: 10},
			wantRow: "A-6,Kolač,\"Opis, s zarezom\",Kolači,4.00,10\n",
		},
		{
			name:    "description with quotes",
			article: Article{Code: "A-7", Name: "Kolač", Description: "Kolač \"Mamin\"", Category: "Kolači", Price: 400, AvailableQuantity: 10},
			wantRow: "A-7,Kolač,\"Kolač \"\"Mamin\"\"\",Kolači,4.00,10\n",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			err := writeExport(&b, []Article{tc.article})
			if err != nil {
				t.Fatalf("writeExport() error = %v", err)
			}

			got := b.String()
			want := exportHeader + tc.wantRow
			if got != want {
				t.Errorf("writeExport() =\n%s\nwant\n%s", got, want)
			}
		})
	}
}
