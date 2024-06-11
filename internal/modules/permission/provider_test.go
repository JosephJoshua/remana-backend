//go:build unit
// +build unit

package permission_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/modules/permission"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type providerRepoPermission struct {
	groupName string
	name      string
}

type providerRepoStub struct {
	roleID           uuid.UUID
	isStoreAdmin     bool
	rolePermissions  []providerRepoPermission
	isStoreAdminErr  error
	hasPermissionErr error
}

func (p *providerRepoStub) IsStoreAdmin(_ context.Context, roleID uuid.UUID) (bool, error) {
	if p.isStoreAdminErr != nil {
		return false, p.isStoreAdminErr
	}

	if p.roleID != roleID {
		return false, nil
	}

	return p.isStoreAdmin, nil
}

func (p *providerRepoStub) HasPermission(
	_ context.Context,
	roleID uuid.UUID,
	permissionGroupName string,
	permissionName string,
) (bool, error) {
	if p.hasPermissionErr != nil {
		return false, p.hasPermissionErr
	}

	if p.roleID != roleID {
		return false, nil
	}

	for _, permission := range p.rolePermissions {
		if permission.groupName == permissionGroupName && permission.name == permissionName {
			return true, nil
		}
	}

	return false, nil
}

func TestCan(t *testing.T) {
	t.Parallel()

	t.Run("returns true when role has permission", func(t *testing.T) {
		t.Parallel()

		var (
			theRoleID      = uuid.New()
			thePermissions = []permission.Permission{
				permission.CreateRole(),
				permission.CreateTechnician(),
			}
		)

		repo := &providerRepoStub{
			roleID:          theRoleID,
			rolePermissions: toProviderRepoPermissions(thePermissions),
			isStoreAdmin:    false,
		}

		s := permission.NewProvider(repo)
		ok, err := s.Can(context.Background(), theRoleID, thePermissions[0])

		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("returns false when role doesn't have permission", func(t *testing.T) {
		t.Parallel()

		var (
			theRoleID      = uuid.New()
			thePermissions = []permission.Permission{
				permission.CreateRole(),
				permission.CreateTechnician(),
			}
			someOtherPermission = permission.AssignPermissionsToRole()
		)

		repo := &providerRepoStub{
			roleID:          theRoleID,
			rolePermissions: toProviderRepoPermissions(thePermissions),
			isStoreAdmin:    false,
		}

		s := permission.NewProvider(repo)
		ok, err := s.Can(context.Background(), theRoleID, someOtherPermission)

		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("returns true when role is store admin", func(t *testing.T) {
		t.Parallel()

		var (
			theRoleID        = uuid.New()
			emptyPermissions = []providerRepoPermission{}
			somePermission   = permission.CreateRole()
		)

		repo := &providerRepoStub{
			roleID:          theRoleID,
			rolePermissions: emptyPermissions,
			isStoreAdmin:    true,
		}

		s := permission.NewProvider(repo)
		ok, err := s.Can(context.Background(), theRoleID, somePermission)

		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("returns error when repository.HasPermission() errors", func(t *testing.T) {
		t.Parallel()

		var (
			theRoleID      = uuid.New()
			thePermissions = []permission.Permission{
				permission.CreateRole(),
				permission.CreateTechnician(),
			}
		)

		repo := &providerRepoStub{
			roleID:           theRoleID,
			rolePermissions:  toProviderRepoPermissions(thePermissions),
			hasPermissionErr: errors.New("oh no!"),
		}

		s := permission.NewProvider(repo)
		_, err := s.Can(context.Background(), theRoleID, thePermissions[0])

		require.Error(t, err)
	})
}

func toProviderRepoPermissions(permissions []permission.Permission) []providerRepoPermission {
	var repoPermissions []providerRepoPermission
	for _, p := range permissions {
		repoPermissions = append(repoPermissions, providerRepoPermission{
			groupName: p.GroupName(),
			name:      p.Name(),
		})
	}

	return repoPermissions
}
