package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/shared"
	"github.com/JosephJoshua/remana-backend/internal/shared/apierror"
	"github.com/JosephJoshua/remana-backend/internal/shared/apperror"
	"github.com/JosephJoshua/remana-backend/internal/shared/readmodel"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type SecurityHandlerSessionManager interface {
	GetUserID(ctx context.Context) (uuid.UUID, error)
}

type SecurityHandlerRepository interface {
	GetUserDetailsByID(ctx context.Context, userID uuid.UUID) (readmodel.UserDetails, error)
}

type SecurityHandler struct {
	sessionManager SecurityHandlerSessionManager
	repo           SecurityHandlerRepository
}

func NewSecurityHandler(sessionManager SecurityHandlerSessionManager, repo SecurityHandlerRepository) *SecurityHandler {
	return &SecurityHandler{
		sessionManager: sessionManager,
		repo:           repo,
	}
}

func (s *SecurityHandler) HandleSessionCookie(
	ctx context.Context,
	_ string,
	_ genapi.SessionCookie,
) (context.Context, error) {
	l := zerolog.Ctx(ctx)

	userID, err := s.sessionManager.GetUserID(ctx)
	if err != nil {
		if errors.Is(err, apperror.ErrMissingSession) {
			return ctx, apierror.ToAPIError(http.StatusUnauthorized, "unauthorized")
		}

		l.Error().Err(err).Msg("failed to get user ID from session")

		return ctx, apierror.ToAPIError(
			http.StatusUnauthorized,
			"invalid session. please try logging out and logging back in",
		)
	}

	user, err := s.repo.GetUserDetailsByID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			return ctx, apierror.ToAPIError(
				http.StatusUnauthorized,
				"invalid session. please try logging out and logging back in",
			)
		}

		l.Error().Err(err).Msg("failed to get user details by ID")
		return ctx, apierror.ToAPIError(http.StatusInternalServerError, "failed to get user details by ID")
	}

	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Interface("user", user)
	})

	ctx = shared.NewContextWithUser(ctx, &user)
	return ctx, nil
}
