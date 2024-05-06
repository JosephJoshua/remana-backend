package core

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

type resourceLocationProvider struct{}

func (r resourceLocationProvider) RepairOrder(id uuid.UUID) url.URL {
	url := url.URL{
		Path: fmt.Sprintf("/repair-orders/%s", id.String()),
	}

	return url
}

func (r resourceLocationProvider) Technician(id uuid.UUID) url.URL {
	url := url.URL{
		Path: fmt.Sprintf("/technicians/%s", id.String()),
	}

	return url
}

func (r resourceLocationProvider) SalesPerson(id uuid.UUID) url.URL {
	url := url.URL{
		Path: fmt.Sprintf("/sales-persons/%s", id.String()),
	}

	return url
}

func (r resourceLocationProvider) DamageType(id uuid.UUID) url.URL {
	url := url.URL{
		Path: fmt.Sprintf("/damage-types/%s", id.String()),
	}

	return url
}
