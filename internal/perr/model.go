package perr

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

// As per RFC7807 (https://tools.ietf.org/html/rfc7807) define a standard error model with a limited set of Flowpipe-specific extensions
// Initial inspiration taken from https://github.com/danielgtaylor/huma/blob/master/error.go
type ErrorDetailModel struct {
	Message  string `json:"message"`
	Location string `json:"location,omitempty"`
}

type ErrorModel struct {
	Instance string `json:"instance" binding:"required"`
	ID       string `json:"-"`
	Type     string `json:"type" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Status   int    `json:"status" binding:"required"`

	// If we don't have required it comes out as pointer and there is a bug in the formatter
	Detail string `json:"detail" binding:"required"`

	ValidationErrors []*ErrorDetailModel `json:"validation_errors,omitempty"`

	// All errors are fatal unless specified
	Retryable bool `json:"retryable,omitempty"`
}

func FromHttpError(err error, statusCode int) ErrorModel {

	var errorMsg string
	if err != nil {
		errorMsg = err.Error()
	}

	switch statusCode {
	case http.StatusNotFound:
		return NotFoundWithMessage(errorMsg)
	case http.StatusForbidden:
		return ForbiddenWithMessage(errorMsg)
	case http.StatusRequestTimeout:
		return TimeoutWithMessage(errorMsg)
	case http.StatusBadRequest:
		return BadRequestWithMessage(errorMsg)
	case http.StatusConflict:
		return ConflictWithMessage(errorMsg)
	case http.StatusPreconditionFailed:
		return PreconditionFailedWithMessage(errorMsg)
	case http.StatusPaymentRequired:
		return UnsupportedPlanValueWithMessage(errorMsg)
	case http.StatusInternalServerError:
		return InternalWithMessage(errorMsg)
	default:
		return InternalWithMessage(errorMsg)
	}
}

func (e ErrorModel) Error() string {
	if e.Detail != "" {
		return e.Title + ": " + e.Detail
	}
	return e.Title
}

func (e ErrorModel) GetStatus() int {
	return e.Status
}

type ValidationError struct {
	Type   string                     `json:"type"`   // Denotes the location where the validation error was encountered.
	Errors validator.ValidationErrors `json:"errors"` // The list of validation errors.
}

func (e ValidationError) Error() string {
	return e.Type + ": " + e.Errors.Error()
}
