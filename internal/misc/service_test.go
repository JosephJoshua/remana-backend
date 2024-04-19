//go:build unit
// +build unit

package misc_test

import (
	"context"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/misc"
	"github.com/stretchr/testify/require"
)

func TestGetHealth(t *testing.T) {
	t.Parallel()

	t.Run("returns no error", func(t *testing.T) {
		t.Parallel()

		s := misc.NewService()
		err := s.GetHealth(context.Background())

		require.NoError(t, err)
	})
}
