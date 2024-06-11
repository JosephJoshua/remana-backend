//go:build unit
// +build unit

package permission_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/appconstant"
	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/apperror"
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

func TestCreateRole(t *testing.T) {
	t.Parallel()

	var (
		theStoreID = uuid.New()
		theRoleID  = uuid.New()
	)

	logger.Init(zerolog.ErrorLevel, appconstant.AppEnvDev)

	requestCtx := appcontext.NewContextWithUser(
		testutil.RequestContextWithLogger(context.Background()),
		testutil.ModifiedUserDetails(func(details *readmodel.UserDetails) {
			details.Store.ID = theStoreID
			details.Role.ID = theRoleID
		}),
	)

	qualifyingPermissionProvider := testutil.NewPermissionProviderStub(theRoleID, []permission.Permission{
		permission.CreateRole(),
	}, nil)

	t.Run("tries to create role when request is valid", func(t *testing.T) {
		t.Parallel()

		repo := &serviceRepoStub{storeID: theStoreID}
		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			repo,
			qualifyingPermissionProvider,
		)

		req := &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: true,
		}

		got, err := s.CreateRole(requestCtx, req)

		require.NoError(t, err)
		require.NotNil(t, got)

		require.NotNil(t, repo.createRoleCalledWith)

		assert.Equal(t, req.Name, repo.createRoleCalledWith.name)
		assert.Equal(t, req.IsStoreAdmin, repo.createRoleCalledWith.isStoreAdmin)
		assert.Equal(t, theStoreID, repo.createRoleCalledWith.storeID)
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
		repo := &serviceRepoStub{storeID: theStoreID}

		s := permission.NewService(resourceLocationProvider, repo, qualifyingPermissionProvider)

		got, err := s.CreateRole(requestCtx, &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: false,
		})

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, theLocation, got.Location)

		require.True(t, resourceLocationProvider.RoleID.IsSet())
		require.NotNil(t, repo.createRoleCalledWith)

		assert.Equal(t, repo.createRoleCalledWith.id, resourceLocationProvider.RoleID.MustGet())
	})

	t.Run("returns unauthorized when user is missing from context", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			&serviceRepoStub{},
			qualifyingPermissionProvider,
		)

		emptyCtx := testutil.RequestContextWithLogger(context.Background())
		_, err := s.CreateRole(emptyCtx, &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: true,
		})

		testutil.AssertAPIStatusCode(t, http.StatusUnauthorized, err)
	})

	t.Run("returns forbidden when user doesn't have permissions", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			&serviceRepoStub{},
			testutil.NewPermissionProviderStub(theRoleID, []permission.Permission{}, nil),
		)

		_, err := s.CreateRole(requestCtx, &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: true,
		})

		testutil.AssertAPIStatusCode(t, http.StatusForbidden, err)
	})

	t.Run("returns bad request when name is empty", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			&serviceRepoStub{storeID: theStoreID},
			qualifyingPermissionProvider,
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

		repo := &serviceRepoStub{existingName: theName, storeID: theStoreID}
		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			repo,
			qualifyingPermissionProvider,
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
			&serviceRepoStub{nameTakenErr: errors.New("oh no!"), storeID: theStoreID},
			qualifyingPermissionProvider,
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
			&serviceRepoStub{createRoleErr: errors.New("oh no!"), storeID: theStoreID},
			qualifyingPermissionProvider,
		)

		_, err := s.CreateRole(requestCtx, &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: false,
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})

	t.Run("returns internal server error when permissionProvider.Can() errors", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			&serviceRepoStub{createRoleErr: errors.New("oh no!"), storeID: theStoreID},
			testutil.NewPermissionProviderStub(theRoleID, nil, errors.New("oh no!")),
		)

		_, err := s.CreateRole(requestCtx, &genapi.CreateRoleRequest{
			Name:         "role 1",
			IsStoreAdmin: false,
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})
}

