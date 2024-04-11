package domain_test

import (
	"testing"

	"github.com/JosephJoshua/repair-management-backend/internal/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	t.Parallel()

	const (
		id   = 1
		name = "store"
		code = "store-code"
	)

	t.Run("new store with negative ID", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewStore(-1, name, code)

		require.ErrorIs(t, err, domain.ErrInvalidID)
		assert.Nil(t, got)
	})

	t.Run("new store with empty name", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewStore(id, "", code)

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Nil(t, got)
	})

	t.Run("new store with empty code", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewStore(id, name, "")

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Nil(t, got)
	})

	t.Run("new store with uppercase code", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewStore(id, name, "ABC")

		require.ErrorIs(t, err, domain.ErrInvalidStoreCode)
		assert.Nil(t, got)
	})

	t.Run("new store with code containing symbols", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewStore(id, name, "a_b_c")

		require.ErrorIs(t, err, domain.ErrInvalidStoreCode)
		assert.Nil(t, got)
	})

	t.Run("new store with code containing numbers", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewStore(id, name, "a123")

		require.ErrorIs(t, err, domain.ErrInvalidStoreCode)
		assert.Nil(t, got)
	})

	t.Run("new store with valid input", func(t *testing.T) {
		t.Parallel()

		got, err := domain.NewStore(id, name, code)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, id, got.ID())
		assert.Equal(t, name, got.Name())
		assert.Equal(t, code, got.Code())
	})

	t.Run("set name with empty name", func(t *testing.T) {
		t.Parallel()

		store, err := domain.NewStore(id, name, code)
		require.NoError(t, err)

		err = store.SetName("")

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Equal(t, name, store.Name())
	})

	t.Run("set name with valid name", func(t *testing.T) {
		t.Parallel()

		const newName = "newname"

		store, err := domain.NewStore(id, name, code)
		require.NoError(t, err)

		err = store.SetName(newName)

		require.NoError(t, err)
		assert.Equal(t, newName, store.Name())
	})

	t.Run("set code with empty code", func(t *testing.T) {
		t.Parallel()

		store, err := domain.NewStore(id, name, code)
		require.NoError(t, err)

		err = store.SetCode("")

		require.ErrorIs(t, err, domain.ErrInputTooShort)
		assert.Equal(t, code, store.Code())
	})

	t.Run("set code with uppercase code", func(t *testing.T) {
		t.Parallel()

		store, err := domain.NewStore(id, name, code)
		require.NoError(t, err)

		err = store.SetCode("ABC")

		require.ErrorIs(t, err, domain.ErrInvalidStoreCode)
		assert.Equal(t, code, store.Code())
	})

	t.Run("set code with code containing symbols", func(t *testing.T) {
		t.Parallel()

		store, err := domain.NewStore(id, name, code)
		require.NoError(t, err)

		err = store.SetCode("a_b_c")

		require.ErrorIs(t, err, domain.ErrInvalidStoreCode)
		assert.Equal(t, code, store.Code())
	})

	t.Run("set code with code containing numbers", func(t *testing.T) {
		t.Parallel()

		store, err := domain.NewStore(id, name, code)
		require.NoError(t, err)

		err = store.SetCode("a123")

		require.ErrorIs(t, err, domain.ErrInvalidStoreCode)
		assert.Equal(t, code, store.Code())
	})

	t.Run("set code with valid code", func(t *testing.T) {
		t.Parallel()

		const newCode = "new-code"

		store, err := domain.NewStore(id, name, code)
		require.NoError(t, err)

		err = store.SetCode(newCode)

		require.NoError(t, err)
		assert.Equal(t, newCode, store.Code())
	})
}
