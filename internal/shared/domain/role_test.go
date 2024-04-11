package domain_test

import (
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRole(t *testing.T) {
	t.Parallel()

	const (
		id        = 1
		name      = "role"
		storeCode = "storecode"
	)

	t.Run("new role with negative ID", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewRole(-1, name, domain.Store{}, false)

		require.ErrorIs(t, err, domain.ErrInvalidID)
		assert.Nil(t, got)
	})

	t.Run("new role with empty name", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewRole(1, "", domain.Store{}, false)

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Nil(t, got)
	})

	t.Run("new role with valid input", func(t *testing.T) {
		t.Parallel()

		store, err := domain.NewStore(1, "store", storeCode)
		require.NoError(t, err)

		got, err := domain.NewRole(id, name, *store, false)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, id, got.ID())
		assert.Equal(t, name, got.Name())
		assert.Equal(t, storeCode, got.StoreCode())
	})

	t.Run("set name with empty name", func(t *testing.T) {
		t.Parallel()

		role, err := domain.NewRole(id, name, domain.Store{}, false)
		require.NoError(t, err)

		err = role.SetName("")

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Equal(t, name, role.Name())
	})

	t.Run("set name with valid name", func(t *testing.T) {
		t.Parallel()

		const newName = "newname"

		role, err := domain.NewRole(id, name, domain.Store{}, false)
		require.NoError(t, err)

		err = role.SetName(newName)

		require.NoError(t, err)
		assert.Equal(t, newName, role.Name())
	})
}
