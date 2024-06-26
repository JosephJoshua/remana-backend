//go:build unit
// +build unit

package repairorder_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/appconstant"
	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/modules/permission"
	"github.com/JosephJoshua/remana-backend/internal/modules/repairorder"
	"github.com/JosephJoshua/remana-backend/internal/modules/repairorder/domain"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRepairOrder(t *testing.T) {
	t.Parallel()

	var (
		theRoleID          = uuid.New()
		theStoreID         = uuid.New()
		theTechnicianID    = uuid.New()
		theSalesPersonID   = uuid.New()
		thePaymentMethodID = uuid.New()

		theDamages = []damage{
			{id: uuid.New(), name: "Screen"},
		}

		thePhoneConditions = []phoneCondition{
			{id: uuid.New(), name: "Screen broken"},
		}

		thePhoneEquipments = []phoneEquipment{
			{id: uuid.New(), name: "Battery"},
		}
	)

	logger.Init(zerolog.ErrorLevel, appconstant.AppEnvDev)

	requestCtx := appcontext.NewContextWithUser(
		testutil.RequestContextWithLogger(context.Background()),
		testutil.ModifiedUserDetails(func(details *readmodel.UserDetails) {
			details.Store.ID = theStoreID
			details.Role.ID = theRoleID
		}),
	)

	qualifyingPermissionProvider := testutil.NewPermissionProviderStub(theRoleID, []permission.Permission{
		permission.CreateRepairOrder(),
	}, nil)

	baseRepo := func() *repositoryStub {
		return &repositoryStub{
			damages:         theDamages,
			phoneConditions: thePhoneConditions,
			phoneEquipments: thePhoneEquipments,
			technicianID:    theTechnicianID,
			salesPersonID:   theSalesPersonID,
			paymentMethodID: thePaymentMethodID,
		}
	}

	validRequest := func() genapi.CreateRepairOrderRequest {
		return genapi.CreateRepairOrderRequest{
			CustomerName:       "John Doe",
			ContactPhoneNumber: "08123456789",
			PhoneType:          "iPhone 12",
			Color:              "Black",
			SalesPersonID:      theSalesPersonID,
			TechnicianID:       theTechnicianID,
			InitialCost:        100,
			DamageTypes:        []uuid.UUID{theDamages[0].id},
			PhoneConditions:    []uuid.UUID{},
			PhoneEquipments:    []uuid.UUID{},
			Photos:             []url.URL{{Host: "example.com", Scheme: "http"}},
			Imei:               genapi.NewOptString("123456789012345"),
			PartsNotCheckedYet: genapi.NewOptString("Battery"),
			Passcode: genapi.NewOptCreateRepairOrderRequestPasscode(genapi.CreateRepairOrderRequestPasscode{
				Value:           "1234",
				IsPatternLocked: true,
			}),
			DownPayment: genapi.NewOptCreateRepairOrderRequestDownPayment(genapi.CreateRepairOrderRequestDownPayment{
				Amount: 100,
				Method: thePaymentMethodID,
			}),
		}
	}

	t.Run("tries to create repair order when request body is valid", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name string
			req  *genapi.CreateRepairOrderRequest
		}{
			{
				name: "only required fields",
				req: &genapi.CreateRepairOrderRequest{
					CustomerName:       "John Doe",
					ContactPhoneNumber: "08123456789",
					PhoneType:          "iPhone 12",
					Color:              "Black",
					SalesPersonID:      theSalesPersonID,
					TechnicianID:       theTechnicianID,
					InitialCost:        100,
					DamageTypes:        []uuid.UUID{theDamages[0].id},
					PhoneConditions:    []uuid.UUID{thePhoneConditions[0].id},
					PhoneEquipments:    []uuid.UUID{thePhoneEquipments[0].id},
					Photos:             []url.URL{{Host: "example.com", Scheme: "http"}},
				},
			},
			{
				name: "with imei",
				req: &genapi.CreateRepairOrderRequest{
					CustomerName:       "John Doe",
					ContactPhoneNumber: "08123456789",
					PhoneType:          "iPhone 12",
					Color:              "Black",
					SalesPersonID:      theSalesPersonID,
					TechnicianID:       theTechnicianID,
					InitialCost:        100,
					DamageTypes:        []uuid.UUID{theDamages[0].id},
					PhoneConditions:    []uuid.UUID{thePhoneConditions[0].id},
					PhoneEquipments:    []uuid.UUID{thePhoneEquipments[0].id},
					Photos:             []url.URL{{Host: "example.com", Scheme: "http"}},
					Imei:               genapi.NewOptString("123456789012345"),
				},
			},
			{
				name: "with parts not checked yet",
				req: &genapi.CreateRepairOrderRequest{
					CustomerName:       "John Doe",
					ContactPhoneNumber: "08123456789",
					PhoneType:          "iPhone 12",
					Color:              "Black",
					SalesPersonID:      theSalesPersonID,
					TechnicianID:       theTechnicianID,
					InitialCost:        100,
					DamageTypes:        []uuid.UUID{theDamages[0].id},
					PhoneConditions:    []uuid.UUID{thePhoneConditions[0].id},
					PhoneEquipments:    []uuid.UUID{thePhoneEquipments[0].id},
					Photos:             []url.URL{{Host: "example.com", Scheme: "http"}},
					PartsNotCheckedYet: genapi.NewOptString("Battery"),
				},
			},
			{
				name: "with passcode",
				req: &genapi.CreateRepairOrderRequest{
					CustomerName:       "John Doe",
					ContactPhoneNumber: "08123456789",
					PhoneType:          "iPhone 12",
					Color:              "Black",
					SalesPersonID:      theSalesPersonID,
					TechnicianID:       theTechnicianID,
					InitialCost:        100,
					DamageTypes:        []uuid.UUID{theDamages[0].id},
					PhoneConditions:    []uuid.UUID{thePhoneConditions[0].id},
					PhoneEquipments:    []uuid.UUID{thePhoneEquipments[0].id},
					Photos:             []url.URL{{Host: "example.com", Scheme: "http"}},
					Passcode: genapi.NewOptCreateRepairOrderRequestPasscode(genapi.CreateRepairOrderRequestPasscode{
						Value:           "1234",
						IsPatternLocked: false,
					}),
				},
			},
			{
				name: "with pattern lock",
				req: &genapi.CreateRepairOrderRequest{
					CustomerName:       "John Doe",
					ContactPhoneNumber: "08123456789",
					PhoneType:          "iPhone 12",
					Color:              "Black",
					SalesPersonID:      theSalesPersonID,
					TechnicianID:       theTechnicianID,
					InitialCost:        100,
					DamageTypes:        []uuid.UUID{theDamages[0].id},
					PhoneConditions:    []uuid.UUID{thePhoneConditions[0].id},
					PhoneEquipments:    []uuid.UUID{thePhoneEquipments[0].id},
					Photos:             []url.URL{{Host: "example.com", Scheme: "http"}},
					Passcode: genapi.NewOptCreateRepairOrderRequestPasscode(genapi.CreateRepairOrderRequestPasscode{
						Value:           "1234",
						IsPatternLocked: true,
					}),
				},
			},
			{
				name: "with down payment",
				req: &genapi.CreateRepairOrderRequest{
					CustomerName:       "John Doe",
					ContactPhoneNumber: "08123456789",
					PhoneType:          "iPhone 12",
					Color:              "Black",
					SalesPersonID:      theSalesPersonID,
					TechnicianID:       theTechnicianID,
					InitialCost:        100,
					DamageTypes:        []uuid.UUID{theDamages[0].id},
					PhoneConditions:    []uuid.UUID{thePhoneConditions[0].id},
					PhoneEquipments:    []uuid.UUID{thePhoneEquipments[0].id},
					Photos:             []url.URL{{Host: "example.com", Scheme: "http"}},
					DownPayment: genapi.NewOptCreateRepairOrderRequestDownPayment(genapi.CreateRepairOrderRequestDownPayment{
						Amount: 100,
						Method: thePaymentMethodID,
					}),
				},
			},
		}

		for _, tc := range testCases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				repo := baseRepo()
				s := repairorder.NewService(
					testutil.NewTimeProviderStub(time.Now()),
					testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
					repo,
					qualifyingPermissionProvider,
					testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
				)

				_, err := s.CreateRepairOrder(requestCtx, tc.req)

				require.NoError(t, err)
				require.NotNil(t, repo.calledWithOrder)

				assert.Equal(t, theStoreID, repo.calledWithOrder.StoreID())
				assert.Equal(t, tc.req.CustomerName, repo.calledWithOrder.CustomerName())
				assert.NotEmpty(t, repo.calledWithOrder.ContactNumber().Value())
				assert.Equal(t, tc.req.PhoneType, repo.calledWithOrder.PhoneType())
				assert.Equal(t, tc.req.Color, repo.calledWithOrder.Color())
				assert.Equal(t, tc.req.SalesPersonID, repo.calledWithOrder.SalesPersonID())
				assert.Equal(t, tc.req.TechnicianID, repo.calledWithOrder.TechnicianID())
				assert.Equal(t, 1, len(repo.calledWithOrder.Costs()))
				assert.Equal(t, tc.req.InitialCost, repo.calledWithOrder.Costs()[0].Amount())
				assert.Equal(t, true, repo.calledWithOrder.Costs()[0].IsInitial())
				assert.Equal(t, len(tc.req.DamageTypes), len(repo.calledWithOrder.Damages()))
				assert.Equal(t, len(tc.req.PhoneConditions), len(repo.calledWithOrder.PhoneConditions()))
				assert.Equal(t, len(tc.req.PhoneEquipments), len(repo.calledWithOrder.PhoneEquipments()))
				assert.Equal(t, len(tc.req.Photos), len(repo.calledWithOrder.Photos()))

				gotPhotoURLs := []url.URL{}
				for _, photo := range repo.calledWithOrder.Photos() {
					gotPhotoURLs = append(gotPhotoURLs, photo.URL())
				}

				assert.Equal(t, tc.req.Photos, gotPhotoURLs)

				if tc.req.Imei.IsSet() {
					require.True(t, repo.calledWithOrder.IMEI().PointerValue().IsSet())
					assert.Equal(t, tc.req.Imei.Value, repo.calledWithOrder.IMEI().PointerValue().MustGet())
				}

				if tc.req.PartsNotCheckedYet.IsSet() {
					require.True(t, repo.calledWithOrder.PartsNotCheckedYet().PointerValue().IsSet())
					assert.Equal(t, tc.req.PartsNotCheckedYet.Value, repo.calledWithOrder.PartsNotCheckedYet().PointerValue().MustGet())
				}

				if tc.req.Passcode.IsSet() {
					require.True(t, repo.calledWithOrder.PhoneSecurityDetails().PointerValue().IsSet())

					gotIsPatternLocked := repo.calledWithOrder.PhoneSecurityDetails().PointerValue().MustGet().Type() == domain.PhoneSecurityTypePattern

					assert.Equal(t, tc.req.Passcode.Value.Value, repo.calledWithOrder.PhoneSecurityDetails().PointerValue().MustGet().Value())
					assert.Equal(t, tc.req.Passcode.Value.IsPatternLocked, gotIsPatternLocked)
				}

				if tc.req.DownPayment.IsSet() {
					require.True(t, repo.calledWithOrder.DownPayment().PointerValue().IsSet())

					assert.Equal(t, uint(tc.req.DownPayment.Value.Amount), repo.calledWithOrder.DownPayment().PointerValue().MustGet().Amount())
					assert.Equal(t, tc.req.DownPayment.Value.Method, repo.calledWithOrder.DownPayment().PointerValue().MustGet().PaymentMethodID())
				}
			})
		}
	})

	t.Run("creates order with random slug", func(t *testing.T) {
		t.Parallel()

		const theSlug = "random-slug"

		slugProvider := testutil.NewRepairOrderSlugProviderStub(theSlug, nil)
		repo := baseRepo()

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			slugProvider,
		)

		req := validRequest()
		_, err := s.CreateRepairOrder(requestCtx, &req)

		require.NoError(t, err)
		require.NotNil(t, repo.calledWithOrder)

		assert.Equal(t, theSlug, repo.calledWithOrder.Slug())
	})

	t.Run("creates order with current time", func(t *testing.T) {
		t.Parallel()

		now := time.Unix(1713917762, 0)
		repo := baseRepo()

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(now),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		_, err := s.CreateRepairOrder(requestCtx, &req)

		require.NoError(t, err)
		require.NotNil(t, repo.calledWithOrder)

		assert.Equal(t, now, repo.calledWithOrder.CreationTime())
	})

	t.Run("creates order with damage names", func(t *testing.T) {
		t.Parallel()

		damages := []damage{
			{id: uuid.New(), name: "Screen"},
			{id: uuid.New(), name: "Battery"},
			{id: uuid.New(), name: "Camera"},
		}

		repo := baseRepo()
		repo.damages = damages

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		req.DamageTypes = []uuid.UUID{damages[0].id, damages[1].id}

		_, err := s.CreateRepairOrder(requestCtx, &req)

		require.NoError(t, err)
		require.NotNil(t, repo.calledWithOrder)

		assert.Equal(t, len(req.DamageTypes), len(repo.calledWithOrder.Damages()))

		var containsFirstDamage, containsSecondDamage bool

		for _, damage := range repo.calledWithOrder.Damages() {
			if damage.Name() == damages[0].name {
				containsFirstDamage = true
			} else if damage.Name() == damages[1].name {
				containsSecondDamage = true
			}
		}

		assert.True(t, containsFirstDamage && containsSecondDamage, "expected order to contain both damages '%s' and '%s', got %#v", damages[0].name, damages[1].name, repo.calledWithOrder.Damages())
	})

	t.Run("returns bad request when damage type is not found", func(t *testing.T) {
		t.Parallel()

		damages := []damage{
			{id: uuid.New(), name: "Screen"},
			{id: uuid.New(), name: "Battery"},
			{id: uuid.New(), name: "Camera"},
		}

		randomID := uuid.New()

		repo := baseRepo()
		repo.damages = damages

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		req.DamageTypes = []uuid.UUID{damages[0].id, randomID}

		_, err := s.CreateRepairOrder(requestCtx, &req)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("creates order with phone condition names", func(t *testing.T) {
		t.Parallel()

		phoneConditions := []phoneCondition{
			{id: uuid.New(), name: "Screen broken"},
			{id: uuid.New(), name: "Camera cracked"},
		}

		repo := baseRepo()
		repo.phoneConditions = phoneConditions

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		req.PhoneConditions = []uuid.UUID{phoneConditions[0].id, phoneConditions[1].id}

		_, err := s.CreateRepairOrder(requestCtx, &req)

		require.NoError(t, err)
		require.NotNil(t, repo.calledWithOrder)

		assert.Equal(t, len(req.PhoneConditions), len(repo.calledWithOrder.PhoneConditions()))

		var containsFirstCondition, containsSecondCondition bool

		for _, condition := range repo.calledWithOrder.PhoneConditions() {
			if condition.Name() == phoneConditions[0].name {
				containsFirstCondition = true
			} else if condition.Name() == phoneConditions[1].name {
				containsSecondCondition = true
			}
		}

		assert.True(t, containsFirstCondition && containsSecondCondition, "expected order to contain both phone conditions '%s' and '%s', got %#v", phoneConditions[0].name, phoneConditions[1].name, repo.calledWithOrder.PhoneConditions())
	})

	t.Run("returns bad request when phone condition is not found", func(t *testing.T) {
		t.Parallel()

		phoneConditions := []phoneCondition{
			{id: uuid.New(), name: "Screen broken"},
			{id: uuid.New(), name: "Camera cracked"},
		}

		randomID := uuid.New()

		repo := baseRepo()
		repo.phoneConditions = phoneConditions

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		req.PhoneConditions = []uuid.UUID{phoneConditions[0].id, randomID}

		_, err := s.CreateRepairOrder(requestCtx, &req)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("creates order with phone equipment names", func(t *testing.T) {
		t.Parallel()

		equipments := []phoneEquipment{
			{id: uuid.New(), name: "Battery"},
			{id: uuid.New(), name: "SIM Card"},
		}

		repo := baseRepo()
		repo.phoneEquipments = equipments

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		req.PhoneEquipments = []uuid.UUID{equipments[0].id, equipments[1].id}

		_, err := s.CreateRepairOrder(requestCtx, &req)

		require.NoError(t, err)
		require.NotNil(t, repo.calledWithOrder)

		assert.Equal(t, len(req.PhoneEquipments), len(repo.calledWithOrder.PhoneEquipments()))

		var containsFirstEquipment, containsSecondEquipment bool

		for _, equipment := range repo.calledWithOrder.PhoneEquipments() {
			if equipment.Name() == equipments[0].name {
				containsFirstEquipment = true
			} else if equipment.Name() == equipments[1].name {
				containsSecondEquipment = true
			}
		}

		assert.True(t, containsFirstEquipment && containsSecondEquipment, "expected order to contain both equipments '%s' and '%s', got %#v", equipments[0].name, equipments[1].name, repo.calledWithOrder.PhoneEquipments())
	})

	t.Run("returns bad request when phone equipment is not found", func(t *testing.T) {
		t.Parallel()

		equipments := []phoneEquipment{
			{id: uuid.New(), name: "Battery"},
			{id: uuid.New(), name: "SIM Card"},
		}

		randomID := uuid.New()

		repo := baseRepo()
		repo.phoneEquipments = equipments

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		req.PhoneEquipments = []uuid.UUID{equipments[0].id, randomID}

		_, err := s.CreateRepairOrder(requestCtx, &req)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns resource location when repair order is created", func(t *testing.T) {
		t.Parallel()

		theLocation := url.URL{Scheme: "http", Host: "example.com", Path: "/repair-orders"}

		locationProvider := testutil.NewResourceLocationProviderStubForRepairOrder(theLocation)
		repo := baseRepo()

		s := repairorder.NewService(testutil.NewTimeProviderStub(time.Now()),
			locationProvider,
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		got, err := s.CreateRepairOrder(requestCtx, &req)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, theLocation, got.Location)

		require.True(t, locationProvider.RepairOrderID.IsSet())
		require.NotNil(t, repo.calledWithOrder)

		assert.Equal(t, repo.calledWithOrder.ID(), locationProvider.RepairOrderID.MustGet())
	})

	t.Run("returns bad request when technician does not exist", func(t *testing.T) {
		t.Parallel()

		technicianID := uuid.New()
		randomID := uuid.New()

		repo := baseRepo()
		repo.technicianID = technicianID

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		req.TechnicianID = randomID

		_, err := s.CreateRepairOrder(requestCtx, &req)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns bad request when sales person does not exist", func(t *testing.T) {
		t.Parallel()

		salesPersonID := uuid.New()
		randomID := uuid.New()

		repo := baseRepo()
		repo.salesPersonID = salesPersonID

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		req.SalesPersonID = randomID

		_, err := s.CreateRepairOrder(requestCtx, &req)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns bad request when payment method does not exist", func(t *testing.T) {
		t.Parallel()

		thePaymentMethodID := uuid.New()
		someRandomID := uuid.New()

		repo := baseRepo()
		repo.paymentMethodID = thePaymentMethodID

		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		req.DownPayment = genapi.NewOptCreateRepairOrderRequestDownPayment(genapi.CreateRepairOrderRequestDownPayment{
			Amount: 100,
			Method: someRandomID,
		})

		_, err := s.CreateRepairOrder(requestCtx, &req)
		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns bad request when contact phone number is invalid", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name          string
			contactNumber string
		}{
			{
				name:          "empty",
				contactNumber: "",
			},
			{
				name:          "too short",
				contactNumber: "123",
			},
			{
				name:          "too long",
				contactNumber: "12345678901234567890",
			},
		}

		for _, tc := range testCases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				repo := baseRepo()
				s := repairorder.NewService(
					testutil.NewTimeProviderStub(time.Now()),
					testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
					repo,
					qualifyingPermissionProvider,
					testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
				)

				req := validRequest()
				req.ContactPhoneNumber = tc.contactNumber

				_, err := s.CreateRepairOrder(requestCtx, &req)
				testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
			})
		}
	})

	t.Run("returns unauthorized when user is missing from context", func(t *testing.T) {
		t.Parallel()

		repo := baseRepo()
		s := repairorder.NewService(
			testutil.NewTimeProviderStub(time.Now()),
			testutil.NewResourceLocationProviderStubForRepairOrder(url.URL{}),
			repo,
			qualifyingPermissionProvider,
			testutil.NewRepairOrderSlugProviderStub("random-slug", nil),
		)

		req := validRequest()
		emptyCtx := testutil.RequestContextWithLogger(context.Background())

		_, err := s.CreateRepairOrder(emptyCtx, &req)

		testutil.AssertAPIStatusCode(t, http.StatusUnauthorized, err)
	})

	t.Run("returns internal server error", func(t *testing.T) {
		testCases := []struct {
			name  string
			setup func(
				repo *repositoryStub,
				locationProvider *testutil.ResourceLocationProviderStub,
				permissionProvider *testutil.PermissionProviderStub,
				slugProvider *testutil.OrderSlugProviderStub,
			)
		}{
			{
				name: "when repository.DoesTechnicianExist() errors",
				setup: func(repo *repositoryStub,
					_ *testutil.ResourceLocationProviderStub,
					_ *testutil.PermissionProviderStub,
					_ *testutil.OrderSlugProviderStub) {
					repo.technicianExistsErr = errors.New("oh no!")
				},
			},
			{
				name: "when repository.DoesSalesPersonIDExist() errors",
				setup: func(repo *repositoryStub,
					_ *testutil.ResourceLocationProviderStub,
					_ *testutil.PermissionProviderStub,
					_ *testutil.OrderSlugProviderStub) {
					repo.salesPersonExistsErr = errors.New("oh no!")
				},
			},
			{
				name: "when repository.DoesPaymentMethodExist() errors",
				setup: func(repo *repositoryStub,
					_ *testutil.ResourceLocationProviderStub,
					_ *testutil.PermissionProviderStub,
					_ *testutil.OrderSlugProviderStub) {
					repo.paymentMethodExistsErr = errors.New("oh no!")
				},
			},
			{
				name: "when repository.CreateRepairOrder() errors",
				setup: func(repo *repositoryStub,
					_ *testutil.ResourceLocationProviderStub,
					_ *testutil.PermissionProviderStub,
					_ *testutil.OrderSlugProviderStub) {
					repo.createErr = errors.New("oh no!")
				},
			},
			{
				name: "when repository.GetDamageNamesByID() errors",
				setup: func(repo *repositoryStub,
					_ *testutil.ResourceLocationProviderStub,
					_ *testutil.PermissionProviderStub,
					_ *testutil.OrderSlugProviderStub) {
					repo.damageNameErr = errors.New("oh no!")
				},
			},
			{
				name: "when repository.GetPhoneConditionNamesByID() errors",
				setup: func(repo *repositoryStub,
					_ *testutil.ResourceLocationProviderStub,
					_ *testutil.PermissionProviderStub,
					_ *testutil.OrderSlugProviderStub) {
					repo.phoneConditionNameErr = errors.New("oh no!")
				},
			},
			{
				name: "when repository.GetPhoneEquipmentNamesByID() errors",
				setup: func(repo *repositoryStub,
					_ *testutil.ResourceLocationProviderStub,
					_ *testutil.PermissionProviderStub,
					_ *testutil.OrderSlugProviderStub) {
					repo.phoneEquipmentNameErr = errors.New("oh no!")
				},
			},
			{
				name: "when repository.GetPhoneEquipmentNamesByID() errors",
				setup: func(repo *repositoryStub,
					_ *testutil.ResourceLocationProviderStub,
					_ *testutil.PermissionProviderStub,
					_ *testutil.OrderSlugProviderStub) {
					repo.phoneEquipmentNameErr = errors.New("oh no!")
				},
			},
			{
				name: "when order slug provider errors",
				setup: func(_ *repositoryStub,
					_ *testutil.ResourceLocationProviderStub,
					_ *testutil.PermissionProviderStub,
					slugProvider *testutil.OrderSlugProviderStub) {
					slugProvider.SetError(errors.New("oh no!"))
				},
			},
			{
				name: "when permissionProvider.Can() errors",
				setup: func(repo *repositoryStub,
					locationProvider *testutil.ResourceLocationProviderStub,
					permissionProvider *testutil.PermissionProviderStub,
					slugProvider *testutil.OrderSlugProviderStub) {
					permissionProvider.SetError(errors.New("oh no!"))
				},
			},
		}

		for _, tc := range testCases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				repo := baseRepo()
				locationProvider := testutil.NewResourceLocationProviderStubForRepairOrder(
					url.URL{Scheme: "http", Host: "example.com", Path: "/repair-orders"},
				)
				slugProvider := testutil.NewRepairOrderSlugProviderStub("random-slug", nil)
				permissionProvider := testutil.NewPermissionProviderStub(theRoleID, []permission.Permission{
					permission.CreateRepairOrder(),
				}, nil)

				tc.setup(repo, locationProvider, permissionProvider, slugProvider)

				s := repairorder.NewService(
					testutil.NewTimeProviderStub(time.Now()),
					locationProvider,
					repo,
					permissionProvider,
					slugProvider,
				)

				req := validRequest()
				_, err := s.CreateRepairOrder(requestCtx, &req)

				testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
			})
		}
	})
}

