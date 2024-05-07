package testutil

import (
	"net/url"

	"github.com/JosephJoshua/remana-backend/internal/optional"
	"github.com/google/uuid"
)

type ResourceLocationProviderStub struct {
	repairOrderLocation    url.URL
	technicianLocation     url.URL
	salesPersonLocation    url.URL
	damageTypeLocation     url.URL
	phoneConditionLocation url.URL
	phoneEquipmentLocation url.URL
	paymentMethodLocation  url.URL
	roleLocation           url.URL

	RepairOrderID    optional.Optional[uuid.UUID]
	TechnicianID     optional.Optional[uuid.UUID]
	SalesPersonID    optional.Optional[uuid.UUID]
	DamageTypeID     optional.Optional[uuid.UUID]
	PhoneConditionID optional.Optional[uuid.UUID]
	PhoneEquipmentID optional.Optional[uuid.UUID]
	PaymentMethodID  optional.Optional[uuid.UUID]
	RoleID           optional.Optional[uuid.UUID]
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

func NewResourceLocationProviderStubForPhoneCondition(location url.URL) *ResourceLocationProviderStub {
	return &ResourceLocationProviderStub{
		phoneConditionLocation: location,
		PhoneConditionID:       optional.None[uuid.UUID](),
	}
}

func NewResourceLocationProviderStubForPhoneEquipment(location url.URL) *ResourceLocationProviderStub {
	return &ResourceLocationProviderStub{
		phoneEquipmentLocation: location,
		PhoneEquipmentID:       optional.None[uuid.UUID](),
	}
}

func NewResourceLocationProviderStubForPaymentMethod(location url.URL) *ResourceLocationProviderStub {
	return &ResourceLocationProviderStub{
		paymentMethodLocation: location,
		PaymentMethodID:       optional.None[uuid.UUID](),
	}
}

func NewResourceLocationProviderStubForRole(location url.URL) *ResourceLocationProviderStub {
	return &ResourceLocationProviderStub{
		roleLocation: location,
		RoleID:       optional.None[uuid.UUID](),
	}
}

func (r *ResourceLocationProviderStub) RepairOrder(id uuid.UUID) url.URL {
	r.RepairOrderID = optional.Some(id)
	return r.repairOrderLocation
}

func (r *ResourceLocationProviderStub) Technician(id uuid.UUID) url.URL {
	r.TechnicianID = optional.Some(id)
	return r.technicianLocation
}

func (r *ResourceLocationProviderStub) SalesPerson(id uuid.UUID) url.URL {
	r.SalesPersonID = optional.Some(id)
	return r.salesPersonLocation
}

func (r *ResourceLocationProviderStub) DamageType(id uuid.UUID) url.URL {
	r.DamageTypeID = optional.Some(id)
	return r.damageTypeLocation
}

func (r *ResourceLocationProviderStub) PhoneCondition(id uuid.UUID) url.URL {
	r.PhoneConditionID = optional.Some(id)
	return r.phoneConditionLocation
}

func (r *ResourceLocationProviderStub) PhoneEquipment(id uuid.UUID) url.URL {
	r.PhoneEquipmentID = optional.Some(id)
	return r.phoneEquipmentLocation
}

func (r *ResourceLocationProviderStub) PaymentMethod(id uuid.UUID) url.URL {
	r.PaymentMethodID = optional.Some(id)
	return r.paymentMethodLocation
}

func (r *ResourceLocationProviderStub) Role(id uuid.UUID) url.URL {
	r.RoleID = optional.Some(id)
	return r.roleLocation
}
