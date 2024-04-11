package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/shared/apierror"
	"github.com/JosephJoshua/remana-backend/internal/shared/apperror"
	"github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type SessionManager interface {
	NewSession(ctx context.Context, userID uuid.UUID) error
}

type LoginCodePromptManager interface {
	NewPrompt(ctx context.Context, userID uuid.UUID) error
	GetUserID(ctx context.Context) (uuid.UUID, error)
	DeletePrompt(ctx context.Context) error
}

type Repository interface {
	GetUserByUsernameAndStoreCode(ctx context.Context, username string, storeCode string) (domain.User, error)
	CheckUserLoginCode(ctx context.Context, userID uuid.UUID, loginCode string) error
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
}

func NewService(
	sessionManager SessionManager,
	loginCodePromptManager LoginCodePromptManager,
	repo Repository,
	hasher PasswordHasher,
) *Service {
	return &Service{
		sessionManager:         sessionManager,
		loginCodePromptManager: loginCodePromptManager,
		repo:                   repo,
		hasher:                 hasher,
	}
}

func (s *Service) Login(ctx context.Context, req *genapi.LoginCredentials) (*genapi.LoginResponse, error) {
	l := zerolog.Ctx(ctx)

	user, err := s.repo.GetUserByUsernameAndStoreCode(ctx, req.Username, req.StoreCode)
	if err != nil {
		// Security measure to prevent timing attacks.
		_ = s.hasher.Check("$2a$14$7IotmYZSWWVoGd.D5xaMLOi2W0bBbHZfNZ0NxX.BpphGmNd9IbC/u", req.Password)

		l.
			Info().
			Str("username", req.GetUsername()).
			Str("store_code", req.GetStoreCode()).
			Msg("wrong username or store code")

		return nil, apierror.ToAPIError(http.StatusUnauthorized, "invalid credentials")
	}

	err = s.hasher.Check(user.Password(), req.Password)
	if err != nil {
		if errors.Is(err, apperror.ErrPasswordMismatch) {
			l.Info().Str("user_id", user.ID().String()).Msg("wrong password")
			return nil, apierror.ToAPIError(http.StatusUnauthorized, "invalid credentials")
		}

		l.Error().Err(err).Msg("PasswordHasher.Check(); failed to check password")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to check password")
	}

	if user.Role().IsStoreAdmin() {
		l.Info().Str("user_id", user.ID().String()).Msg("store admin logged in")

		if err = s.sessionManager.NewSession(ctx, user.ID()); err != nil {
			l.Error().Err(err).Msg("SessionManager.NewSession(); failed to create session")
			return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to create session")
		}

		return &genapi.LoginResponse{
			Type: genapi.LoginResponseTypeAdmin,
		}, nil
	}

	l.Info().Str("user_id", user.ID().String()).Msg("store employee login code prompt initiated")

	if err = s.loginCodePromptManager.NewPrompt(ctx, user.ID()); err != nil {
		l.Error().Err(err).Msg("LoginCodePromptManager.NewPrompt(); failed to create login code prompt")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to create login code prompt")
	}

	return &genapi.LoginResponse{
		Type: genapi.LoginResponseTypeEmployee,
	}, nil
}

func (s *Service) LoginCodePrompt(ctx context.Context, req *genapi.LoginCodePrompt) error {
	l := zerolog.Ctx(ctx)

	userID, err := s.loginCodePromptManager.GetUserID(ctx)
	if err != nil {
		if errors.Is(err, apperror.ErrMisingLoginCodePrompt) {
			return apierror.ToAPIError(http.StatusBadRequest, "missing login code prompt ID. please call /login first")
		}

		l.Error().Err(err).Msg("LoginCodePromptManager.GetUserID(); failed to get user ID")
		return apierror.ToAPIError(http.StatusInternalServerError, "failed to get user ID")
	}

	err = s.repo.CheckUserLoginCode(ctx, userID, req.GetLoginCode())
	if err != nil {
		if errors.Is(err, apperror.ErrLoginCodeMismatch) {
			l.Info().Str("user_id", userID.String()).Msg("wrong login code")
			return apierror.ToAPIError(http.StatusBadRequest, "wrong login code")
		}

		l.Error().Err(err).Msg("Repository.CheckUserLoginCode(); failed to check login code")
		return apierror.ToAPIError(http.StatusInternalServerError, "failed to check login code")
	}

	l.Info().Str("user_id", userID.String()).Msg("store employee logged in")

	if err = s.sessionManager.NewSession(ctx, userID); err != nil {
		l.Error().Err(err).Msg("SessionManager.NewSession(); failed to create session")
		return apierror.ToAPIError(http.StatusInternalServerError, "failed to create session")
	}

	if err = s.loginCodePromptManager.DeletePrompt(ctx); err != nil {
		l.Error().Err(err).Msg("LoginCodePromptManager.DeletePrompt(); failed to delete login code prompt")
		return apierror.ToAPIError(http.StatusInternalServerError, "failed to delete login code prompt")
	}

	return nil
}
