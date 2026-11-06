package article

// JSON imena odgovaraju Kotlin Article modelu u Android aplikaciji.

type Article struct {
	ID                int     `json:"id"`
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Price             float64 `json:"price"`
	AvailableQuantity int     `json:"availableQuantity"`
}
