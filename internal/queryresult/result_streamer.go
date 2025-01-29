package queryresult

import "github.com/turbot/pipe-fittings/v2/queryresult"

type ResultStreamer = queryresult.ResultStreamer[*TimingMetadata]

func NewResultStreamer() *ResultStreamer {
	return queryresult.NewResultStreamer[*TimingMetadata]()
}
