package controldisplay

import (
	"fmt"
	"slices"
	"strings"
)

const severityMaxLen = len("CRITICAL")

type SeverityRenderer struct {
	severity string
}

func NewSeverityRenderer(severity string) *SeverityRenderer {
	return &SeverityRenderer{
		strings.ToUpper(severity),
	}
}

// Render returns ther severity oin upper case, got 'critical' and 'high' severities
// for all other values an empty string is returned
// NOTE: adds a trailing space
func (r SeverityRenderer) Render() string {
	if slices.Contains([]string{"CRITICAL", "HIGH"}, r.severity) {
		return fmt.Sprintf("%s ", ControlColors.Severity(r.severity))
	}
	return ""
}
