package domain

import "fmt"

const (
	RoleNameMinLength = 3
	RoleNameMaxLength = 255
)

type Role struct {
	id           int
	name         string
	storeCode    string
	isStoreAdmin bool
}

func NewRole(id int, name string, store Store, isStoreAdmin bool) (*Role, error) {
	role := new(Role)

	role.SetIsStoreAdmin(isStoreAdmin)
	role.SetStore(store)

	if err := role.setID(id); err != nil {
		return nil, fmt.Errorf("failed to create new role: %w", err)
	}

	if err := role.SetName(name); err != nil {
		return nil, fmt.Errorf("failed to create new role: %w", err)
	}

	return role, nil
}

func (r *Role) SetName(name string) error {
	if len(name) < RoleNameMinLength {
		return fmt.Errorf("error setting name of role: %w", ErrInputTooShort)
	}

	if len(name) > RoleNameMaxLength {
		return fmt.Errorf("error setting name of role: %w", ErrInputTooLong)
	}

	r.name = name
	return nil
}

func (r *Role) SetStore(store Store) {
	r.storeCode = store.Code()
}

func (r *Role) SetIsStoreAdmin(isStoreAdmin bool) {
	r.isStoreAdmin = isStoreAdmin
}

func (r *Role) ID() int {
	return r.id
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) StoreCode() string {
	return r.storeCode
}

func (r *Role) IsStoreAdmin() bool {
	return r.isStoreAdmin
}

func (r *Role) setID(id int) error {
	if id < 0 {
		return fmt.Errorf("error setting id of role: %w", ErrInvalidID)
	}

	r.id = id
	return nil
}
