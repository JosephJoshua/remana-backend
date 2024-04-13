// Code generated by ogen, DO NOT EDIT.

package genapi

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
)

func recordError(string, error) {}

// handleGetMyUserDetailsRequest handles getMyUserDetails operation.
//
// Returns details of the currently logged in user.
//
// GET /users/me
func (s *Server) handleGetMyUserDetailsRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetMyUserDetails",
			ID:   "getMyUserDetails",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionCookie(ctx, "GetMyUserDetails", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "SessionCookie",
					Err:              err,
				}
				if encodeErr := encodeErrorResponse(s.h.NewError(ctx, err), w); encodeErr != nil {
					recordError("Security:SessionCookie", err)
				}
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			if encodeErr := encodeErrorResponse(s.h.NewError(ctx, err), w); encodeErr != nil {
				recordError("Security", err)
			}
			return
		}
	}

	var response *UserDetails
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetMyUserDetails",
			OperationSummary: "Returns details of the currently logged in user",
			OperationID:      "getMyUserDetails",
			Body:             nil,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = struct{}
			Params   = struct{}
			Response = *UserDetails
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			nil,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetMyUserDetails(ctx)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetMyUserDetails(ctx)
	}
	if err != nil {
		if errRes, ok := errors.Into[*ErrorStatusCode](err); ok {
			if err := encodeErrorResponse(errRes, w); err != nil {
				recordError("Internal", err)
			}
			return
		}
		if errors.Is(err, ht.ErrNotImplemented) {
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
		if err := encodeErrorResponse(s.h.NewError(ctx, err), w); err != nil {
			recordError("Internal", err)
		}
		return
	}

	if err := encodeGetMyUserDetailsResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleLoginRequest handles login operation.
//
// Logs in with credentials.
//
// POST /auth/login
func (s *Server) handleLoginRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "Login",
			ID:   "login",
		}
	)
	request, close, err := s.decodeLoginRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response *LoginResponse
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "Login",
			OperationSummary: "Logs in with credentials",
			OperationID:      "login",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *LoginCredentials
			Params   = struct{}
			Response = *LoginResponse
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			nil,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.Login(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.Login(ctx, request)
	}
	if err != nil {
		if errRes, ok := errors.Into[*ErrorStatusCode](err); ok {
			if err := encodeErrorResponse(errRes, w); err != nil {
				recordError("Internal", err)
			}
			return
		}
		if errors.Is(err, ht.ErrNotImplemented) {
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
		if err := encodeErrorResponse(s.h.NewError(ctx, err), w); err != nil {
			recordError("Internal", err)
		}
		return
	}

	if err := encodeLoginResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleLoginCodePromptRequest handles loginCodePrompt operation.
//
// Logs store employees in with the login code given by the store admin. Should only be called after
// [/auth/login](#/auth/login) has been called.
//
// POST /auth/login-code
func (s *Server) handleLoginCodePromptRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "LoginCodePrompt",
			ID:   "loginCodePrompt",
		}
	)
	request, close, err := s.decodeLoginCodePromptRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response *LoginCodePromptNoContent
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "LoginCodePrompt",
			OperationSummary: "Logs store employees in with login code",
			OperationID:      "loginCodePrompt",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *LoginCodePrompt
			Params   = struct{}
			Response = *LoginCodePromptNoContent
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			nil,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				err = s.h.LoginCodePrompt(ctx, request)
				return response, err
			},
		)
	} else {
		err = s.h.LoginCodePrompt(ctx, request)
	}
	if err != nil {
		if errRes, ok := errors.Into[*ErrorStatusCode](err); ok {
			if err := encodeErrorResponse(errRes, w); err != nil {
				recordError("Internal", err)
			}
			return
		}
		if errors.Is(err, ht.ErrNotImplemented) {
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
		if err := encodeErrorResponse(s.h.NewError(ctx, err), w); err != nil {
			recordError("Internal", err)
		}
		return
	}

	if err := encodeLoginCodePromptResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleLogoutRequest handles logout operation.
//
// Logs out current session.
//
// POST /auth/logout
func (s *Server) handleLogoutRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "Logout",
			ID:   "logout",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionCookie(ctx, "Logout", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "SessionCookie",
					Err:              err,
				}
				if encodeErr := encodeErrorResponse(s.h.NewError(ctx, err), w); encodeErr != nil {
					recordError("Security:SessionCookie", err)
				}
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			if encodeErr := encodeErrorResponse(s.h.NewError(ctx, err), w); encodeErr != nil {
				recordError("Security", err)
			}
			return
		}
	}

	var response *LogoutResetContent
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "Logout",
			OperationSummary: "Logs out current session",
			OperationID:      "logout",
			Body:             nil,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = struct{}
			Params   = struct{}
			Response = *LogoutResetContent
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			nil,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				err = s.h.Logout(ctx)
				return response, err
			},
		)
	} else {
		err = s.h.Logout(ctx)
	}
	if err != nil {
		if errRes, ok := errors.Into[*ErrorStatusCode](err); ok {
			if err := encodeErrorResponse(errRes, w); err != nil {
				recordError("Internal", err)
			}
			return
		}
		if errors.Is(err, ht.ErrNotImplemented) {
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
		if err := encodeErrorResponse(s.h.NewError(ctx, err), w); err != nil {
			recordError("Internal", err)
		}
		return
	}

	if err := encodeLogoutResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}
