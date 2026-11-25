package catalogue

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

var exportColumns = []string{"code", "name", "description", "category", "price", "availableQuantity"}

func writeExport(w io.Writer, articles []Article) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(exportColumns); err != nil {
		return fmt.Errorf("write export header: %w", err)
	}

	for _, a := range articles {
		row := []string{a.Code, a.Name, a.Description, a.Category, formatPrice(a.Price), strconv.Itoa(a.AvailableQuantity)}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write article %s: %w", a.Code, err)
		}
	}

	cw.Flush()
	if err := cw.Error(); err != nil {
		return fmt.Errorf("flush export: %w", err)
	}
	return nil
}

func formatPrice(cents int64) string {
	return fmt.Sprintf("%d.%02d", cents/100, cents%100)
}
