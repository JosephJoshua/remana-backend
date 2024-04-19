package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/JosephJoshua/remana-backend/internal/auth"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/misc"
	"github.com/JosephJoshua/remana-backend/internal/shared/apierror"
	"github.com/JosephJoshua/remana-backend/internal/shared/repository"
	"github.com/JosephJoshua/remana-backend/internal/user"
	"github.com/go-faster/jx"
	"github.com/jackc/pgx/v5/pgxpool"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
	"github.com/rs/zerolog"
)

type authService = auth.Service
type userService = user.Service
type miscService = misc.Service

type server struct {
	*authService
	*userService
	*miscService
}

type Middleware func(next http.Handler) http.Handler

func NewAPIServer(db *pgxpool.Pool) (*genapi.Server, []Middleware, error) {
	sm := newAuthSessionManager()
	pm := newLoginCodePromptManager()

	middlewares := []Middleware{requestLoggerMiddleware, sm.middleware, pm.middleware}

	authService := auth.NewService(
		sm,
		pm,
		repository.NewSQLAuthRepository(db),
		&PasswordHasher{},
	)

	userService := user.NewService()
	miscService := misc.NewService()

	srv := server{
		authService: authService,
		userService: userService,
		miscService: miscService,
	}

	securityHandler := auth.NewSecurityHandler(sm, repository.NewSQLUserRepository(db))

	oasSrv, err := genapi.NewServer(srv, securityHandler, genapi.WithErrorHandler(handleServerError))
	if err != nil {
		return nil, []Middleware{}, fmt.Errorf("error creating oas server: %w", err)
	}

	return oasSrv, middlewares, nil
}

func (s server) NewError(_ context.Context, err error) *genapi.ErrorStatusCode {
	var apiErr *genapi.ErrorStatusCode
	if errors.As(err, &apiErr) {
		return apiErr
	}

	if errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied) {
		return apierror.ToAPIError(http.StatusUnauthorized, "unauthorized")
	}

	return apierror.ToAPIError(http.StatusInternalServerError, "unexpected internal error")
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
