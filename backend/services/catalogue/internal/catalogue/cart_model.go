package catalogue

// CartItem is one article in a cart. Prices are in cents.
type CartItem struct {
	ArticleID  int    `json:"articleId"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	UnitPrice  int64  `json:"unitPrice"`
	TotalPrice int64  `json:"totalPrice"`
}

// Cart is the cart of one user. TotalPrice is in cents.
type Cart struct {
	Items      []CartItem `json:"items"`
	TotalPrice int64      `json:"totalPrice"`
}
