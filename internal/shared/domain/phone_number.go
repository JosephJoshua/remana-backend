package domain

import (
	"errors"
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

type PhoneNumber interface {
	Value() string
}

type phoneNumber struct {
	value string
}

func NewPhoneNumber(value string) (PhoneNumber, error) {
	num, err := phonenumbers.Parse(value, "ID")
	if err != nil {
		return nil, fmt.Errorf("failed to parse phone number: %w", err)
	}

	if !phonenumbers.IsValidNumber(num) {
		return nil, errors.New("invalid phone number")
	}

	return phoneNumber{value: phonenumbers.Format(num, phonenumbers.E164)}, nil
}

func (p phoneNumber) Value() string {
	return p.value
}
