package domain

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	UsernameMinLength = 5
	UsernameMaxLength = 100
)

type User interface {
	SetUsername(username string) error
	SetPassword(password string)
	SetRole(role Role) error
	ID() uuid.UUID
	Username() string
	Password() string
	StoreCode() string
	Role() Role
}

type user struct {
	id        uuid.UUID
	username  string
	password  string
	storeCode string
	role      Role
}

func NewUser(id uuid.UUID, username string, password string, store Store, role Role) (User, error) {
	user := new(user)

	user.id = id
	user.storeCode = store.Code()

	user.SetPassword(password)

	if err := user.SetRole(role); err != nil {
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	if err := user.SetUsername(username); err != nil {
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	return user, nil
}

func (u *user) SetUsername(username string) error {
	if len(username) > UsernameMaxLength {
		return fmt.Errorf("error setting username of user: %w", ErrInputTooLong)
	}

	if len(username) < UsernameMinLength {
		return fmt.Errorf("error setting username of user: %w", ErrInputTooShort)
	}

	u.username = username
	return nil
}

func (u *user) SetPassword(password string) {
	u.password = password
}

func (u *user) SetRole(role Role) error {
	if role.StoreCode() != u.StoreCode() {
		return fmt.Errorf("error setting role of user: %w", ErrInvalidStoreCode)
	}

	u.role = role
	return nil
}

func (u *user) ID() uuid.UUID {
	return u.id
}

func (u *user) Username() string {
	return u.username
}

func (u *user) Password() string {
	return u.password
}

func (u *user) StoreCode() string {
	return u.storeCode
}

func (u *user) Role() Role {
	return u.role
}
