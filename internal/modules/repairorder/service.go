package repairorder

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/apierror"
	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/modules/repairorder/domain"
	"github.com/JosephJoshua/remana-backend/internal/shared"
	shareddomain "github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Repository interface {
	CreateRepairOrder(ctx context.Context, order domain.Order) error
	GetDamageNamesByIDs(ctx context.Context, storeID uuid.UUID, ids []uuid.UUID) ([]string, error)
	GetPhoneConditionNamesByIDs(ctx context.Context, storeID uuid.UUID, ids []uuid.UUID) ([]string, error)
	GetPhoneEquipmentNamesByIDs(ctx context.Context, storeID uuid.UUID, ids []uuid.UUID) ([]string, error)
	DoesTechnicianExist(ctx context.Context, storeID uuid.UUID, technicianID uuid.UUID) (bool, error)
	DoesSalesExist(ctx context.Context, storeID uuid.UUID, salesID uuid.UUID) (bool, error)
	DoesPaymentMethodExist(ctx context.Context, storeID uuid.UUID, paymentMethodID uuid.UUID) (bool, error)
}

type OrderSlugProvider interface {
	Generate(ctx context.Context, storeID uuid.UUID) (string, error)
}

type TimeProvider interface {
	Now() time.Time
}

type ResourceLocationProvider interface {
	RepairOrder(orderID uuid.UUID) (url.URL, error)
}

type Service struct {
	timeProvider      TimeProvider
	locationProvider  ResourceLocationProvider
	repo              Repository
	orderSlugProvider OrderSlugProvider
}

func NewService(
	timeProvider TimeProvider,
	locationProvider ResourceLocationProvider,
	repo Repository,
	orderSlugProvider OrderSlugProvider,
) *Service {
	return &Service{
		timeProvider:      timeProvider,
		locationProvider:  locationProvider,
		repo:              repo,
		orderSlugProvider: orderSlugProvider,
	}
}

func (s *Service) CreateRepairOrder(
	ctx context.Context,
	req *genapi.CreateRepairOrderRequest,
) (*genapi.CreateRepairOrderCreated, error) {
	l := zerolog.Ctx(ctx)
	creationTime := s.timeProvider.Now()

	user, ok := shared.GetUserFromContext(ctx)
	if !ok {
		return nil, apierror.ToAPIError(http.StatusUnauthorized, "unauthorized")
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

	opts, err := s.buildOrderOptions(req)
	if err != nil {
		return nil, err
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

	repairOrder, err := domain.NewOrder(
		creationTime,
		slug,
		storeID,
		req.CustomerName,
		contactNumber,
		req.PhoneType,
		req.Color,
		uint(req.InitialCost),
		phoneConditions,
		phoneEquipments,
		damages,
		req.Photos,
		req.SalesID,
		req.TechnicianID,
		opts...,
	)

	if err != nil {
		return nil, apierror.ToAPIError(http.StatusBadRequest, err.Error())
	}

	err = s.repo.CreateRepairOrder(ctx, repairOrder)
	if err != nil {
		l.Error().Err(err).Msg("failed to create repair order")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to create repair order")
	}

	location, err := s.locationProvider.RepairOrder(repairOrder.ID())
	if err != nil {
		l.Error().Err(err).Msg("failed to get resource location")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to get resource location")
	}

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

	ok, err = s.repo.DoesSalesExist(ctx, storeID, req.SalesID)

	if err != nil {
		l.Error().Err(err).Msg("failed to check if sales exists")
		return apierror.ToAPIError(http.StatusInternalServerError, "failed to check if sales exists")
	}

	if !ok {
		return apierror.ToAPIError(http.StatusBadRequest, "sales does not exist")
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

func (s *Service) buildOrderOptions(req *genapi.CreateRepairOrderRequest) ([]domain.OrderOption, error) {
	opts := []domain.OrderOption{}

	if req.Imei.IsSet() {
		opts = append(opts, domain.WithIMEI(req.Imei.Value))
	}

	if req.PartsNotCheckedYet.IsSet() {
		opts = append(opts, domain.WithPartsNotCheckedYet(req.PartsNotCheckedYet.Value))
	}

	if req.Passcode.IsSet() {
		var securityDetails domain.PhoneSecurityDetails

		if req.Passcode.Value.IsPatternLocked {
			patternSecurity, err := domain.NewPatternSecurity(req.Passcode.Value.Value)
			if err != nil {
				return nil, apierror.ToAPIError(http.StatusBadRequest, err.Error())
			}

			securityDetails = patternSecurity
		} else {
			securityDetails = domain.NewPasscodeSecurity(req.Passcode.Value.Value)
		}

		opts = append(opts, domain.WithPhoneSecurityDetails(securityDetails))
	}

	if req.DownPayment.IsSet() {
		if req.DownPayment.Value.Amount <= 0 {
			return nil, apierror.ToAPIError(http.StatusBadRequest, "down payment amount must be greater than 0")
		}

		opts = append(opts, domain.WithDownPayment(uint(req.DownPayment.Value.Amount), req.DownPayment.Value.Method))
	}

	return opts, nil
}
