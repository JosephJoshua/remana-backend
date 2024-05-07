package permission_test

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
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type repositoryStub struct {
	createCalledWith struct {
		id           uuid.UUID
		storeID      uuid.UUID
		name         string
		isStoreAdmin bool
	}
	storeID      uuid.UUID
	existingName string
	createErr    error
	nameTakenErr error
}

func (r *repositoryStub) CreateRole(
	_ context.Context,
	id uuid.UUID,
	storeID uuid.UUID,
	name string,
	isStoreAdmin bool,
) error {
	if r.createErr != nil {
		return r.createErr
	}

	r.createCalledWith.id = id
	r.createCalledWith.storeID = storeID
	r.createCalledWith.name = name
	r.createCalledWith.isStoreAdmin = isStoreAdmin

	return nil
}

func (r *repositoryStub) IsRoleNameTaken(_ context.Context, storeID uuid.UUID, name string) (bool, error) {
	if r.nameTakenErr != nil {
		return false, r.nameTakenErr
	}

	return r.storeID == storeID && r.existingName == name, nil
}

func TestCreateRole(t *testing.T) {
	t.Parallel()

	var (
		theStoreID = uuid.New()
	)

	logger.Init(zerolog.ErrorLevel, appconstant.AppEnvDev)

	requestCtx := appcontext.NewContextWithUser(
		testutil.RequestContextWithLogger(context.Background()),
		testutil.ModifiedUserDetails(func(details *readmodel.UserDetails) {
			details.Store.ID = theStoreID
		}),
	)

	t.Run("tries to create role when request is valid", func(t *testing.T) {
		t.Parallel()

		repo := &repositoryStub{storeID: theStoreID}
		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			repo,
		)

		req := &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: true,
		}

		got, err := s.CreateRole(requestCtx, req)

		require.NoError(t, err)
		require.NotNil(t, got)

		require.NotNil(t, repo.createCalledWith)

		assert.Equal(t, req.Name, repo.createCalledWith.name)
		assert.Equal(t, req.IsStoreAdmin, repo.createCalledWith.isStoreAdmin)
		assert.Equal(t, theStoreID, repo.createCalledWith.storeID)
	})

	t.Run("returns resource location when role is created", func(t *testing.T) {
		t.Parallel()

		var (
			theLocation = url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/roles/ef21dc9e-c364-41cd-8c03-fa289d11e3a7",
			}
		)

		resourceLocationProvider := testutil.NewResourceLocationProviderStubForRole(theLocation)
		repo := &repositoryStub{storeID: theStoreID}

		s := permission.NewService(resourceLocationProvider, repo)

		got, err := s.CreateRole(requestCtx, &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: false,
		})

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, theLocation, got.Location)

		require.True(t, resourceLocationProvider.RoleID.IsSet())
		require.NotNil(t, repo.createCalledWith)

		assert.Equal(t, repo.createCalledWith.id, resourceLocationProvider.RoleID.MustGet())
	})

	t.Run("returns unauthorized when user is missing from context", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			&repositoryStub{},
		)

		emptyCtx := testutil.RequestContextWithLogger(context.Background())
		_, err := s.CreateRole(emptyCtx, &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: true,
		})

		testutil.AssertAPIStatusCode(t, http.StatusUnauthorized, err)
	})

	t.Run("returns bad request when name is empty", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			&repositoryStub{storeID: theStoreID},
		)

		_, err := s.CreateRole(requestCtx, &genapi.CreateRoleRequest{
			Name:         "",
			IsStoreAdmin: true,
		})

		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns conflict when name is taken", func(t *testing.T) {
		t.Parallel()

		const (
			theName = "role 1"
		)

		repo := &repositoryStub{existingName: theName, storeID: theStoreID}
		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			repo,
		)

		req := &genapi.CreateRoleRequest{
			Name:         theName,
			IsStoreAdmin: true,
		}

		_, err := s.CreateRole(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})

	t.Run("returns internal server error when repository.IsRoleNameTaken() errors", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			&repositoryStub{nameTakenErr: errors.New("oh no!"), storeID: theStoreID},
		)

		_, err := s.CreateRole(requestCtx, &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: true,
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})

	t.Run("returns internal server error when repository.CreateRole() errors", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			&repositoryStub{createErr: errors.New("oh no!"), storeID: theStoreID},
		)

		_, err := s.CreateRole(requestCtx, &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: false,
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})
}
