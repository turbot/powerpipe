package perr

import (
	"fmt"
	"net/http"
)

const (
	ErrorCodeTooManyRequests = "error_too_many_requests"
)

func TooManyRequests(itemType string, id string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeTooManyRequests,
		Title:    "Too Many Requests",
		Status:   http.StatusTooManyRequests,
	}
	if id != "" {
		e.Detail = fmt.Sprintf("%s = %s.", itemType, id)
	}
	return e
}

func TooManyRequestsWithMessage(msg string) ErrorModel {
	return TooManyRequestsWithTypeAndMessage(ErrorCodeTooManyRequests, msg)
}

func TooManyRequestsWithTypeAndMessage(errorType string, msg string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     errorType,
		Title:    "Too Many Requests",
		Status:   http.StatusTooManyRequests,
		Detail:   msg,
	}
	return e
}

func IsTooManyRequests(err error) bool {
	e, ok := err.(ErrorModel)
	return ok && e.Status == http.StatusTooManyRequests
}
