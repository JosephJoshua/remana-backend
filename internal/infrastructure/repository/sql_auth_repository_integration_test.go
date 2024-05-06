//go:build integration
// +build integration

package repository_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/infrastructure/repository"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
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

func TestLogin(t *testing.T) {
	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)
	requestCtx := testutil.RequestContextWithLogger(context.Background())

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

	const (
		theStoreCode = "store-code"

		theAdminUsername = "admin"
		theAdminPassword = "AdminPassword12345"

		theEmployeeUsername = "employee"
		theEmployeePassword = "EmployeePassword12345"
	)

	queries := gensql.New(db)

	seedLogin(
		context.Background(),
		t,
		queries,
		theStoreCode,
		theAdminUsername,
		theAdminPassword,
		theEmployeeUsername,
		theEmployeePassword,
	)

	t.Run("returns type 'admin' when logging in as admin", func(t *testing.T) {
		req := &genapi.LoginCredentials{
			Username:  theAdminUsername,
			Password:  theAdminPassword,
			StoreCode: theStoreCode,
		}

		repo := repository.NewSQLAuthRepository(db)
		s := auth.NewService(
			serviceSessionManagerStub{},
			loginCodePromptManagerStub{},
			repo,
			testutil.PasswordHasherStub{},
		)

		got, err := s.Login(requestCtx, req)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, genapi.LoginResponseTypeAdmin, got.Type)
	})

	t.Run("returns type 'employee' when logging in as employee", func(t *testing.T) {
		req := &genapi.LoginCredentials{
			Username:  theEmployeeUsername,
			Password:  theEmployeePassword,
			StoreCode: theStoreCode,
		}

		repo := repository.NewSQLAuthRepository(db)
		s := auth.NewService(
			serviceSessionManagerStub{},
			loginCodePromptManagerStub{},
			repo,
			testutil.PasswordHasherStub{},
		)

		got, err := s.Login(requestCtx, req)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, genapi.LoginResponseTypeEmployee, got.Type)
	})

	t.Run("returns unauthorized", func(t *testing.T) {
		testCases := []struct {
			name  string
			setup func(req *genapi.LoginCredentials)
		}{
			{
				name: "when username is wrong",
				setup: func(req *genapi.LoginCredentials) {
					req.Username = "random-username"
				},
			},
			{
				name: "when password is wrong",
				setup: func(req *genapi.LoginCredentials) {
					req.Password = "random-password"
				},
			},
			{
				name: "when store code is wrong",
				setup: func(req *genapi.LoginCredentials) {
					req.StoreCode = "random-store-code"
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := &genapi.LoginCredentials{
					Username:  theAdminUsername,
					Password:  theAdminPassword,
					StoreCode: theStoreCode,
				}

				tc.setup(req)

				repo := repository.NewSQLAuthRepository(db)
				s := auth.NewService(
					serviceSessionManagerStub{},
					loginCodePromptManagerStub{},
					repo,
					testutil.PasswordHasherStub{},
				)

				_, err := s.Login(requestCtx, req)
				testutil.AssertAPIStatusCode(t, http.StatusUnauthorized, err)
			})
		}
	})
}

