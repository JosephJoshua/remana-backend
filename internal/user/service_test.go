//go:build unit
// +build unit

package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/shared"
	"github.com/JosephJoshua/remana-backend/internal/shared/logger"
	"github.com/JosephJoshua/remana-backend/internal/shared/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/JosephJoshua/remana-backend/internal/user"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMyUserDetails(t *testing.T) {
	t.Parallel()

	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)
	requestCtx := testutil.RequestContextWithLogger(context.Background())

	t.Run("returns internal server error if user is missing from context", func(t *testing.T) {
		t.Parallel()

		s := user.NewService()
		_, err := s.GetMyUserDetails(requestCtx)

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})

	t.Run("returns user details", func(t *testing.T) {
		t.Parallel()

		s := user.NewService()

		user := readmodel.UserDetails{
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

		ctx := shared.NewContextWithUser(requestCtx, &user)
		got, err := s.GetMyUserDetails(ctx)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, user.ID.String(), got.ID.String())
		assert.Equal(t, user.Username, got.Username)
		assert.Equal(t, user.Role.ID.String(), got.Role.ID.String())
		assert.Equal(t, user.Role.Name, got.Role.Name)
		assert.Equal(t, user.Role.IsStoreAdmin, got.Role.IsStoreAdmin)
		assert.Equal(t, user.Store.ID.String(), got.Store.ID.String())
		assert.Equal(t, user.Store.Name, got.Store.Name)
		assert.Equal(t, user.Store.Code, got.Store.Code)
	})
}
