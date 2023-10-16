package utils

import "time"

// TimePtr returns a pointer to the time.Time value passed in.
func TimePtr(t time.Time) *time.Time {
	return &t
}

// TimeNowPtr returns a pointer to the current time in UTC.
func TimeNow() *time.Time {
	return TimePtr(time.Now().UTC())
}
