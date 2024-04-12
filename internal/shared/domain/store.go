package domain

import (
	"fmt"
	"unicode"

	"github.com/google/uuid"
)

type Store interface {
	SetName(name string) error
	SetCode(code string) error
	ID() uuid.UUID
	Name() string
	Code() string
}

type store struct {
	id   uuid.UUID
	name string
	code string
}

func NewStore(id uuid.UUID, name string, code string) (Store, error) {
	store := new(store)

	store.id = id

	if err := store.SetName(name); err != nil {
		return nil, fmt.Errorf("failed to create new store: %w", err)
	}

	if err := store.SetCode(code); err != nil {
		return nil, fmt.Errorf("failed to create new store: %w", err)
	}

	return store, nil
}

func (s *store) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("error setting name of store: %w", ErrInputTooShort)
	}

	s.name = name
	return nil
}

func (s *store) SetCode(code string) error {
	if code == "" {
		return fmt.Errorf("error setting code of store: %w", ErrInputTooShort)
	}

	for _, c := range code {
		if c == '-' || unicode.IsLower(c) {
			continue
		}

		return fmt.Errorf("error setting store code of user: %w", ErrInvalidStoreCode)
	}

	s.code = code
	return nil
}

func (s *store) ID() uuid.UUID {
	return s.id
}

func (s *store) Name() string {
	return s.name
}

func (s *store) Code() string {
	return s.code
}
