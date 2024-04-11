package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/JosephJoshua/remana-backend/internal/auth"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/JosephJoshua/remana-backend/internal/shared/repository"
	"github.com/go-faster/jx"
	"github.com/google/uuid"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
	"github.com/rs/zerolog"
)

type server struct {
	*auth.Service
}

type Middleware func(next http.Handler) http.Handler

func NewAPIServer() (*genapi.Server, []Middleware, error) {
	dummyStore, err := domain.NewStore(1, "store", "store")
	if err != nil {
		return nil, []Middleware{}, err
	}

	adminRole, err := domain.NewRole(1, "admin", *dummyStore, true)
	if err != nil {
		return nil, []Middleware{}, err
	}

	userRole, err := domain.NewRole(2, "user", *dummyStore, false)
	if err != nil {
		return nil, []Middleware{}, err
	}

	password, err := (&PasswordHasher{}).Hash("password")
	if err != nil {
		return nil, []Middleware{}, err
	}

	adminUser, err := domain.NewUser(uuid.New(), "username", password, *dummyStore, *adminRole)
	if err != nil {
		return nil, []Middleware{}, err
	}

	employeeUser, err := domain.NewUser(uuid.New(), "username2", password, *dummyStore, *userRole)
	if err != nil {
		return nil, []Middleware{}, err
	}

	sm := newAuthSessionManager()
	pm := newLoginCodePromptManager()

	middlewares := []Middleware{requestLoggerMiddleware, sm.middleware, pm.middleware}

	authService := auth.NewService(
		sm,
		pm,
		repository.NewMemoryAuthRepository([]domain.User{*adminUser, *employeeUser}),
		&PasswordHasher{},
	)

	srv := server{
		Service: authService,
	}

	oasSrv, err := genapi.NewServer(srv, genapi.WithErrorHandler(handleServerError))
	if err != nil {
		return nil, []Middleware{}, fmt.Errorf("error creating oas server: %w", err)
	}

	return oasSrv, middlewares, nil
}

func (s server) NewError(_ context.Context, _ error) *genapi.ErrorStatusCode {
	return nil
}

func handleServerError(_ context.Context, w http.ResponseWriter, r *http.Request, err error) {
	code := http.StatusInternalServerError
	message := "unexpected internal server error"

	var (
		ctError *validate.InvalidContentTypeError
		ogenErr ogenerrors.Error
	)

	switch {
	case errors.Is(err, ht.ErrNotImplemented):
		code = http.StatusNotImplemented
		message = "operation not implemented"

	case errors.As(err, &ctError):
		// Takes precedence over ogenerrors.Error.
		code = http.StatusUnsupportedMediaType
		message = "invalid content type"

	case errors.As(err, &ogenErr):
		code = ogenErr.Code()
		message = ogenErr.Error()

	default:
		l := zerolog.Ctx(r.Context())
		l.Error().Err(err).Msg("handleServerError(); unexpected internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	apiError := genapi.Error{
		Message: message,
	}

	e := jx.GetEncoder()
	e.Obj(func(e *jx.Encoder) {
		e.Field("message", func(e *jx.Encoder) {
			e.StrEscape(apiError.GetMessage())
		})
	})

	_, _ = w.Write(e.Bytes())
}
