package readmodel

import "github.com/google/uuid"

type AuthnUser struct {
	ID           uuid.UUID
	Password     string
	IsStoreAdmin bool
}

type UserDetailsRole struct {
	ID           uuid.UUID
	Name         string
	IsStoreAdmin bool
}

type UserDetailsStore struct {
	ID   uuid.UUID
	Name string
	Code string
}

type UserDetails struct {
	ID       uuid.UUID
	Username string
	Role     UserDetailsRole
	Store    UserDetailsStore
}
