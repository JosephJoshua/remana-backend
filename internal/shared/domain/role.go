package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type Role interface {
	SetName(name string) error
	SetIsStoreAdmin(isStoreAdmin bool)
	ID() uuid.UUID
	Name() string
	StoreCode() string
	IsStoreAdmin() bool
}

type role struct {
	id           uuid.UUID
	name         string
	storeCode    string
	isStoreAdmin bool
}

func NewRole(id uuid.UUID, name string, store Store, isStoreAdmin bool) (Role, error) {
	role := new(role)

	role.id = id
	role.storeCode = store.Code()

	role.SetIsStoreAdmin(isStoreAdmin)

	if err := role.SetName(name); err != nil {
		return nil, fmt.Errorf("failed to create new role: %w", err)
	}

	return role, nil
}

func (r *role) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("error setting name of role: %w", ErrInputTooShort)
	}

	r.name = name
	return nil
}

func (r *role) SetIsStoreAdmin(isStoreAdmin bool) {
	r.isStoreAdmin = isStoreAdmin
}

func (r *role) ID() uuid.UUID {
	return r.id
}

func (r *role) Name() string {
	return r.name
}

func (r *role) StoreCode() string {
	return r.storeCode
}

func (r *role) IsStoreAdmin() bool {
	return r.isStoreAdmin
}
