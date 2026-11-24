package catalogue

type CartItem struct {
	ArticleID  int    `json:"articleId"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	UnitPrice  int64  `json:"unitPrice"`
	TotalPrice int64  `json:"totalPrice"`
}

type Cart struct {
	ID         int        `json:"id"`
	Items      []CartItem `json:"items"`
	TotalPrice int64      `json:"totalPrice"`
}
