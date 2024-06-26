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
	"github.com/JosephJoshua/remana-backend/internal/modules/technician"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTechnician(t *testing.T) {
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

	seedCreateTechnician(
		context.Background(),
		t,
		queries,
		theStoreID,
	)

	t.Run("creates technician in db", func(t *testing.T) {
		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLTechnicianRepository(db)

		s := technician.NewService(
			locationProvider,
			permissionProviderStub{},
			repo,
		)

		req := &genapi.CreateTechnicianRequest{
			Name: "technician 1",
		}

		_, err := s.CreateTechnician(requestCtx, req)

		require.NoError(t, err)
		require.True(t, locationProvider.TechnicianID.IsSet(), "location provider not called with technician id")

		technicianID := locationProvider.TechnicianID.MustGet()
		technician, err := queries.GetTechnicianForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(technicianID),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("technician with ID %s not found in db", technicianID.String())
		}

		require.NoError(t, err)

		assert.Equal(t, technicianID, typemapper.MustPgtypeUUIDToUUID(technician.TechnicianID))
		assert.Equal(t, theStoreID, typemapper.MustPgtypeUUIDToUUID(technician.StoreID))
		assert.Equal(t, req.Name, technician.TechnicianName)
	})

	t.Run("returns bad request when name is taken", func(t *testing.T) {
		const (
			theName = "technician 1"
		)

		_, err := queries.SeedTechnician(context.Background(), gensql.SeedTechnicianParams{
			TechnicianID:   typemapper.UUIDToPgtypeUUID(uuid.New()),
			TechnicianName: theName,
			StoreID:        typemapper.UUIDToPgtypeUUID(theStoreID),
		})
		require.NoError(t, err)

		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLTechnicianRepository(db)

		s := technician.NewService(
			locationProvider,
			permissionProviderStub{},
			repo,
		)

		req := &genapi.CreateTechnicianRequest{
			Name: theName,
		}

		_, err = s.CreateTechnician(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})

	t.Run("returns bad request when name is taken (case insensitive)", func(t *testing.T) {
		const (
			theNameInDB = "technician 1"
			theNewName  = "Technician 1"
		)

		_, err := queries.SeedTechnician(context.Background(), gensql.SeedTechnicianParams{
			TechnicianID:   typemapper.UUIDToPgtypeUUID(uuid.New()),
			TechnicianName: theNameInDB,
			StoreID:        typemapper.UUIDToPgtypeUUID(theStoreID),
		})
		require.NoError(t, err)

		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLTechnicianRepository(db)

		s := technician.NewService(
			locationProvider,
			permissionProviderStub{},
			repo,
		)

		req := &genapi.CreateTechnicianRequest{
			Name: theNewName,
		}

		_, err = s.CreateTechnician(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})
}

func seedCreateTechnician(
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
