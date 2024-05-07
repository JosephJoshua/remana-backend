package permission

import (
	"context"
	"net/http"
	"net/url"

	"github.com/JosephJoshua/remana-backend/internal/apierror"
	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Repository interface {
	CreateRole(ctx context.Context, id uuid.UUID, storeID uuid.UUID, name string, isStoreAdmin bool) error
	IsRoleNameTaken(ctx context.Context, storeID uuid.UUID, name string) (bool, error)
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
