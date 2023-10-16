package perr

import (
	"net/http"
)

const (
	ErrorCodeUnauthorized = "error_unauthorized"
)

func Unauthorized() ErrorModel {
	return UnauthorizedWithMessage("Unauthorized.")
}

func UnauthorizedWithMessage(msg string) ErrorModel {
	return ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeUnauthorized,
		Title:    "Unauthorized",
		Status:   http.StatusUnauthorized,
		Detail:   msg,
	}
}

func IsUnauthorized(err error) bool {
	e, ok := err.(ErrorModel)
	return ok && e.Status == http.StatusUnauthorized
}
