package queryresult

import (
	"time"

	"github.com/turbot/pipe-fittings/v2/queryresult"
)

type CheckTimingMetadata struct {
	Duration time.Duration
}

// GetTiming implements TimingContainer - we implement this interface
// to allow CheckTimingMetadata to be used to parameterize the ResultStreamer
func (t CheckTimingMetadata) GetTiming() any {
	return t
}

type Result = queryresult.Result[*queryresult.QueryTimingMetadata]

func NewResult(cols []*queryresult.ColumnDef) *Result {
	return queryresult.NewResult[*queryresult.QueryTimingMetadata](cols, &queryresult.QueryTimingMetadata{})
}
