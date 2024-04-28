package testutil

import (
	"net/url"

	"github.com/JosephJoshua/remana-backend/internal/optional"
	"github.com/google/uuid"
)

type ResourceLocationProviderStub struct {
	repairOrderLocation url.URL
	repairOrderErr      error
	RepairOrderID       optional.Optional[uuid.UUID]
}

func NewResourceLocationProviderStubForRepairOrder(
	location url.URL,
	err error,
) *ResourceLocationProviderStub {
	return &ResourceLocationProviderStub{
		repairOrderLocation: location,
		repairOrderErr:      err,
		RepairOrderID:       optional.None[uuid.UUID](),
	}
}

func (r *ResourceLocationProviderStub) SetRepairOrderErr(err error) {
	r.repairOrderErr = err
}

func (r *ResourceLocationProviderStub) RepairOrder(orderID uuid.UUID) (url.URL, error) {
	if r.repairOrderErr != nil {
		return url.URL{}, r.repairOrderErr
	}

	r.RepairOrderID = optional.Some(orderID)
	return r.repairOrderLocation, nil
}
