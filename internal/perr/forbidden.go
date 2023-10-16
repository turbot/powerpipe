package perr

import (
	"net/http"
)

const (
	ErrorCodeForbidden = "error_forbidden"
)

func Forbidden() ErrorModel {
	return ForbiddenWithMessage("Forbidden.")
}

func ForbiddenWithMessage(msg string) ErrorModel {
	return ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeForbidden,
		Title:    "Forbidden",
		Detail:   msg,
		Status:   http.StatusForbidden,
	}
}

func IsForbidden(err error) bool {
	e, ok := err.(ErrorModel)
	return ok && e.Status == http.StatusForbidden
}
