package constants

import (
	"github.com/thediveo/enumflag/v2"
	"github.com/turbot/pipe-fittings/constants"
)

type OutputMode enumflag.Flag

const (
	OutputModePretty OutputMode = iota
	OutputModePlain
	OutputModeYaml
	OutputModeJson
)

// Map enumeration values to their textual representations (value identifiers).

var OutputModeIds = map[OutputMode][]string{
	OutputModePretty: {constants.OutputFormatPretty},
	OutputModePlain:  {constants.OutputFormatPlain},
	OutputModeYaml:   {constants.OutputFormatYAML},
	OutputModeJson:   {constants.OutputFormatJSON},
}

type QueryOutputMode enumflag.Flag

const (
	QueryOutputModeCsv QueryOutputMode = iota
	QueryOutputModeJson
	QueryOutputModeLine
	QueryOutputModeSnapshot
	QueryOutputModeSnapshotShort
	QueryOutputModeTable
)

var QueryOutputModeIds = map[QueryOutputMode][]string{
	QueryOutputModeCsv:           {constants.OutputFormatCSV},
	QueryOutputModeJson:          {constants.OutputFormatJSON},
	QueryOutputModeLine:          {constants.OutputFormatLine},
	QueryOutputModeSnapshot:      {constants.OutputFormatSnapshot},
	QueryOutputModeSnapshotShort: {constants.OutputFormatSnapshotShort},
	QueryOutputModeTable:         {constants.OutputFormatTable},
}

type DashboardOutputMode enumflag.Flag

const (
	DashboardOutputModeSnapshot DashboardOutputMode = iota
	DashboardOutputModeSnapshotShort
	DashboardOutputModeNone
)

var DashboardOutputModeIds = map[DashboardOutputMode][]string{
	DashboardOutputModeSnapshot:      {constants.OutputFormatSnapshot},
	DashboardOutputModeSnapshotShort: {constants.OutputFormatSnapshotShort},
	DashboardOutputModeNone:          {constants.OutputFormatNone},
}

type CheckOutputMode enumflag.Flag

const (
	CheckOutputModeText  CheckOutputMode = iota
	CheckOutputModeBrief CheckOutputMode = iota
	CheckOutputModeCsv
	CheckOutputModeHTM
	CheckOutputModeJSON
	CheckOutputModeMd
	CheckOutputModeTest
	CheckOutputModeSnapshot
	CheckOutputModeNone
)

var CheckOutputModeIds = map[CheckOutputMode][]string{
	CheckOutputModeText:     {constants.OutputFormatText},
	CheckOutputModeBrief:    {constants.OutputFormatBrief},
	CheckOutputModeCsv:      {constants.OutputFormatCSV},
	CheckOutputModeHTM:      {constants.OutputFormatHTML},
	CheckOutputModeJSON:     {constants.OutputFormatJSON},
	CheckOutputModeMd:       {constants.OutputFormatMD},
	CheckOutputModeTest:     {constants.OutputFormatText},
	CheckOutputModeSnapshot: {constants.OutputFormatSnapshot},
	CheckOutputModeNone:     {constants.OutputFormatNone},
}

func FlagValues[T comparable](mappings map[T][]string) []string {
	var res = make([]string, 0, len(mappings))
	for _, v := range mappings {
		res = append(res, v[0])
	}
	return res

}
