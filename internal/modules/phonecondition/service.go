package phonecondition

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
	CreatePhoneCondition(ctx context.Context, id uuid.UUID, storeID uuid.UUID, name string) error
	IsNameTaken(ctx context.Context, storeID uuid.UUID, name string) (bool, error)
}

type ResourceLocationProvider interface {
	PhoneCondition(phoneConditionID uuid.UUID) url.URL
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

func (s *Service) CreatePhoneCondition(
	ctx context.Context,
	req *genapi.CreatePhoneConditionRequest,
) (*genapi.CreatePhoneConditionCreated, error) {
	l := zerolog.Ctx(ctx)

	user, ok := appcontext.GetUserFromContext(ctx)
	if !ok {
		l.Error().Msg("user is missing from context")
		return nil, apierror.ToAPIError(http.StatusUnauthorized, "unauthorized")
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

	if err := s.repo.CreatePhoneCondition(ctx, id, user.Store.ID, req.Name); err != nil {
		l.Error().Err(err).Msg("failed to create phone condition")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to create phone condition")
	}

	location := s.resourceLocationProvider.PhoneCondition(id)
	return &genapi.CreatePhoneConditionCreated{
		Location: location,
	}, nil
}
