//go:build unit
// +build unit

package auth_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type securityHandlerSessionManagerStub struct {
	userID *uuid.UUID
	err    error
}

func (s *securityHandlerSessionManagerStub) GetUserID(_ context.Context) (uuid.UUID, error) {
	if s.err != nil {
		return uuid.UUID{}, s.err
	}

	if s.userID == nil {
		return uuid.UUID{}, apperror.ErrMissingSession
	}

	return *s.userID, nil
}

type securityHandlerRepositoryStub struct {
	userDetails *readmodel.UserDetails
	err         error
}

func (s *securityHandlerRepositoryStub) GetUserDetailsByID(
	_ context.Context,
	_ uuid.UUID,
) (readmodel.UserDetails, error) {
	var emptyUserDetails readmodel.UserDetails

	if s.err != nil {
		return emptyUserDetails, s.err
	}

	if s.userDetails == nil {
		return emptyUserDetails, apperror.ErrUserNotFound
	}

	return *s.userDetails, nil
}

func TestHandleSessionCookie(t *testing.T) {
	t.Parallel()

	t.Run("returns unauthorized when session manager errors", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()

		sh := auth.NewSecurityHandler(
			&securityHandlerSessionManagerStub{userID: &userID, err: errors.New("oh no error")},
			&securityHandlerRepositoryStub{userDetails: nil, err: nil},
		)

		ctx := context.Background()
		ctx, err := sh.HandleSessionCookie(ctx, "", genapi.SessionCookie{APIKey: ""})

		var apiErr *genapi.ErrorStatusCode
		require.ErrorAs(t, err, &apiErr)

		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)

		_, ok := appcontext.GetUserFromContext(ctx)
		assert.False(t, ok)
	})

	t.Run("returns unauthorized when session is missing", func(t *testing.T) {
		t.Parallel()

		sh := auth.NewSecurityHandler(
			&securityHandlerSessionManagerStub{userID: nil, err: nil},
			&securityHandlerRepositoryStub{userDetails: nil, err: nil},
		)

		ctx := context.Background()
		ctx, err := sh.HandleSessionCookie(ctx, "", genapi.SessionCookie{APIKey: ""})

		var apiErr *genapi.ErrorStatusCode
		require.ErrorAs(t, err, &apiErr)

		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)

		_, ok := appcontext.GetUserFromContext(ctx)
		assert.False(t, ok)
	})

	t.Run("returns unauthorized when user not found", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()

		sh := auth.NewSecurityHandler(
			&securityHandlerSessionManagerStub{userID: &userID, err: nil},
			&securityHandlerRepositoryStub{userDetails: nil, err: nil},
		)

		ctx := context.Background()
		ctx, err := sh.HandleSessionCookie(ctx, "", genapi.SessionCookie{APIKey: ""})

		var apiErr *genapi.ErrorStatusCode
		require.ErrorAs(t, err, &apiErr)

		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)

		_, ok := appcontext.GetUserFromContext(ctx)
		assert.False(t, ok)
	})

	t.Run("returns internal server error when repo errors", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()

		var userDetails readmodel.UserDetails

		sh := auth.NewSecurityHandler(
			&securityHandlerSessionManagerStub{userID: &userID, err: nil},
			&securityHandlerRepositoryStub{userDetails: &userDetails, err: errors.New("oh no error")},
		)

		ctx := context.Background()
		ctx, err := sh.HandleSessionCookie(ctx, "", genapi.SessionCookie{APIKey: ""})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)

		_, ok := appcontext.GetUserFromContext(ctx)
		assert.False(t, ok)
	})

	t.Run("adds user details to context", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()

		userDetails := readmodel.UserDetails{
			ID:       uuid.New(),
			Username: "username",
			Role: readmodel.UserDetailsRole{
				ID:           uuid.New(),
				Name:         "role",
				IsStoreAdmin: true,
			},
			Store: readmodel.UserDetailsStore{
				ID:   uuid.New(),
				Name: "store",
				Code: "code",
			},
		}

		sh := auth.NewSecurityHandler(
			&securityHandlerSessionManagerStub{userID: &userID, err: nil},
			&securityHandlerRepositoryStub{userDetails: &userDetails, err: nil},
		)

		ctx := context.Background()
		ctx, err := sh.HandleSessionCookie(ctx, "", genapi.SessionCookie{APIKey: ""})

		require.NoError(t, err)

		got, ok := appcontext.GetUserFromContext(ctx)

		require.True(t, ok)
		assert.EqualExportedValues(t, userDetails, *got)
	})
}
