package permission

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/JosephJoshua/remana-backend/internal/apierror"
	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type GetPermissionIDDetail struct {
	GroupName string
	Name      string
}

type Repository interface {
	CreateRole(ctx context.Context, id uuid.UUID, storeID uuid.UUID, name string, isStoreAdmin bool) error
	IsRoleNameTaken(ctx context.Context, storeID uuid.UUID, name string) (bool, error)
	GetPermissionIDs(ctx context.Context, permissions []GetPermissionIDDetail) ([]uuid.UUID, error)
	DoesRoleExist(ctx context.Context, roleID uuid.UUID) (bool, error)
	AssignPermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
}

type ResourceLocationProvider interface {
	Role(roleID uuid.UUID) url.URL
}

type Service struct {
	resourceLocationProvider ResourceLocationProvider
	repo                     Repository
}

func NewService(resourceLocationProvider ResourceLocationProvider, repo Repository) *Service {
	return &Service{
		resourceLocationProvider: resourceLocationProvider,
		repo:                     repo,
	}
}

func (s *Service) CreateRole(
	ctx context.Context,
	req *genapi.CreateRoleRequest,
) (*genapi.CreateRoleCreated, error) {
	l := zerolog.Ctx(ctx)

	user, ok := appcontext.GetUserFromContext(ctx)
	if !ok {
		l.Error().Msg("user is missing from context")
		return nil, apierror.ToAPIError(http.StatusUnauthorized, "unauthorized")
	}

	if req.Name == "" {
		return nil, apierror.ToAPIError(http.StatusBadRequest, "name is required and cannot be empty")
	}

	if taken, err := s.repo.IsRoleNameTaken(ctx, user.Store.ID, req.Name); taken {
		return nil, apierror.ToAPIError(http.StatusConflict, "name is taken")
	} else if err != nil {
		l.Error().Err(err).Msg("failed to check if role name is taken")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to check if role name is taken")
	}

	id := uuid.New()

	if err := s.repo.CreateRole(ctx, id, user.Store.ID, req.Name, req.IsStoreAdmin); err != nil {
		l.Error().Err(err).Msg("failed to create role")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to create role")
	}

	location := s.resourceLocationProvider.Role(id)
	return &genapi.CreateRoleCreated{
		Location: location,
	}, nil
}

func (s *Service) AssignPermissionsToRole(
	ctx context.Context,
	req *genapi.AssignPermissionsToRoleRequest,
	params genapi.AssignPermissionsToRoleParams,
) error {
	l := zerolog.Ctx(ctx)

	_, ok := appcontext.GetUserFromContext(ctx)
	if !ok {
		l.Error().Msg("user is missing from context")
		return apierror.ToAPIError(http.StatusUnauthorized, "unauthorized")
	}

	if exists, err := s.repo.DoesRoleExist(ctx, params.RoleId); err != nil {
		l.Error().Err(err).Msg("failed to check if role exists")
		return apierror.ToAPIError(http.StatusInternalServerError, "failed to check if role exists")
	} else if !exists {
		return apierror.ToAPIError(http.StatusBadRequest, "role does not exist")
	}

	permissions := make([]GetPermissionIDDetail, 0, len(req.Permissions))
	for _, p := range req.Permissions {
		permissions = append(permissions, GetPermissionIDDetail{
			GroupName: p.GroupName,
			Name:      p.Name,
		})
	}

	ids, err := s.repo.GetPermissionIDs(ctx, permissions)
	if errors.Is(err, apperror.ErrPermissionNotFound) {
		return apierror.ToAPIError(http.StatusBadRequest, "permission does not exist")
	} else if err != nil {
		l.Error().Err(err).Msg("failed to check if permissions exist")
		return apierror.ToAPIError(http.StatusInternalServerError, "failed to check if permissions exist")
	}

	if err = s.repo.AssignPermissionsToRole(ctx, params.RoleId, ids); err != nil {
		l.Error().Err(err).Msg("failed to assign permissions to role")
		return apierror.ToAPIError(http.StatusInternalServerError, "failed to assign permissions to role")
	}

	return nil
}
