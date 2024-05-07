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

	"github.com/JosephJoshua/remana-backend/internal/appconstant"
	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/infrastructure/repository"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/modules/permission"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRole(t *testing.T) {
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

		theLocation = url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/roles/bc80e136-12cb-46f3-8b4f-5ec2b00802d3",
		}
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

	t.Run("creates role in db", func(t *testing.T) {
		locationProvider := testutil.NewResourceLocationProviderStubForPhoneCondition(theLocation)
		repo := repository.NewSQLPermissionRepository(db)

		s := permission.NewService(
			locationProvider,
			repo,
		)

		req := &genapi.CreateRoleRequest{
			Name:         "Admin",
			IsStoreAdmin: true,
		}

		_, err := s.CreateRole(requestCtx, req)

		require.NoError(t, err)
		require.True(t, locationProvider.RoleID.IsSet(), "location provider not called with role id")

		roleID := locationProvider.RoleID.MustGet()
		role, err := queries.GetRoleForTesting(
			context.Background(),
			typemapper.UUIDToPgtypeUUID(roleID),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("role with ID %s not found in db", roleID.String())
		}

		require.NoError(t, err)

		assert.Equal(t, roleID, typemapper.MustPgtypeUUIDToUUID(role.RoleID))
		assert.Equal(t, theStoreID, typemapper.MustPgtypeUUIDToUUID(role.StoreID))
		assert.Equal(t, req.Name, role.RoleName)
		assert.Equal(t, req.IsStoreAdmin, role.IsStoreAdmin)
	})

	t.Run("returns bad request when name is taken", func(t *testing.T) {
		const (
			theName = "admin"
		)

		_, err := queries.SeedRole(context.Background(), gensql.SeedRoleParams{
			RoleID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
			RoleName:     theName,
			IsStoreAdmin: false,
			StoreID:      typemapper.UUIDToPgtypeUUID(theStoreID),
		})
		require.NoError(t, err)

		locationProvider := testutil.NewResourceLocationProviderStubForRole(theLocation)
		repo := repository.NewSQLPermissionRepository(db)

		s := permission.NewService(
			locationProvider,
			repo,
		)

		req := &genapi.CreateRoleRequest{
			Name:         theName,
			IsStoreAdmin: true,
		}

		_, err = s.CreateRole(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})

	t.Run("returns bad request when name is taken (case insensitive)", func(t *testing.T) {
		const (
			theNameInDB = "admin"
			theNewName  = "ADMIN"
		)

		_, err := queries.SeedRole(context.Background(), gensql.SeedRoleParams{
			RoleID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
			RoleName:     theNameInDB,
			IsStoreAdmin: false,
			StoreID:      typemapper.UUIDToPgtypeUUID(theStoreID),
		})
		require.NoError(t, err)

		locationProvider := testutil.NewResourceLocationProviderStubForRole(theLocation)
		repo := repository.NewSQLPermissionRepository(db)

		s := permission.NewService(
			locationProvider,
			repo,
		)

		req := &genapi.CreateRoleRequest{
			Name:         theNewName,
			IsStoreAdmin: true,
		}

		_, err = s.CreateRole(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})
}

func seedCreateRole(
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
