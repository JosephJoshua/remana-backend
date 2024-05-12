package apperror

type appError string

func (e appError) Error() string {
	return string(e)
}

const (
	ErrValueAlreadySet        appError = appError("value already set")
	ErrInvalidInput           appError = appError("invalid input")
	ErrPasswordTooLong        appError = appError("password too long")
	ErrPasswordMismatch       appError = appError("password mismatch")
	ErrMisingLoginCodePrompt  appError = appError("missing login code prompt")
	ErrMissingSession         appError = appError("missing session")
	ErrUserNotFound           appError = appError("user not found")
	ErrDamageNotFound         appError = appError("damage not found")
	ErrPhoneConditionNotFound appError = appError("phone condition not found")
	ErrPhoneEquipmentNotFound appError = appError("phone equipment not found")
	ErrPermissionNotFound     appError = appError("permission not found")
	ErrLoginCodeMismatch      appError = appError("login code mismatch")
)