type damage struct {
	id   uuid.UUID
	name string
}

type phoneCondition struct {
	id   uuid.UUID
	name string
}

type phoneEquipment struct {
	id   uuid.UUID
	name string
}

type repositoryStub struct {
	damages                []damage
	phoneConditions        []phoneCondition
	phoneEquipments        []phoneEquipment
	technicianID           uuid.UUID
	salesPersonID          uuid.UUID
	paymentMethodID        uuid.UUID
	calledWithOrder        domain.Order
	createErr              error
	damageNameErr          error
	phoneConditionNameErr  error
	phoneEquipmentNameErr  error
	storeExistsErr         error
	technicianExistsErr    error
	salesPersonExistsErr   error
	paymentMethodExistsErr error
}

func (r *repositoryStub) CreateRepairOrder(_ context.Context, order domain.Order) error {
	if r.createErr != nil {
		return r.createErr
	}

	r.calledWithOrder = order
	return nil
}

func (r *repositoryStub) GetDamageNamesByIDs(_ context.Context, storeID uuid.UUID, ids []uuid.UUID) ([]string, error) {
	if r.damageNameErr != nil {
		return []string{}, r.damageNameErr
	}

	names := []string{}

	for _, id := range ids {
		found := false

		for _, damage := range r.damages {
			if damage.id == id {
				names = append(names, damage.name)
				found = true

				break
			}
		}

		if !found {
			return []string{}, apperror.ErrDamageNotFound
		}
	}

	return names, nil
}

