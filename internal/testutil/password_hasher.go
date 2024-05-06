package testutil

import "github.com/JosephJoshua/remana-backend/internal/apperror"

type PasswordHasherStub struct{}

func (p PasswordHasherStub) Hash(password string) (string, error) {
	return password, nil
}

func (p PasswordHasherStub) Check(hashedPassword, password string) error {
	if hashedPassword != password {
		return apperror.ErrPasswordMismatch
	}

	return nil
}
