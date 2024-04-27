package domain

import (
	"fmt"
	"net/url"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/shared/apperror"
	shareddomain "github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/JosephJoshua/remana-backend/internal/shared/optional"
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
	SalesID() uuid.UUID
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

	setDownPayment(amount uint, paymentMethodID uuid.UUID) error
	setPhoneSecurityDetails(details PhoneSecurityDetails) error
	setPartsNotCheckedYet(parts string) error
	setIMEI(imei string) error
}

type OrderOption func(Order) error

type order struct {
	id                   uuid.UUID
	creationTime         time.Time
	slug                 string
	storeID              uuid.UUID
	customerName         string
	contactNumber        shareddomain.PhoneNumber
	phoneType            string
	color                string
	salesID              uuid.UUID
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

func NewOrder(
	creationTime time.Time,
	slug string,
	storeID uuid.UUID,
	customerName string,
	contactNumber shareddomain.PhoneNumber,
	phoneType string,
	color string,
	initialCost uint,
	phoneConditions []string,
	phoneEquipments []string,
	damages []string,
	photos []url.URL,
	salesID uuid.UUID,
	technicianID uuid.UUID,
	options ...OrderOption,
) (Order, error) {
	if slug == "" {
		return nil, fmt.Errorf("%w: slug is empty", apperror.ErrInvalidInput)
	}

	if customerName == "" {
		return nil, fmt.Errorf("%w: customerName is empty", apperror.ErrInvalidInput)
	}

	if phoneType == "" {
		return nil, fmt.Errorf("%w: phoneType is empty", apperror.ErrInvalidInput)
	}

	if color == "" {
		return nil, fmt.Errorf("%w: color is empty", apperror.ErrInvalidInput)
	}

	cost, costErr := newInitialOrderCost(uuid.New(), initialCost, creationTime)
	if costErr != nil {
		return nil, fmt.Errorf("failed to create initial order cost: %w", costErr)
	}

	phoneConditionVOs := make([]PhoneCondition, 0, len(phoneConditions))
	for _, condition := range phoneConditions {
		phoneCondition, err := newPhoneCondition(uuid.New(), condition)
		if err != nil {
			return nil, fmt.Errorf("failed to create phone condition: %w", err)
		}

		phoneConditionVOs = append(phoneConditionVOs, phoneCondition)
	}

	phoneEquipmentVOs := make([]PhoneEquipment, 0, len(phoneEquipments))
	for _, equipment := range phoneEquipments {
		phoneEquipment, err := newPhoneEquipment(uuid.New(), equipment)
		if err != nil {
			return nil, fmt.Errorf("failed to create phone equipment: %w", err)
		}

		phoneEquipmentVOs = append(phoneEquipmentVOs, phoneEquipment)
	}

	if len(damages) == 0 {
		return nil, fmt.Errorf("%w: damages is empty", apperror.ErrInvalidInput)
	}

	damageVOs := make([]Damage, 0, len(damages))
	for _, damage := range damages {
		damageVO, err := newDamage(uuid.New(), damage)
		if err != nil {
			return nil, fmt.Errorf("failed to create damage: %w", err)
		}

		damageVOs = append(damageVOs, damageVO)
	}

	if len(photos) == 0 {
		return nil, fmt.Errorf("%w: photos is empty", apperror.ErrInvalidInput)
	}

	photoVOs := make([]OrderPhoto, 0, len(photos))
	for _, photo := range photos {
		photoVO := newOrderPhoto(uuid.New(), photo)
		photoVOs = append(photoVOs, photoVO)
	}

	o := &order{
		id:                   uuid.New(),
		creationTime:         creationTime,
		slug:                 slug,
		storeID:              storeID,
		customerName:         customerName,
		contactNumber:        contactNumber,
		phoneType:            phoneType,
		color:                color,
		costs:                []OrderCost{cost},
		phoneConditions:      phoneConditionVOs,
		phoneEquipments:      phoneEquipmentVOs,
		damages:              damageVOs,
		photos:               photoVOs,
		salesID:              salesID,
		technicianID:         technicianID,
		imei:                 optional.None[string](),
		partsNotCheckedYet:   optional.None[string](),
		phoneSecurityDetails: optional.None[PhoneSecurityDetails](),
		confirmationTime:     optional.None[time.Time](),
		confirmationContents: optional.None[string](),
		pickUpTime:           optional.None[time.Time](),
		completionTime:       optional.None[time.Time](),
		cancellationTime:     optional.None[time.Time](),
		cancellationReason:   optional.None[string](),
		downPayment:          optional.None[OrderPayment](),
		repayment:            optional.None[OrderPayment](),
	}

	for _, option := range options {
		if err := option(o); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return o, nil
}

func WithDownPayment(amount uint, paymentMethodID uuid.UUID) OrderOption {
	return func(o Order) error {
		if err := o.setDownPayment(amount, paymentMethodID); err != nil {
			return fmt.Errorf("failed to set down payment: %w", err)
		}

		return nil
	}
}

func WithPhoneSecurityDetails(details PhoneSecurityDetails) OrderOption {
	return func(o Order) error {
		if err := o.setPhoneSecurityDetails(details); err != nil {
			return fmt.Errorf("failed to set phone security details: %w", err)
		}

		return nil
	}
}

func WithPartsNotCheckedYet(parts string) OrderOption {
	return func(o Order) error {
		if err := o.setPartsNotCheckedYet(parts); err != nil {
			return fmt.Errorf("failed to set parts not checked yet: %w", err)
		}

		return nil
	}
}

func WithIMEI(imei string) OrderOption {
	return func(o Order) error {
		if err := o.setIMEI(imei); err != nil {
			return fmt.Errorf("failed to set imei: %w", err)
		}

		return nil
	}
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

func (o *order) SalesID() uuid.UUID {
	return o.salesID
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

func (o *order) setDownPayment(amount uint, paymentMethodID uuid.UUID) error {
	if o.downPayment.IsSet() {
		return apperror.ErrValueAlreadySet
	}

	op, err := newOrderPayment(amount, paymentMethodID)
	if err != nil {
		return fmt.Errorf("failed to create order payment: %w", err)
	}

	o.downPayment = optional.Some(op)
	return nil
}

func (o *order) setPhoneSecurityDetails(details PhoneSecurityDetails) error {
	if o.phoneSecurityDetails.IsSet() {
		return apperror.ErrValueAlreadySet
	}

	o.phoneSecurityDetails = optional.Some(details)
	return nil
}

func (o *order) setPartsNotCheckedYet(parts string) error {
	if o.partsNotCheckedYet.IsSet() {
		return apperror.ErrValueAlreadySet
	}

	if parts == "" {
		return fmt.Errorf("%w: parts is empty", apperror.ErrInvalidInput)
	}

	o.partsNotCheckedYet = optional.Some(parts)
	return nil
}

func (o *order) setIMEI(imei string) error {
	if o.imei.IsSet() {
		return apperror.ErrValueAlreadySet
	}

	if imei == "" {
		return fmt.Errorf("%w: imei is empty", apperror.ErrInvalidInput)
	}

	o.imei = optional.Some(imei)
	return nil
}
