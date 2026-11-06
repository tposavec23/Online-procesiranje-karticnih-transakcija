package account

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// errInvalidCredentials se vraća kada korisničko ime ili lozinka nisu ispravni.
// Ista greška se vraća i za nepostojećeg korisnika, kako odgovor ne bi otkrio postoji li račun.
var errInvalidCredentials = errors.New("invalid username or password")

// AuthService provjerava korisničke podatke pri prijavi.
type AuthService struct {
	repo *UserRepository
}

func NewAuthService(repo *UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// login vraća korisnika ako su korisničko ime i lozinka ispravni.
func (s *AuthService) login(ctx context.Context, username, password string) (User, error) {
	u, err := s.repo.findByUsername(ctx, username)
	if errors.Is(err, errUserNotFound) {
		return User{}, errInvalidCredentials
	}
	if err != nil {
		return User{}, fmt.Errorf("login: %w", err)
	}

	// bcrypt uspoređuje upisanu lozinku s hashom iz baze, bez spremanja lozinke u čitljivom obliku.
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return User{}, errInvalidCredentials
	}

	return u, nil
}
