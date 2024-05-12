//go:build integration
// +build integration

package repository_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/appconstant"
	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/infrastructure/repository"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/modules/paymentmethod"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePaymentMethod(t *testing.T) {
	logger.Init(zerolog.ErrorLevel, appconstant.AppEnvDev)

	pool, initErr := testutil.StartDockerPool()
	require.NoError(t, initErr, "error starting docker pool")

	postgresResource, db, initErr := testutil.StartPostgresContainer(pool)
	require.NoError(t, initErr, "error starting postgres container")

	t.Cleanup(func() {
		if purgeErr := testutil.PurgeDockerResources(pool, []*dockertest.Resource{postgresResource}); purgeErr != nil {
			t.Fatalf("failed to purge docker resources: %v", purgeErr)
		}
	})

	initErr = testutil.MigratePostgres(context.Background(), db)
	require.NoError(t, initErr, "error migrating database")

	var (
		theStoreID = uuid.New()
	)

	requestCtx := appcontext.NewContextWithUser(
		testutil.RequestContextWithLogger(context.Background()),
		testutil.ModifiedUserDetails(func(details *readmodel.UserDetails) {
			details.Store.ID = theStoreID
		}),
	)

	queries := gensql.New(db)

	seedCreatePaymentMethod(
		context.Background(),
		t,
		queries,
		theStoreID,
	)

	t.Run("creates payment method in db", func(t *testing.T) {
		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPaymentMethodRepository(db)

		s := paymentmethod.NewService(
			locationProvider,
			repo,
		)

		req := &genapi.CreatePaymentMethodRequest{
			Name: "payment method 1",
		}

		_, err := s.CreatePaymentMethod(requestCtx, req)

		require.NoError(t, err)
		require.True(t, locationProvider.PaymentMethodID.IsSet(), "location provider not called with payment method id")

		paymentMethodID := locationProvider.PaymentMethodID.MustGet()
		paymentMethod, err := queries.GetPaymentMethodForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(paymentMethodID),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("payment method with ID %s not found in db", paymentMethodID.String())
		}

		require.NoError(t, err)

		assert.Equal(t, paymentMethodID, typemapper.MustPgtypeUUIDToUUID(paymentMethod.PaymentMethodID))
		assert.Equal(t, theStoreID, typemapper.MustPgtypeUUIDToUUID(paymentMethod.StoreID))
		assert.Equal(t, req.Name, paymentMethod.PaymentMethodName)
	})

	t.Run("returns bad request when name is taken", func(t *testing.T) {
		const (
			theName = "payment method 1"
		)

		_, err := queries.SeedPaymentMethod(context.Background(), gensql.SeedPaymentMethodParams{
			PaymentMethodID:   typemapper.UUIDToPgtypeUUID(uuid.New()),
			PaymentMethodName: theName,
			StoreID:           typemapper.UUIDToPgtypeUUID(theStoreID),
		})
		require.NoError(t, err)

		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPaymentMethodRepository(db)

		s := paymentmethod.NewService(
			locationProvider,
			repo,
		)

		req := &genapi.CreatePaymentMethodRequest{
			Name: theName,
		}

		_, err = s.CreatePaymentMethod(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})

	t.Run("returns bad request when name is taken (case insensitive)", func(t *testing.T) {
		const (
			theNameInDB = "payment method 1"
			theNewName  = "Payment Method 1"
		)

		_, err := queries.SeedPaymentMethod(context.Background(), gensql.SeedPaymentMethodParams{
			PaymentMethodID:   typemapper.UUIDToPgtypeUUID(uuid.New()),
			PaymentMethodName: theNameInDB,
			StoreID:           typemapper.UUIDToPgtypeUUID(theStoreID),
		})
		require.NoError(t, err)

		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPaymentMethodRepository(db)

		s := paymentmethod.NewService(
			locationProvider,
			repo,
		)

		req := &genapi.CreatePaymentMethodRequest{
			Name: theNewName,
		}

		_, err = s.CreatePaymentMethod(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})
}

func seedCreatePaymentMethod(
	ctx context.Context,
	t *testing.T,
	queries *gensql.Queries,
	theStoreID uuid.UUID,
) {
	t.Helper()

	const maxWait = 2 * time.Second

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
}
