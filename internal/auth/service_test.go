package auth_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/auth"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/shared/apperror"
	"github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/JosephJoshua/remana-backend/internal/shared/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sessionManagerStub struct {
	userID *uuid.UUID
}

func (s *sessionManagerStub) NewSession(_ context.Context, userID uuid.UUID) error {
	s.userID = &userID
	return nil
}

type loginCodePromptManagerStub struct {
	userID *uuid.UUID
}

func (l *loginCodePromptManagerStub) NewPrompt(_ context.Context, userID uuid.UUID) error {
	l.userID = &userID
	return nil
}

func (l *loginCodePromptManagerStub) GetUserID(_ context.Context) (uuid.UUID, error) {
	if l.userID == nil {
		return uuid.UUID{}, apperror.ErrMisingLoginCodePrompt
	}

	return *l.userID, nil
}

func (l *loginCodePromptManagerStub) DeletePrompt(_ context.Context) error {
	l.userID = nil
	return nil
}

type passwordHasherStub struct{}

func (p *passwordHasherStub) Hash(password string) (string, error) {
	return password, nil
}

func (p *passwordHasherStub) Check(hashedPassword, password string) error {
	if hashedPassword != password {
		return apperror.ErrPasswordMismatch
	}

	return nil
}

func TestLogin(t *testing.T) {
	t.Parallel()

	const (
		correctUsername  = "testuser"
		correctPassword  = "testpassword"
		correctStoreCode = "teststore"

		wrongUsername  = "wronguser"
		wrongPassword  = "wrongpassword"
		wrongStoreCode = "wrongstore"
	)

	id, initError := uuid.NewUUID()
	require.NoError(t, initError)

	store, initError := domain.NewStore(1, correctStoreCode, correctStoreCode)
	require.NoError(t, initError)

	adminRole, initError := domain.NewRole(1, "admin", store, true)
	require.NoError(t, initError)

	employeeRole, initError := domain.NewRole(2, "employee", store, false)
	require.NoError(t, initError)

	adminUser, initError := domain.NewUser(id, correctUsername, correctPassword, store, adminRole)
	require.NoError(t, initError)

	employeeUser, initError := domain.NewUser(id, correctUsername, correctPassword, store, employeeRole)
	require.NoError(t, initError)

	t.Run("failed login with wrong password", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := repository.NewMemoryAuthRepository([]domain.User{adminUser}, []domain.LoginCode{})

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		_, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  wrongPassword,
			StoreCode: correctStoreCode,
		})

		var apiError *genapi.ErrorStatusCode

		require.ErrorAs(t, err, &apiError)
		assert.Equal(t, http.StatusUnauthorized, apiError.GetStatusCode())

		assert.Nil(t, sessionManager.userID)
		assert.Nil(t, loginCodePromptManager.userID)
	})

	t.Run("failed login with wrong store code", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := repository.NewMemoryAuthRepository([]domain.User{adminUser}, []domain.LoginCode{})

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		_, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  correctPassword,
			StoreCode: wrongStoreCode,
		})

		var apiError *genapi.ErrorStatusCode

		require.ErrorAs(t, err, &apiError)
		assert.Equal(t, http.StatusUnauthorized, apiError.GetStatusCode())

		assert.Nil(t, sessionManager.userID)
		assert.Nil(t, loginCodePromptManager.userID)
	})

	t.Run("failed login with wrong username", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := repository.NewMemoryAuthRepository([]domain.User{adminUser}, []domain.LoginCode{})

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		_, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  wrongUsername,
			Password:  correctPassword,
			StoreCode: correctStoreCode,
		})

		var apiError *genapi.ErrorStatusCode

		require.ErrorAs(t, err, &apiError)
		assert.Equal(t, http.StatusUnauthorized, apiError.GetStatusCode())

		assert.Nil(t, sessionManager.userID)
		assert.Nil(t, loginCodePromptManager.userID)
	})

	t.Run("admin successful login", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := repository.NewMemoryAuthRepository([]domain.User{adminUser}, []domain.LoginCode{})

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		got, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  correctPassword,
			StoreCode: correctStoreCode,
		})

		require.NoError(t, err)
		assert.Equal(t, genapi.LoginResponseTypeAdmin, got.Type)
		assert.NotNil(t, sessionManager.userID)
		assert.Equal(t, adminUser.ID().String(), sessionManager.userID.String())
		assert.Nil(t, loginCodePromptManager.userID)
	})

	t.Run("employee login code prompt", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := repository.NewMemoryAuthRepository([]domain.User{employeeUser}, []domain.LoginCode{})

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		got, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  correctPassword,
			StoreCode: correctStoreCode,
		})

		require.NoError(t, err)
		assert.Equal(t, genapi.LoginResponseTypeEmployee, got.Type)
		assert.NotNil(t, loginCodePromptManager.userID)
		assert.Equal(t, employeeUser.ID().String(), loginCodePromptManager.userID.String())
		assert.Nil(t, sessionManager.userID)
	})
}

func TestLoginCodePrompt(t *testing.T) {
	t.Parallel()

	var userID = uuid.New()
	const code = "A1B2C3D4"

	store, initErr := domain.NewStore(1, "store", "store")
	require.NoError(t, initErr)

	role, initErr := domain.NewRole(1, "store", store, false)
	require.NoError(t, initErr)

	user, initErr := domain.NewUser(userID, "username", "password", store, role)
	require.NoError(t, initErr)

	loginCode, initErr := domain.NewLoginCode(user, code)
	require.NoError(t, initErr)

	repo := repository.NewMemoryAuthRepository([]domain.User{user}, []domain.LoginCode{loginCode})

	t.Run("empty prompt", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		err := s.LoginCodePrompt(context.Background(), &genapi.LoginCodePrompt{
			LoginCode: "12345678",
		})

		var apiError *genapi.ErrorStatusCode

		require.ErrorAs(t, err, &apiError)
		assert.Equal(t, http.StatusBadRequest, apiError.GetStatusCode())
		assert.Nil(t, sessionManager.userID)
	})

	t.Run("wrong login code", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := &loginCodePromptManagerStub{userID: &userID}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		err := s.LoginCodePrompt(context.Background(), &genapi.LoginCodePrompt{
			LoginCode: "12345678",
		})

		var apiError *genapi.ErrorStatusCode

		require.ErrorAs(t, err, &apiError)
		assert.Equal(t, http.StatusBadRequest, apiError.GetStatusCode())
		assert.Nil(t, sessionManager.userID)
		assert.NotNil(t, loginCodePromptManager.userID)
		assert.Equal(t, userID.String(), loginCodePromptManager.userID.String())
	})

	t.Run("correct login code", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := &loginCodePromptManagerStub{userID: &userID}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		err := s.LoginCodePrompt(context.Background(), &genapi.LoginCodePrompt{
			LoginCode: code,
		})

		require.NoError(t, err)
		require.NotNil(t, sessionManager.userID)
		assert.Equal(t, userID.String(), sessionManager.userID.String())
		assert.Nil(t, loginCodePromptManager.userID)
	})
}
