package perr

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	ErrorCodeQuotaExceeded             = "error_quota_exceeded"
	ErrorCodeAssociationExists         = "association_exists"
	ErrorCodeWorkspaceMembershipExists = "workspace_membership_exists"
)

func QuotaExceeded(itemType string, max int) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeQuotaExceeded,
		Title:    "Quota Exceeded",
		// Closest without being obviously wrong
		// See https://softwareengineering.stackexchange.com/questions/288376/recommended-http-status-code-for-plan-limit-exceeded-response
		Status: http.StatusPaymentRequired,
	}
	if itemType != "" {
		itemTypeElements := strings.Split(itemType, ".")
		e.Detail = fmt.Sprintf("Maximum number of allowed %s is %d.", itemTypeElements[len(itemTypeElements)-2], max)
	}
	return e
}

func AssociationExist(item string, associations string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeAssociationExists,
		Title:    "Association Exist",
		Status:   http.StatusForbidden,
	}
	if item != "" {
		e.Detail = fmt.Sprintf("Connection %s is associated with workspaces: %s", item, associations)
	}
	return e
}

func MembershipExist(item string, associations string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeWorkspaceMembershipExists,
		Title:    "Workspace Membership Exist",
		Status:   http.StatusConflict,
	}
	if item != "" {
		e.Detail = fmt.Sprintf("User %s is associated with workspaces: %s", item, associations)
	}
	return e
}

func QuotaExceededWithMessage(msg string) ErrorModel {
	e := ErrorModel{
		Instance: reference(),
		Type:     ErrorCodeQuotaExceeded,
		Title:    "Quota Exceeded",
		// Closest without being obviously wrong
		// See https://softwareengineering.stackexchange.com/questions/288376/recommended-http-status-code-for-plan-limit-exceeded-response
		Status: http.StatusPaymentRequired,
		Detail: msg,
	}
	return e
}

func IsQuotaExceeded(err error) bool {
	e, ok := err.(ErrorModel)
	return ok && e.Status == http.StatusPaymentRequired
}
