package perr

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	ErrorCodeUnsupportedPlanValue = "error_unsupported_plan_value"
)

func UnsupportedPlanValue(itemType, value string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeUnsupportedPlanValue,
		Title:    "Unsupported Plan Value",
		// Closest without being obviously wrong
		// See https://softwareengineering.stackexchange.com/questions/288376/recommended-http-status-code-for-plan-limit-exceeded-response
		Status: http.StatusPaymentRequired,
	}
	if itemType != "" {
		itemTypeElements := strings.Split(itemType, ".")
		e.Detail = fmt.Sprintf("%s is not a supported value for %s in the current plan.", value, itemTypeElements[len(itemTypeElements)-2])
	}
	return e
}

func UnsupportedPlanValueWithMessage(msg string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeUnsupportedPlanValue,
		Title:    "Unsupported Plan Value",
		// Closest without being obviously wrong
		// See https://softwareengineering.stackexchange.com/questions/288376/recommended-http-status-code-for-plan-limit-exceeded-response
		Status: http.StatusPaymentRequired,
		Detail: msg,
	}
	return e
}
