package domain_test

import (
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRole(t *testing.T) {
	t.Parallel()

	var id = uuid.New()

	const (
		name      = "role"
		storeCode = "storecode"
	)

	store, initErr := domain.NewStore(uuid.New(), "store", storeCode)
	require.NoError(t, initErr)

	t.Run("new role with empty name", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewRole(id, "", store, false)

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Nil(t, got)
	})

	t.Run("new role with valid input", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewRole(id, name, store, false)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, id.String(), got.ID().String())
		assert.Equal(t, name, got.Name())
		assert.Equal(t, storeCode, got.StoreCode())
	})

	t.Run("set name with empty name", func(t *testing.T) {
		t.Parallel()

		role, err := domain.NewRole(id, name, store, false)
		require.NoError(t, err)

		err = role.SetName("")

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Equal(t, name, role.Name())
	})

	t.Run("set name with valid name", func(t *testing.T) {
		t.Parallel()

		const newName = "newname"

		role, err := domain.NewRole(id, name, store, false)
		require.NoError(t, err)

		err = role.SetName(newName)

		require.NoError(t, err)
		assert.Equal(t, newName, role.Name())
	})
}
