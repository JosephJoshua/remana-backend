package testutil

import (
	"context"

	"github.com/google/uuid"
)

type OrderSlugProviderStub struct {
	slug string
	err  error
}

func NewRepairOrderSlugProviderStub(slug string, err error) *OrderSlugProviderStub {
	return &OrderSlugProviderStub{
		slug: slug,
		err:  err,
	}
}

func (o *OrderSlugProviderStub) SetError(err error) {
	o.err = err
}

func (o *OrderSlugProviderStub) Generate(_ context.Context, _ uuid.UUID) (string, error) {
	if o.err != nil {
		return "", o.err
	}

	return o.slug, nil
}
