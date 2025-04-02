package display

import (
	"fmt"
	"time"

	"github.com/turbot/powerpipe/internal/queryresult"
)

func PrintTiming(timingMetadata *queryresult.CheckTimingMetadata) {
	durationString := getDurationString(timingMetadata.Duration)
	fmt.Printf("\nTime: %s\n", durationString) //nolint:forbidigo // intentional use of fmt
}

func getDurationString(duration time.Duration) string {
	// Calculate duration since startTime and round down to the nearest millisecond
	durationInMS := duration / time.Millisecond
	//nolint:durationcheck // we want to print the duration in milliseconds
	duration = durationInMS * time.Millisecond

	durationString := duration.String()
	return durationString
}