func TestLoginCodePrompt(t *testing.T) {
	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)
	requestCtx := testutil.RequestContextWithLogger(context.Background())

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
		theUserID    = uuid.New()
		theLoginCode = "A1B2C3D4"
	)

	queries := gensql.New(db)

	seedLoginCodePrompt(
		context.Background(),
		t,
		queries,
		theUserID,
		theLoginCode,
	)

	// ORDER MATTERS!

	t.Run("returns bad request when user is not found", func(t *testing.T) {
		req := &genapi.LoginCodePrompt{
			LoginCode: theLoginCode,
		}

		var someRandomID = uuid.New()

		repo := repository.NewSQLAuthRepository(db)
		loginCodePromptManager := &loginCodePromptManagerStub{userID: someRandomID}

		s := auth.NewService(
			serviceSessionManagerStub{},
			loginCodePromptManager,
			repo,
			testutil.PasswordHasherStub{},
		)

		err := s.LoginCodePrompt(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns bad request when login code is wrong", func(t *testing.T) {
		req := &genapi.LoginCodePrompt{
			LoginCode: "AAAAAAAA",
		}

		repo := repository.NewSQLAuthRepository(db)
		loginCodePromptManager := &loginCodePromptManagerStub{userID: theUserID}

		s := auth.NewService(
			serviceSessionManagerStub{},
			loginCodePromptManager,
			repo,
			testutil.PasswordHasherStub{},
		)

		err := s.LoginCodePrompt(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("deletes login code when login code is correct", func(t *testing.T) {
		_, err := queries.GetLoginCodeByUserIDAndCode(context.Background(), gensql.GetLoginCodeByUserIDAndCodeParams{
			UserID:    typemapper.UUIDToPgtypeUUID(theUserID),
			LoginCode: theLoginCode,
		})
		require.NoError(t, err, pgx.ErrNoRows)

		req := &genapi.LoginCodePrompt{
			LoginCode: theLoginCode,
		}

		repo := repository.NewSQLAuthRepository(db)
		loginCodePromptManager := &loginCodePromptManagerStub{userID: theUserID}

		s := auth.NewService(
			serviceSessionManagerStub{},
			loginCodePromptManager,
			repo,
			testutil.PasswordHasherStub{},
		)

		err = s.LoginCodePrompt(requestCtx, req)
		require.NoError(t, err)

		_, err = queries.GetLoginCodeByUserIDAndCode(context.Background(), gensql.GetLoginCodeByUserIDAndCodeParams{
			UserID:    typemapper.UUIDToPgtypeUUID(theUserID),
			LoginCode: theLoginCode,
		})
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestHandleSessionCookie(t *testing.T) {
	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)
	requestCtx := testutil.RequestContextWithLogger(context.Background())

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
		theUserDetails = readmodel.UserDetails{
			ID:       uuid.New(),
			Username: "johndoe",
			Role: readmodel.UserDetailsRole{
				ID:           uuid.New(),
				Name:         "Admin",
				IsStoreAdmin: true,
			},
			Store: readmodel.UserDetailsStore{
				ID:   uuid.New(),
				Name: "Store 1",
				Code: "store-one",
			},
		}
	)

	queries := gensql.New(db)

	seedHandleSessionCookie(
		context.Background(),
		t,
		queries,
		theUserDetails,
	)

	t.Run("returns context with user details when session is valid", func(t *testing.T) {
		sm := securityHandlerSessionManagerStub{userID: theUserDetails.ID}
		repo := repository.NewSQLAuthRepository(db)

		s := auth.NewSecurityHandler(sm, repo)

		ctx, err := s.HandleSessionCookie(requestCtx, "", genapi.SessionCookie{APIKey: ""})
		require.NoError(t, err)

		got, ok := appcontext.GetUserFromContext(ctx)
		require.True(t, ok)
		require.NotNil(t, got)

		assert.EqualExportedValues(t, theUserDetails, *got)
	})

	t.Run("returns unauthorized when user is missing", func(t *testing.T) {
		var someRandomID = uuid.New()

		sm := securityHandlerSessionManagerStub{userID: someRandomID}
		repo := repository.NewSQLAuthRepository(db)

		s := auth.NewSecurityHandler(sm, repo)
		_, err := s.HandleSessionCookie(requestCtx, "", genapi.SessionCookie{APIKey: ""})

		testutil.AssertAPIStatusCode(t, http.StatusUnauthorized, err)
	})
}

func seedLogin(
	ctx context.Context,
	t *testing.T,
	queries *gensql.Queries,
	theStoreCode string,
	theAdminUsername string,
	theAdminPassword string,
	theEmployeeUsername string,
	theEmployeePassword string,
) {
	t.Helper()

	const maxWait = 5 * time.Second

	ctx, cancel := context.WithTimeout(ctx, maxWait)
	defer cancel()

	storeID, err := queries.SeedStore(ctx, gensql.SeedStoreParams{
		StoreID:      typemapper.UUIDToPgtypeUUID(uuid.New()),
		StoreName:    "Not important",
		StoreCode:    theStoreCode,
		StoreAddress: "Not important",
		PhoneNumber:  "+6281234567890",
	})
	require.NoError(t, err)

	adminRoleID, err := queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		RoleName:     "Admin",
		StoreID:      storeID,
		IsStoreAdmin: true,
	})
	require.NoError(t, err)

	employeeRoleID, err := queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		RoleName:     "Employee",
		StoreID:      storeID,
		IsStoreAdmin: false,
	})
	require.NoError(t, err)

	_, err = queries.SeedUser(ctx, gensql.SeedUserParams{
		UserID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		Username:     theAdminUsername,
		UserPassword: theAdminPassword,
		RoleID:       adminRoleID,
		StoreID:      storeID,
	})
	require.NoError(t, err)

	_, err = queries.SeedUser(ctx, gensql.SeedUserParams{
		UserID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		Username:     theEmployeeUsername,
		UserPassword: theEmployeePassword,
		RoleID:       employeeRoleID,
		StoreID:      storeID,
	})
	require.NoError(t, err)
}

