package catalogue

// Article is a single item in the catalogue.
type Article struct {
	ID                int    `json:"id"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Price             int64  `json:"price"`
	AvailableQuantity int    `json:"availableQuantity"`
}
