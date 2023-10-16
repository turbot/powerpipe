package perr

import (
	"net/http"
)

const (
	ErrorCodeInternal = "error_internal"
)

func InternalWithMessage(msg string) ErrorModel {
	id := reference()
	e := ErrorModel{
		Instance: id,
		Type:     ErrorCodeInternal,
		Title:    "Internal Error",
		Status:   http.StatusInternalServerError,
		Detail:   msg,
	}
	return e
}

func Internal(err error) ErrorModel {
	id := reference()
	e := ErrorModel{
		Instance: id,
		Type:     ErrorCodeInternal,
		Title:    "Internal Error",
		Status:   http.StatusInternalServerError,
		Detail:   err.Error(),
	}
	return e
}

func IsInternal(err error) bool {
	e, ok := err.(ErrorModel)
	return ok && e.Status == http.StatusInternalServerError
}
