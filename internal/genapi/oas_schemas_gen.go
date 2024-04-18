// Code generated by ogen, DO NOT EDIT.

package genapi

import (
	"fmt"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

func (s *ErrorStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

// Ref: #
type Error struct {
	Message string `json:"message"`
}

// GetMessage returns the value of Message.
func (s *Error) GetMessage() string {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *Error) SetMessage(val string) {
	s.Message = val
}

// ErrorStatusCode wraps Error with StatusCode.
type ErrorStatusCode struct {
	StatusCode int
	Response   Error
}

// GetStatusCode returns the value of StatusCode.
func (s *ErrorStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *ErrorStatusCode) GetResponse() Error {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *ErrorStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *ErrorStatusCode) SetResponse(val Error) {
	s.Response = val
}

// GetHealthNoContent is response for GetHealth operation.
type GetHealthNoContent struct{}

// Ref: #
type LoginCodePrompt struct {
	LoginCode string `json:"login_code"`
}

// GetLoginCode returns the value of LoginCode.
func (s *LoginCodePrompt) GetLoginCode() string {
	return s.LoginCode
}

// SetLoginCode sets the value of LoginCode.
func (s *LoginCodePrompt) SetLoginCode(val string) {
	s.LoginCode = val
}

// LoginCodePromptNoContent is response for LoginCodePrompt operation.
type LoginCodePromptNoContent struct{}

// Ref: #
type LoginCredentials struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	StoreCode string `json:"store_code"`
}

// GetUsername returns the value of Username.
func (s *LoginCredentials) GetUsername() string {
	return s.Username
}

// GetPassword returns the value of Password.
func (s *LoginCredentials) GetPassword() string {
	return s.Password
}

// GetStoreCode returns the value of StoreCode.
func (s *LoginCredentials) GetStoreCode() string {
	return s.StoreCode
}

// SetUsername sets the value of Username.
func (s *LoginCredentials) SetUsername(val string) {
	s.Username = val
}

// SetPassword sets the value of Password.
func (s *LoginCredentials) SetPassword(val string) {
	s.Password = val
}

// SetStoreCode sets the value of StoreCode.
func (s *LoginCredentials) SetStoreCode(val string) {
	s.StoreCode = val
}

// Ref: #
type LoginResponse struct {
	// The type of user that logged in:
	// * `admin` - Store admin. The session ID is returned in a cookie named
	// `session_id`. You need to include this cookie in subsequent requests.
	// * `employee` - Store employee. The user needs to log in with a login
	// code given by the store admin. The login code prompt ID is returned in
	// a cookie named `login_code_prompt_id`. You need to visit [/auth/login-code](#/auth/loginCodePrompt)
	// with the login code to log in.
	Type LoginResponseType `json:"type"`
}

// GetType returns the value of Type.
func (s *LoginResponse) GetType() LoginResponseType {
	return s.Type
}

// SetType sets the value of Type.
func (s *LoginResponse) SetType(val LoginResponseType) {
	s.Type = val
}

// The type of user that logged in:
// * `admin` - Store admin. The session ID is returned in a cookie named
// `session_id`. You need to include this cookie in subsequent requests.
// * `employee` - Store employee. The user needs to log in with a login
// code given by the store admin. The login code prompt ID is returned in
// a cookie named `login_code_prompt_id`. You need to visit [/auth/login-code](#/auth/loginCodePrompt)
// with the login code to log in.
type LoginResponseType string

const (
	LoginResponseTypeAdmin    LoginResponseType = "admin"
	LoginResponseTypeEmployee LoginResponseType = "employee"
)

// AllValues returns all LoginResponseType values.
func (LoginResponseType) AllValues() []LoginResponseType {
	return []LoginResponseType{
		LoginResponseTypeAdmin,
		LoginResponseTypeEmployee,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s LoginResponseType) MarshalText() ([]byte, error) {
	switch s {
	case LoginResponseTypeAdmin:
		return []byte(s), nil
	case LoginResponseTypeEmployee:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *LoginResponseType) UnmarshalText(data []byte) error {
	switch LoginResponseType(data) {
	case LoginResponseTypeAdmin:
		*s = LoginResponseTypeAdmin
		return nil
	case LoginResponseTypeEmployee:
		*s = LoginResponseTypeEmployee
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// LogoutResetContent is response for Logout operation.
type LogoutResetContent struct{}

type SessionCookie struct {
	APIKey string
}

// GetAPIKey returns the value of APIKey.
func (s *SessionCookie) GetAPIKey() string {
	return s.APIKey
}

// SetAPIKey sets the value of APIKey.
func (s *SessionCookie) SetAPIKey(val string) {
	s.APIKey = val
}

// Ref: #
type UserDetails struct {
	ID       uuid.UUID        `json:"id"`
	Username string           `json:"username"`
	Role     UserDetailsRole  `json:"role"`
	Store    UserDetailsStore `json:"store"`
}

// GetID returns the value of ID.
func (s *UserDetails) GetID() uuid.UUID {
	return s.ID
}

// GetUsername returns the value of Username.
func (s *UserDetails) GetUsername() string {
	return s.Username
}

// GetRole returns the value of Role.
func (s *UserDetails) GetRole() UserDetailsRole {
	return s.Role
}

// GetStore returns the value of Store.
func (s *UserDetails) GetStore() UserDetailsStore {
	return s.Store
}

// SetID sets the value of ID.
func (s *UserDetails) SetID(val uuid.UUID) {
	s.ID = val
}

// SetUsername sets the value of Username.
func (s *UserDetails) SetUsername(val string) {
	s.Username = val
}

// SetRole sets the value of Role.
func (s *UserDetails) SetRole(val UserDetailsRole) {
	s.Role = val
}

// SetStore sets the value of Store.
func (s *UserDetails) SetStore(val UserDetailsStore) {
	s.Store = val
}

type UserDetailsRole struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	IsStoreAdmin bool      `json:"is_store_admin"`
}

// GetID returns the value of ID.
func (s *UserDetailsRole) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *UserDetailsRole) GetName() string {
	return s.Name
}

// GetIsStoreAdmin returns the value of IsStoreAdmin.
func (s *UserDetailsRole) GetIsStoreAdmin() bool {
	return s.IsStoreAdmin
}

// SetID sets the value of ID.
func (s *UserDetailsRole) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *UserDetailsRole) SetName(val string) {
	s.Name = val
}

// SetIsStoreAdmin sets the value of IsStoreAdmin.
func (s *UserDetailsRole) SetIsStoreAdmin(val bool) {
	s.IsStoreAdmin = val
}

type UserDetailsStore struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

// GetID returns the value of ID.
func (s *UserDetailsStore) GetID() uuid.UUID {
	return s.ID
}

// GetName returns the value of Name.
func (s *UserDetailsStore) GetName() string {
	return s.Name
}

// GetCode returns the value of Code.
func (s *UserDetailsStore) GetCode() string {
	return s.Code
}

// SetID sets the value of ID.
func (s *UserDetailsStore) SetID(val uuid.UUID) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *UserDetailsStore) SetName(val string) {
	s.Name = val
}

// SetCode sets the value of Code.
func (s *UserDetailsStore) SetCode(val string) {
	s.Code = val
}