func seedLoginCodePrompt(
	ctx context.Context,
	t *testing.T,
	queries *gensql.Queries,
	theUserID uuid.UUID,
	theLoginCode string,
) {
	t.Helper()

	const maxWait = 5 * time.Second

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

	roleID, err := queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		RoleName:     "Not important",
		StoreID:      storeID,
		IsStoreAdmin: false,
	})
	require.NoError(t, err)

	_, err = queries.SeedUser(ctx, gensql.SeedUserParams{
		UserID:       typemapper.UUIDToPgtypeUUID(theUserID),
		Username:     "Not Important",
		UserPassword: "notimportant",
		RoleID:       roleID,
		StoreID:      storeID,
	})
	require.NoError(t, err)

	_, err = queries.SeedLoginCode(ctx, gensql.SeedLoginCodeParams{
		LoginCodeID: typemapper.UUIDToPgtypeUUID(uuid.New()),
		UserID:      typemapper.UUIDToPgtypeUUID(theUserID),
		LoginCode:   theLoginCode,
	})
	require.NoError(t, err)
}

func seedHandleSessionCookie(
	ctx context.Context,
	t *testing.T,
	queries *gensql.Queries,
	theUserDetails readmodel.UserDetails,
) {
	t.Helper()

	const maxWait = 5 * time.Second

	ctx, cancel := context.WithTimeout(ctx, maxWait)
	defer cancel()

	_, err := queries.SeedStore(ctx, gensql.SeedStoreParams{
		StoreID:      typemapper.UUIDToPgtypeUUID(theUserDetails.Store.ID),
		StoreName:    theUserDetails.Store.Name,
		StoreCode:    theUserDetails.Store.Code,
		StoreAddress: "Not Important",
		PhoneNumber:  "081234567890",
	})
	require.NoError(t, err)

	_, err = queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(theUserDetails.Role.ID),
		RoleName:     theUserDetails.Role.Name,
		IsStoreAdmin: theUserDetails.Role.IsStoreAdmin,
		StoreID:      typemapper.UUIDToPgtypeUUID(theUserDetails.Store.ID),
	})
	require.NoError(t, err)

	_, err = queries.SeedUser(ctx, gensql.SeedUserParams{
		UserID:       typemapper.UUIDToPgtypeUUID(theUserDetails.ID),
		Username:     theUserDetails.Username,
		UserPassword: "notimportant",
		RoleID:       typemapper.UUIDToPgtypeUUID(theUserDetails.Role.ID),
		StoreID:      typemapper.UUIDToPgtypeUUID(theUserDetails.Store.ID),
	})
	require.NoError(t, err)
}

type serviceSessionManagerStub struct{}

func (s serviceSessionManagerStub) NewSession(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (s serviceSessionManagerStub) DeleteSession(_ context.Context) error {
	return nil
}

type loginCodePromptManagerStub struct {
	userID uuid.UUID
}

func (l loginCodePromptManagerStub) NewPrompt(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (l loginCodePromptManagerStub) GetUserID(_ context.Context) (uuid.UUID, error) {
	return l.userID, nil
}

func (l loginCodePromptManagerStub) DeletePrompt(_ context.Context) error {
	return nil
}

type securityHandlerSessionManagerStub struct {
	userID uuid.UUID
}

func (s securityHandlerSessionManagerStub) GetUserID(_ context.Context) (uuid.UUID, error) {
	return s.userID, nil
}
