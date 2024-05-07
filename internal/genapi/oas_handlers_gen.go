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

// handleCreateDamageTypeRequest handles createDamageType operation.
//
// Creates a new damage type.
//
// POST /damage-types
func (s *Server) handleCreateDamageTypeRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "CreateDamageType",
			ID:   "createDamageType",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionCookie(ctx, "CreateDamageType", r)
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
	request, close, err := s.decodeCreateDamageTypeRequest(r)
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

	var response *CreateDamageTypeCreated
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "CreateDamageType",
			OperationSummary: "Creates a new damage type",
			OperationID:      "createDamageType",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *CreateDamageTypeRequest
			Params   = struct{}
			Response = *CreateDamageTypeCreated
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
				response, err = s.h.CreateDamageType(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.CreateDamageType(ctx, request)
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

	if err := encodeCreateDamageTypeResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleCreatePaymentMethodRequest handles createPaymentMethod operation.
//
// Creates a new payment method.
//
// POST /payment-methods
func (s *Server) handleCreatePaymentMethodRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "CreatePaymentMethod",
			ID:   "createPaymentMethod",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionCookie(ctx, "CreatePaymentMethod", r)
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
	request, close, err := s.decodeCreatePaymentMethodRequest(r)
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

	var response *CreatePaymentMethodCreated
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "CreatePaymentMethod",
			OperationSummary: "Creates a new payment method",
			OperationID:      "createPaymentMethod",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *CreatePaymentMethodRequest
			Params   = struct{}
			Response = *CreatePaymentMethodCreated
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
				response, err = s.h.CreatePaymentMethod(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.CreatePaymentMethod(ctx, request)
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

	if err := encodeCreatePaymentMethodResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleCreatePhoneConditionRequest handles createPhoneCondition operation.
//
// Creates a new phone condition.
//
// POST /phone-conditions
func (s *Server) handleCreatePhoneConditionRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "CreatePhoneCondition",
			ID:   "createPhoneCondition",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionCookie(ctx, "CreatePhoneCondition", r)
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
	request, close, err := s.decodeCreatePhoneConditionRequest(r)
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

	var response *CreatePhoneConditionCreated
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "CreatePhoneCondition",
			OperationSummary: "Creates a new phone condition",
			OperationID:      "createPhoneCondition",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *CreatePhoneConditionRequest
			Params   = struct{}
			Response = *CreatePhoneConditionCreated
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
				response, err = s.h.CreatePhoneCondition(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.CreatePhoneCondition(ctx, request)
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

	if err := encodeCreatePhoneConditionResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleCreatePhoneEquipmentRequest handles createPhoneEquipment operation.
//
// Creates a new phone equipment.
//
// POST /phone-equipments
func (s *Server) handleCreatePhoneEquipmentRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "CreatePhoneEquipment",
			ID:   "createPhoneEquipment",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionCookie(ctx, "CreatePhoneEquipment", r)
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
	request, close, err := s.decodeCreatePhoneEquipmentRequest(r)
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

	var response *CreatePhoneEquipmentCreated
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "CreatePhoneEquipment",
			OperationSummary: "Creates a new phone equipment",
			OperationID:      "createPhoneEquipment",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *CreatePhoneEquipmentRequest
			Params   = struct{}
			Response = *CreatePhoneEquipmentCreated
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
				response, err = s.h.CreatePhoneEquipment(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.CreatePhoneEquipment(ctx, request)
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

	if err := encodeCreatePhoneEquipmentResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleCreateRepairOrderRequest handles createRepairOrder operation.
//
// Creates a new repair order.
//
// POST /repair-orders
func (s *Server) handleCreateRepairOrderRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "CreateRepairOrder",
			ID:   "createRepairOrder",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionCookie(ctx, "CreateRepairOrder", r)
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
	request, close, err := s.decodeCreateRepairOrderRequest(r)
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

	var response *CreateRepairOrderCreated
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "CreateRepairOrder",
			OperationSummary: "Creates a new repair order",
			OperationID:      "createRepairOrder",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *CreateRepairOrderRequest
			Params   = struct{}
			Response = *CreateRepairOrderCreated
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
				response, err = s.h.CreateRepairOrder(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.CreateRepairOrder(ctx, request)
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

	if err := encodeCreateRepairOrderResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleCreateRoleRequest handles createRole operation.
//
// Creates a role.
//
// POST /roles
func (s *Server) handleCreateRoleRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "CreateRole",
			ID:   "createRole",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionCookie(ctx, "CreateRole", r)
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
	request, close, err := s.decodeCreateRoleRequest(r)
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

	var response *CreateRoleCreated
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "CreateRole",
			OperationSummary: "Creates a role",
			OperationID:      "createRole",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *CreateRoleRequest
			Params   = struct{}
			Response = *CreateRoleCreated
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
				response, err = s.h.CreateRole(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.CreateRole(ctx, request)
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

	if err := encodeCreateRoleResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleCreateSalesPersonRequest handles createSalesPerson operation.
//
// Creates a new sales person.
//
// POST /sales-persons
func (s *Server) handleCreateSalesPersonRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "CreateSalesPerson",
			ID:   "createSalesPerson",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionCookie(ctx, "CreateSalesPerson", r)
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
	request, close, err := s.decodeCreateSalesPersonRequest(r)
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

	var response *CreateSalesPersonCreated
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "CreateSalesPerson",
			OperationSummary: "Creates a new sales person",
			OperationID:      "createSalesPerson",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *CreateSalesPersonRequest
			Params   = struct{}
			Response = *CreateSalesPersonCreated
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
				response, err = s.h.CreateSalesPerson(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.CreateSalesPerson(ctx, request)
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

	if err := encodeCreateSalesPersonResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleCreateTechnicianRequest handles createTechnician operation.
//
// Creates a new technician.
//
// POST /technicians
func (s *Server) handleCreateTechnicianRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "CreateTechnician",
			ID:   "createTechnician",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionCookie(ctx, "CreateTechnician", r)
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
	request, close, err := s.decodeCreateTechnicianRequest(r)
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

	var response *CreateTechnicianCreated
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "CreateTechnician",
			OperationSummary: "Creates a new technician",
			OperationID:      "createTechnician",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *CreateTechnicianRequest
			Params   = struct{}
			Response = *CreateTechnicianCreated
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
				response, err = s.h.CreateTechnician(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.CreateTechnician(ctx, request)
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

	if err := encodeCreateTechnicianResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetHealthRequest handles getHealth operation.
//
// Returns the health status of the service.
//
// GET /healthz
func (s *Server) handleGetHealthRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err error
	)

	var response *GetHealthNoContent
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetHealth",
			OperationSummary: "Returns the health status of the service",
			OperationID:      "getHealth",
			Body:             nil,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = struct{}
			Params   = struct{}
			Response = *GetHealthNoContent
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
				err = s.h.GetHealth(ctx)
				return response, err
			},
		)
	} else {
		err = s.h.GetHealth(ctx)
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

	if err := encodeGetHealthResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

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
