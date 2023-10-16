package perr

import (
	"fmt"
	"net/http"
)

const (
	ErrorCodeBadRequest        = "error_bad_request"
	ErrorCodeInvalidData       = "error_invalid_data"
	ErrorCodeDependencyFailure = "error_dependency_failure"
)

func BadRequest(itemType string, id string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeBadRequest,
		Title:    "Bad Request",
		Status:   http.StatusBadRequest,
	}
	if id != "" {
		e.Detail = fmt.Sprintf("%s = %s.", itemType, id)
	}
	return e
}

func BadRequestWithMessage(msg string) ErrorModel {
	return BadRequestWithTypeAndMessage(ErrorCodeBadRequest, msg)
}

func BadRequestWithTypeAndMessage(errorType string, msg string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     errorType,
		Title:    "Bad Request",
		Status:   http.StatusBadRequest,
		Detail:   msg,
	}
	return e
}

func IsBadRequest(err error) bool {
	e, ok := err.(ErrorModel)
	return ok && e.Status == http.StatusBadRequest
}
