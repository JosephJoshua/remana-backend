package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/optional"
	"github.com/google/uuid"
)

type OrderCost interface {
	ID() uuid.UUID
	Amount() int
	Reason() optional.Optional[string]
	IsInitial() bool
	CreationTime() time.Time
}

type orderCost struct {
	id           uuid.UUID
	amount       int
	reason       optional.Optional[string]
	creationTime time.Time
}

func newInitialOrderCost(id uuid.UUID, amount uint, creationTime time.Time) (OrderCost, error) {
	if amount == 0 {
		return nil, fmt.Errorf("%w: value is zero", apperror.ErrInvalidInput)
	}

	if amount > math.MaxInt {
		return nil, fmt.Errorf("%w: value is greater than MaxInt", apperror.ErrInvalidInput)
	}

	return orderCost{
		id:           id,
		amount:       int(amount),
		reason:       optional.None[string](),
		creationTime: creationTime,
	}, nil
}

func newAdditionalOrderCost(id uuid.UUID, amount int, reason string, creationTime time.Time) (OrderCost, error) {
	if amount == 0 {
		return nil, fmt.Errorf("%w: value is zero", apperror.ErrInvalidInput)
	}

	if reason == "" {
		return nil, fmt.Errorf("%w: reason is empty", apperror.ErrInvalidInput)
	}

	return orderCost{
		id:           id,
		amount:       amount,
		reason:       optional.Some(reason),
		creationTime: creationTime,
	}, nil
}

func (o orderCost) ID() uuid.UUID {
	return o.id
}

func (o orderCost) Amount() int {
	return o.amount
}

func (o orderCost) Reason() optional.Optional[string] {
	return o.reason
}

func (o orderCost) IsInitial() bool {
	return !o.reason.IsSet()
}

func (o orderCost) CreationTime() time.Time {
	return o.creationTime
}
