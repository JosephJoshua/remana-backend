//go:build integration
// +build integration

package repository_test

import (
	"context"
	"errors"
	"fmt"
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
		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPermissionRepository(db)

		s := permission.NewService(
			locationProvider,
			repo,
			permissionProviderStub{},
		)

		req := &genapi.CreateRoleRequest{
			Name:         "Admin",
			IsStoreAdmin: true,
		}

		_, err := s.CreateRole(requestCtx, req)

		fmt.Println(locationProvider.RoleID.MustGet())

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

		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPermissionRepository(db)

		s := permission.NewService(
			locationProvider,
			repo,
			permissionProviderStub{},
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

		locationProvider := &testutil.ResourceLocationProviderStub{}
		repo := repository.NewSQLPermissionRepository(db)

		s := permission.NewService(
			locationProvider,
			repo,
			permissionProviderStub{},
		)

		req := &genapi.CreateRoleRequest{
			Name:         theNewName,
			IsStoreAdmin: true,
		}

		_, err = s.CreateRole(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})
}

func TestAssignPermissionsToRole(t *testing.T) {
	logger.Init(zerolog.DebugLevel, appconstant.AppEnvDev)

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
		theStoreID     = uuid.New()
		theRoleID      = uuid.New()
		thePermissions = []struct {
			ID        uuid.UUID
			GroupName string
			Name      string
		}{
			{ID: uuid.New(), GroupName: "permission_one", Name: "test_one"},
			{ID: uuid.New(), GroupName: "permission_two", Name: "test_two"},
		}
		thePermissionsReqItems = []genapi.AssignPermissionsToRoleRequestPermissionsItem{
			{GroupName: "permission_one", Name: "test_one"},
			{GroupName: "permission_two", Name: "test_two"},
		}
	)

	requestCtx := appcontext.NewContextWithUser(
		testutil.RequestContextWithLogger(context.Background()),
		testutil.ModifiedUserDetails(func(details *readmodel.UserDetails) {
			details.Store.ID = theStoreID
		}),
	)

	queries := gensql.New(db)

	seedAssignPermissionsToRole(
		context.Background(),
		t,
		queries,
		theStoreID,
		theRoleID,
		thePermissions,
	)

	t.Run("returns bad request when role doesn't exist", func(t *testing.T) {
		var (
			someRandomID = uuid.New()
		)

		s := permission.NewService(
			&testutil.ResourceLocationProviderStub{},
			repository.NewSQLPermissionRepository(db),
			permissionProviderStub{},
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissionsReqItems,
		}

		params := genapi.AssignPermissionsToRoleParams{
			RoleId: someRandomID,
		}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns bad request when permission doesn't exist", func(t *testing.T) {
		s := permission.NewService(
			&testutil.ResourceLocationProviderStub{},
			repository.NewSQLPermissionRepository(db),
			permissionProviderStub{},
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: []genapi.AssignPermissionsToRoleRequestPermissionsItem{
				thePermissionsReqItems[0],
				{GroupName: "permission_two", Name: "non_existent"},
			},
		}

		params := genapi.AssignPermissionsToRoleParams{
			RoleId: theRoleID,
		}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("assigns permissions to role", func(t *testing.T) {
		s := permission.NewService(
			&testutil.ResourceLocationProviderStub{},
			repository.NewSQLPermissionRepository(db),
			permissionProviderStub{},
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissionsReqItems,
		}

		params := genapi.AssignPermissionsToRoleParams{
			RoleId: theRoleID,
		}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		require.NoError(t, err)

		permissionIDs := make([]uuid.UUID, 0, len(thePermissions))
		for _, p := range thePermissions {
			permissionIDs = append(permissionIDs, p.ID)
		}

		n, err := queries.DoesRoleHavePermissions(
			context.Background(),
			gensql.DoesRoleHavePermissionsParams{
				RoleID:        typemapper.UUIDToPgtypeUUID(theRoleID),
				PermissionIds: typemapper.UUIDsToPgtypeUUIDs(permissionIDs),
			},
		)

		require.NoError(t, err)
		assert.Equal(
			t,
			int64(len(permissionIDs)), n,
			"expected role to have %d permissions assigned, got %d", len(permissionIDs), n,
		)
	})
}

func TestCan(t *testing.T) {
	logger.Init(zerolog.DebugLevel, appconstant.AppEnvDev)

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
		theAdminRoleID    = uuid.New()
		theEmployeeRoleID = uuid.New()
		thePermissions    = []permission.Permission{
			permission.CreateRole(),
			permission.CreateDamageType(),
		}
	)

	requestCtx := testutil.RequestContextWithLogger(context.Background())
	queries := gensql.New(db)

	seedCan(
		context.Background(),
		t,
		queries,
		theAdminRoleID,
		theEmployeeRoleID,
		thePermissions,
	)

	t.Run("returns true when role is store admin", func(t *testing.T) {
		var (
			someRandomPermission = permission.AssignPermissionsToRole()
		)

		p := permission.NewProvider(
			repository.NewSQLPermissionRepository(db),
		)

		ok, err := p.Can(requestCtx, theAdminRoleID, someRandomPermission)

		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("returns true when role has permission", func(t *testing.T) {
		p := permission.NewProvider(
			repository.NewSQLPermissionRepository(db),
		)

		ok, err := p.Can(requestCtx, theEmployeeRoleID, thePermissions[0])

		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("returns false when role doesn't have permission", func(t *testing.T) {
		var (
			someOtherPermission = permission.AssignPermissionsToRole()
		)

		p := permission.NewProvider(
			repository.NewSQLPermissionRepository(db),
		)

		ok, err := p.Can(requestCtx, theEmployeeRoleID, someOtherPermission)

		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("returns error when role doesn't exist", func(t *testing.T) {
		var (
			someRandomID = uuid.New()
		)

		p := permission.NewProvider(
			repository.NewSQLPermissionRepository(db),
		)

		_, err := p.Can(requestCtx, someRandomID, thePermissions[0])
		require.Error(t, err)
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

func seedAssignPermissionsToRole(
	ctx context.Context,
	t *testing.T,
	queries *gensql.Queries,
	theStoreID uuid.UUID,
	theRoleID uuid.UUID,
	thePermissions []struct {
		ID        uuid.UUID
		GroupName string
		Name      string
	},
) {
	t.Helper()

	const maxWait = 3 * time.Second

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

	_, err = queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(theRoleID),
		RoleName:     "Not important",
		IsStoreAdmin: false,
		StoreID:      typemapper.UUIDToPgtypeUUID(theStoreID),
	})
	require.NoError(t, err)

	permissionGroups := make(map[string]uuid.UUID)
	for _, permission := range thePermissions {
		if _, ok := permissionGroups[permission.GroupName]; !ok {
			id, err := queries.SeedPermissionGroup(ctx, gensql.SeedPermissionGroupParams{
				PermissionGroupID:   typemapper.UUIDToPgtypeUUID(uuid.New()),
				PermissionGroupName: permission.GroupName,
			})
			require.NoError(t, err)

			permissionGroups[permission.GroupName] = typemapper.MustPgtypeUUIDToUUID(id)
		}
	}

	for _, permission := range thePermissions {
		_, err = queries.SeedPermission(ctx, gensql.SeedPermissionParams{
			PermissionID:          typemapper.UUIDToPgtypeUUID(permission.ID),
			PermissionName:        permission.Name,
			PermissionDisplayName: permission.Name,
			PermissionGroupID:     typemapper.UUIDToPgtypeUUID(permissionGroups[permission.GroupName]),
		})
		require.NoError(t, err)
	}
}

func seedCan(
	ctx context.Context,
	t *testing.T,
	queries *gensql.Queries,
	theAdminRoleID uuid.UUID,
	theEmployeeRoleID uuid.UUID,
	thePermissions []permission.Permission,
) {
	t.Helper()

	const maxWait = 3 * time.Second

	ctx, cancel := context.WithTimeout(ctx, maxWait)
	defer cancel()

	storeID, err := queries.SeedStore(ctx, gensql.SeedStoreParams{
		StoreID:      typemapper.UUIDToPgtypeUUID(uuid.New()),
		StoreName:    "Not important",
		StoreCode:    "not-important",
		StoreAddress: "Not important",
		PhoneNumber:  "+6281234567890",
	})
	require.NoError(t, err)

	_, err = queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(theAdminRoleID),
		RoleName:     "Not important",
		IsStoreAdmin: true,
		StoreID:      storeID,
	})
	require.NoError(t, err)

	_, err = queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(theEmployeeRoleID),
		RoleName:     "Not important",
		IsStoreAdmin: false,
		StoreID:      storeID,
	})
	require.NoError(t, err)

	permissionGroups := make(map[string]uuid.UUID)
	for _, permission := range thePermissions {
		groupName := permission.GroupName()

		if _, ok := permissionGroups[groupName]; !ok {
			id, err := queries.SeedPermissionGroup(ctx, gensql.SeedPermissionGroupParams{
				PermissionGroupID:   typemapper.UUIDToPgtypeUUID(uuid.New()),
				PermissionGroupName: groupName,
			})
			require.NoError(t, err)

			permissionGroups[groupName] = typemapper.MustPgtypeUUIDToUUID(id)
		}
	}

	for _, permission := range thePermissions {
		permissionID := uuid.New()

		_, err = queries.SeedPermission(ctx, gensql.SeedPermissionParams{
			PermissionID:          typemapper.UUIDToPgtypeUUID(permissionID),
			PermissionName:        permission.Name(),
			PermissionDisplayName: permission.Name(),
			PermissionGroupID:     typemapper.UUIDToPgtypeUUID(permissionGroups[permission.GroupName()]),
		})
		require.NoError(t, err)

		err = queries.SeedRolePermission(ctx, gensql.SeedRolePermissionParams{
			RoleID:       typemapper.UUIDToPgtypeUUID(theEmployeeRoleID),
			PermissionID: typemapper.UUIDToPgtypeUUID(permissionID),
		})
		require.NoError(t, err)
	}
}
