package perr

import (
	"fmt"
	"net/http"
)

const (
	ErrorCodePreconditionFailed = "error_precondition_failed"
)

func PreconditionFailed(expected string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodePreconditionFailed,
		Title:    "Precondition Failed",
		Status:   http.StatusPreconditionFailed,
		// Deliberately exclude the actual version for security reasons
		Detail: fmt.Sprintf("If-Match = %s does not match the current resource version.", expected),
	}
	return e
}

func PreconditionFailedWithMessage(msg string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodePreconditionFailed,
		Title:    "Precondition Failed",
		Status:   http.StatusPreconditionFailed,
		Detail:   msg,
	}
	return e
}

func IsPreconditionFailed(err error) bool {
	e, ok := err.(ErrorModel)
	return ok && e.Status == http.StatusPreconditionFailed
}
