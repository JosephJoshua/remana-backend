// Code generated by ogen, DO NOT EDIT.

package genapi

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// Login implements login operation.
//
// Logs in with credentials.
//
// POST /auth/login
func (UnimplementedHandler) Login(ctx context.Context, req *LoginCredentials) (r *LoginResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// LoginCodePrompt implements loginCodePrompt operation.
//
// Logs store employees in with the login code given by the store admin. Should only be called after
// [/auth/login](#/auth/login) has been called.
//
// POST /auth/login-code
func (UnimplementedHandler) LoginCodePrompt(ctx context.Context, req *LoginCodePrompt, params LoginCodePromptParams) error {
	return ht.ErrNotImplemented
}

// NewError creates *ErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (UnimplementedHandler) NewError(ctx context.Context, err error) (r *ErrorStatusCode) {
	r = new(ErrorStatusCode)
	return r
}
