package dashboardexecute

import "time"

type TimeRange struct {
	From *time.Time `json:"from"`
	To   *time.Time `json:"to"`
}

type InputValues struct {
	Inputs map[string]interface{} `json:"inputs"`
	// map of time ranges, keyed by target benchmark/detection
	DetectionTimeRange TimeRange `json:"detection_time_ranges"`
}

// ctor
func NewInputValues() *InputValues {
	return &InputValues{
		Inputs: make(map[string]interface{}),
	}
}

func (v *InputValues) Empty() bool {
	if v == nil {
		return true
	}
	return len(v.Inputs) > 0
}
