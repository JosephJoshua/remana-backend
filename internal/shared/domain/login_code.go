package domain

import (
	"fmt"
	"unicode"

	"github.com/google/uuid"
)

const (
	LoginCodeLength = 8
)

type LoginCode interface {
	UserID() uuid.UUID
	Code() string
}

type loginCode struct {
	userID uuid.UUID
	code   string
}

func NewLoginCode(user User, code string) (LoginCode, error) {
	loginCode := new(loginCode)

	loginCode.userID = user.ID()

	if err := loginCode.setCode(code); err != nil {
		return nil, fmt.Errorf("error creating login code: %w", err)
	}

	return loginCode, nil
}

func (l *loginCode) UserID() uuid.UUID {
	return l.userID
}

func (l *loginCode) Code() string {
	return l.code
}

func (l *loginCode) setCode(code string) error {
	if len(code) < LoginCodeLength {
		return fmt.Errorf("error setting code of login code: %w", ErrInputTooShort)
	}

	if len(code) > LoginCodeLength {
		return fmt.Errorf("error setting code of login code: %w", ErrInputTooLong)
	}

	newCode := []rune(code)

	for i, c := range code {
		if unicode.IsDigit(c) {
			continue
		}

		if unicode.IsLetter(c) {
			if unicode.IsLower(c) {
				newCode[i] = unicode.ToUpper(c)
			}

			continue
		}

		return fmt.Errorf("error setting code of login code: %w", ErrInvalidLoginCode)
	}

	l.code = string(newCode)
	return nil
}
