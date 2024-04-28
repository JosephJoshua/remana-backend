//go:build unit
// +build unit

package domain_test

import (
	"fmt"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/modules/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPhoneNumber(t *testing.T) {
	testCases := []struct {
		input       string
		valid       bool
		expectedVal string
	}{
		{"+628123456789", true, "+628123456789"},
		{"08123456789", true, "+628123456789"},
		{"0812345678a", true, "+62812345678"},
		{"+62812345678a", true, "+62812345678"},
		{"+15417543010", true, "+15417543010"},
		{"0812345", false, ""},
		{"+62812345", false, ""},
		{"0812345678901234", false, ""},
		{"+62812345678901234", false, ""},
	}

	for _, tc := range testCases {
		tc := tc

		var validLabel string

		if tc.valid {
			validLabel = "valid"
		} else {
			validLabel = "invalid"
		}

		t.Run(fmt.Sprintf("'%s' is %v", tc.input, validLabel), func(t *testing.T) {
			t.Parallel()

			got, err := domain.NewPhoneNumber(tc.input)

			if tc.valid {
				require.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tc.expectedVal, got.Value())
			} else {
				require.Error(t, err)
			}
		})
	}
}
