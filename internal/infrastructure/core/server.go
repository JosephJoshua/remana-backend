package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/JosephJoshua/remana-backend/internal/apierror"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/infrastructure/repository"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth"
	"github.com/JosephJoshua/remana-backend/internal/modules/damagetype"
	"github.com/JosephJoshua/remana-backend/internal/modules/misc"
	"github.com/JosephJoshua/remana-backend/internal/modules/paymentmethod"
	"github.com/JosephJoshua/remana-backend/internal/modules/permission"
	"github.com/JosephJoshua/remana-backend/internal/modules/phonecondition"
	"github.com/JosephJoshua/remana-backend/internal/modules/phoneequipment"
	"github.com/JosephJoshua/remana-backend/internal/modules/repairorder"
	"github.com/JosephJoshua/remana-backend/internal/modules/salesperson"
	"github.com/JosephJoshua/remana-backend/internal/modules/technician"
	"github.com/JosephJoshua/remana-backend/internal/modules/user"
	"github.com/go-faster/jx"
	"github.com/jackc/pgx/v5/pgxpool"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
	"github.com/rs/zerolog"
)

type authService = auth.Service
type userService = user.Service
type permissionService = permission.Service
type technicianService = technician.Service
type salesPersonService = salesperson.Service
type damageTypeService = damagetype.Service
type phoneConditionService = phonecondition.Service
type phoneEquipmentService = phoneequipment.Service
type paymentMethodService = paymentmethod.Service
type repairOrderService = repairorder.Service
type miscService = misc.Service

type server struct {
	*authService
	*userService
	*permissionService
	*technicianService
	*salesPersonService
	*damageTypeService
	*phoneConditionService
	*phoneEquipmentService
	*paymentMethodService
	*repairOrderService
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

	permissionService := permission.NewService(
		resourceLocationProvider{},
		repository.NewSQLPermissionRepository(db),
	)

	repairOrderService := repairorder.NewService(
		timeProvider{},
		resourceLocationProvider{},
		repository.NewSQLRepairOrderRepository(db),
		newRepairOrderSlugProvider(db),
	)

	technicianService := technician.NewService(
		resourceLocationProvider{},
		repository.NewSQLTechnicianRepository(db),
	)

	salesPersonService := salesperson.NewService(
		resourceLocationProvider{},
		repository.NewSQLSalesPersonRepository(db),
	)

	damageTypeService := damagetype.NewService(
		resourceLocationProvider{},
		repository.NewSQLDamageTypeRepository(db),
	)

	phoneConditionService := phonecondition.NewService(
		resourceLocationProvider{},
		repository.NewSQLPhoneConditionRepository(db),
	)

	phoneEquipmentService := phoneequipment.NewService(
		resourceLocationProvider{},
		repository.NewSQLPhoneEquipmentRepository(db),
	)

	paymentMethodService := paymentmethod.NewService(
		resourceLocationProvider{},
		repository.NewSQLPaymentMethodRepository(db),
	)

	userService := user.NewService()
	miscService := misc.NewService()

	srv := server{
		authService:           authService,
		userService:           userService,
		permissionService:     permissionService,
		technicianService:     technicianService,
		salesPersonService:    salesPersonService,
		damageTypeService:     damageTypeService,
		phoneConditionService: phoneConditionService,
		phoneEquipmentService: phoneEquipmentService,
		paymentMethodService:  paymentMethodService,
		repairOrderService:    repairOrderService,
		miscService:           miscService,
	}

	securityHandler := auth.NewSecurityHandler(sm, repository.NewSQLAuthRepository(db))

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
