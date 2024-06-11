package damagetype

import (
	"context"
	"net/http"
	"net/url"

	"github.com/JosephJoshua/remana-backend/internal/apierror"
	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/modules/permission"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Repository interface {
	CreateDamageType(ctx context.Context, id uuid.UUID, storeID uuid.UUID, name string) error
	IsNameTaken(ctx context.Context, storeID uuid.UUID, name string) (bool, error)
}

type ResourceLocationProvider interface {
	DamageType(damageTypeID uuid.UUID) url.URL
}

type Service struct {
	resourceLocationProvider ResourceLocationProvider
	permissionProvider       permission.Provider
	repo                     Repository
}

func NewService(
	resourceLocationProvider ResourceLocationProvider,
	permissionProvider permission.Provider,
	repo Repository,
) *Service {
	return &Service{
		resourceLocationProvider: resourceLocationProvider,
		permissionProvider:       permissionProvider,
		repo:                     repo,
	}
}

func (s *Service) CreateDamageType(
	ctx context.Context,
	req *genapi.CreateDamageTypeRequest,
) (*genapi.CreateDamageTypeCreated, error) {
	l := zerolog.Ctx(ctx)

	user, ok := appcontext.GetUserFromContext(ctx)
	if !ok {
		l.Error().Msg("user is missing from context")
		return nil, apierror.ToAPIError(http.StatusUnauthorized, "unauthorized")
	}

	if can, err := s.permissionProvider.Can(ctx, user.Role.ID, permission.CreateDamageType()); err != nil {
		l.Error().Err(err).Msg("failed to check permission")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to check permission")
	} else if !can {
		return nil, apierror.ToAPIError(http.StatusForbidden, "insufficient permissions")
	}

	if req.Name == "" {
		return nil, apierror.ToAPIError(http.StatusBadRequest, "name is required and cannot be empty")
	}

	if taken, err := s.repo.IsNameTaken(ctx, user.Store.ID, req.Name); taken {
		return nil, apierror.ToAPIError(http.StatusConflict, "name is taken")
	} else if err != nil {
		l.Error().Err(err).Msg("failed to check if name is taken")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to check if name is taken")
	}

	id := uuid.New()

	if err := s.repo.CreateDamageType(ctx, id, user.Store.ID, req.Name); err != nil {
		l.Error().Err(err).Msg("failed to create damage type")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to create damage type")
	}

	location := s.resourceLocationProvider.DamageType(id)
	return &genapi.CreateDamageTypeCreated{
		Location: location,
	}, nil
}
