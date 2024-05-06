//go:build integration
// +build integration

package repository_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/infrastructure/repository"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/modules/repairorder"
	"github.com/JosephJoshua/remana-backend/internal/modules/shared"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type damage struct {
	id   uuid.UUID
	name string
}

type phoneCondition struct {
	id   uuid.UUID
	name string
}

type phoneEquipment struct {
	id   uuid.UUID
	name string
}

func TestCreateRepairOrder(t *testing.T) {
	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)

	pool, initErr := testutil.StartDockerPool()
	require.NoError(t, initErr, "error starting docker pool")

	postgresResource, db, initErr := testutil.StartPostgresContainer(pool)
	require.NoError(t, initErr, "error starting postgres container")

	t.Cleanup(func() {
		if purgeErr := testutil.PurgeDockerResources(pool, []*dockertest.Resource{postgresResource}); purgeErr != nil {
			t.Fatalf("failed to purge docker resources: %v", initErr)
		}
	})

	initErr = testutil.MigratePostgres(context.Background(), db)
	require.NoError(t, initErr, "error migrating database")

	var (
		theCreationTime = time.Unix(1713917762, 0)
		theLocation     = url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/repair-orders/61821d3a-dbdd-4f41-a2c6-14faf7b57c67",
		}

		theStoreID         = uuid.New()
		theSalesID         = uuid.New()
		theTechnicianID    = uuid.New()
		thePaymentMethodID = uuid.New()

		theDamage = damage{
			id:   uuid.New(),
			name: "Broken Screen",
		}

		thePhoneCondition = phoneCondition{
			id:   uuid.New(),
			name: "Screen scratched",
		}

		theEquipment = phoneEquipment{
			id:   uuid.New(),
			name: "Battery",
		}

		otherStoreID              = uuid.New()
		otherStoreSalesID         = uuid.New()
		otherStoreTechnicianID    = uuid.New()
		otherStorePaymentMethodID = uuid.New()

		otherStoreDamage = damage{
			id:   uuid.New(),
			name: "Broken Screen",
		}

		otherStorePhoneCondition = phoneCondition{
			id:   uuid.New(),
			name: "Screen scratched",
		}

		otherStoreEquipment = phoneEquipment{
			id:   uuid.New(),
			name: "Battery",
		}
	)

	requestCtx := appcontext.NewContextWithUser(
		testutil.RequestContextWithLogger(context.Background()),
		&readmodel.UserDetails{
			ID:       uuid.New(),
			Username: "not important",
			Role: readmodel.UserDetailsRole{
				ID:           uuid.New(),
				Name:         "not important",
				IsStoreAdmin: true,
			},
			Store: readmodel.UserDetailsStore{
				ID:   theStoreID,
				Name: "not important",
				Code: "not-important",
			},
		})

	queries := gensql.New(db)

	seedCreateRepairOrder(
		context.Background(),
		t,
		queries,
		theStoreID,
		otherStoreID,
		theSalesID,
		otherStoreSalesID,
		theTechnicianID,
		otherStoreTechnicianID,
		thePaymentMethodID,
		otherStorePaymentMethodID,
		theDamage,
		otherStoreDamage,
		thePhoneCondition,
		otherStorePhoneCondition,
		theEquipment,
		otherStoreEquipment,
	)

	validRequest := func() genapi.CreateRepairOrderRequest {
		return genapi.CreateRepairOrderRequest{
			CustomerName:       "John Doe",
			ContactPhoneNumber: "08123456789",
			PhoneType:          "iPhone 12",
			Color:              "Black",
			SalesID:            theSalesID,
			TechnicianID:       theTechnicianID,
			InitialCost:        100,
			DamageTypes:        []uuid.UUID{theDamage.id},
			PhoneConditions:    []uuid.UUID{thePhoneCondition.id},
			PhoneEquipments:    []uuid.UUID{theEquipment.id},
			Photos:             []url.URL{{Host: "example.com", Scheme: "http"}},
			Imei:               genapi.NewOptString("123456789012345"),
			PartsNotCheckedYet: genapi.NewOptString("Battery"),
			Passcode: genapi.NewOptCreateRepairOrderRequestPasscode(genapi.CreateRepairOrderRequestPasscode{
				Value:           "1234",
				IsPatternLocked: true,
			}),
			DownPayment: genapi.NewOptCreateRepairOrderRequestDownPayment(genapi.CreateRepairOrderRequestDownPayment{
				Amount: 100,
				Method: thePaymentMethodID,
			}),
		}
	}

	t.Run("creates repair order in db", func(t *testing.T) {
		timeProvider := testutil.NewTimeProviderStub(theCreationTime)
		locationProvider := testutil.NewResourceLocationProviderStubForRepairOrder(theLocation)
		slugProvider := testutil.NewRepairOrderSlugProviderStub("some-slug", nil)

		repo := repository.NewSQLRepairOrderRepository(db)
		s := repairorder.NewService(timeProvider, locationProvider, repo, slugProvider)

		req := validRequest()

		_, err := s.CreateRepairOrder(requestCtx, &req)

		require.NoError(t, err)
		require.True(t, locationProvider.RepairOrderID.IsSet(), "location provider not called with repair order id")

		orderID := locationProvider.RepairOrderID.MustGet()
		order, err := queries.GetRepairOrderForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(orderID),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("repair order with ID %s not found in db", orderID.String())
		}

		require.NoError(t, err)

		assert.NotEmpty(t, order.Slug)
		assert.Equal(t, theStoreID, typemapper.MustPgtypeUUIDToUUID(order.StoreID))
		assert.Equal(t, req.CustomerName, order.CustomerName)
		assert.Equal(t, req.Color, order.Color)
		assert.NotEmpty(t, order.ContactNumber)
		assert.Equal(t, req.PhoneType, order.PhoneType)

		require.True(t, order.SalesID.Valid)
		assert.Equal(t, req.SalesID, typemapper.MustPgtypeUUIDToUUID(order.SalesID))

		require.True(t, order.TechnicianID.Valid)
		assert.Equal(t, req.TechnicianID, typemapper.MustPgtypeUUIDToUUID(order.TechnicianID))

		assert.Equal(t, req.Imei.Value, order.Imei.String)
		assert.Equal(t, req.PartsNotCheckedYet.Value, order.PartsNotCheckedYet.String)
		assert.Equal(t, req.Passcode.Value.IsPatternLocked, order.IsPatternLocked.Bool)
		assert.Equal(t, req.Passcode.Value.Value, order.PasscodeOrPattern.String)
		assert.Equal(t, int32(req.DownPayment.Value.Amount), order.DownPaymentAmount.Int32)

		require.True(t, order.DownPaymentMethodID.Valid)
		assert.Equal(t, req.DownPayment.Value.Method, typemapper.MustPgtypeUUIDToUUID(order.DownPaymentMethodID))

		assert.False(t, order.PickUpTime.Valid)
		assert.False(t, order.CompletionTime.Valid)
		assert.False(t, order.CancellationTime.Valid)
		assert.False(t, order.CancellationReason.Valid)
		assert.False(t, order.WarrantyDays.Valid)
		assert.False(t, order.ConfirmationTime.Valid)
		assert.False(t, order.ConfirmationContent.Valid)
		assert.False(t, order.RepaymentAmount.Valid)
		assert.False(t, order.RepaymentMethodID.Valid)

		damages, err := queries.GetRepairOrderDamagesForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(orderID),
		)

		require.NoError(t, err)

		assert.Equal(t, 1, len(damages))
		assert.Equal(t, theDamage.name, damages[0].DamageName)

		phoneConditions, err := queries.GetRepairOrderPhoneConditionsForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(orderID),
		)

		require.NoError(t, err)

		assert.Equal(t, 1, len(phoneConditions))
		assert.Equal(t, thePhoneCondition.name, phoneConditions[0].PhoneConditionName)

		equipments, err := queries.GetRepairOrderPhoneEquipmentsForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(orderID),
		)

		require.NoError(t, err)

		assert.Equal(t, 1, len(equipments))
		assert.Equal(t, theEquipment.name, equipments[0].PhoneEquipmentName)

		costs, err := queries.GetRepairOrderCostsForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(orderID),
		)

		require.NoError(t, err)

		assert.Equal(t, 1, len(costs))
		assert.Equal(t, int32(req.InitialCost), costs[0].Amount)
		assert.Empty(t, costs[0].Reason.String)

		photos, err := queries.GetRepairOrderPhotosForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(orderID),
		)

		require.NoError(t, err)
		assert.Equal(t, 1, len(photos))

		_, err = url.Parse(photos[0].PhotoUrl)
		require.NoError(t, err, "invalid photo url")
	})

	t.Run("returns bad request", func(t *testing.T) {
		var someRandomID = uuid.New()

		testCases := []struct {
			name  string
			setup func(req *genapi.CreateRepairOrderRequest)
		}{
			{
				name: "when damage type does not exist",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.DamageTypes = []uuid.UUID{someRandomID}
				},
			},
			{
				name: "when phone condition does not exist",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.PhoneConditions = []uuid.UUID{someRandomID}
				},
			},
			{
				name: "when phone equipment does not exist",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.PhoneEquipments = []uuid.UUID{someRandomID}
				},
			},
			{
				name: "when payment method does not exist",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.DownPayment = genapi.NewOptCreateRepairOrderRequestDownPayment(genapi.CreateRepairOrderRequestDownPayment{
						Amount: 100,
						Method: someRandomID,
					})
				},
			},
			{
				name: "when technician does not exist",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.TechnicianID = someRandomID
				},
			},
			{
				name: "when sales does not exist",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.SalesID = someRandomID
				},
			},
			{
				name: "when damage type is from different store",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.DamageTypes = []uuid.UUID{otherStoreDamage.id}
				},
			},
			{
				name: "when phone condition is from different store",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.PhoneConditions = []uuid.UUID{otherStorePhoneCondition.id}
				},
			},
			{
				name: "when phone equipment is from different store",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.PhoneEquipments = []uuid.UUID{otherStoreEquipment.id}
				},
			},
			{
				name: "when payment method is from different store",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.DownPayment = genapi.NewOptCreateRepairOrderRequestDownPayment(genapi.CreateRepairOrderRequestDownPayment{
						Amount: 100,
						Method: otherStorePaymentMethodID,
					})
				},
			},
			{
				name: "when technician is from different store",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.TechnicianID = otherStoreTechnicianID
				},
			},
			{
				name: "when sales is from different store",
				setup: func(req *genapi.CreateRepairOrderRequest) {
					req.SalesID = otherStoreSalesID
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				timeProvider := testutil.NewTimeProviderStub(theCreationTime)
				locationProvider := testutil.NewResourceLocationProviderStubForRepairOrder(theLocation)
				slugProvider := testutil.NewRepairOrderSlugProviderStub("some-slug", nil)
				repo := repository.NewSQLRepairOrderRepository(db)

				s := repairorder.NewService(timeProvider, locationProvider, repo, slugProvider)

				req := validRequest()
				tc.setup(&req)

				_, err := s.CreateRepairOrder(requestCtx, &req)
				testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
			})
		}
	})
}

