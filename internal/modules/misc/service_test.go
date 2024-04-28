//go:build unit
// +build unit

package misc_test

import (
	"context"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/misc"
	"github.com/JosephJoshua/remana-backend/internal/modules/shared"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestGetHealth(t *testing.T) {
	t.Parallel()

	logger.Init(zerolog.ErrorLevel, shared.AppEnvDev)
	requestCtx := testutil.RequestContextWithLogger(context.Background())

	t.Run("returns no error", func(t *testing.T) {
		t.Parallel()

		s := misc.NewService()
		err := s.GetHealth(requestCtx)

		require.NoError(t, err)
	})
}
