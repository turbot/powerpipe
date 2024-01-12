package types

import "github.com/thediveo/enumflag/v2"

// ① Define your new enum flag type. It can be derived from enumflag.Flag,
// but it doesn't need to be as long as it satisfies constraints.Integer.

type OutputMode enumflag.Flag

// ② Define the enumeration values

const (
	OutputModePretty OutputMode = iota
	OutputModePlain
	OutputModeYaml
	OutputModeJson
)

// ③ Map enumeration values to their textual representations (value identifiers).

var OutputModeIds = map[OutputMode][]string{
	OutputModePretty: {"pretty"},
	OutputModePlain:  {"plain"},
	OutputModeYaml:   {"yaml"},
	OutputModeJson:   {"json"},
}
