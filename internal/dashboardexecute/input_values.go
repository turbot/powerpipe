package dashboardexecute

import (
	"github.com/turbot/pipe-fittings/v2/utils"
)

type InputValues struct {
	Inputs map[string]interface{} `json:"inputs"`
	// map of time ranges, keyed by target benchmark/detection
	DateTimeRange utils.TimeRange `json:"detection_time_ranges"`
}

func NewInputValues() *InputValues {
	return &InputValues{
		Inputs: make(map[string]interface{}),
	}
}

func (v *InputValues) Empty() bool {
	if v == nil {
		return true
	}
	return len(v.Inputs) == 0
}
