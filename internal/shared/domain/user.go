package domain

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	UsernameMinLength = 5
	UsernameMaxLength = 100
)

type User struct {
	id        uuid.UUID
	username  string
	password  string
	storeCode string
	role      Role
}

func NewUser(id uuid.UUID, username string, password string, store Store, role Role) (*User, error) {
	user := new(User)

	user.id = id

	user.SetStore(store)
	user.SetRole(role)
	user.SetPassword(password)

	if err := user.SetUsername(username); err != nil {
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	return user, nil
}

func (u *User) SetUsername(username string) error {
	if len(username) > UsernameMaxLength {
		return fmt.Errorf("error setting username of user: %w", ErrInputTooLong)
	}

	if len(username) < UsernameMinLength {
		return fmt.Errorf("error setting username of user: %w", ErrInputTooShort)
	}

	u.username = username
	return nil
}

func (u *User) SetPassword(password string) {
	u.password = password
}

func (u *User) SetStore(store Store) {
	u.storeCode = store.Code()
}

func (u *User) SetRole(role Role) {
	u.role = role
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Password() string {
	return u.password
}

func (u *User) StoreCode() string {
	return u.storeCode
}

func (u *User) Role() *Role {
	return &u.role
}
