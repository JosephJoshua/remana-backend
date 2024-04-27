package domain

import (
	"fmt"

	"github.com/JosephJoshua/remana-backend/internal/shared/apperror"
	"github.com/google/uuid"
)

type PhoneEquipment interface {
	ID() uuid.UUID
	Name() string
}

type phoneEquipment struct {
	id   uuid.UUID
	name string
}

func newPhoneEquipment(id uuid.UUID, name string) (PhoneEquipment, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name is empty", apperror.ErrInvalidInput)
	}

	return phoneEquipment{
		id:   id,
		name: name,
	}, nil
}

func (p phoneEquipment) ID() uuid.UUID {
	return p.id
}

func (p phoneEquipment) Name() string {
	return p.name
}
