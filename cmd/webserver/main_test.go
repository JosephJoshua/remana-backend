//go:build e2e
// +build e2e

package main_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/appconstant"
	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/infrastructure/core"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestAuthnFlow(t *testing.T) {
	t.Parallel()

	logger.Init(zerolog.ErrorLevel, appconstant.AppEnvDev)

	const (
		theAdminUsername    = "admin"
		theAdminPassword    = "Password123"
		theEmployeeUsername = "employee"
		theEmployeePassword = "Password123"
		theStoreCode        = "store-one"
		theLoginCode        = "A1B2C3D4"
	)

	db := setupTest(t)

	seedAuthnFlow(
		t,
		db,
		theAdminUsername,
		mustHashPassword(t, theAdminPassword),
		theEmployeeUsername,
		mustHashPassword(t, theEmployeePassword),
		theStoreCode,
		theLoginCode,
	)

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
			"username":   theAdminUsername,
			"password":   theAdminPassword,
			"store_code": theStoreCode,
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
			"username":   theEmployeeUsername,
			"password":   theEmployeePassword,
			"store_code": theStoreCode,
		}).
		Expect().
		Status(http.StatusOK)

	loginAsEmployee.JSON().Object().ContainsKey("type").HasValue("type", "employee")
	loginAsEmployee.Cookie("login_code_prompt_id")

	e.POST("/auth/login-code").WithName("supply invalid employee login code").
		WithJSON(map[string]interface{}{
			"login_code": "1234567890",
		}).
		Expect().
		Status(http.StatusBadRequest)

	e.POST("/auth/login-code").WithName("supply incorrect employee login code").
		WithJSON(map[string]interface{}{
			"login_code": "12345678",
		}).
		Expect().
		Status(http.StatusBadRequest)

	e.POST("/auth/login-code").WithName("supply correct employee login code").
		WithJSON(map[string]interface{}{
			"login_code": theLoginCode,
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

func TestCreateRepairOrderFlow(t *testing.T) {
	t.Parallel()

	logger.Init(zerolog.ErrorLevel, appconstant.AppEnvDev)

	const (
		theUsername  = "admin"
		thePassword  = "Password123"
		theStoreCode = "store-one"
	)

	db := setupTest(t)
	seedCreateRepairOrderFlow(t, db, theUsername, mustHashPassword(t, thePassword), theStoreCode)

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

	e.POST("/auth/login").WithName("login as admin").
		WithJSON(map[string]interface{}{
			"username":   theUsername,
			"password":   thePassword,
			"store_code": theStoreCode,
		}).
		Expect().
		Status(http.StatusOK)

	technicianLocation := e.POST("/technicians").WithName("create technician").
		WithJSON(map[string]interface{}{
			"name": "John Doe",
		}).
		Expect().
		Status(http.StatusCreated).
		Header("Location").NotEmpty().Raw()

	parts := strings.Split(technicianLocation, "/")
	technicianID := parts[len(parts)-1]

	e.POST("/technicians").WithName("create technician with same name (conflict)").
		WithJSON(map[string]interface{}{
			"name": "john doe",
		}).
		Expect().
		Status(http.StatusConflict)

	salesPersonLocation := e.POST("/sales-persons").WithName("create sales person").
		WithJSON(map[string]interface{}{
			"name": "John Doe (Sales ver)",
		}).
		Expect().
		Status(http.StatusCreated).
		Header("Location").NotEmpty().Raw()

	parts = strings.Split(salesPersonLocation, "/")
	salesPersonID := parts[len(parts)-1]

	e.POST("/sales-persons").WithName("create sales person with same name (conflict)").
		WithJSON(map[string]interface{}{
			"name": "john doe (Sales ver)",
		}).
		Expect().
		Status(http.StatusConflict)

	damageTypeLocation := e.POST("/damage-types").WithName("create damage type").
		WithJSON(map[string]interface{}{
			"name": "Broken screen",
		}).
		Expect().
		Status(http.StatusCreated).
		Header("Location").NotEmpty().Raw()

	parts = strings.Split(damageTypeLocation, "/")
	damageTypeID := parts[len(parts)-1]

	e.POST("/damage-types").WithName("create damage type with same name (conflict)").
		WithJSON(map[string]interface{}{
			"name": "broken screen",
		}).
		Expect().
		Status(http.StatusConflict)

	phoneConditionLocation := e.POST("/phone-conditions").WithName("create phone condition").
		WithJSON(map[string]interface{}{
			"name": "Missing button",
		}).
		Expect().
		Status(http.StatusCreated).
		Header("Location").NotEmpty().Raw()

	parts = strings.Split(phoneConditionLocation, "/")
	phoneConditionID := parts[len(parts)-1]

	e.POST("/phone-conditions").WithName("create phone condition with same name (conflict)").
		WithJSON(map[string]interface{}{
			"name": "missing button",
		}).
		Expect().
		Status(http.StatusConflict)

	phoneEquipmentLocation := e.POST("/phone-equipments").WithName("create phone equipment").
		WithJSON(map[string]interface{}{
			"name": "Battery",
		}).
		Expect().
		Status(http.StatusCreated).
		Header("Location").NotEmpty().Raw()

	parts = strings.Split(phoneEquipmentLocation, "/")
	phoneEquipmentID := parts[len(parts)-1]

	e.POST("/phone-equipments").WithName("create phone equipment with same name (conflict)").
		WithJSON(map[string]interface{}{
			"name": "battery",
		}).
		Expect().
		Status(http.StatusConflict)

	paymentMethodLocation := e.POST("/payment-methods").WithName("create payment method").
		WithJSON(map[string]interface{}{
			"name": "Cash",
		}).
		Expect().
		Status(http.StatusCreated).
		Header("Location").NotEmpty().Raw()

	parts = strings.Split(paymentMethodLocation, "/")
	paymentMethodID := parts[len(parts)-1]

	e.POST("/payment-methods").WithName("create payment method with same name (conflict)").
		WithJSON(map[string]interface{}{
			"name": "cash",
		}).
		Expect().
		Status(http.StatusConflict)

	e.POST("/repair-orders").WithName("create minimal repair order").
		WithJSON(map[string]interface{}{
			"customer_name":        "John Doe",
			"contact_phone_number": "+6281234567890",
			"phone_type":           "Advan G4",
			"color":                "Black",
			"initial_cost":         100000,
			"sales_person_id":      salesPersonID,
			"technician_id":        technicianID,
			"damage_types":         []string{damageTypeID},
			"photos":               []string{"https://example.com/photo1.jpg", "https://example.com/photo2.jpg"},
		}).
		Expect().
		Status(http.StatusCreated).
		Header("Location").NotEmpty()

	e.POST("/repair-orders").WithName("create full repair order").
		WithJSON(map[string]interface{}{
			"customer_name":         "John Doe",
			"contact_phone_number":  "+6281234567890",
			"phone_type":            "Advan G4",
			"imei":                  "123456789012345",
			"parts_not_checked_yet": "Back cover",
			"passcode": map[string]interface{}{
				"value":             "12345678",
				"is_pattern_locked": true,
			},
			"color":        "Black",
			"initial_cost": 100000,
			"down_payment": map[string]interface{}{
				"amount": 50000,
				"method": paymentMethodID,
			},
			"sales_person_id":  salesPersonID,
			"technician_id":    technicianID,
			"damage_types":     []string{damageTypeID},
			"phone_conditions": []string{phoneConditionID},
			"phone_equipments": []string{phoneEquipmentID},
			"photos":           []string{"https://example.com/photo1.jpg", "https://example.com/photo2.jpg"},
		}).
		Expect().
		Status(http.StatusCreated).
		Header("Location").NotEmpty()

	var someRandomID = uuid.New()

	e.POST("/repair-orders").WithName("create repair order with invalid IDs").
		WithJSON(map[string]interface{}{
			"customer_name":         "John Doe",
			"contact_phone_number":  "+6281234567890",
			"phone_type":            "Advan G4",
			"imei":                  "123456789012345",
			"parts_not_checked_yet": "Back cover",
			"passcode": map[string]interface{}{
				"value":             "12345678",
				"is_pattern_locked": true,
			},
			"color":        "Black",
			"initial_cost": 100000,
			"down_payment": map[string]interface{}{
				"amount": 50000,
				"method": paymentMethodID,
			},
			"sales_person_id":  salesPersonID,
			"technician_id":    technicianID,
			"damage_types":     []string{someRandomID.String()},
			"phone_conditions": []string{someRandomID.String()},
			"phone_equipments": []string{phoneEquipmentID},
			"photos":           []string{"https://example.com/photo1.jpg", "https://example.com/photo2.jpg"},
		}).
		Expect().
		Status(http.StatusBadRequest)

	e.POST("/repair-orders").WithName("create repair order with invalid input").
		WithJSON(map[string]interface{}{
			"customer_name":         "John Doe",
			"contact_phone_number":  "+6281234567890",
			"phone_type":            "Advan G4",
			"imei":                  "123456789012345",
			"parts_not_checked_yet": "Back cover",
			"passcode": map[string]interface{}{
				"value":             "12345678",
				"is_pattern_locked": true,
			},
			"color":        "Black",
			"initial_cost": 100000,
			"down_payment": map[string]interface{}{
				"amount": 150000,
				"method": paymentMethodID,
			},
			"sales_person_id":  salesPersonID,
			"technician_id":    technicianID,
			"damage_types":     []string{damageTypeID},
			"phone_conditions": []string{phoneConditionID},
			"phone_equipments": []string{phoneEquipmentID},
			"photos":           []string{"https://example.com/photo1.jpg", "https://example.com/photo2.jpg"},
		}).
		Expect().
		Status(http.StatusBadRequest)
}

func mustHashPassword(t *testing.T, password string) string {
	ph := core.PasswordHasher{}

	hashed, err := ph.Hash(password)
	require.NoError(t, err)

	return hashed
}

func seedAuthnFlow(
	t *testing.T,
	db *pgxpool.Pool,
	theAdminUsername string,
	theAdminPassword string,
	theEmployeeUsername string,
	theEmployeePassword string,
	theStoreCode string,
	theLoginCode string,
) {
	t.Helper()

	const maxWait = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), maxWait)
	defer cancel()

	queries := gensql.New(db)

	storeID, err := queries.SeedStore(ctx, gensql.SeedStoreParams{
		StoreID:      typemapper.UUIDToPgtypeUUID(uuid.New()),
		StoreName:    "Not important",
		StoreCode:    theStoreCode,
		StoreAddress: "Not important",
		PhoneNumber:  "081234567890",
	})
	require.NoError(t, err)

	adminRoleID, err := queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		RoleName:     "Not important",
		StoreID:      storeID,
		IsStoreAdmin: true,
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

	employeeRoleID, err := queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		RoleName:     "Not important",
		StoreID:      storeID,
		IsStoreAdmin: false,
	})
	require.NoError(t, err)

	employeeUserID, err := queries.SeedUser(ctx, gensql.SeedUserParams{
		UserID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		Username:     theEmployeeUsername,
		UserPassword: theEmployeePassword,
		RoleID:       employeeRoleID,
		StoreID:      storeID,
	})
	require.NoError(t, err)

	_, err = queries.SeedLoginCode(ctx, gensql.SeedLoginCodeParams{
		LoginCodeID: typemapper.UUIDToPgtypeUUID(uuid.New()),
		UserID:      employeeUserID,
		LoginCode:   theLoginCode,
	})
	require.NoError(t, err)
}

func seedCreateRepairOrderFlow(
	t *testing.T,
	db *pgxpool.Pool,
	theUsername string,
	thePassword string,
	theStoreCode string,
) {
	t.Helper()

	const maxWait = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), maxWait)
	defer cancel()

	queries := gensql.New(db)

	storeID, err := queries.SeedStore(ctx, gensql.SeedStoreParams{
		StoreID:      typemapper.UUIDToPgtypeUUID(uuid.New()),
		StoreName:    "Not important",
		StoreCode:    theStoreCode,
		StoreAddress: "Not important",
		PhoneNumber:  "081234567890",
	})
	require.NoError(t, err)

	roleID, err := queries.SeedRole(ctx, gensql.SeedRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		RoleName:     "Not important",
		StoreID:      storeID,
		IsStoreAdmin: true,
	})
	require.NoError(t, err)

	_, err = queries.SeedUser(ctx, gensql.SeedUserParams{
		UserID:       typemapper.UUIDToPgtypeUUID(uuid.New()),
		Username:     theUsername,
		UserPassword: thePassword,
		RoleID:       roleID,
		StoreID:      storeID,
	})
	require.NoError(t, err)
}
