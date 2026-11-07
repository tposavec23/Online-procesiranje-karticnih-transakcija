package account

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Returned for both an unknown username and a wrong password,
// so the response does not reveal whether the account exists.
var errInvalidCredentials = errors.New("invalid username or password")

type AuthService struct {
	repo *UserRepository
}

func NewAuthService(repo *UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) login(ctx context.Context, username, password string) (User, error) {
	u, err := s.repo.findByUsername(ctx, username)
	if errors.Is(err, errUserNotFound) {
		return User{}, errInvalidCredentials
	}
	if err != nil {
		return User{}, fmt.Errorf("login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return User{}, errInvalidCredentials
	}

	if !u.Active {
		return User{}, errInvalidCredentials
	}

	return u, nil
}
