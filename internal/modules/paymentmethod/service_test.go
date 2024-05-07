//go:build unit
// +build unit

package paymentmethod_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/JosephJoshua/remana-backend/internal/appconstant"
	"github.com/JosephJoshua/remana-backend/internal/appcontext"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/logger"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/modules/paymentmethod"
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type repositoryStub struct {
	createCalledWith struct {
		id      uuid.UUID
		storeID uuid.UUID
		name    string
	}
	storeID      uuid.UUID
	existingName string
	createErr    error
	nameTakenErr error
}

func (r *repositoryStub) CreatePaymentMethod(_ context.Context, id uuid.UUID, storeID uuid.UUID, name string) error {
	if r.createErr != nil {
		return r.createErr
	}

	r.createCalledWith.id = id
	r.createCalledWith.storeID = storeID
	r.createCalledWith.name = name

	return nil
}

func (r *repositoryStub) IsNameTaken(_ context.Context, storeID uuid.UUID, name string) (bool, error) {
	if r.nameTakenErr != nil {
		return false, r.nameTakenErr
	}

	return r.storeID == storeID && r.existingName == name, nil
}

func TestCreatePaymentMethod(t *testing.T) {
	t.Parallel()

	var (
		theStoreID = uuid.New()
	)

	logger.Init(zerolog.ErrorLevel, appconstant.AppEnvDev)
	requestCtx := appcontext.NewContextWithUser(
		testutil.RequestContextWithLogger(context.Background()),
		&readmodel.UserDetails{
			ID:       uuid.New(),
			Username: "not important",
			Role: readmodel.UserDetailsRole{
				ID:           uuid.New(),
				Name:         "not important",
				IsStoreAdmin: true,
			},
			Store: readmodel.UserDetailsStore{
				ID:   theStoreID,
				Name: "not important",
				Code: "not-important",
			},
		},
	)

	t.Run("tries to create payment method when request is valid", func(t *testing.T) {
		t.Parallel()

		repo := &repositoryStub{storeID: theStoreID}
		s := paymentmethod.NewService(
			testutil.NewResourceLocationProviderStubForPaymentMethod(url.URL{}),
			repo,
		)

		req := &genapi.CreatePaymentMethodRequest{
			Name: "payment method 1",
		}

		got, err := s.CreatePaymentMethod(requestCtx, req)

		require.NoError(t, err)
		require.NotNil(t, got)

		require.NotNil(t, repo.createCalledWith)

		assert.Equal(t, req.Name, repo.createCalledWith.name)
		assert.Equal(t, theStoreID, repo.createCalledWith.storeID)
	})

	t.Run("returns resource location when payment method is created", func(t *testing.T) {
		t.Parallel()

		var (
			theLocation = url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/payment-methods/ef21dc9e-c364-41cd-8c03-fa289d11e3a7",
			}
		)

		resourceLocationProvider := testutil.NewResourceLocationProviderStubForPaymentMethod(theLocation)
		repo := &repositoryStub{storeID: theStoreID}

		s := paymentmethod.NewService(resourceLocationProvider, repo)

		got, err := s.CreatePaymentMethod(requestCtx, &genapi.CreatePaymentMethodRequest{
			Name: "payment method 1",
		})

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, theLocation, got.Location)

		require.True(t, resourceLocationProvider.PaymentMethodID.IsSet())
		require.NotNil(t, repo.createCalledWith)

		assert.Equal(t, repo.createCalledWith.id, resourceLocationProvider.PaymentMethodID.MustGet())
	})

	t.Run("returns unauthorized when user is missing from context", func(t *testing.T) {
		t.Parallel()

		s := paymentmethod.NewService(
			testutil.NewResourceLocationProviderStubForPaymentMethod(url.URL{}),
			&repositoryStub{},
		)

		emptyCtx := testutil.RequestContextWithLogger(context.Background())
		_, err := s.CreatePaymentMethod(emptyCtx, &genapi.CreatePaymentMethodRequest{
			Name: "payment method 1",
		})

		testutil.AssertAPIStatusCode(t, http.StatusUnauthorized, err)
	})

	t.Run("returns bad request when name is empty", func(t *testing.T) {
		t.Parallel()

		s := paymentmethod.NewService(
			testutil.NewResourceLocationProviderStubForPaymentMethod(url.URL{}),
			&repositoryStub{storeID: theStoreID},
		)

		_, err := s.CreatePaymentMethod(requestCtx, &genapi.CreatePaymentMethodRequest{
			Name: "",
		})

		testutil.AssertAPIStatusCode(t, http.StatusBadRequest, err)
	})

	t.Run("returns conflict when name is taken", func(t *testing.T) {
		t.Parallel()

		const (
			theName = "payment method 1"
		)

		repo := &repositoryStub{existingName: theName, storeID: theStoreID}
		s := paymentmethod.NewService(
			testutil.NewResourceLocationProviderStubForPaymentMethod(url.URL{}),
			repo,
		)

		req := &genapi.CreatePaymentMethodRequest{
			Name: theName,
		}

		_, err := s.CreatePaymentMethod(requestCtx, req)
		testutil.AssertAPIStatusCode(t, http.StatusConflict, err)
	})

	t.Run("returns internal server error when repository.IsNameTaken() errors", func(t *testing.T) {
		t.Parallel()

		s := paymentmethod.NewService(
			testutil.NewResourceLocationProviderStubForPaymentMethod(url.URL{}),
			&repositoryStub{nameTakenErr: errors.New("oh no!"), storeID: theStoreID},
		)

		_, err := s.CreatePaymentMethod(requestCtx, &genapi.CreatePaymentMethodRequest{
			Name: "payment method 1",
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})

	t.Run("returns internal server error when repository.CreatePaymentMethod() errors", func(t *testing.T) {
		t.Parallel()

		s := paymentmethod.NewService(
			testutil.NewResourceLocationProviderStubForPaymentMethod(url.URL{}),
			&repositoryStub{createErr: errors.New("oh no!"), storeID: theStoreID},
		)

		_, err := s.CreatePaymentMethod(requestCtx, &genapi.CreatePaymentMethodRequest{
			Name: "payment method 1",
		})

		testutil.AssertAPIStatusCode(t, http.StatusInternalServerError, err)
	})
}