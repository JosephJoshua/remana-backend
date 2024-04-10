package auth_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/JosephJoshua/repair-management-backend/internal/auth"
	"github.com/JosephJoshua/repair-management-backend/internal/genapi"
	"github.com/JosephJoshua/repair-management-backend/internal/shared/domain"
	"github.com/JosephJoshua/repair-management-backend/internal/shared/repository"
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
		return uuid.UUID{}, errors.New("invalid user ID")
	}

	return *l.userID, nil
}

type passwordHasherStub struct{}

func (p *passwordHasherStub) Hash(password string) (string, error) {
	return password, nil
}

func (p *passwordHasherStub) Check(hashedPassword, password string) error {
	if hashedPassword != password {
		return errors.New("passwords do not match")
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

	adminRole, initError := domain.NewRole(1, "admin", *store, true)
	require.NoError(t, initError)

	employeeRole, initError := domain.NewRole(2, "employee", *store, false)
	require.NoError(t, initError)

	adminUser, initError := domain.NewUser(id, correctUsername, correctPassword, *store, *adminRole)
	require.NoError(t, initError)

	employeeUser, initError := domain.NewUser(id, correctUsername, correctPassword, *store, *employeeRole)
	require.NoError(t, initError)

	t.Run("failed login with wrong password", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := repository.NewMemoryAuthRepository([]domain.User{*adminUser})

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{}, url.URL{})

		_, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  wrongPassword,
			StoreCode: correctStoreCode,
		})

		require.Error(t, err)
		assert.Nil(t, sessionManager.userID)
		assert.Nil(t, loginCodePromptManager.userID)
	})

	t.Run("failed login with wrong store code", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := repository.NewMemoryAuthRepository([]domain.User{*adminUser})

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{}, url.URL{})

		_, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  correctPassword,
			StoreCode: wrongStoreCode,
		})

		require.Error(t, err)
		assert.Nil(t, sessionManager.userID)
		assert.Nil(t, loginCodePromptManager.userID)
	})

	t.Run("failed login with wrong username", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := repository.NewMemoryAuthRepository([]domain.User{*adminUser})

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{}, url.URL{})

		_, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  wrongUsername,
			Password:  correctPassword,
			StoreCode: correctStoreCode,
		})

		require.Error(t, err)
		assert.Nil(t, sessionManager.userID)
		assert.Nil(t, loginCodePromptManager.userID)
	})

	t.Run("admin successful login", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := repository.NewMemoryAuthRepository([]domain.User{*adminUser})

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{}, url.URL{})

		got, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  correctPassword,
			StoreCode: correctStoreCode,
		})

		require.NoError(t, err)
		assert.IsType(t, new(genapi.LoginNoContent), got)
		assert.EqualValues(t, adminUser.ID(), *sessionManager.userID)
		assert.Nil(t, loginCodePromptManager.userID)
	})

	t.Run("employee login code prompt", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := repository.NewMemoryAuthRepository([]domain.User{*employeeUser})

		promptURL := url.URL{
			Scheme: "https",
			Host:   "example.com",
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{}, promptURL)

		got, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  correctPassword,
			StoreCode: correctStoreCode,
		})

		require.NoError(t, err)

		if got2, ok := got.(*genapi.LoginCodePromptRedirection); ok {
			assert.Equal(t, promptURL.String(), got2.PromptURL.String())
		} else {
			t.Errorf("expected *genapi.LoginCodePromptRedirection as response, got %T", got)
		}

		assert.Equal(t, employeeUser.ID(), *loginCodePromptManager.userID)
		assert.Nil(t, sessionManager.userID)
	})
}
