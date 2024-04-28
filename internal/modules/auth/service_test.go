//go:build unit
// +build unit

package auth_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth"
	"github.com/JosephJoshua/remana-backend/internal/shared"
	"github.com/JosephJoshua/remana-backend/internal/shared/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type serviceSessionManagerStub struct {
	userID *uuid.UUID
}

func (s *serviceSessionManagerStub) NewSession(_ context.Context, userID uuid.UUID) error {
	s.userID = &userID
	return nil
}

func (s *serviceSessionManagerStub) DeleteSession(_ context.Context) error {
	s.userID = nil
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

type serviceRepositoryStub struct {
	user             readmodel.AuthnUser
	username         string
	storeCode        string
	loginCode        string
	loginCodeDeleted bool
}

func (a *serviceRepositoryStub) GetUserByUsernameAndStoreCode(
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

func (a *serviceRepositoryStub) CheckAndDeleteUserLoginCode(
	_ context.Context,
	userID uuid.UUID,
	loginCode string,
) error {
	if a.user.ID.String() != userID.String() || a.loginCode != loginCode {
		return apperror.ErrLoginCodeMismatch
	}

	a.loginCodeDeleted = true
	return nil
}

func TestLogin(t *testing.T) {
	t.Parallel()

	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)
	requestCtx := testutil.RequestContextWithLogger(context.Background())

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

	t.Run("returns unauthorized", func(t *testing.T) {
		testCases := []struct {
			name string
			req  *genapi.LoginCredentials
		}{
			{
				name: "when password is wrong",
				req: &genapi.LoginCredentials{
					Username:  correctUsername,
					Password:  wrongPassword,
					StoreCode: correctStoreCode,
				},
			},
			{
				name: "when store code is wrong",
				req: &genapi.LoginCredentials{
					Username:  correctUsername,
					Password:  correctPassword,
					StoreCode: wrongStoreCode,
				},
			},
			{
				name: "when username is wrong",
				req: &genapi.LoginCredentials{
					Username:  wrongUsername,
					Password:  correctPassword,
					StoreCode: correctStoreCode,
				},
			},
		}

		for _, tc := range testCases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				sessionManager := new(serviceSessionManagerStub)
				loginCodePromptManager := new(loginCodePromptManagerStub)

				repo := &serviceRepositoryStub{
					user:             adminUser,
					username:         correctUsername,
					storeCode:        correctStoreCode,
					loginCode:        "",
					loginCodeDeleted: false,
				}

				s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})
				_, err := s.Login(requestCtx, tc.req)

				testutil.AssertAPIStatusCode(t, http.StatusUnauthorized, err)

				assert.Nil(t, sessionManager.userID)
				assert.Nil(t, loginCodePromptManager.userID)
			})
		}
	})

	t.Run("creates new session when user is store admin and credentials are correct", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(serviceSessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := &serviceRepositoryStub{
			user:             adminUser,
			username:         correctUsername,
			storeCode:        correctStoreCode,
			loginCode:        "",
			loginCodeDeleted: false,
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		got, err := s.Login(requestCtx, &genapi.LoginCredentials{
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

		sessionManager := new(serviceSessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := &serviceRepositoryStub{
			user:             employeeUser,
			username:         correctUsername,
			storeCode:        correctStoreCode,
			loginCode:        "",
			loginCodeDeleted: false,
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		got, err := s.Login(requestCtx, &genapi.LoginCredentials{
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

	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)
	requestCtx := testutil.RequestContextWithLogger(context.Background())

	var userID = uuid.New()
	const loginCode = "A1B2C3D4"

	user := readmodel.AuthnUser{
		ID:           userID,
		Password:     "password",
		IsStoreAdmin: false,
	}

	t.Run("returns bad request when prompt hasn't been initiated yet", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(serviceSessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)
		repo := &serviceRepositoryStub{
			user:             user,
			loginCode:        loginCode,
			username:         "testuser",
			storeCode:        "teststore",
			loginCodeDeleted: false,
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		err := s.LoginCodePrompt(requestCtx, &genapi.LoginCodePrompt{
			LoginCode: "12345678",
		})

		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)

		assert.Nil(t, sessionManager.userID)
		assert.False(t, repo.loginCodeDeleted)
	})

	t.Run("returns bad request when login code is wrong", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(serviceSessionManagerStub)
		loginCodePromptManager := &loginCodePromptManagerStub{userID: &userID}
		repo := &serviceRepositoryStub{
			user:             user,
			loginCode:        loginCode,
			username:         "testuser",
			storeCode:        "teststore",
			loginCodeDeleted: false,
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		err := s.LoginCodePrompt(requestCtx, &genapi.LoginCodePrompt{
			LoginCode: "12345678",
		})

		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)

		assert.Nil(t, sessionManager.userID)
		assert.NotNil(t, loginCodePromptManager.userID)
		assert.Equal(t, userID.String(), loginCodePromptManager.userID.String())
		assert.False(t, repo.loginCodeDeleted)
	})

	t.Run("creates new session when login code is correct", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(serviceSessionManagerStub)
		loginCodePromptManager := &loginCodePromptManagerStub{userID: &userID}
		repo := &serviceRepositoryStub{
			user:             user,
			loginCode:        loginCode,
			username:         "testuser",
			storeCode:        "teststore",
			loginCodeDeleted: false,
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, &passwordHasherStub{})

		err := s.LoginCodePrompt(requestCtx, &genapi.LoginCodePrompt{
			LoginCode: loginCode,
		})

		require.NoError(t, err)
		require.NotNil(t, sessionManager.userID)
		assert.Equal(t, userID.String(), sessionManager.userID.String())
		assert.Nil(t, loginCodePromptManager.userID)
		assert.True(t, repo.loginCodeDeleted)
	})
}

func TestLogout(t *testing.T) {
	t.Parallel()

	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)
	requestCtx := testutil.RequestContextWithLogger(context.Background())

	t.Run("deletes session", func(t *testing.T) {
		t.Parallel()

		var userID = uuid.New()

		sessionManager := &serviceSessionManagerStub{userID: &userID}

		s := auth.NewService(
			sessionManager,
			new(loginCodePromptManagerStub),
			new(serviceRepositoryStub),
			new(passwordHasherStub),
		)

		err := s.Logout(requestCtx)

		require.NoError(t, err)
		assert.Nil(t, sessionManager.userID)
	})
}
