package apperror

type appError string

func (e appError) Error() string {
	return string(e)
}

const (
	ErrPasswordTooLong  appError = appError("password too long")
	ErrPasswordMismatch appError = appError("password mismatch")
	ErrUserNotFound     appError = appError("user not found")
)