func TestAssignPermissionsToRole(t *testing.T) {
	t.Parallel()

	var (
		theStoreID     = uuid.New()
		theRoleID      = uuid.New()
		thePermissions = []genapi.AssignPermissionsToRoleRequestPermissionsItem{
			{GroupName: "permission", Name: "test"},
			{GroupName: "permission_two", Name: "test_two"},
		}
		theRepoPermissions = []serviceRepoPermission{
			{id: uuid.New(), groupName: "permission", name: "test"},
			{id: uuid.New(), groupName: "permission_two", name: "test_two"},
		}
	)

	logger.Init(zerolog.ErrorLevel, appconstant.AppEnvDev)

	requestCtx := appcontext.NewContextWithUser(
		testutil.RequestContextWithLogger(context.Background()),
		testutil.ModifiedUserDetails(func(details *readmodel.UserDetails) {
			details.Store.ID = theStoreID
			details.Role.ID = theRoleID
		}),
	)

	qualifyingPermissionProvider := testutil.NewPermissionProviderStub(theRoleID, []permission.Permission{
		permission.AssignPermissionsToRole(),
	}, nil)

	baseRepo := func() *serviceRepoStub {
		return &serviceRepoStub{storeID: theStoreID, roleID: theRoleID, permissions: theRepoPermissions}
	}

	t.Run("tries to assign permissions when request is valid", func(t *testing.T) {
		t.Parallel()

		repo := baseRepo()
		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			repo,
			qualifyingPermissionProvider,
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissions,
		}

		params := genapi.AssignPermissionsToRoleParams{RoleId: theRoleID}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		require.NoError(t, err)

		assert.Equal(t, theRoleID, repo.assignPermissionsCalledWith.roleID)

		permissionIDs := make([]uuid.UUID, 0, len(theRepoPermissions))
		for _, p := range theRepoPermissions {
			permissionIDs = append(permissionIDs, p.id)
		}

		assert.ElementsMatch(t, permissionIDs, repo.assignPermissionsCalledWith.permissionIDs)
	})

	t.Run("returns no error when request is valid", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			baseRepo(),
			qualifyingPermissionProvider,
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissions,
		}

		params := genapi.AssignPermissionsToRoleParams{RoleId: theRoleID}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		require.NoError(t, err)
	})

	t.Run("returns unauthorized when user is missing from context", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			baseRepo(),
			qualifyingPermissionProvider,
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissions,
		}

		params := genapi.AssignPermissionsToRoleParams{RoleId: theRoleID}
		emptyCtx := testutil.RequestContextWithLogger(context.Background())

		err := s.AssignPermissionsToRole(emptyCtx, req, params)
		testutil.AssertAPIStatusCode(t, http.StatusUnauthorized, err)
	})

	t.Run("returns forbidden when user doesn't have permissions", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			baseRepo(),
			testutil.NewPermissionProviderStub(theRoleID, []permission.Permission{}, nil),
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissions,
		}

		params := genapi.AssignPermissionsToRoleParams{RoleId: theRoleID}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		testutil.AssertAPIStatusCode(t, http.StatusForbidden, err)
	})

	t.Run("returns bad request when a permission doesn't exist", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			baseRepo(),
			qualifyingPermissionProvider,
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: []genapi.AssignPermissionsToRoleRequestPermissionsItem{
				thePermissions[0],
				{GroupName: "some_random_group", Name: "some_random_permission"},
			},
		}

		params := genapi.AssignPermissionsToRoleParams{RoleId: theRoleID}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns bad request when the role doesn't exist", func(t *testing.T) {
		t.Parallel()

		var (
			someRandomID = uuid.New()
		)

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			baseRepo(),
			qualifyingPermissionProvider,
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissions,
		}

		params := genapi.AssignPermissionsToRoleParams{RoleId: someRandomID}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns internal server error when repository.GetPermissionIDs() errors", func(t *testing.T) {
		t.Parallel()

		repo := baseRepo()
		repo.getPermissionIDsErr = errors.New("oh no!")

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			repo,
			qualifyingPermissionProvider,
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissions,
		}

		params := genapi.AssignPermissionsToRoleParams{RoleId: theRoleID}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})

	t.Run("returns internal server error when repository.DoesRoleExist() errors", func(t *testing.T) {
		t.Parallel()

		repo := baseRepo()
		repo.roleExistsErr = errors.New("oh no!")

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			repo,
			qualifyingPermissionProvider,
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissions,
		}

		params := genapi.AssignPermissionsToRoleParams{RoleId: theRoleID}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})

	t.Run("returns internal server error when repository.AssignPermissionsToRole() errors", func(t *testing.T) {
		t.Parallel()

		repo := baseRepo()
		repo.assignPermissionsErr = errors.New("oh no!")

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			repo,
			qualifyingPermissionProvider,
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissions,
		}

		params := genapi.AssignPermissionsToRoleParams{RoleId: theRoleID}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})

	t.Run("returns internal server error when permissionProvider.Can() errors", func(t *testing.T) {
		t.Parallel()

		s := permission.NewService(
			testutil.NewResourceLocationProviderStubForRole(url.URL{}),
			baseRepo(),
			testutil.NewPermissionProviderStub(theRoleID, nil, errors.New("oh no!")),
		)

		req := &genapi.AssignPermissionsToRoleRequest{
			Permissions: thePermissions,
		}

		params := genapi.AssignPermissionsToRoleParams{RoleId: theRoleID}

		err := s.AssignPermissionsToRole(requestCtx, req, params)
		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})
}

