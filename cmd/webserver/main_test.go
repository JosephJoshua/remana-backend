//go:build e2e
// +build e2e

package main_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/infrastructure/core"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/shared"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestAuthnFlow(t *testing.T) {
	t.Parallel()

	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)

	db, cleanup := setupTest(t)
	t.Cleanup(func() {
		l := logger.MustGet()

		err := cleanup()
		l.Error().Err(err).Msg("error cleaning up test resources")
	})

	const seedMaxWait = 5 * time.Second

	seedCtx, cancelSeed := context.WithTimeout(context.Background(), seedMaxWait)
	t.Cleanup(cancelSeed)

	seedUsers(seedCtx, t, db)

	addr := runServer(context.Background(), t, db)
	client := createHTTPClient(t)

	waitForReady(context.Background(), t, client, addr, 5*time.Second)

	e := httpexpect.WithConfig(httpexpect.Config{
		Reporter: httpexpect.NewFatalReporter(t),
		Client:   &client,
		BaseURL: (&url.URL{
			Scheme: "https",
			Host:   addr,
		}).String(),
	})

	e.POST("/auth/login").WithName("login with incorrect credentials").
		WithJSON(map[string]interface{}{
			"username":   "wrongusername",
			"password":   "wrongpassword",
			"store_code": "wrongstorecode",
		}).
		Expect().
		Status(http.StatusUnauthorized)

	loginAsAdmin := e.POST("/auth/login").WithName("login as admin").
		WithJSON(map[string]interface{}{
			"username":   "admin",
			"password":   "Password123",
			"store_code": "store-one",
		}).
		Expect().
		Status(http.StatusOK)

	loginAsAdmin.JSON().Object().ContainsKey("type").HasValue("type", "admin")
	loginAsAdmin.Cookies().NotEmpty()

	e.GET("/users/me").WithName("get user details of admin").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("id").NotEmpty()

	e.POST("/auth/logout").WithName("logout").
		Expect().
		Status(http.StatusResetContent).
		NoContent().
		Cookies().NotEmpty()

	e.GET("/users/me").WithName("verify user is logged out").
		Expect().
		Status(http.StatusUnauthorized)

	loginAsEmployee := e.POST("/auth/login").WithName("log in as employee").
		WithJSON(map[string]interface{}{
			"username":   "employee",
			"password":   "Password123",
			"store_code": "store-one",
		}).
		Expect().
		Status(http.StatusOK)

	loginAsEmployee.JSON().Object().ContainsKey("type").HasValue("type", "employee")
	loginAsEmployee.Cookie("login_code_prompt_id")

	e.POST("/auth/login-code").WithName("supply incorrect employee login code").
		WithJSON(map[string]interface{}{
			"login_code": "12345678",
		}).
		Expect().
		Status(http.StatusBadRequest)

	e.POST("/auth/login-code").WithName("supply correct employee login code").
		WithJSON(map[string]interface{}{
			"login_code": "A1B2C3D4",
		}).
		Expect().
		Status(http.StatusNoContent).
		NoContent().
		Cookies().NotEmpty()

	e.GET("/users/me").WithName("get user details of employee").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("id").NotEmpty()

	e.POST("/auth/logout").WithName("logout").
		Expect().
		Status(http.StatusResetContent).
		NoContent().
		Cookies().NotEmpty()

	e.GET("/users/me").WithName("verify user is logged out").
		Expect().
		Status(http.StatusUnauthorized)
}

func seedUsers(ctx context.Context, t testing.TB, db *pgxpool.Pool) {
	t.Helper()

	queries := gensql.New(db)

	storeID, err := queries.SeedStore(ctx, gensql.SeedStoreParams{
		StoreID:      typemapper.UUIDToPgtypeUUID(uuid.New()),
		StoreName:    "Store 1",
		StoreCode:    "store-one",
		StoreAddress: "123 Main St",
		PhoneNumber:  "081234567890",
	})
	require.NoError(t, err)

	adminRoleID, err := queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		RoleName:     "Admin",
		StoreID:      storeID,
		IsStoreAdmin: true,
	})
	require.NoError(t, err)

	passwordHasher := &core.PasswordHasher{}

	hashedPassword, err := passwordHasher.Hash("Password123")
	require.NoError(t, err)

	_, err = queries.SeedUser(ctx, gensql.SeedUserParams{
		UserID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		Username:     "admin",
		UserPassword: hashedPassword,
		RoleID:       adminRoleID,
		StoreID:      storeID,
	})
	require.NoError(t, err)

	employeeRoleID, err := queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		RoleName:     "Employee",
		StoreID:      storeID,
		IsStoreAdmin: false,
	})
	require.NoError(t, err)

	employeeUserID, err := queries.SeedUser(ctx, gensql.SeedUserParams{
		UserID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		Username:     "employee",
		UserPassword: hashedPassword,
		RoleID:       employeeRoleID,
		StoreID:      storeID,
	})
	require.NoError(t, err)

	_, err = queries.SeedLoginCode(ctx, gensql.SeedLoginCodeParams{
		LoginCodeID: typemapper.UUIDToPgtypeUUID(uuid.New()),
		UserID:      employeeUserID,
		LoginCode:   "A1B2C3D4",
	})
	require.NoError(t, err)
}
