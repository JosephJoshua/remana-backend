//go:build unit
// +build unit

package domain_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/repairorder/domain"
	shareddomain "github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrder(t *testing.T) {
	dummyURL := url.URL{}
	dummyTime := time.Now()
	dummyID := uuid.New()
	dummyContactNumber, initErr := shareddomain.NewPhoneNumber("081234567890")

	require.NoError(t, initErr)

	var invalidDownPaymentAmount uint
	var invalidIMEI string
	var invalidPartsNotCheckedYet string

	testCases := []struct {
		testName           string
		slug               string
		customerName       string
		phoneType          string
		color              string
		initialCost        uint
		damages            []string
		photos             []url.URL
		downPaymentAmount  *uint
		imei               *string
		partsNotCheckedYet *string
		valid              bool
	}{
		{
			"empty slug",
			"",
			"customer 1",
			"phone type",
			"white",
			100,
			[]string{"damage 1"},
			[]url.URL{dummyURL},
			nil,
			nil,
			nil,
			false,
		},
		{
			"empty customerName",
			"slug",
			"",
			"phone type",
			"white",
			100,
			[]string{"damage 1"},
			[]url.URL{dummyURL},
			nil,
			nil,
			nil,
			false,
		},
		{
			"empty phoneType",
			"slug",
			"customer 1",
			"",
			"white",
			100,
			[]string{"damage 1"},
			[]url.URL{dummyURL},
			nil,
			nil,
			nil,
			false,
		},
		{
			"empty color",
			"slug",
			"customer 1",
			"phone type",
			"",
			100,
			[]string{"damage 1"},
			[]url.URL{dummyURL},
			nil,
			nil,
			nil,
			false,
		},
		{
			"empty initialCost",
			"slug",
			"customer 1",
			"phone type",
			"white",
			0,
			[]string{"damage 1"},
			[]url.URL{dummyURL},
			nil,
			nil,
			nil,
			false,
		},
		{
			"empty damages",
			"slug",
			"customer 1",
			"phone type",
			"white",
			100,
			[]string{},
			[]url.URL{dummyURL},
			nil,
			nil,
			nil,
			false,
		},
		{
			"empty photos",
			"slug",
			"customer 1",
			"phone type",
			"white",
			100,
			[]string{"damage 1"},
			[]url.URL{},
			nil,
			nil,
			nil,
			false,
		},
		{
			"empty downPaymentAmount",
			"slug",
			"customer 1",
			"phone type",
			"white",
			100,
			[]string{"damage 1"},
			[]url.URL{dummyURL},
			&invalidDownPaymentAmount,
			nil,
			nil,
			false,
		},
		{
			"empty imei",
			"slug",
			"customer 1",
			"phone type",
			"white",
			100,
			[]string{"damage 1"},
			[]url.URL{dummyURL},
			nil,
			&invalidIMEI,
			nil,
			false,
		},
		{
			"empty partsNotCheckedYet",
			"slug",
			"customer 1",
			"phone type",
			"white",
			100,
			[]string{"damage 1"},
			[]url.URL{dummyURL},
			nil,
			nil,
			&invalidPartsNotCheckedYet,
			false,
		},
		{
			"valid",
			"slug",
			"customer 1",
			"phone type",
			"white",
			100,
			[]string{"damage 1"},
			[]url.URL{dummyURL},
			nil,
			nil,
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			opts := []domain.OrderOption{}

			if tc.downPaymentAmount != nil {
				opts = append(opts, domain.WithDownPayment(*tc.downPaymentAmount, dummyID))
			}

			if tc.imei != nil {
				opts = append(opts, domain.WithIMEI(*tc.imei))
			}

			if tc.partsNotCheckedYet != nil {
				opts = append(opts, domain.WithPartsNotCheckedYet(*tc.partsNotCheckedYet))
			}

			got, err := domain.NewOrder(
				dummyTime,
				tc.slug,
				dummyID,
				tc.customerName,
				dummyContactNumber,
				tc.phoneType,
				tc.color,
				tc.initialCost,
				[]string{},
				[]string{},
				tc.damages,
				tc.photos,
				dummyID,
				dummyID,
				opts...,
			)

			if tc.valid {
				require.NoError(t, err)
				assert.NotNil(t, got)
			} else {
				require.Error(t, err)
			}
		})
	}
}
