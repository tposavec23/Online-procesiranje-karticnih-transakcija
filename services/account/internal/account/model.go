package account

// Role odgovara enumu Role u Android aplikaciji.
type Role string

const (
	RoleMerchant Role = "MERCHANT"
	RoleAdmin    Role = "ADMIN"
)

// User je korisnik koji se može prijaviti u sustav.
// JSON imena odgovaraju Kotlin modelu User u Android aplikaciji.
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Role         Role   `json:"role"`
	Active       bool   `json:"active"`
	PasswordHash string `json:"-"` // hash lozinke se nikad ne šalje u JSON odgovoru
}
