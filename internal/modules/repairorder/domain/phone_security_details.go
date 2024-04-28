package domain

import (
	"fmt"
	"unicode"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
)

type phoneSecurityType string

const (
	PhoneSecurityTypeNone     = phoneSecurityType("none")
	PhoneSecurityTypePasscode = phoneSecurityType("passcode")
	PhoneSecurityTypePattern  = phoneSecurityType("pattern")
)

type PhoneSecurityDetails interface {
	Value() string
	Type() phoneSecurityType
}

type phoneSecurityDetails struct {
	value        string
	securityType phoneSecurityType
}

func NewNoSecurity() PhoneSecurityDetails {
	return phoneSecurityDetails{
		value:        "",
		securityType: PhoneSecurityTypeNone,
	}
}

func NewPasscodeSecurity(value string) PhoneSecurityDetails {
	return phoneSecurityDetails{
		value:        value,
		securityType: PhoneSecurityTypePasscode,
	}
}

func NewPatternSecurity(value string) (PhoneSecurityDetails, error) {
	for _, c := range value {
		if !unicode.IsDigit(c) {
			return nil, fmt.Errorf("%w: value contains non-digit", apperror.ErrInvalidInput)
		}
	}

	return &phoneSecurityDetails{
		value:        value,
		securityType: PhoneSecurityTypePattern,
	}, nil
}

func (p phoneSecurityDetails) Value() string {
	return p.value
}

func (p phoneSecurityDetails) Type() phoneSecurityType {
	return p.securityType
}
