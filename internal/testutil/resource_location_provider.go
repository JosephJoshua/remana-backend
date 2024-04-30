package testutil

import (
	"net/url"

	"github.com/JosephJoshua/remana-backend/internal/optional"
	"github.com/google/uuid"
)

type ResourceLocationProviderStub struct {
	repairOrderLocation url.URL
	technicianLocation  url.URL

	RepairOrderID optional.Optional[uuid.UUID]
	TechnicianID  optional.Optional[uuid.UUID]
}

func NewResourceLocationProviderStubForRepairOrder(location url.URL) *ResourceLocationProviderStub {
	return &ResourceLocationProviderStub{
		repairOrderLocation: location,
		RepairOrderID:       optional.None[uuid.UUID](),
	}
}

func NewResourceLocationProviderStubForTechnician(location url.URL) *ResourceLocationProviderStub {
	return &ResourceLocationProviderStub{
		technicianLocation: location,
		TechnicianID:       optional.None[uuid.UUID](),
	}
}

func (r *ResourceLocationProviderStub) RepairOrder(orderID uuid.UUID) url.URL {
	r.RepairOrderID = optional.Some(orderID)
	return r.repairOrderLocation
}

func (r *ResourceLocationProviderStub) Technician(technicianID uuid.UUID) url.URL {
	r.TechnicianID = optional.Some(technicianID)
	return r.technicianLocation
}
