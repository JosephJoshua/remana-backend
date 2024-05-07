package domain

import (
	"fmt"
	"net/url"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	shareddomain "github.com/JosephJoshua/remana-backend/internal/modules/shared/domain"
	"github.com/JosephJoshua/remana-backend/internal/optional"
	"github.com/google/uuid"
)

type Order interface {
	// ChangeContactPhoneNumber(newPhoneNumber shareddomain.PhoneNumber)
	// AddDamage(damage string)
	// RemoveDamage(damage string)
	// AddPhoneCondition(condition string)
	// RemovePhoneCondition(condition string)
	// MutateCost(amount int, reason string)
	// AddPhoto(photoID uuid.UUID)
	// ChangeTechnician(newTechnicianID uuid.UUID)
	// ConfirmToCustomer(time time.Time, contents string)
	// PickUpByCustomer(repayment OrderPayment)
	// CompleteRepair()
	// Cancel()

	ID() uuid.UUID
	CreationTime() time.Time
	Slug() string
	StoreID() uuid.UUID
	CustomerName() string
	ContactNumber() shareddomain.PhoneNumber
	PhoneType() string
	Color() string
	SalesPersonID() uuid.UUID
	TechnicianID() uuid.UUID
	Costs() []OrderCost
	PhoneConditions() []PhoneCondition
	PhoneEquipments() []PhoneEquipment
	Damages() []Damage
	Photos() []OrderPhoto
	IMEI() optional.Optional[string]
	PartsNotCheckedYet() optional.Optional[string]
	PhoneSecurityDetails() optional.Optional[PhoneSecurityDetails]
	ConfirmationTime() optional.Optional[time.Time]
	ConfirmationContents() optional.Optional[string]
	PickUpTime() optional.Optional[time.Time]
	CompletionTime() optional.Optional[time.Time]
	CancellationTime() optional.Optional[time.Time]
	CancellationReason() optional.Optional[string]
	DownPayment() optional.Optional[OrderPayment]
	Repayment() optional.Optional[OrderPayment]
}

type order struct {
	id                   uuid.UUID
	creationTime         time.Time
	slug                 string
	storeID              uuid.UUID
	customerName         string
	contactNumber        shareddomain.PhoneNumber
	phoneType            string
	color                string
	salesPersonID        uuid.UUID
	technicianID         uuid.UUID
	costs                []OrderCost
	phoneConditions      []PhoneCondition
	phoneEquipments      []PhoneEquipment
	damages              []Damage
	photos               []OrderPhoto
	imei                 optional.Optional[string]
	partsNotCheckedYet   optional.Optional[string]
	phoneSecurityDetails optional.Optional[PhoneSecurityDetails]
	confirmationTime     optional.Optional[time.Time]
	confirmationContents optional.Optional[string]
	pickUpTime           optional.Optional[time.Time]
	completionTime       optional.Optional[time.Time]
	cancellationTime     optional.Optional[time.Time]
	cancellationReason   optional.Optional[string]
	downPayment          optional.Optional[OrderPayment]
	repayment            optional.Optional[OrderPayment]
}

type NewOrderParams struct {
	CreationTime         time.Time
	Slug                 string
	StoreID              uuid.UUID
	CustomerName         string
	ContactNumber        shareddomain.PhoneNumber
	PhoneType            string
	Color                string
	InitialCost          uint
	PhoneConditions      []string
	PhoneEquipments      []string
	Damages              []string
	Photos               []url.URL
	SalesPersonID        uuid.UUID
	TechnicianID         uuid.UUID
	Imei                 optional.Optional[string]
	PartsNotCheckedYet   optional.Optional[string]
	DownPayment          optional.Optional[OrderPayment]
	PhoneSecurityDetails optional.Optional[PhoneSecurityDetails]
}

