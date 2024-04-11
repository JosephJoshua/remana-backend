package apperror

type appError string

func (e appError) Error() string {
	return string(e)
}

const (
	ErrPasswordTooLong       appError = appError("password too long")
	ErrPasswordMismatch      appError = appError("password mismatch")
	ErrMisingLoginCodePrompt appError = appError("missing login code prompt")
	ErrUserNotFound          appError = appError("user not found")
	ErrLoginCodeMismatch     appError = appError("login code mismatch")
)
