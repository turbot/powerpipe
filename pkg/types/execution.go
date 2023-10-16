package types

import (
	"encoding/json"
	"time"
)

type EventLogEntry struct {
	EventType string          `json:"event_type"`
	Timestamp *time.Time      `json:"ts"`
	Payload   json.RawMessage `json:"payload"`
}

type Stack struct {
	ID         string            `json:"id"`
	Status     string            `json:"status"`
	StepStatus map[int]string    `json:"pipeline_step_status"`
	Stacks     map[string]*Stack `json:"children"`
}
