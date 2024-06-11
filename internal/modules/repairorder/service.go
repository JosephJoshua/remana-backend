package repairorder

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/apierror"
	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/modules/permission"
	"github.com/JosephJoshua/remana-backend/internal/modules/repairorder/domain"
	shareddomain "github.com/JosephJoshua/remana-backend/internal/modules/shared/domain"
	"github.com/JosephJoshua/remana-backend/internal/optional"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Repository interface {
	CreateRepairOrder(ctx context.Context, order domain.Order) error
	GetDamageNamesByIDs(ctx context.Context, storeID uuid.UUID, ids []uuid.UUID) ([]string, error)
	GetPhoneConditionNamesByIDs(ctx context.Context, storeID uuid.UUID, ids []uuid.UUID) ([]string, error)
	GetPhoneEquipmentNamesByIDs(ctx context.Context, storeID uuid.UUID, ids []uuid.UUID) ([]string, error)
	DoesTechnicianExist(ctx context.Context, storeID uuid.UUID, technicianID uuid.UUID) (bool, error)
	DoesSalesPersonExist(ctx context.Context, storeID uuid.UUID, salesPersonID uuid.UUID) (bool, error)
	DoesPaymentMethodExist(ctx context.Context, storeID uuid.UUID, paymentMethodID uuid.UUID) (bool, error)
}

type OrderSlugProvider interface {
	Generate(ctx context.Context, storeID uuid.UUID) (string, error)
}

type TimeProvider interface {
	Now() time.Time
}

type ResourceLocationProvider interface {
	RepairOrder(orderID uuid.UUID) url.URL
}

type Service struct {
	timeProvider       TimeProvider
	locationProvider   ResourceLocationProvider
	repo               Repository
	orderSlugProvider  OrderSlugProvider
	permissionProvider permission.Provider
}

func NewService(
	timeProvider TimeProvider,
	locationProvider ResourceLocationProvider,
	repo Repository,
	permissionProvider permission.Provider,
	orderSlugProvider OrderSlugProvider,
) *Service {
	return &Service{
		timeProvider:       timeProvider,
		locationProvider:   locationProvider,
		repo:               repo,
		permissionProvider: permissionProvider,
		orderSlugProvider:  orderSlugProvider,
	}
}

