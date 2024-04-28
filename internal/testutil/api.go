package testutil

import (
	"context"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RequestContextWithLogger(ctx context.Context) context.Context {
	l := logger.MustGet()
	return l.WithContext(ctx)
}

func AssertAPIStatusCode(t *testing.T, expected int, got error) {
	t.Helper()

	require.Error(t, got)

	var apiErr *genapi.ErrorStatusCode
	require.ErrorAs(t, got, &apiErr, "expected error to be of type *genapi.ErrorStatusCode")

	assert.Equal(t, expected, apiErr.StatusCode, "expected status code to be %d, got %+v", expected, apiErr)
}
