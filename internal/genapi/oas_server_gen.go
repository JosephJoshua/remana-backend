// Code generated by ogen, DO NOT EDIT.

package genapi

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// CreateRepairOrder implements createRepairOrder operation.
	//
	// Creates a new repair order.
	//
	// POST /repair-orders
	CreateRepairOrder(ctx context.Context, req *CreateRepairOrderRequest) (*CreateRepairOrderCreated, error)
	// CreateSalesPerson implements createSalesPerson operation.
	//
	// Creates a new sales person.
	//
	// POST /sales-persons
	CreateSalesPerson(ctx context.Context, req *CreateSalesPersonRequest) (*CreateSalesPersonCreated, error)
	// CreateTechnician implements createTechnician operation.
	//
	// Creates a new technician.
	//
	// POST /technicians
	CreateTechnician(ctx context.Context, req *CreateTechnicianRequest) (*CreateTechnicianCreated, error)
	// GetHealth implements getHealth operation.
	//
	// Returns the health status of the service.
	//
	// GET /healthz
	GetHealth(ctx context.Context) error
	// GetMyUserDetails implements getMyUserDetails operation.
	//
	// Returns details of the currently logged in user.
	//
	// GET /users/me
	GetMyUserDetails(ctx context.Context) (*UserDetails, error)
	// Login implements login operation.
	//
	// Logs in with credentials.
	//
	// POST /auth/login
	Login(ctx context.Context, req *LoginCredentials) (*LoginResponse, error)
	// LoginCodePrompt implements loginCodePrompt operation.
	//
	// Logs store employees in with the login code given by the store admin. Should only be called after
	// [/auth/login](#/auth/login) has been called.
	//
	// POST /auth/login-code
	LoginCodePrompt(ctx context.Context, req *LoginCodePrompt) error
	// Logout implements logout operation.
	//
	// Logs out current session.
	//
	// POST /auth/logout
	Logout(ctx context.Context) error
	// NewError creates *ErrorStatusCode from error returned by handler.
	//
	// Used for common default response.
	NewError(ctx context.Context, err error) *ErrorStatusCode
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}
