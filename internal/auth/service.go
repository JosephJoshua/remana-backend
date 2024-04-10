package auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/JosephJoshua/repair-management-backend/internal/genapi"
	"github.com/JosephJoshua/repair-management-backend/internal/shared/apierror"
	"github.com/JosephJoshua/repair-management-backend/internal/shared/apperror"
	"github.com/JosephJoshua/repair-management-backend/internal/shared/domain"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type SessionManager interface {
	NewSession(ctx context.Context, userID uuid.UUID) error
}

type LoginCodePromptManager interface {
	NewPrompt(ctx context.Context, userID uuid.UUID) error
	GetUserID(ctx context.Context) (uuid.UUID, error)
}

type Repository interface {
	GetUserByUsernameAndStoreCode(ctx context.Context, username string, storeCode string) (domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(hashedPassword, password string) error
}

type Service struct {
	sessionManager         SessionManager
	loginCodePromptManager LoginCodePromptManager
	repo                   Repository
	hasher                 PasswordHasher
	loginCodePromptURL     url.URL
}

func NewService(
	sessionManager SessionManager,
	loginCodePromptManager LoginCodePromptManager,
	repo Repository,
	hasher PasswordHasher,
	loginCodePromptURL url.URL,
) *Service {
	return &Service{
		sessionManager:         sessionManager,
		loginCodePromptManager: loginCodePromptManager,
		repo:                   repo,
		hasher:                 hasher,
		loginCodePromptURL:     loginCodePromptURL,
	}
}

func (s *Service) Login(ctx context.Context, req *genapi.LoginCredentials) (genapi.LoginRes, error) {
	l := zerolog.Ctx(ctx)

	user, err := s.repo.GetUserByUsernameAndStoreCode(ctx, req.Username, req.StoreCode)
	if err != nil {
		return nil, apierror.ToAPIError(http.StatusUnauthorized, "invalid credentials")
	}

	err = s.hasher.Check(user.Password(), req.Password)
	if err != nil {
		if errors.Is(err, apperror.ErrPasswordMismatch) {
			return nil, apierror.ToAPIError(http.StatusUnauthorized, "invalid credentials")
		}

		l.Error().Err(err).Msg("PasswordHasher.Check(); failed to check password")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to check password")
	}

	if user.Role().IsStoreAdmin() {
		if err = s.sessionManager.NewSession(ctx, user.ID()); err != nil {
			l.Error().Err(err).Msg("SessionManager.NewSession(); failed to create session")
			return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to create session")
		}

		return &genapi.LoginNoContent{}, nil
	}

	if err = s.loginCodePromptManager.NewPrompt(ctx, user.ID()); err != nil {
		l.Error().Err(err).Msg("LoginCodePromptManager.NewPrompt(); failed to create login code prompt")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to create login code prompt")
	}

	return &genapi.LoginCodePromptRedirection{
		PromptURL: s.loginCodePromptURL,
	}, nil
}

func (s *Service) LoginCodePrompt(ctx context.Context, req *genapi.LoginCodePrompt, params genapi.LoginCodePromptParams) error {
	return nil
}
