//go:build unit
// +build unit

package domain_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/modules/repairorder/domain"
	shareddomain "github.com/JosephJoshua/remana-backend/internal/modules/shared/domain"
	"github.com/JosephJoshua/remana-backend/internal/optional"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrder(t *testing.T) {
	t.Run("returns new order", func(t *testing.T) {
		theContactNumber, initErr := shareddomain.NewPhoneNumber("081234567890")
		require.NoError(t, initErr)

		params := domain.NewOrderParams{
			CreationTime:    time.Now(),
			Slug:            "slug",
			StoreID:         uuid.New(),
			CustomerName:    "John Doe",
			ContactNumber:   theContactNumber,
			PhoneType:       "Advan G5",
			Color:           "White",
			InitialCost:     100,
			PhoneConditions: []string{"condition 1"},
			PhoneEquipments: []string{"equipment 1"},
			Damages:         []string{"damage 1"},
			Photos:          []url.URL{{Host: "example.com"}},
			SalesPersonID:   uuid.New(),
			TechnicianID:    uuid.New(),
		}

		got, err := domain.NewOrder(params)
		require.NoError(t, err)

		assert.Equal(t, params.CreationTime, got.CreationTime())
		assert.Equal(t, params.Slug, got.Slug())
		assert.Equal(t, params.StoreID, got.StoreID())
		assert.Equal(t, params.CustomerName, got.CustomerName())
		assert.Equal(t, params.ContactNumber.Value(), got.ContactNumber().Value())
		assert.Equal(t, params.PhoneType, got.PhoneType())
		assert.Equal(t, params.Color, got.Color())
		assert.Equal(t, params.SalesPersonID, got.SalesPersonID())
		assert.Equal(t, params.TechnicianID, got.TechnicianID())

		require.NotEmpty(t, got.Costs())
		assert.Equal(t, int(params.InitialCost), got.Costs()[0].Amount())

		require.NotEmpty(t, got.PhoneConditions())
		assert.Equal(t, params.PhoneConditions[0], got.PhoneConditions()[0].Name())

		require.NotEmpty(t, got.Damages())
		assert.Equal(t, params.Damages[0], got.Damages()[0].Name())

		require.NotEmpty(t, got.PhoneEquipments())
		assert.Equal(t, params.PhoneEquipments[0], got.PhoneEquipments()[0].Name())

		require.NotEmpty(t, got.Photos())
		assert.Equal(t, params.Photos[0], got.Photos()[0].URL())
	})

	t.Run("returns invalid input error", func(t *testing.T) {
		dummyURL := url.URL{Host: "example.com"}
		dummyTime := time.Now()
		dummyID := uuid.New()

		dummyContactNumber, initErr := shareddomain.NewPhoneNumber("081234567890")
		require.NoError(t, initErr)

		testCases := []struct {
			name  string
			setup func(params *domain.NewOrderParams)
		}{
			{
				name: "empty slug",
				setup: func(params *domain.NewOrderParams) {
					params.Slug = ""
				},
			},
			{
				name: "empty customer name",
				setup: func(params *domain.NewOrderParams) {
					params.CustomerName = ""
				},
			},
			{
				name: "empty phone type",
				setup: func(params *domain.NewOrderParams) {
					params.PhoneType = ""
				},
			},
			{
				name: "empty color",
				setup: func(params *domain.NewOrderParams) {
					params.Color = ""
				},
			},
			{
				name: "zero initial cost",
				setup: func(params *domain.NewOrderParams) {
					params.InitialCost = 0
				},
			},
			{
				name: "down payment greater than initial cost",
				setup: func(params *domain.NewOrderParams) {
					payment, err := domain.NewOrderPayment(500, dummyID)
					require.NoError(t, err)

					params.InitialCost = 100
					params.DownPayment = optional.Some(payment)
				},
			},
			{
				name: "empty damages",
				setup: func(params *domain.NewOrderParams) {
					params.Damages = []string{}
				},
			},
			{
				name: "damages set, but is empty string",
				setup: func(params *domain.NewOrderParams) {
					params.Damages = []string{""}
				},
			},
			{
				name: "phone equipments set, but is empty string",
				setup: func(params *domain.NewOrderParams) {
					params.PhoneEquipments = []string{""}
				},
			},
			{
				name: "phone conditions set, but is empty string",
				setup: func(params *domain.NewOrderParams) {
					params.PhoneConditions = []string{""}
				},
			},
			{
				name: "empty photos",
				setup: func(params *domain.NewOrderParams) {
					params.Photos = []url.URL{}
				},
			},
			{
				name: "imei set, but empty",
				setup: func(params *domain.NewOrderParams) {
					params.Imei = optional.Some("")
				},
			},
			{
				name: "parts not checked yet set, but empty",
				setup: func(params *domain.NewOrderParams) {
					params.PartsNotCheckedYet = optional.Some("")
				},
			},
		}

		for _, tc := range testCases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				params := domain.NewOrderParams{
					CreationTime:  dummyTime,
					Slug:          "slug",
					StoreID:       dummyID,
					CustomerName:  "John Doe",
					ContactNumber: dummyContactNumber,
					PhoneType:     "Advan G5",
					Color:         "White",
					InitialCost:   100,
					Damages:       []string{"damage 1"},
					Photos:        []url.URL{dummyURL},
					SalesPersonID: dummyID,
					TechnicianID:  dummyID,
				}

				tc.setup(&params)

				_, err := domain.NewOrder(params)
				require.ErrorIs(t, err, apperror.ErrInvalidInput)
			})
		}
	})

}
