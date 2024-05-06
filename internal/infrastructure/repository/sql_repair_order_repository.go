package repository

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/modules/repairorder/domain"
	"github.com/JosephJoshua/remana-backend/internal/optional"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLRepairOrderRepository struct {
	queries *gensql.Queries
	db      *pgxpool.Pool
}

func NewSQLRepairOrderRepository(db *pgxpool.Pool) *SQLRepairOrderRepository {
	return &SQLRepairOrderRepository{
		queries: gensql.New(db),
		db:      db,
	}
}

func (r *SQLRepairOrderRepository) CreateRepairOrder(ctx context.Context, order domain.Order) (err error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			if errors.Is(rollbackErr, pgx.ErrTxClosed) {
				return
			}

			err = fmt.Errorf("failed to rollback transaction: %w", rollbackErr)
		}
	}()

	qtx := r.queries.WithTx(tx)

	params, err := r.buildCreateRepairOrderParams(order)
	if err != nil {
		return fmt.Errorf("failed to build create repair order params: %w", err)
	}

	if err = qtx.CreateRepairOrder(ctx, params); err != nil {
		return fmt.Errorf("failed to create repair order: %w", err)
	}

	if err = r.attachRepairOrderDamages(ctx, qtx, order); err != nil {
		return fmt.Errorf("failed to attach repair order damages: %w", err)
	}

	if err = r.attachRepairOrderPhoneConditions(ctx, qtx, order); err != nil {
		return fmt.Errorf("failed to attach repair order phone conditions: %w", err)
	}

	if err = r.attachRepairOrderPhoneEquipments(ctx, qtx, order); err != nil {
		return fmt.Errorf("failed to attach repair order phone equipments: %w", err)
	}

	if err = r.attachRepairOrderPhotos(ctx, qtx, order); err != nil {
		return fmt.Errorf("failed to attach repair order photos: %w", err)
	}

	if err = r.attachRepairOrderCosts(ctx, qtx, order); err != nil {
		return fmt.Errorf("failed to attach repair order costs: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *SQLRepairOrderRepository) GetDamageNamesByIDs(
	ctx context.Context,
	storeID uuid.UUID,
	ids []uuid.UUID,
) ([]string, error) {
	damageNames, err := r.queries.GetDamageNamesByIDs(
		ctx,
		gensql.GetDamageNamesByIDsParams{
			StoreID: typemapper.UUIDToPgtypeUUID(storeID),
			Ids:     typemapper.UUIDsToPgtypeUUIDs(ids),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get damage names by IDs: %w", err)
	}

	if len(damageNames) < len(ids) {
		return nil, apperror.ErrDamageNotFound
	}

	return damageNames, nil
}

func (r *SQLRepairOrderRepository) GetPhoneConditionNamesByIDs(
	ctx context.Context,
	storeID uuid.UUID,
	ids []uuid.UUID,
) ([]string, error) {
	phoneConditionNames, err := r.queries.GetPhoneConditionNamesByIDs(
		ctx,
		gensql.GetPhoneConditionNamesByIDsParams{
			StoreID: typemapper.UUIDToPgtypeUUID(storeID),
			Ids:     typemapper.UUIDsToPgtypeUUIDs(ids),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get phone condition names by IDs: %w", err)
	}

	if len(phoneConditionNames) < len(ids) {
		return nil, apperror.ErrPhoneConditionNotFound
	}

	return phoneConditionNames, nil
}

func (r *SQLRepairOrderRepository) GetPhoneEquipmentNamesByIDs(
	ctx context.Context,
	storeID uuid.UUID,
	ids []uuid.UUID,
) ([]string, error) {
	phoneEquipmentNames, err := r.queries.GetPhoneEquipmentNamesByIDs(
		ctx,
		gensql.GetPhoneEquipmentNamesByIDsParams{
			StoreID: typemapper.UUIDToPgtypeUUID(storeID),
			Ids:     typemapper.UUIDsToPgtypeUUIDs(ids),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get phone equipment names by IDs: %w", err)
	}

	if len(phoneEquipmentNames) < len(ids) {
		return nil, apperror.ErrPhoneEquipmentNotFound
	}

	return phoneEquipmentNames, nil
}

func (r *SQLRepairOrderRepository) DoesSalesPersonExist(
	ctx context.Context,
	storeID uuid.UUID,
	salesPersonID uuid.UUID,
) (bool, error) {
	_, err := r.queries.DoesSalesPersonExist(
		ctx,
		gensql.DoesSalesPersonExistParams{
			StoreID:       typemapper.UUIDToPgtypeUUID(storeID),
			SalesPersonID: typemapper.UUIDToPgtypeUUID(salesPersonID),
		},
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if sales person exists: %w", err)
	}

	return true, nil
}

func (r *SQLRepairOrderRepository) DoesTechnicianExist(
	ctx context.Context,
	storeID uuid.UUID,
	technicianID uuid.UUID,
) (bool, error) {
	_, err := r.queries.DoesTechnicianExist(
		ctx,
		gensql.DoesTechnicianExistParams{
			StoreID:      typemapper.UUIDToPgtypeUUID(storeID),
			TechnicianID: typemapper.UUIDToPgtypeUUID(technicianID),
		},
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if technician exists: %w", err)
	}

	return true, nil
}

func (r *SQLRepairOrderRepository) DoesPaymentMethodExist(
	ctx context.Context,
	storeID uuid.UUID,
	paymentMethodID uuid.UUID,
) (bool, error) {
	_, err := r.queries.DoesPaymentMethodExist(
		ctx,
		gensql.DoesPaymentMethodExistParams{
			StoreID:         typemapper.UUIDToPgtypeUUID(storeID),
			PaymentMethodID: typemapper.UUIDToPgtypeUUID(paymentMethodID),
		},
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if payment method exists: %w", err)
	}

	return true, nil
}

func (r *SQLRepairOrderRepository) buildCreateRepairOrderParams(
	order domain.Order,
) (gensql.CreateRepairOrderParams, error) {
	securityDetails := order.PhoneSecurityDetails()

	passcodeOrPattern := typemapper.OptionalStringToPgtypeText(optional.None[string]())
	isPatternLocked := typemapper.OptionalBoolToPgtypeBool(optional.None[bool]())

	if securityDetails.IsSet() && securityDetails.MustGet().Type() != domain.PhoneSecurityTypeNone {
		passcodeOrPattern = typemapper.StringToPgtypeText(securityDetails.MustGet().Value())
		isPatternLocked = typemapper.BoolToPgtypeBool(securityDetails.MustGet().Type() == domain.PhoneSecurityTypePattern)
	}

	downPaymentAmount := typemapper.OptionalInt32ToPgtypeInt4(optional.None[int32]())
	downPaymentMethodID := typemapper.OptionalUUIDToPgtypeUUID(optional.None[uuid.UUID]())

	downPayment := order.DownPayment()

	if downPayment.IsSet() {
		if downPayment.MustGet().Amount() > math.MaxInt {
			return gensql.CreateRepairOrderParams{}, errors.New("down payment amount is greater than MaxInt")
		}

		downPaymentAmount = typemapper.Int32ToPgtypeInt4(int32(downPayment.MustGet().Amount()))
		downPaymentMethodID = typemapper.UUIDToPgtypeUUID(downPayment.MustGet().PaymentMethodID())
	}

	return gensql.CreateRepairOrderParams{
		RepairOrderID:       typemapper.UUIDToPgtypeUUID(order.ID()),
		CreationTime:        typemapper.TimeToPgtypeTimestamptz(order.CreationTime()),
		Slug:                order.Slug(),
		StoreID:             typemapper.UUIDToPgtypeUUID(order.StoreID()),
		CustomerName:        order.CustomerName(),
		ContactNumber:       order.ContactNumber().Value(),
		PhoneType:           order.PhoneType(),
		Color:               order.Color(),
		SalesPersonID:       typemapper.UUIDToPgtypeUUID(order.SalesPersonID()),
		TechnicianID:        typemapper.UUIDToPgtypeUUID(order.TechnicianID()),
		Imei:                typemapper.OptionalStringToPgtypeText(order.IMEI()),
		PartsNotCheckedYet:  typemapper.OptionalStringToPgtypeText(order.PartsNotCheckedYet()),
		PasscodeOrPattern:   passcodeOrPattern,
		IsPatternLocked:     isPatternLocked,
		DownPaymentAmount:   downPaymentAmount,
		DownPaymentMethodID: downPaymentMethodID,
	}, nil
}

func (r *SQLRepairOrderRepository) attachRepairOrderDamages(
	ctx context.Context,
	qtx *gensql.Queries,
	order domain.Order,
) error {
	params := make([]gensql.AddDamagesToRepairOrderParams, 0, len(order.Damages()))
	for _, damage := range order.Damages() {
		params = append(params, gensql.AddDamagesToRepairOrderParams{
			RepairOrderDamageID: typemapper.UUIDToPgtypeUUID(damage.ID()),
			RepairOrderID:       typemapper.UUIDToPgtypeUUID(order.ID()),
			DamageName:          damage.Name(),
		})
	}

	n, err := qtx.AddDamagesToRepairOrder(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to add damages to repair order: %w", err)
	}

	if n < int64(len(order.Damages())) {
		return errors.New("failed to add all damages to repair order")
	}

	return nil
}

func (r *SQLRepairOrderRepository) attachRepairOrderPhoneConditions(
	ctx context.Context,
	qtx *gensql.Queries,
	order domain.Order,
) error {
	params := make([]gensql.AddPhoneConditionsToRepairOrderParams, 0, len(order.PhoneConditions()))
	for _, phoneCondition := range order.PhoneConditions() {
		params = append(params, gensql.AddPhoneConditionsToRepairOrderParams{
			RepairOrderPhoneConditionID: typemapper.UUIDToPgtypeUUID(phoneCondition.ID()),
			RepairOrderID:               typemapper.UUIDToPgtypeUUID(order.ID()),
			PhoneConditionName:          phoneCondition.Name(),
		})
	}

	n, err := qtx.AddPhoneConditionsToRepairOrder(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to add phone conditions to repair order: %w", err)
	}

	if n < int64(len(order.PhoneConditions())) {
		return errors.New("failed to add all phone conditions to repair order")
	}

	return nil
}

func (r *SQLRepairOrderRepository) attachRepairOrderPhoneEquipments(
	ctx context.Context,
	qtx *gensql.Queries,
	order domain.Order,
) error {
	params := make([]gensql.AddPhoneEquipmentsToRepairOrderParams, 0, len(order.PhoneEquipments()))
	for _, phoneEquipment := range order.PhoneEquipments() {
		params = append(params, gensql.AddPhoneEquipmentsToRepairOrderParams{
			RepairOrderPhoneEquipmentID: typemapper.UUIDToPgtypeUUID(phoneEquipment.ID()),
			RepairOrderID:               typemapper.UUIDToPgtypeUUID(order.ID()),
			PhoneEquipmentName:          phoneEquipment.Name(),
		})
	}

	n, err := qtx.AddPhoneEquipmentsToRepairOrder(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to add phone equipments to repair order: %w", err)
	}

	if n < int64(len(order.PhoneEquipments())) {
		return errors.New("failed to add all phone equipments to repair order")
	}

	return nil
}

func (r *SQLRepairOrderRepository) attachRepairOrderCosts(
	ctx context.Context,
	qtx *gensql.Queries,
	order domain.Order,
) error {
	params := make([]gensql.AddCostsToRepairOrderParams, 0, len(order.Costs()))
	for _, cost := range order.Costs() {
		params = append(params, gensql.AddCostsToRepairOrderParams{
			RepairOrderCostID: typemapper.UUIDToPgtypeUUID(cost.ID()),
			RepairOrderID:     typemapper.UUIDToPgtypeUUID(order.ID()),
			Amount:            int32(cost.Amount()),
			Reason:            typemapper.OptionalStringToPgtypeText(cost.Reason()),
			CreationTime:      typemapper.TimeToPgtypeTimestamptz(cost.CreationTime()),
		})
	}

	n, err := qtx.AddCostsToRepairOrder(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to add costs to repair order: %w", err)
	}

	if n < int64(len(order.Costs())) {
		return errors.New("failed to add all costs to repair order")
	}

	return nil
}

func (r *SQLRepairOrderRepository) attachRepairOrderPhotos(
	ctx context.Context,
	qtx *gensql.Queries,
	order domain.Order,
) error {
	params := make([]gensql.AddPhotosToRepairOrderParams, 0, len(order.Photos()))
	for _, photo := range order.Photos() {
		url := photo.URL()

		params = append(params, gensql.AddPhotosToRepairOrderParams{
			RepairOrderPhotoID: typemapper.UUIDToPgtypeUUID(photo.ID()),
			RepairOrderID:      typemapper.UUIDToPgtypeUUID(order.ID()),
			PhotoUrl:           url.String(),
		})
	}

	n, err := qtx.AddPhotosToRepairOrder(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to add photos to repair order: %w", err)
	}

	if n < int64(len(order.Photos())) {
		return errors.New("failed to add all photos to repair order")
	}

	return nil
}
