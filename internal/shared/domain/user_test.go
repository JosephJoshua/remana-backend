package domain_test

import (
	"testing"

	"github.com/JosephJoshua/repair-management-backend/internal/shared/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	t.Parallel()

	var id = uuid.New()

	const (
		username = "username"
		password = "password"
	)

	t.Run("new user with username too short", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewUser(id, "123", password, domain.Store{}, domain.Role{})

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Nil(t, got)
	})

	t.Run("new user with username too long", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewUser(
			id,
			"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901",
			password,
			domain.Store{},
			domain.Role{},
		)

		require.ErrorIs(t, err, domain.ErrInputTooLong)
		assert.Nil(t, got)
	})

	t.Run("new user with valid input", func(t *testing.T) {
		t.Parallel()

		const storeCode = "storecode"

		store, err := domain.NewStore(1, "store", storeCode)
		require.NoError(t, err)

		role, err := domain.NewRole(1, "role", *store, false)
		require.NoError(t, err)

		got, err := domain.NewUser(id, username, password, *store, *role)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, id.String(), got.ID().String())
		assert.Equal(t, username, got.Username())
		assert.Equal(t, password, got.Password())
		assert.Equal(t, storeCode, got.StoreCode())
		assert.EqualValues(t, role, got.Role())
	})

	t.Run("set username with too short input", func(t *testing.T) {
		t.Parallel()

		user, err := domain.NewUser(id, username, password, domain.Store{}, domain.Role{})
		require.NoError(t, err)

		err = user.SetUsername("123")

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Equal(t, username, user.Username())
	})

	t.Run("set username with too long input", func(t *testing.T) {
		t.Parallel()

		user, err := domain.NewUser(id, username, password, domain.Store{}, domain.Role{})
		require.NoError(t, err)

		err = user.SetUsername(
			"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901",
		)

		require.ErrorIs(t, err, domain.ErrInputTooLong)
		assert.Equal(t, username, user.Username())
	})

	t.Run("set username with valid input", func(t *testing.T) {
		t.Parallel()

		const newUsername = "newusername"

		user, err := domain.NewUser(id, username, password, domain.Store{}, domain.Role{})
		require.NoError(t, err)

		err = user.SetUsername(newUsername)

		require.NoError(t, err)
		assert.Equal(t, newUsername, user.Username())
	})
}
