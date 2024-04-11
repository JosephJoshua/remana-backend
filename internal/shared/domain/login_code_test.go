package domain_test

import (
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginCode(t *testing.T) {
	t.Parallel()

	var userID = uuid.New()

	store, initErr := domain.NewStore(1, "store", "storecode")
	require.NoError(t, initErr)

	role, initErr := domain.NewRole(1, "role", store, false)
	require.NoError(t, initErr)

	user, initErr := domain.NewUser(userID, "username", "password", store, role)
	require.NoError(t, initErr)

	t.Run("new login code with code too short", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewLoginCode(user, "1234")

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Nil(t, got)
	})

	t.Run("new login code with code too long", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewLoginCode(user, "123456789")

		require.ErrorIs(t, err, domain.ErrInputTooLong)
		assert.Nil(t, got)
	})

	t.Run("new login code with code containing symbols", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewLoginCode(user, "1234567_")

		require.ErrorIs(t, err, domain.ErrInvalidLoginCode)
		assert.Nil(t, got)
	})

	t.Run("new login code with valid input", func(t *testing.T) {
		t.Parallel()

		const code = "A1B2C3D4"

		got, err := domain.NewLoginCode(user, code)

		require.NoError(t, err)

		assert.Equal(t, code, got.Code())
		assert.Equal(t, userID.String(), got.UserID().String())
	})

	t.Run("new login code with lowercase code", func(t *testing.T) {
		t.Parallel()

		const code = "a1b2c3d4"

		got, err := domain.NewLoginCode(user, code)

		require.NoError(t, err)

		assert.Equal(t, "A1B2C3D4", got.Code())
		assert.Equal(t, userID.String(), got.UserID().String())
	})
}
