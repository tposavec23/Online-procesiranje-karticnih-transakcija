package account

type Role string

const (
	RoleMerchant Role = "MERCHANT"
	RoleAdmin    Role = "ADMIN"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Role         Role   `json:"role"`
	Active       bool   `json:"active"`
	PasswordHash string `json:"-"`
}
