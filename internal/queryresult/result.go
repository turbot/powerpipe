package queryresult

import (
	"github.com/turbot/pipe-fittings/v2/queryresult"
	"time"
)

type TimingMetadata struct {
	Duration time.Duration
}

// GetTiming implements TimingContainer - we implement this interface
// to allow TimingMetadata to be used to parameterize the ResultStreamer
func (t TimingMetadata) GetTiming() any {
	return t
}

type Result = queryresult.Result[*TimingMetadata]

func NewResult(cols []*queryresult.ColumnDef) *Result {
	return queryresult.NewResult[*TimingMetadata](cols, &TimingMetadata{})
}
