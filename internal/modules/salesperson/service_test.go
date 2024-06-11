//go:build unit
// +build unit

package salesperson_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/appconstant"
	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/modules/permission"
	"github.com/JosephJoshua/remana-backend/internal/modules/salesperson"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSalesPerson(t *testing.T) {
	t.Parallel()

	var (
		theStoreID                   = uuid.New()
		theRoleID                    = uuid.New()
		qualifyingPermissionProvider = testutil.NewPermissionProviderStub(
			theRoleID,
			[]permission.Permission{
				permission.CreateSalesPerson(),
			},
			nil,
		)
	)

	logger.Init(zerolog.ErrorLevel, appconstant.AppEnvDev)
	requestCtx := appcontext.NewContextWithUser(
		testutil.RequestContextWithLogger(context.Background()),
		testutil.ModifiedUserDetails(func(details *readmodel.UserDetails) {
			details.Role.ID = theRoleID
			details.Store.ID = theStoreID
		}),
	)

	t.Run("tries to create sales person when request is valid", func(t *testing.T) {
		t.Parallel()

		repo := &repositoryStub{storeID: theStoreID}
		s := salesperson.NewService(
			testutil.NewResourceLocationProviderStubForSalesPerson(url.URL{}),
			qualifyingPermissionProvider,
			repo,
		)

		req := &genapi.CreateSalesPersonRequest{
			Name: "sales person 1",
		}

		got, err := s.CreateSalesPerson(requestCtx, req)

		require.NoError(t, err)
		require.NotNil(t, got)

		require.NotNil(t, repo.createCalledWith)

		assert.Equal(t, req.Name, repo.createCalledWith.name)
		assert.Equal(t, theStoreID, repo.createCalledWith.storeID)
	})

	t.Run("returns resource location when sales person is created", func(t *testing.T) {
		t.Parallel()

		var (
			theLocation = url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/sales-persons/ef21dc9e-c364-41cd-8c03-fa289d11e3a7",
			}
		)

		resourceLocationProvider := testutil.NewResourceLocationProviderStubForSalesPerson(theLocation)
		repo := &repositoryStub{storeID: theStoreID}

		s := salesperson.NewService(resourceLocationProvider, qualifyingPermissionProvider, repo)

		got, err := s.CreateSalesPerson(requestCtx, &genapi.CreateSalesPersonRequest{
			Name: "sales person 1",
		})

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, theLocation, got.Location)

		require.True(t, resourceLocationProvider.SalesPersonID.IsSet())
		require.NotNil(t, repo.createCalledWith)

		assert.Equal(t, repo.createCalledWith.id, resourceLocationProvider.SalesPersonID.MustGet())
	})

	t.Run("returns unauthorized when user is missing from context", func(t *testing.T) {
		t.Parallel()

		s := salesperson.NewService(
			testutil.NewResourceLocationProviderStubForSalesPerson(url.URL{}),
			qualifyingPermissionProvider,
			&repositoryStub{},
		)

		emptyCtx := testutil.RequestContextWithLogger(context.Background())
		_, err := s.CreateSalesPerson(emptyCtx, &genapi.CreateSalesPersonRequest{
			Name: "sales person 1",
		})

		testutil.AssertAPIStatusCode(t, http.StatusUnauthorized, err)
	})

	t.Run("returns forbidden when role doesn't have permission", func(t *testing.T) {
		t.Parallel()

		s := salesperson.NewService(
			testutil.NewResourceLocationProviderStubForSalesPerson(url.URL{}),
			testutil.NewPermissionProviderStub(theRoleID, []permission.Permission{}, nil),
			&repositoryStub{},
		)

		_, err := s.CreateSalesPerson(requestCtx, &genapi.CreateSalesPersonRequest{
			Name: "sales person 1",
		})

		testutil.AssertAPIStatusCode(t, http.StatusForbidden, err)
	})

	t.Run("returns bad request when name is empty", func(t *testing.T) {
		t.Parallel()

		s := salesperson.NewService(
			testutil.NewResourceLocationProviderStubForSalesPerson(url.URL{}),
			qualifyingPermissionProvider,
			&repositoryStub{storeID: theStoreID},
		)

		_, err := s.CreateSalesPerson(requestCtx, &genapi.CreateSalesPersonRequest{
			Name: "",
		})

		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns conflict when name is taken", func(t *testing.T) {
		t.Parallel()

		const (
			theName = "sales person 1"
		)

		repo := &repositoryStub{existingName: theName, storeID: theStoreID}
		s := salesperson.NewService(
			testutil.NewResourceLocationProviderStubForSalesPerson(url.URL{}),
			qualifyingPermissionProvider,
			repo,
		)

		req := &genapi.CreateSalesPersonRequest{
			Name: theName,
		}

		_, err := s.CreateSalesPerson(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})

	t.Run("returns internal server error when repository.IsNameTaken() errors", func(t *testing.T) {
		t.Parallel()

		s := salesperson.NewService(
			testutil.NewResourceLocationProviderStubForSalesPerson(url.URL{}),
			qualifyingPermissionProvider,
			&repositoryStub{nameTakenErr: errors.New("oh no!"), storeID: theStoreID},
		)

		_, err := s.CreateSalesPerson(requestCtx, &genapi.CreateSalesPersonRequest{
			Name: "sales person 1",
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})

	t.Run("returns internal server error when repository.CreateSalesPerson() errors", func(t *testing.T) {
		t.Parallel()

		s := salesperson.NewService(
			testutil.NewResourceLocationProviderStubForSalesPerson(url.URL{}),
			qualifyingPermissionProvider,
			&repositoryStub{createErr: errors.New("oh no!"), storeID: theStoreID},
		)

		_, err := s.CreateSalesPerson(requestCtx, &genapi.CreateSalesPersonRequest{
			Name: "sales person 1",
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})

	t.Run("returns internal server error when permissionProvider.Can() errors", func(t *testing.T) {
		t.Parallel()

		s := salesperson.NewService(
			testutil.NewResourceLocationProviderStubForSalesPerson(url.URL{}),
			testutil.NewPermissionProviderStub(theRoleID, []permission.Permission{}, errors.New("oh no!")),
			&repositoryStub{createErr: errors.New("oh no!"), storeID: theStoreID},
		)

		_, err := s.CreateSalesPerson(requestCtx, &genapi.CreateSalesPersonRequest{
			Name: "sales person 1",
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})
}

type repositoryStub struct {
	createCalledWith struct {
		id      uuid.UUID
		storeID uuid.UUID
		name    string
	}
	storeID      uuid.UUID
	existingName string
	createErr    error
	nameTakenErr error
}

func (r *repositoryStub) CreateSalesPerson(_ context.Context, id uuid.UUID, storeID uuid.UUID, name string) error {
	if r.createErr != nil {
		return r.createErr
	}

	r.createCalledWith.id = id
	r.createCalledWith.storeID = storeID
	r.createCalledWith.name = name

	return nil
}

func (r *repositoryStub) IsNameTaken(_ context.Context, storeID uuid.UUID, name string) (bool, error) {
	if r.nameTakenErr != nil {
		return false, r.nameTakenErr
	}

	return r.storeID == storeID && r.existingName == name, nil
}