func (r *repositoryStub) GetPhoneConditionNamesByIDs(_ context.Context, storeID uuid.UUID, ids []uuid.UUID) ([]string, error) {
	if r.phoneConditionNameErr != nil {
		return []string{}, r.phoneConditionNameErr
	}

	names := []string{}

	for _, id := range ids {
		found := false

		for _, phoneCondition := range r.phoneConditions {
			if phoneCondition.id == id {
				names = append(names, phoneCondition.name)
				found = true

				break
			}
		}

		if !found {
			return []string{}, apperror.ErrPhoneConditionNotFound
		}
	}

	return names, nil
}

func (r *repositoryStub) GetPhoneEquipmentNamesByIDs(_ context.Context, storeID uuid.UUID, ids []uuid.UUID) ([]string, error) {
	if r.phoneEquipmentNameErr != nil {
		return []string{}, r.phoneEquipmentNameErr
	}

	names := []string{}

	for _, id := range ids {
		found := false

		for _, equipment := range r.phoneEquipments {
			if equipment.id == id {
				names = append(names, equipment.name)
				found = true

				break
			}
		}

		if !found {
			return []string{}, apperror.ErrPhoneEquipmentNotFound
		}
	}

	return names, nil
}

func (r *repositoryStub) DoesTechnicianExist(_ context.Context, storeID uuid.UUID, technicianID uuid.UUID) (bool, error) {
	if r.technicianExistsErr != nil {
		return false, r.technicianExistsErr
	}

	return technicianID == r.technicianID, nil
}

func (r *repositoryStub) DoesSalesPersonExist(_ context.Context, storeID uuid.UUID, salesPersonID uuid.UUID) (bool, error) {
	if r.salesPersonExistsErr != nil {
		return false, r.salesPersonExistsErr
	}

	return salesPersonID == r.salesPersonID, nil
}

func (r *repositoryStub) DoesPaymentMethodExist(_ context.Context, storeID uuid.UUID, paymentMethodID uuid.UUID) (bool, error) {
	if r.paymentMethodExistsErr != nil {
		return false, r.paymentMethodExistsErr
	}

	return paymentMethodID == r.paymentMethodID, nil
}
