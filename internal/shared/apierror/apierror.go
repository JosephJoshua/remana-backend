package apierror

import (
	"github.com/JosephJoshua/repair-management-backend/internal/genapi"
)

func ToAPIError(statusCode int, message string) *genapi.ErrorStatusCode {
	return &genapi.ErrorStatusCode{
		StatusCode: statusCode,
		Response: genapi.Error{
			Message: message,
		},
	}
}
