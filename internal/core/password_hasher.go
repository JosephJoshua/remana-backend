package core

import (
	"errors"
	"fmt"

	"github.com/JosephJoshua/remana-backend/internal/shared/apperror"
	"golang.org/x/crypto/bcrypt"
)

// We're using a cost of 14 here, which is leaning more towards the secure side,
// even if it will compromise a bit on performance.
const (
	bcryptCost = 14
)

type PasswordHasher struct{}

func (p *PasswordHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return "", apperror.ErrPasswordTooLong
	}

	if err != nil {
		return "", fmt.Errorf("error generating hash from password: %w", err)
	}

	return string(hashedPassword), nil
}

func (p *PasswordHasher) Check(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, bcrypt.ErrHashTooShort) {
		return apperror.ErrPasswordMismatch
	}

	if err != nil {
		return fmt.Errorf("error comparing hash and password: %w", err)
	}

	return nil
}