func NewOrder(
	params NewOrderParams,
) (Order, error) {
	if err := validateNewOrderParams(params); err != nil {
		return nil, fmt.Errorf("failed to validate new order params: %w", err)
	}

	cost, costErr := newInitialOrderCost(uuid.New(), params.InitialCost, params.CreationTime)
	if costErr != nil {
		return nil, fmt.Errorf("failed to create initial order cost: %w", costErr)
	}

	phoneConditionVOs := make([]PhoneCondition, 0, len(params.PhoneConditions))
	for _, condition := range params.PhoneConditions {
		phoneCondition, err := newPhoneCondition(uuid.New(), condition)
		if err != nil {
			return nil, fmt.Errorf("failed to create phone condition: %w", err)
		}

		phoneConditionVOs = append(phoneConditionVOs, phoneCondition)
	}

	phoneEquipmentVOs := make([]PhoneEquipment, 0, len(params.PhoneEquipments))
	for _, equipment := range params.PhoneEquipments {
		phoneEquipment, err := newPhoneEquipment(uuid.New(), equipment)
		if err != nil {
			return nil, fmt.Errorf("failed to create phone equipment: %w", err)
		}

		phoneEquipmentVOs = append(phoneEquipmentVOs, phoneEquipment)
	}

	damageVOs := make([]Damage, 0, len(params.Damages))
	for _, damage := range params.Damages {
		damageVO, err := newDamage(uuid.New(), damage)
		if err != nil {
			return nil, fmt.Errorf("failed to create damage: %w", err)
		}

		damageVOs = append(damageVOs, damageVO)
	}

	photoVOs := make([]OrderPhoto, 0, len(params.Photos))
	for _, photo := range params.Photos {
		photoVO := newOrderPhoto(uuid.New(), photo)
		photoVOs = append(photoVOs, photoVO)
	}

	o := &order{
		id:                   uuid.New(),
		creationTime:         params.CreationTime,
		slug:                 params.Slug,
		storeID:              params.StoreID,
		customerName:         params.CustomerName,
		contactNumber:        params.ContactNumber,
		phoneType:            params.PhoneType,
		color:                params.Color,
		costs:                []OrderCost{cost},
		phoneConditions:      phoneConditionVOs,
		phoneEquipments:      phoneEquipmentVOs,
		damages:              damageVOs,
		photos:               photoVOs,
		salesPersonID:        params.SalesPersonID,
		technicianID:         params.TechnicianID,
		imei:                 params.Imei,
		partsNotCheckedYet:   params.PartsNotCheckedYet,
		phoneSecurityDetails: params.PhoneSecurityDetails,
		downPayment:          params.DownPayment,
		confirmationTime:     optional.None[time.Time](),
		confirmationContents: optional.None[string](),
		pickUpTime:           optional.None[time.Time](),
		completionTime:       optional.None[time.Time](),
		cancellationTime:     optional.None[time.Time](),
		cancellationReason:   optional.None[string](),
		repayment:            optional.None[OrderPayment](),
	}

	return o, nil
}

func (o *order) ID() uuid.UUID {
	return o.id
}

func (o *order) CreationTime() time.Time {
	return o.creationTime
}

func (o *order) Slug() string {
	return o.slug
}

func (o *order) StoreID() uuid.UUID {
	return o.storeID
}

func (o *order) CustomerName() string {
	return o.customerName
}

func (o *order) ContactNumber() shareddomain.PhoneNumber {
	return o.contactNumber
}

func (o *order) PhoneType() string {
	return o.phoneType
}

func (o *order) Color() string {
	return o.color
}

func (o *order) SalesPersonID() uuid.UUID {
	return o.salesPersonID
}

func (o *order) TechnicianID() uuid.UUID {
	return o.technicianID
}

func (o *order) Costs() []OrderCost {
	return o.costs
}

func (o *order) PhoneConditions() []PhoneCondition {
	return o.phoneConditions
}

func (o *order) PhoneEquipments() []PhoneEquipment {
	return o.phoneEquipments
}

func (o *order) Damages() []Damage {
	return o.damages
}

func (o *order) Photos() []OrderPhoto {
	return o.photos
}

func (o *order) IMEI() optional.Optional[string] {
	return o.imei
}

func (o *order) PartsNotCheckedYet() optional.Optional[string] {
	return o.partsNotCheckedYet
}

func (o *order) PhoneSecurityDetails() optional.Optional[PhoneSecurityDetails] {
	return o.phoneSecurityDetails
}

func (o *order) ConfirmationTime() optional.Optional[time.Time] {
	return o.confirmationTime
}

func (o *order) ConfirmationContents() optional.Optional[string] {
	return o.confirmationContents
}

func (o *order) PickUpTime() optional.Optional[time.Time] {
	return o.pickUpTime
}

func (o *order) CompletionTime() optional.Optional[time.Time] {
	return o.completionTime
}

func (o *order) CancellationTime() optional.Optional[time.Time] {
	return o.cancellationTime
}

func (o *order) CancellationReason() optional.Optional[string] {
	return o.cancellationReason
}

func (o *order) DownPayment() optional.Optional[OrderPayment] {
	return o.downPayment
}

func (o *order) Repayment() optional.Optional[OrderPayment] {
	return o.repayment
}

func validateNewOrderParams(params NewOrderParams) error {
	if params.Slug == "" {
		return fmt.Errorf("%w: slug is empty", apperror.ErrInvalidInput)
	}

	if params.CustomerName == "" {
		return fmt.Errorf("%w: customerName is empty", apperror.ErrInvalidInput)
	}

	if params.PhoneType == "" {
		return fmt.Errorf("%w: phoneType is empty", apperror.ErrInvalidInput)
	}

	if params.Color == "" {
		return fmt.Errorf("%w: color is empty", apperror.ErrInvalidInput)
	}

	if params.Imei.IsSet() && params.Imei.MustGet() == "" {
		return fmt.Errorf("%w: imei is empty", apperror.ErrInvalidInput)
	}

	if params.PartsNotCheckedYet.IsSet() && params.PartsNotCheckedYet.MustGet() == "" {
		return fmt.Errorf("%w: imei is empty", apperror.ErrInvalidInput)
	}

	if params.DownPayment.IsSet() && params.DownPayment.MustGet().Amount() > params.InitialCost {
		return fmt.Errorf("%w: down payment amount is greater than initial cost", apperror.ErrInvalidInput)
	}

	if len(params.Damages) == 0 {
		return fmt.Errorf("%w: damages is empty", apperror.ErrInvalidInput)
	}

	if len(params.Photos) == 0 {
		return fmt.Errorf("%w: photos is empty", apperror.ErrInvalidInput)
	}

	return nil
}