type serviceRepoPermission struct {
	id        uuid.UUID
	groupName string
	name      string
}

type serviceRepoStub struct {
	createRoleCalledWith struct {
		id           uuid.UUID
		storeID      uuid.UUID
		name         string
		isStoreAdmin bool
	}
	assignPermissionsCalledWith struct {
		roleID        uuid.UUID
		permissionIDs []uuid.UUID
	}
	storeID              uuid.UUID
	roleID               uuid.UUID
	permissions          []serviceRepoPermission
	existingName         string
	createRoleErr        error
	assignPermissionsErr error
	nameTakenErr         error
	roleExistsErr        error
	getPermissionIDsErr  error
}

func (r *serviceRepoStub) CreateRole(
	_ context.Context,
	id uuid.UUID,
	storeID uuid.UUID,
	name string,
	isStoreAdmin bool,
) error {
	if r.createRoleErr != nil {
		return r.createRoleErr
	}

	r.createRoleCalledWith.id = id
	r.createRoleCalledWith.storeID = storeID
	r.createRoleCalledWith.name = name
	r.createRoleCalledWith.isStoreAdmin = isStoreAdmin

	return nil
}

func (r *serviceRepoStub) IsRoleNameTaken(_ context.Context, storeID uuid.UUID, name string) (bool, error) {
	if r.nameTakenErr != nil {
		return false, r.nameTakenErr
	}

	return r.storeID == storeID && r.existingName == name, nil
}

func (r *serviceRepoStub) GetPermissionIDs(
	_ context.Context,
	permissions []permission.GetPermissionIDDetail,
) ([]uuid.UUID, error) {
	if r.getPermissionIDsErr != nil {
		return []uuid.UUID{}, r.getPermissionIDsErr
	}

	ids := make([]uuid.UUID, 0, len(permissions))
	for _, toFind := range permissions {
		var found bool

		for _, p := range r.permissions {
			if p.groupName == toFind.GroupName && p.name == toFind.Name {
				ids = append(ids, p.id)
				found = true

				break
			}
		}

		if !found {
			return []uuid.UUID{}, apperror.ErrPermissionNotFound
		}
	}

	return ids, nil
}

func (r *serviceRepoStub) DoesRoleExist(_ context.Context, roleID uuid.UUID) (bool, error) {
	if r.roleExistsErr != nil {
		return false, r.roleExistsErr
	}

	return roleID == r.roleID, nil
}

func (r *serviceRepoStub) AssignPermissionsToRole(_ context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	if r.assignPermissionsErr != nil {
		return r.assignPermissionsErr
	}

	r.assignPermissionsCalledWith.roleID = roleID
	r.assignPermissionsCalledWith.permissionIDs = permissionIDs

	return nil
}
