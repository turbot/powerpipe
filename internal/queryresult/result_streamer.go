package queryresult

import "github.com/turbot/pipe-fittings/v2/queryresult"

type ResultStreamer = queryresult.ResultStreamer[*CheckTimingMetadata]

func NewResultStreamer() *ResultStreamer {
	return queryresult.NewResultStreamer[*CheckTimingMetadata]()
}
