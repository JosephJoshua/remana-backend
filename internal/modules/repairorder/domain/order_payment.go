package domain

import (
	"fmt"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/google/uuid"
)

type OrderPayment interface {
	Amount() uint
	PaymentMethodID() uuid.UUID
}

type orderPayment struct {
	amount          uint
	paymentMethodID uuid.UUID
}

func NewOrderPayment(amount uint, paymentMethodID uuid.UUID) (OrderPayment, error) {
	if amount == 0 {
		return nil, fmt.Errorf("%w: amount is zero", apperror.ErrInvalidInput)
	}

	return orderPayment{amount: amount, paymentMethodID: paymentMethodID}, nil
}

func (o orderPayment) Amount() uint {
	return o.amount
}

func (o orderPayment) PaymentMethodID() uuid.UUID {
	return o.paymentMethodID
}
