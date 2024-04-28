package domain

import (
	"fmt"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/google/uuid"
)

type Damage interface {
	ID() uuid.UUID
	Name() string
}

type damage struct {
	id   uuid.UUID
	name string
}

func newDamage(id uuid.UUID, name string) (Damage, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name is empty", apperror.ErrInvalidInput)
	}

	return damage{
		id:   id,
		name: name,
	}, nil
}

func (d damage) ID() uuid.UUID {
	return d.id
}

func (d damage) Name() string {
	return d.name
}
