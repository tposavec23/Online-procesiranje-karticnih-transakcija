package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// errUserNotFound se vraća kada u bazi ne postoji korisnik s traženim korisničkim imenom.
var errUserNotFound = errors.New("user not found")

// UserRepository dohvaća korisnike iz baze podataka.
type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// findByUsername dohvaća korisnika prema korisničkom imenu.
func (r *UserRepository) findByUsername(ctx context.Context, username string) (User, error) {
	var u User
	err := r.db.QueryRow(ctx, `
		SELECT id, username, password_hash, role, active
		FROM account.users
		WHERE username = $1`, username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Active)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, errUserNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("find user by username: %w", err)
	}
	return u, nil
}
