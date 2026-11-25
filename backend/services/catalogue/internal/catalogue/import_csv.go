package catalogue

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

// Reasons why a CSV row is skipped. They are shown to the user as
// "redak N: <reason>".
const (
	reasonInvalidPrice     = "neispravna cijena"
	reasonInvalidQuantity  = "neispravna količina"
	reasonMissingField     = "nedostaje obavezno polje"
	reasonWrongColumnCount = "neispravan broj stupaca"
	reasonMalformedRow     = "neispravan format retka"
)

const csvColumnCount = 6

// csvHeader is the exact header every import file must start with.
var csvHeader = [csvColumnCount]string{"code", "name", "description", "category", "price", "quantity"}

var errInvalidHeader = errors.New("invalid CSV header")

// RowError describes one skipped CSV row. Row counts data rows only,
// so the first row after the header is row 1.
type RowError struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
}

// parseArticlesCSV reads an import file and returns the valid articles and
// one RowError for every row that was skipped. It returns errInvalidHeader
// when the header is missing or wrong, and then no rows are read.
func parseArticlesCSV(r io.Reader) ([]Article, []RowError, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1 // the column count is checked per row
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil || !isValidHeader(header) {
		return nil, nil, errInvalidHeader
	}

	articles := []Article{}
	rowErrors := []RowError{}

	for row := 1; ; row++ {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			rowErrors = append(rowErrors, RowError{Row: row, Reason: reasonMalformedRow})
			continue
		}

		article, reason := parseArticleRecord(record)
		if reason != "" {
			rowErrors = append(rowErrors, RowError{Row: row, Reason: reason})
			continue
		}
		articles = append(articles, article)
	}

	return articles, rowErrors, nil
}

func isValidHeader(header []string) bool {
	if len(header) != csvColumnCount {
		return false
	}
	for i, name := range header {
		name = strings.TrimPrefix(name, "\ufeff") // BOM written by Excel
		if strings.ToLower(strings.TrimSpace(name)) != csvHeader[i] {
			return false
		}
	}
	return true
}

func parseArticleRecord(record []string) (Article, string) {
	if len(record) != csvColumnCount {
		return Article{}, reasonWrongColumnCount
	}

	code := strings.TrimSpace(record[0])
	name := strings.TrimSpace(record[1])
	description := strings.TrimSpace(record[2])
	category := strings.TrimSpace(record[3])

	if code == "" || name == "" || category == "" {
		return Article{}, reasonMissingField
	}

	price, err := parsePriceCents(record[4])
	if err != nil || price <= 0 {
		return Article{}, reasonInvalidPrice
	}

	quantity, err := parseQuantity(record[5])
	if err != nil {
		return Article{}, reasonInvalidQuantity
	}

	return Article{
		Code:              code,
		Name:              name,
		Description:       description,
		Category:          category,
		Price:             price,
		AvailableQuantity: quantity,
	}, ""
}

// parseQuantity reads a whole number >= 0. An empty value means 0.
func parseQuantity(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0, errors.New("invalid quantity")
	}
	return n, nil
}

func parsePriceCents(s string) (int64, error) {
	s = strings.ReplaceAll(strings.TrimSpace(s), ",", ".")

	negative := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")

	whole, fraction, _ := strings.Cut(s, ".")
	if whole == "" || !isDigits(whole) || len(whole) > 8 {
		return 0, errors.New("invalid price")
	}
	if len(fraction) > 2 || (fraction != "" && !isDigits(fraction)) {
		return 0, errors.New("invalid price")
	}
	fraction += strings.Repeat("0", 2-len(fraction)) // "5" -> "50"

	cents, err := strconv.ParseInt(whole+fraction, 10, 64)
	if err != nil {
		return 0, errors.New("invalid price")
	}
	if negative {
		cents = -cents
	}
	return cents, nil
}

func isDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return s != ""
}
