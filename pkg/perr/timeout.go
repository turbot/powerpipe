package perr

import (
	"fmt"
	"net/http"
)

const (
	ErrorCodeRequestTimeout = "error_request_timeout"
)

func Timeout(itemType string, id string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeRequestTimeout,
		Title:    "Timeout",
		Status:   http.StatusRequestTimeout,
	}
	if id != "" {
		e.Detail = fmt.Sprintf("%s = %s.", itemType, id)
	}
	return e
}

func TimeoutWithMessage(msg string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeRequestTimeout,
		Title:    "Timeout",
		Status:   http.StatusRequestTimeout,
		Detail:   msg,
	}
	return e
}

func IsTimeout(err error) bool {
	e, ok := err.(ErrorModel)
	return ok && (e.Status == http.StatusRequestTimeout)
}
