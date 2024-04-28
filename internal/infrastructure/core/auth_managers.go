package core

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
)

const (
	userIDKey = "user_id"
)

type authSessionManager struct {
	sm *scs.SessionManager
}

func newAuthSessionManager() *authSessionManager {
	sm := scs.New()

	sm.Cookie.Name = "session_id"
	sm.Cookie.Secure = true

	return &authSessionManager{
		sm: sm,
	}
}

func (a *authSessionManager) NewSession(ctx context.Context, userID uuid.UUID) error {
	if err := a.sm.RenewToken(ctx); err != nil {
		return fmt.Errorf("failed to renew session token: %w", err)
	}

	a.sm.Put(ctx, userIDKey, userID.String())
	return nil
}

func (a *authSessionManager) DeleteSession(ctx context.Context) error {
	return a.sm.Destroy(ctx)
}

func (a *authSessionManager) GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID := a.sm.GetString(ctx, userIDKey)
	if userID == "" {
		return uuid.UUID{}, fmt.Errorf(
			"failed to get user ID from session: %w",
			apperror.ErrMissingSession,
		)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to parse user ID from session: %w", err)
	}

	return parsedUserID, nil
}

func (a *authSessionManager) middleware(next http.Handler) http.Handler {
	return a.sm.LoadAndSave(next)
}

const (
	loginCodePromptCookieLifetime    = 1 * time.Hour
	loginCodePromptCookieIdleTimeout = 30 * time.Minute
)

type loginCodePromptManager struct {
	sm *scs.SessionManager
}

func newLoginCodePromptManager() *loginCodePromptManager {
	sm := scs.New()

	sm.Lifetime = loginCodePromptCookieLifetime
	sm.IdleTimeout = loginCodePromptCookieIdleTimeout
	sm.Cookie.Name = "login_code_prompt_id"
	sm.Cookie.Secure = true

	return &loginCodePromptManager{
		sm: sm,
	}
}

func (l *loginCodePromptManager) NewPrompt(ctx context.Context, userID uuid.UUID) error {
	if err := l.sm.RenewToken(ctx); err != nil {
		return fmt.Errorf("failed to renew login code prompt token: %w", err)
	}

	l.sm.Put(ctx, userIDKey, userID.String())
	return nil
}

func (l *loginCodePromptManager) GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID := l.sm.GetString(ctx, userIDKey)
	if userID == "" {
		return uuid.UUID{}, fmt.Errorf(
			"failed to get user ID from login code prompt session: %w",
			apperror.ErrMisingLoginCodePrompt,
		)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to parse user ID from login code prompt session: %w", err)
	}

	return parsedUserID, nil
}

func (l *loginCodePromptManager) DeletePrompt(ctx context.Context) error {
	return l.sm.Destroy(ctx)
}

func (l *loginCodePromptManager) middleware(next http.Handler) http.Handler {
	return l.sm.LoadAndSave(next)
}
