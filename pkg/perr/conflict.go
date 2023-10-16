package perr

import (
	"fmt"
	"net/http"
)

const (
	ErrorCodeConflict = "error_conflict"
)

func Conflict(itemType string, id string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeConflict,
		Title:    "Conflict",
		Status:   http.StatusConflict,
	}
	if id != "" {
		e.Detail = fmt.Sprintf("%s %s already in use. If you have just deleted the resource, it may take a few minutes for the resource to be fully deleted from the system.", itemType, id)
	}
	return e
}

func ConflictWithMessage(msg string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeConflict,
		Title:    "Conflict",
		Status:   http.StatusConflict,
		Detail:   msg,
	}
	return e
}

func IsConflict(err error) bool {
	e, ok := err.(ErrorModel)
	return ok && e.Status == http.StatusConflict
}
