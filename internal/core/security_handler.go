package core

import (
	"context"
	"errors"
	"net/http"

	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/shared"
	"github.com/JosephJoshua/remana-backend/internal/shared/apierror"
	"github.com/JosephJoshua/remana-backend/internal/shared/apperror"
	"github.com/JosephJoshua/remana-backend/internal/shared/readmodel"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type securityHandlerRepository interface {
	GetUserDetailsByID(ctx context.Context, userID uuid.UUID) (readmodel.UserDetails, error)
}

type securityHandler struct {
	sm   *scs.SessionManager
	repo securityHandlerRepository
}

func newSecurityHandler(sm *scs.SessionManager, repo securityHandlerRepository) *securityHandler {
	return &securityHandler{
		sm:   sm,
		repo: repo,
	}
}

func (s *securityHandler) HandleSessionCookie(
	ctx context.Context,
	_ string,
	_ genapi.SessionCookie,
) (context.Context, error) {
	l := zerolog.Ctx(ctx)

	userIDStr := s.sm.GetString(ctx, userIDKey)
	if userIDStr == "" {
		return ctx, apierror.ToAPIError(http.StatusUnauthorized, "unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		l.Error().Err(err).Msg("invalid user ID found in session")

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
