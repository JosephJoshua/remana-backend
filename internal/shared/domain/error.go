package domain

type domainError string

func (e domainError) Error() string {
	return string(e)
}

const (
	ErrInputTooLong     domainError = domainError("input too long")
	ErrInputTooShort    domainError = domainError("input too short")
	ErrInvalidStoreCode domainError = domainError("invalid store code")
	ErrInvalidLoginCode domainError = domainError("invalid login code")
)
