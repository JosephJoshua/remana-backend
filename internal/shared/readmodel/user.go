package readmodel

import "github.com/google/uuid"

type AuthnUser struct {
	ID           uuid.UUID
	Password     string
	IsStoreAdmin bool
}
