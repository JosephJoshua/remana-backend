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
	"github.com/JosephJoshua/remana-backend/internal/modules/phonecondition"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePhoneCondition(t *testing.T) {
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

	seedCreateRole(
		context.Background(),
		t,
		queries,
		theStoreID,
	)

	t.Run("creates phone condition in db", func(t *testing.T) {
		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPhoneConditionRepository(db)

		s := phonecondition.NewService(
			locationProvider,
			repo,
		)

		req := &genapi.CreatePhoneConditionRequest{
			Name: "phone condition 1",
		}

		_, err := s.CreatePhoneCondition(requestCtx, req)

		require.NoError(t, err)
		require.True(t, locationProvider.PhoneConditionID.IsSet(), "location provider not called with phone condition id")

		phoneConditionID := locationProvider.PhoneConditionID.MustGet()
		phoneCondition, err := queries.GetPhoneConditionForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(phoneConditionID),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("phone condition with ID %s not found in db", phoneConditionID.String())
		}

		require.NoError(t, err)

		assert.Equal(t, phoneConditionID, typemapper.MustPgtypeUUIDToUUID(phoneCondition.PhoneConditionID))
		assert.Equal(t, theStoreID, typemapper.MustPgtypeUUIDToUUID(phoneCondition.StoreID))
		assert.Equal(t, req.Name, phoneCondition.PhoneConditionName)
	})

	t.Run("returns bad request when name is taken", func(t *testing.T) {
		const (
			theName = "phone condition 1"
		)

		_, err := queries.SeedPhoneCondition(context.Background(), gensql.SeedPhoneConditionParams{
			PhoneConditionID:   typemapper.UUIDToPgtypeUUID(uuid.New()),
			PhoneConditionName: theName,
			StoreID:            typemapper.UUIDToPgtypeUUID(theStoreID),
		})
		require.NoError(t, err)

		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPhoneConditionRepository(db)

		s := phonecondition.NewService(
			locationProvider,
			repo,
		)

		req := &genapi.CreatePhoneConditionRequest{
			Name: theName,
		}

		_, err = s.CreatePhoneCondition(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})

	t.Run("returns bad request when name is taken (case insensitive)", func(t *testing.T) {
		const (
			theNameInDB = "phone condition 1"
			theNewName  = "Phone Condition 1"
		)

		_, err := queries.SeedPhoneCondition(context.Background(), gensql.SeedPhoneConditionParams{
			PhoneConditionID:   typemapper.UUIDToPgtypeUUID(uuid.New()),
			PhoneConditionName: theNameInDB,
			StoreID:            typemapper.UUIDToPgtypeUUID(theStoreID),
		})
		require.NoError(t, err)

		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPhoneConditionRepository(db)

		s := phonecondition.NewService(
			locationProvider,
			repo,
		)

		req := &genapi.CreatePhoneConditionRequest{
			Name: theNewName,
		}

		_, err = s.CreatePhoneCondition(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})
}

func seedCreatePhoneCondition(
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