func seedCreateRepairOrder(
	ctx context.Context,
	t *testing.T,
	queries *gensql.Queries,
	theStoreID uuid.UUID,
	otherStoreID uuid.UUID,
	theSalesID uuid.UUID,
	otherStoreSalesID uuid.UUID,
	theTechnicianID uuid.UUID,
	otherStoreTechnicianID uuid.UUID,
	thePaymentMethodID uuid.UUID,
	otherStorePaymentMethodID uuid.UUID,
	theDamage damage,
	otherStoreDamage damage,
	thePhoneCondition phoneCondition,
	otherStorePhoneCondition phoneCondition,
	theEquipment phoneEquipment,
	otherStoreEquipment phoneEquipment,
) {
	t.Helper()

	const maxWait = 5 * time.Second

	ctx, cancel := context.WithTimeout(ctx, maxWait)
	defer cancel()

	_, err := queries.SeedStore(ctx, gensql.SeedStoreParams{
		StoreID:      typemapper.UUIDToPgtypeUUID(theStoreID),
		StoreName:    "Not important",
		StoreCode:    "not-important",
		StoreAddress: "Not important",
		PhoneNumber:  "+6281234567890",
	})
	require.NoError(t, err)

	_, err = queries.SeedStore(ctx, gensql.SeedStoreParams{
		StoreID:      typemapper.UUIDToPgtypeUUID(otherStoreID),
		StoreName:    "Not important",
		StoreCode:    "not-important",
		StoreAddress: "Not important",
		PhoneNumber:  "+6281234567890",
	})
	require.NoError(t, err)

	_, err = queries.SeedSales(ctx, gensql.SeedSalesParams{
		SalesID:   typemapper.UUIDToPgtypeUUID(theSalesID),
		SalesName: "Not important",
		StoreID:   typemapper.UUIDToPgtypeUUID(theStoreID),
	})
	require.NoError(t, err)

	_, err = queries.SeedSales(ctx, gensql.SeedSalesParams{
		SalesID:   typemapper.UUIDToPgtypeUUID(otherStoreSalesID),
		SalesName: "Not important",
		StoreID:   typemapper.UUIDToPgtypeUUID(otherStoreID),
	})
	require.NoError(t, err)

	_, err = queries.SeedTechnician(ctx, gensql.SeedTechnicianParams{
		TechnicianID:   typemapper.UUIDToPgtypeUUID(theTechnicianID),
		TechnicianName: "Not important",
		StoreID:        typemapper.UUIDToPgtypeUUID(theStoreID),
	})
	require.NoError(t, err)

	_, err = queries.SeedTechnician(ctx, gensql.SeedTechnicianParams{
		TechnicianID:   typemapper.UUIDToPgtypeUUID(otherStoreTechnicianID),
		TechnicianName: "Not important",
		StoreID:        typemapper.UUIDToPgtypeUUID(otherStoreID),
	})
	require.NoError(t, err)

	_, err = queries.SeedPaymentMethod(ctx, gensql.SeedPaymentMethodParams{
		PaymentMethodID:   typemapper.UUIDToPgtypeUUID(thePaymentMethodID),
		PaymentMethodName: "Not important",
		StoreID:           typemapper.UUIDToPgtypeUUID(theStoreID),
	})
	require.NoError(t, err)

	_, err = queries.SeedPaymentMethod(ctx, gensql.SeedPaymentMethodParams{
		PaymentMethodID:   typemapper.UUIDToPgtypeUUID(otherStorePaymentMethodID),
		PaymentMethodName: "Not important",
		StoreID:           typemapper.UUIDToPgtypeUUID(otherStoreID),
	})
	require.NoError(t, err)

	_, err = queries.SeedDamageType(ctx, gensql.SeedDamageTypeParams{
		DamageTypeID:   typemapper.UUIDToPgtypeUUID(theDamage.id),
		StoreID:        typemapper.UUIDToPgtypeUUID(theStoreID),
		DamageTypeName: theDamage.name,
	})
	require.NoError(t, err)

	_, err = queries.SeedDamageType(ctx, gensql.SeedDamageTypeParams{
		DamageTypeID:   typemapper.UUIDToPgtypeUUID(otherStoreDamage.id),
		StoreID:        typemapper.UUIDToPgtypeUUID(otherStoreID),
		DamageTypeName: theDamage.name,
	})
	require.NoError(t, err)

	_, err = queries.SeedPhoneCondition(ctx, gensql.SeedPhoneConditionParams{
		PhoneConditionID:   typemapper.UUIDToPgtypeUUID(thePhoneCondition.id),
		StoreID:            typemapper.UUIDToPgtypeUUID(theStoreID),
		PhoneConditionName: thePhoneCondition.name,
	})
	require.NoError(t, err)

	_, err = queries.SeedPhoneCondition(ctx, gensql.SeedPhoneConditionParams{
		PhoneConditionID:   typemapper.UUIDToPgtypeUUID(otherStorePhoneCondition.id),
		StoreID:            typemapper.UUIDToPgtypeUUID(otherStoreID),
		PhoneConditionName: thePhoneCondition.name,
	})
	require.NoError(t, err)

	_, err = queries.SeedPhoneEquipment(ctx, gensql.SeedPhoneEquipmentParams{
		PhoneEquipmentID:   typemapper.UUIDToPgtypeUUID(theEquipment.id),
		StoreID:            typemapper.UUIDToPgtypeUUID(theStoreID),
		PhoneEquipmentName: theEquipment.name,
	})
	require.NoError(t, err)

	_, err = queries.SeedPhoneEquipment(ctx, gensql.SeedPhoneEquipmentParams{
		PhoneEquipmentID:   typemapper.UUIDToPgtypeUUID(otherStoreEquipment.id),
		StoreID:            typemapper.UUIDToPgtypeUUID(otherStoreID),
		PhoneEquipmentName: theEquipment.name,
	})
	require.NoError(t, err)
}
