//go:build unit
// +build unit

package auth_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/modules/shared"
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

type serviceRepositoryStub struct {
	user              readmodel.User
	username          string
	storeCode         string
	loginCode         string
	loginCodeDeleted  bool
	getUserErr        error
	checkLoginCodeErr error
}

func (a *serviceRepositoryStub) GetUserByUsernameAndStoreCode(
	_ context.Context,
	username string,
	storeCode string,
) (readmodel.User, error) {
	var emptyUser readmodel.User

	if a.getUserErr != nil {
		return emptyUser, a.getUserErr
	}

	if a.username != username || a.storeCode != storeCode {
		return emptyUser, apperror.ErrUserNotFound
	}

	return a.user, nil
}

func (a *serviceRepositoryStub) CheckAndDeleteUserLoginCode(
	_ context.Context,
	userID uuid.UUID,
	loginCode string,
) error {
	if a.checkLoginCodeErr != nil {
		return a.checkLoginCodeErr
	}

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

	adminUser := readmodel.User{
		ID:           userID,
		Password:     correctPassword,
		IsStoreAdmin: true,
	}

	employeeUser := readmodel.User{
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

				s := auth.NewService(sessionManager, loginCodePromptManager, repo, testutil.PasswordHasherStub{})
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

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, testutil.PasswordHasherStub{})

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

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, testutil.PasswordHasherStub{})

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

	t.Run("returns internal server error when repository.GetUserByUsernameAndStoreCode() errors", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(serviceSessionManagerStub)
		loginCodePromptManager := new(loginCodePromptManagerStub)

		repo := &serviceRepositoryStub{
			user:             adminUser,
			username:         correctUsername,
			storeCode:        correctStoreCode,
			loginCode:        "",
			loginCodeDeleted: false,
			getUserErr:       errors.New("oh no!"),
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, testutil.PasswordHasherStub{})
		_, err := s.Login(requestCtx, &genapi.LoginCredentials{
			Username:  correctUsername,
			Password:  correctPassword,
			StoreCode: correctStoreCode,
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)

		assert.Nil(t, sessionManager.userID)
		assert.Nil(t, loginCodePromptManager.userID)
	})
}

func TestLoginCodePrompt(t *testing.T) {
	t.Parallel()

	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)
	requestCtx := testutil.RequestContextWithLogger(context.Background())

	var userID = uuid.New()
	const loginCode = "A1B2C3D4"

	user := readmodel.User{
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

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, testutil.PasswordHasherStub{})

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

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, testutil.PasswordHasherStub{})

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

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, testutil.PasswordHasherStub{})

		err := s.LoginCodePrompt(requestCtx, &genapi.LoginCodePrompt{
			LoginCode: loginCode,
		})

		require.NoError(t, err)
		require.NotNil(t, sessionManager.userID)
		assert.Equal(t, userID.String(), sessionManager.userID.String())
		assert.Nil(t, loginCodePromptManager.userID)
		assert.True(t, repo.loginCodeDeleted)
	})

	t.Run("returns internal server error when repository.CheckAndDeleteLoginCode() errors", func(t *testing.T) {
		t.Parallel()

		sessionManager := new(serviceSessionManagerStub)
		loginCodePromptManager := &loginCodePromptManagerStub{userID: &userID}
		repo := &serviceRepositoryStub{
			user:              user,
			loginCode:         loginCode,
			username:          "testuser",
			storeCode:         "teststore",
			loginCodeDeleted:  false,
			checkLoginCodeErr: errors.New("oh no!"),
		}

		s := auth.NewService(sessionManager, loginCodePromptManager, repo, testutil.PasswordHasherStub{})

		err := s.LoginCodePrompt(requestCtx, &genapi.LoginCodePrompt{
			LoginCode: "12345678",
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)

		assert.Nil(t, sessionManager.userID)
		assert.NotNil(t, loginCodePromptManager.userID)
		assert.Equal(t, userID.String(), loginCodePromptManager.userID.String())
		assert.False(t, repo.loginCodeDeleted)
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
			new(testutil.PasswordHasherStub),
		)

		err := s.Logout(requestCtx)

		require.NoError(t, err)
		assert.Nil(t, sessionManager.userID)
	})
}
