package domain

import (
	"net/url"

	"github.com/google/uuid"
)

type OrderPhoto interface {
	ID() uuid.UUID
	URL() url.URL
}

type orderPhoto struct {
	id  uuid.UUID
	url url.URL
}

func newOrderPhoto(id uuid.UUID, url url.URL) OrderPhoto {
	return orderPhoto{
		id:  id,
		url: url,
	}
}

func (o orderPhoto) ID() uuid.UUID {
	return o.id
}

func (o orderPhoto) URL() url.URL {
	return o.url
}
