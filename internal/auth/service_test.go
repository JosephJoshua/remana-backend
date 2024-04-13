package auth_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/auth"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/shared/apperror"
	"github.com/JosephJoshua/remana-backend/internal/shared/readmodel"
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

type authRepositoryStub struct {
	user             readmodel.AuthnUser
	username         string
	storeCode        string
	loginCode        string
	loginCodeDeleted bool
}

func (a *authRepositoryStub) GetUserByUsernameAndStoreCode(
	_ context.Context,
	username string,
	storeCode string,
) (readmodel.AuthnUser, error) {
	if a.username != username || a.storeCode != storeCode {
		var user readmodel.AuthnUser
		return user, apperror.ErrUserNotFound
	}

	return a.user, nil
}

func (a *authRepositoryStub) CheckAndDeleteUserLoginCode(_ context.Context, userID uuid.UUID, loginCode string) error {
	if a.user.ID.String() != userID.String() || a.loginCode != loginCode {
		return apperror.ErrLoginCodeMismatch
	}

	a.loginCodeDeleted = true
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

	var userID = uuid.New()

	adminUser := readmodel.AuthnUser{
		ID:           userID,
		Password:     correctPassword,
		IsStoreAdmin: true,
	}

	employeeUser := readmodel.AuthnUser{
		ID:           userID,
		Password:     correctPassword,
		IsStoreAdmin: false,
	}

	t.Run("returns unauthorized when password is wrong", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := &authRepositoryStub{
			user:             adminUser,
			username:         correctUsername,
			storeCode:        correctStoreCode,
			loginCode:        "",
			loginCodeDeleted: false,
		}
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

	t.Run("returns unauthorized when store code is wrong", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := &authRepositoryStub{
			user:             adminUser,
			username:         correctUsername,
			storeCode:        correctStoreCode,
			loginCode:        "",
			loginCodeDeleted: false,
		}

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

	t.Run("returns unauthorized when username is wrong", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := &authRepositoryStub{
			user:             adminUser,
			username:         correctUsername,
			storeCode:        correctStoreCode,
			loginCode:        "",
			loginCodeDeleted: false,
		}

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

	t.Run("creates new session when user is store admin and credentials are correct", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := &authRepositoryStub{
			user:             adminUser,
			username:         correctUsername,
			storeCode:        correctStoreCode,
			loginCode:        "",
			loginCodeDeleted: false,
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		got, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  correctPassword,
			StoreCode: correctStoreCode,
		})

		require.NoError(t, err)
		assert.Equal(t, genapi.LoginResponseTypeAdmin, got.Type)
		assert.NotNil(t, sessionManager.userID)
		assert.Equal(t, adminUser.ID.String(), sessionManager.userID.String())
		assert.Nil(t, loginCodePromptManager.userID)
	})

	t.Run("creates new login code prompt when user is store employee and credentials are correct", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := &authRepositoryStub{
			user:             employeeUser,
			username:         correctUsername,
			storeCode:        correctStoreCode,
			loginCode:        "",
			loginCodeDeleted: false,
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		got, err := s.Login(context.Background(), &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  correctPassword,
			StoreCode: correctStoreCode,
		})

		require.NoError(t, err)
		assert.Equal(t, genapi.LoginResponseTypeEmployee, got.Type)
		assert.NotNil(t, loginCodePromptManager.userID)
		assert.Equal(t, employeeUser.ID.String(), loginCodePromptManager.userID.String())
		assert.Nil(t, sessionManager.userID)
	})
}

func TestLoginCodePrompt(t *testing.T) {
	t.Parallel()

	var userID = uuid.New()
	const loginCode = "A1B2C3D4"

	user := readmodel.AuthnUser{
		ID:           userID,
		Password:     "password",
		IsStoreAdmin: false,
	}

	t.Run("returns bad request when prompt hasn't been initiated yet", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := &authRepositoryStub{
			user:             user,
			loginCode:        loginCode,
			username:         "testuser",
			storeCode:        "teststore",
			loginCodeDeleted: false,
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		err := s.LoginCodePrompt(context.Background(), &genapi.LoginCodePrompt{
			LoginCode: "12345678",
		})

		var apiError *genapi.ErrorStatusCode

		require.ErrorAs(t, err, &apiError)
		assert.Equal(t, http.StatusBadRequest, apiError.GetStatusCode())
		assert.Nil(t, sessionManager.userID)
		assert.False(t, repo.loginCodeDeleted)
	})

	t.Run("returns bad request when login code is wrong", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := &loginCodePromptManagerStub{userID: &userID}
		repo := &authRepositoryStub{
			user:             user,
			loginCode:        loginCode,
			username:         "testuser",
			storeCode:        "teststore",
			loginCodeDeleted: false,
		}

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
		assert.False(t, repo.loginCodeDeleted)
	})

	t.Run("creates new session when login code is correct", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(sessionManagerStub)
		loginCodePromptManager := &loginCodePromptManagerStub{userID: &userID}
		repo := &authRepositoryStub{
			user:             user,
			loginCode:        loginCode,
			username:         "testuser",
			storeCode:        "teststore",
			loginCodeDeleted: false,
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		err := s.LoginCodePrompt(context.Background(), &genapi.LoginCodePrompt{
			LoginCode: loginCode,
		})

		require.NoError(t, err)
		require.NotNil(t, sessionManager.userID)
		assert.Equal(t, userID.String(), sessionManager.userID.String())
		assert.Nil(t, loginCodePromptManager.userID)
		assert.True(t, repo.loginCodeDeleted)
	})
}
