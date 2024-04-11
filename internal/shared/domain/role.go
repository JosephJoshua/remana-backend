package domain

import "fmt"

type Role interface {
	SetName(name string) error
	SetIsStoreAdmin(isStoreAdmin bool)
	ID() int
	Name() string
	StoreCode() string
	IsStoreAdmin() bool
}

type role struct {
	id           int
	name         string
	storeCode    string
	isStoreAdmin bool
}

func NewRole(id int, name string, store Store, isStoreAdmin bool) (Role, error) {
	role := new(role)

	role.storeCode = store.Code()

	role.SetIsStoreAdmin(isStoreAdmin)

	if err := role.setID(id); err != nil {
		return nil, fmt.Errorf("failed to create new role: %w", err)
	}

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

func (r *role) ID() int {
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

func (r *role) setID(id int) error {
	if id < 0 {
		return fmt.Errorf("error setting id of role: %w", ErrInvalidID)
	}

	r.id = id
	return nil
}
