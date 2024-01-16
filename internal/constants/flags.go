package constants

import "github.com/thediveo/enumflag/v2"

type OutputMode enumflag.Flag

const (
	OutputModePretty OutputMode = iota
	OutputModePlain
	OutputModeYaml
	OutputModeJson
)

// Map enumeration values to their textual representations (value identifiers).

var OutputModeIds = map[OutputMode][]string{
	OutputModePretty: {"pretty"},
	OutputModePlain:  {"plain"},
	OutputModeYaml:   {"yaml"},
	OutputModeJson:   {"json"},
}

type QueryOutputMode enumflag.Flag

const (
	QueryOutputModePretty QueryOutputMode = iota
	QueryOutputModePlain
	QueryOutputModeYaml
	QueryOutputModeJson
)

var QueryOutputModeIds = map[QueryOutputMode][]string{
	QueryOutputModePretty: {"pretty"},
	QueryOutputModePlain:  {"plain"},
	QueryOutputModeYaml:   {"yaml"},
	QueryOutputModeJson:   {"json"},
}

type CheckOutputMode enumflag.Flag

const (
	CheckOutputModePretty CheckOutputMode = iota
	CheckOutputModePlain
	CheckOutputModeYaml
	CheckOutputModeJson
)

var CheckOutputModeIds = map[CheckOutputMode][]string{
	CheckOutputModePretty: {"pretty"},
	CheckOutputModePlain:  {"plain"},
	CheckOutputModeYaml:   {"yaml"},
	CheckOutputModeJson:   {"json"},
}

type DashboardOutputMode enumflag.Flag

const (
	DashboardOutputModePretty DashboardOutputMode = iota
	DashboardOutputModePlain
	DashboardOutputModeYaml
	DashboardOutputModeJson
)

var DashboardOutputModeIds = map[DashboardOutputMode][]string{
	DashboardOutputModePretty: {"pretty"},
	DashboardOutputModePlain:  {"plain"},
	DashboardOutputModeYaml:   {"yaml"},
	DashboardOutputModeJson:   {"json"},
}

func FlagValues[T comparable](mappings map[T][]string) []string {
	var res = make([]string, 0, len(mappings))
	for _, v := range mappings {
		res = append(res, v[0])
	}
	return res

}
