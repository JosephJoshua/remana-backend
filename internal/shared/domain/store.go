package domain

import (
	"fmt"
	"unicode"
)

type Store struct {
	id   int
	name string
	code string
}

func NewStore(id int, name string, code string) (*Store, error) {
	store := new(Store)

	if err := store.setID(id); err != nil {
		return nil, fmt.Errorf("failed to create new store: %w", err)
	}

	if err := store.SetName(name); err != nil {
		return nil, fmt.Errorf("failed to create new store: %w", err)
	}

	if err := store.SetCode(code); err != nil {
		return nil, fmt.Errorf("failed to create new store: %w", err)
	}

	return store, nil
}

func (s *Store) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("error setting name of store: %w", ErrInputTooShort)
	}

	s.name = name
	return nil
}

func (s *Store) SetCode(code string) error {
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

func (s *Store) ID() int {
	return s.id
}

func (s *Store) Name() string {
	return s.name
}

func (s *Store) Code() string {
	return s.code
}

func (s *Store) setID(id int) error {
	if id < 0 {
		return fmt.Errorf("error setting id of store: %w", ErrInvalidID)
	}

	s.id = id
	return nil
}
