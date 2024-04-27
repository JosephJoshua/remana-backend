//go:build unit
// +build unit

package domain_test

import (
	"fmt"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/repairorder/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPatternSecurity(t *testing.T) {
	testCases := []struct {
		input string
		valid bool
	}{
		{"0123456789", true},
		{"1234567890", true},
		{"12a3456789", false},
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

			got, err := domain.NewPatternSecurity(tc.input)

			if tc.valid {
				require.NoError(t, err)
				assert.NotNil(t, got)
			} else {
				require.Error(t, err)
			}
		})
	}
}
