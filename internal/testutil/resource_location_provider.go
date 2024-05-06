package testutil

import (
	"net/url"

	"github.com/JosephJoshua/remana-backend/internal/optional"
	"github.com/google/uuid"
)

type ResourceLocationProviderStub struct {
	repairOrderLocation url.URL
	technicianLocation  url.URL
	salesPersonLocation url.URL
	damageTypeLocation  url.URL

	RepairOrderID optional.Optional[uuid.UUID]
	TechnicianID  optional.Optional[uuid.UUID]
	SalesPersonID optional.Optional[uuid.UUID]
	DamageTypeID  optional.Optional[uuid.UUID]
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

func NewResourceLocationProviderStubForSalesPerson(location url.URL) *ResourceLocationProviderStub {
	return &ResourceLocationProviderStub{
		salesPersonLocation: location,
		SalesPersonID:       optional.None[uuid.UUID](),
	}
}

func NewResourceLocationProviderStubForDamageType(location url.URL) *ResourceLocationProviderStub {
	return &ResourceLocationProviderStub{
		damageTypeLocation: location,
		DamageTypeID:       optional.None[uuid.UUID](),
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

func (r *ResourceLocationProviderStub) SalesPerson(salesPersonID uuid.UUID) url.URL {
	r.SalesPersonID = optional.Some(salesPersonID)
	return r.salesPersonLocation
}

func (r *ResourceLocationProviderStub) DamageType(damageTypeID uuid.UUID) url.URL {
	r.DamageTypeID = optional.Some(damageTypeID)
	return r.damageTypeLocation
}
