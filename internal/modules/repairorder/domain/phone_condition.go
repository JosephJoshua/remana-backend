package domain

import (
	"fmt"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/google/uuid"
)

type PhoneCondition interface {
	ID() uuid.UUID
	Name() string
}

type phoneCondition struct {
	id   uuid.UUID
	name string
}

func newPhoneCondition(id uuid.UUID, name string) (PhoneCondition, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: value is empty", apperror.ErrInvalidInput)
	}

	return phoneCondition{
		id:   id,
		name: name,
	}, nil
}

func (p phoneCondition) ID() uuid.UUID {
	return p.id
}

func (p phoneCondition) Name() string {
	return p.name
}