func (s *Service) CreateRepairOrder(
	ctx context.Context,
	req *genapi.CreateRepairOrderRequest,
) (*genapi.CreateRepairOrderCreated, error) {
	l := zerolog.Ctx(ctx)
	creationTime := s.timeProvider.Now()

	user, ok := appcontext.GetUserFromContext(ctx)
	if !ok {
		l.Error().Msg("user is missing from context")
		return nil, apierror.ToAPIError(http.StatusUnauthorized, "unauthorized")
	}

	if can, err := s.permissionProvider.Can(ctx, user.Role.ID, permission.CreateRepairOrder()); err != nil {
		l.Error().Err(err).Msg("failed to check permission")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to check permission")
	} else if !can {
		return nil, apierror.ToAPIError(http.StatusForbidden, "insufficient permissions")
	}

	storeID := user.Store.ID

	contactNumber, err := shareddomain.NewPhoneNumber(req.ContactPhoneNumber)
	if err != nil {
		return nil, apierror.ToAPIError(http.StatusBadRequest, "invalid contact phone number")
	}

	if req.InitialCost <= 0 {
		return nil, apierror.ToAPIError(http.StatusBadRequest, "initial cost must be greater than 0")
	}

	slug, err := s.orderSlugProvider.Generate(ctx, storeID)
	if err != nil {
		l.Error().Err(err).Msg("failed to generate repair order slug")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to generate repair order slug")
	}

	var phoneSecurityDetails optional.Optional[domain.PhoneSecurityDetails]

	if req.Passcode.IsSet() {
		if req.Passcode.Value.IsPatternLocked {
			tmp, securityErr := domain.NewPatternSecurity(req.Passcode.Value.Value)
			if securityErr != nil {
				return nil, apierror.ToAPIError(http.StatusBadRequest, securityErr.Error())
			}

			phoneSecurityDetails = optional.Some(tmp)
		} else {
			phoneSecurityDetails = optional.Some(domain.NewPasscodeSecurity(req.Passcode.Value.Value))
		}
	}

	var downPayment optional.Optional[domain.OrderPayment]

	if req.DownPayment.IsSet() {
		if req.DownPayment.Value.Amount <= 0 {
			return nil, apierror.ToAPIError(http.StatusBadRequest, "down payment amount must be greater than 0")
		}

		tmp, paymentErr := domain.NewOrderPayment(uint(req.DownPayment.Value.Amount), req.DownPayment.Value.Method)
		if paymentErr != nil {
			return nil, apierror.ToAPIError(http.StatusBadRequest, paymentErr.Error())
		}

		downPayment = optional.Some(tmp)
	}

	if err = s.checkReferentialIntegrity(ctx, l, storeID, req); err != nil {
		return nil, err
	}

	damages, err := s.repo.GetDamageNamesByIDs(ctx, storeID, req.DamageTypes)
	if err != nil {
		if errors.Is(err, apperror.ErrDamageNotFound) {
			return nil, apierror.ToAPIError(http.StatusBadRequest, "one or more damage types do not exist")
		}

		l.Error().Err(err).Msg("failed to get damage names by IDs")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to get damage names by IDs")
	}

	phoneConditions, err := s.repo.GetPhoneConditionNamesByIDs(ctx, storeID, req.PhoneConditions)
	if err != nil {
		if errors.Is(err, apperror.ErrPhoneConditionNotFound) {
			return nil, apierror.ToAPIError(http.StatusBadRequest, "one or more phone conditions do not exist")
		}

		l.Error().Err(err).Msg("failed to get phone condition names by IDs")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to get phone condition names by IDs")
	}

	phoneEquipments, err := s.repo.GetPhoneEquipmentNamesByIDs(ctx, storeID, req.PhoneEquipments)
	if err != nil {
		if errors.Is(err, apperror.ErrPhoneEquipmentNotFound) {
			return nil, apierror.ToAPIError(http.StatusBadRequest, "one or more phone equipments do not exist")
		}

		l.Error().Err(err).Msg("failed to get phone equipment names by IDs")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to get phone equipment names by IDs")
	}

	var imei optional.Optional[string]
	if req.Imei.IsSet() {
		imei = optional.Some(req.Imei.Value)
	}

	var partsNotCheckedYet optional.Optional[string]
	if req.PartsNotCheckedYet.IsSet() {
		partsNotCheckedYet = optional.Some(req.PartsNotCheckedYet.Value)
	}

	params := domain.NewOrderParams{
		CreationTime:         creationTime,
		Slug:                 slug,
		StoreID:              storeID,
		CustomerName:         req.CustomerName,
		ContactNumber:        contactNumber,
		PhoneType:            req.PhoneType,
		Color:                req.Color,
		InitialCost:          uint(req.InitialCost),
		PhoneConditions:      phoneConditions,
		PhoneEquipments:      phoneEquipments,
		Damages:              damages,
		Photos:               req.Photos,
		SalesPersonID:        req.SalesPersonID,
		TechnicianID:         req.TechnicianID,
		Imei:                 imei,
		PartsNotCheckedYet:   partsNotCheckedYet,
		DownPayment:          downPayment,
		PhoneSecurityDetails: phoneSecurityDetails,
	}

	repairOrder, err := domain.NewOrder(params)

	if err != nil {
		return nil, apierror.ToAPIError(http.StatusBadRequest, err.Error())
	}

	err = s.repo.CreateRepairOrder(ctx, repairOrder)
	if err != nil {
		l.Error().Err(err).Msg("failed to create repair order")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to create repair order")
	}

	location := s.locationProvider.RepairOrder(repairOrder.ID())
	return &genapi.CreateRepairOrderCreated{
		Location: location,
	}, nil
}

func (s *Service) checkReferentialIntegrity(
	ctx context.Context,
	l *zerolog.Logger,
	storeID uuid.UUID,
	req *genapi.CreateRepairOrderRequest,
) error {
	ok, err := s.repo.DoesTechnicianExist(ctx, storeID, req.TechnicianID)

	if err != nil {
		l.Error().Err(err).Msg("failed to check if technician exists")
		return apierror.ToAPIError(http.StatusInternalServerError, "failed to check if technician exists")
	}

	if !ok {
		return apierror.ToAPIError(http.StatusBadRequest, "technician does not exist")
	}

	ok, err = s.repo.DoesSalesPersonExist(ctx, storeID, req.SalesPersonID)

	if err != nil {
		l.Error().Err(err).Msg("failed to check if sales person exists")
		return apierror.ToAPIError(http.StatusInternalServerError, "failed to check if sales person exists")
	}

	if !ok {
		return apierror.ToAPIError(http.StatusBadRequest, "sales person does not exist")
	}

	if req.DownPayment.IsSet() {
		ok, err = s.repo.DoesPaymentMethodExist(ctx, storeID, req.DownPayment.Value.GetMethod())
		if err != nil {
			l.Error().Err(err).Msg("failed to check if payment method exists")
			return apierror.ToAPIError(http.StatusInternalServerError, "failed to check if payment method exists")
		}

		if !ok {
			return apierror.ToAPIError(http.StatusBadRequest, "payment method does not exist")
		}
	}

	return nil
}
