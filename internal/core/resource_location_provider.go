package core

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

type resourceLocationProvider struct{}

func (r resourceLocationProvider) RepairOrder(orderID uuid.UUID) (url.URL, error) {
	url := url.URL{
		Path: fmt.Sprintf("/repair-orders/%s", orderID.String()),
	}

	return url, nil
}
