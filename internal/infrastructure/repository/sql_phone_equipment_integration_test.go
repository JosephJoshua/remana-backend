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
	"github.com/JosephJoshua/remana-backend/internal/modules/phoneequipment"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePhoneEquipment(t *testing.T) {
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

	seedCreatePhoneEquipment(
		context.Background(),
		t,
		queries,
		theStoreID,
	)

	t.Run("creates phone equipment in db", func(t *testing.T) {
		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPhoneEquipmentRepository(db)

		s := phoneequipment.NewService(
			locationProvider,
			permissionProviderStub{},
			repo,
		)

		req := &genapi.CreatePhoneEquipmentRequest{
			Name: "phone equipment 1",
		}

		_, err := s.CreatePhoneEquipment(requestCtx, req)

		require.NoError(t, err)
		require.True(t, locationProvider.PhoneEquipmentID.IsSet(), "location provider not called with phone equipment id")

		phoneEquipmentID := locationProvider.PhoneEquipmentID.MustGet()
		phoneEquipment, err := queries.GetPhoneEquipmentForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(phoneEquipmentID),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("phone equipment with ID %s not found in db", phoneEquipmentID.String())
		}

		require.NoError(t, err)

		assert.Equal(t, phoneEquipmentID, typemapper.MustPgtypeUUIDToUUID(phoneEquipment.PhoneEquipmentID))
		assert.Equal(t, theStoreID, typemapper.MustPgtypeUUIDToUUID(phoneEquipment.StoreID))
		assert.Equal(t, req.Name, phoneEquipment.PhoneEquipmentName)
	})

	t.Run("returns bad request when name is taken", func(t *testing.T) {
		const (
			theName = "phone equipment 1"
		)

		_, err := queries.SeedPhoneEquipment(context.Background(), gensql.SeedPhoneEquipmentParams{
			PhoneEquipmentID:   typemapper.UUIDToPgtypeUUID(uuid.New()),
			PhoneEquipmentName: theName,
			StoreID:            typemapper.UUIDToPgtypeUUID(theStoreID),
		})
		require.NoError(t, err)

		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPhoneEquipmentRepository(db)

		s := phoneequipment.NewService(
			locationProvider,
			permissionProviderStub{},
			repo,
		)

		req := &genapi.CreatePhoneEquipmentRequest{
			Name: theName,
		}

		_, err = s.CreatePhoneEquipment(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})

	t.Run("returns bad request when name is taken (case insensitive)", func(t *testing.T) {
		const (
			theNameInDB = "phone equipment 1"
			theNewName  = "Phone Equipment 1"
		)

		_, err := queries.SeedPhoneEquipment(context.Background(), gensql.SeedPhoneEquipmentParams{
			PhoneEquipmentID:   typemapper.UUIDToPgtypeUUID(uuid.New()),
			PhoneEquipmentName: theNameInDB,
			StoreID:            typemapper.UUIDToPgtypeUUID(theStoreID),
		})
		require.NoError(t, err)

		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPhoneEquipmentRepository(db)

		s := phoneequipment.NewService(
			locationProvider,
			permissionProviderStub{},
			repo,
		)

		req := &genapi.CreatePhoneEquipmentRequest{
			Name: theNewName,
		}

		_, err = s.CreatePhoneEquipment(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})
}

func seedCreatePhoneEquipment(
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
