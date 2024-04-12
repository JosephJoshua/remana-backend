package domain

type domainError string

func (e domainError) Error() string {
	return string(e)
}

const (
	ErrInputTooLong     domainError = domainError("input too long")
	ErrInputTooShort    domainError = domainError("input too short")
	ErrInvalidStoreCode domainError = domainError("store code must only contain lowercase letters separated by hyphens")
	ErrInvalidLoginCode domainError = domainError("login code must only contain alphanumeric characters")
)
