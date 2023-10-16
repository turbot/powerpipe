package perr

import (
	"fmt"
	"net/http"
)

const (
	ErrorCodeNotFound = "error_not_found"
)

func NotFound(itemType string, id string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeNotFound,
		Title:    "Not Found",
		Status:   http.StatusNotFound,
	}
	if id != "" {
		e.Detail = fmt.Sprintf("%s = %s.", itemType, id)
	}
	return e
}

func NotFoundWithMessage(msg string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeNotFound,
		Title:    "Not Found",
		Status:   http.StatusNotFound,
		Detail:   msg,
	}
	return e
}

func IsNotFound(err error) bool {
	e, ok := err.(ErrorModel)
	return ok && (e.Status == http.StatusNotFound)
}
