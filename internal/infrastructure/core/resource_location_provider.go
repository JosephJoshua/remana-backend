package core

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

type resourceLocationProvider struct{}

func (r resourceLocationProvider) RepairOrder(orderID uuid.UUID) url.URL {
	url := url.URL{
		Path: fmt.Sprintf("/repair-orders/%s", orderID.String()),
	}

	return url
}

func (r resourceLocationProvider) Technician(technicianID uuid.UUID) url.URL {
	url := url.URL{
		Path: fmt.Sprintf("/technicians/%s", technicianID.String()),
	}

	return url
}
