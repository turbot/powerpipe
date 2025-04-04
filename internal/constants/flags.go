package constants

import (
	"github.com/thediveo/enumflag/v2"
	"github.com/turbot/pipe-fittings/v2/constants"
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
	QueryOutputModeNone
)

// powerpipe snapshot
const OutputFormatPpSnapshotShort = "pps"

var QueryOutputModeIds = map[QueryOutputMode][]string{
	QueryOutputModeCsv:           {constants.OutputFormatCSV},
	QueryOutputModeJson:          {constants.OutputFormatJSON},
	QueryOutputModeLine:          {constants.OutputFormatLine},
	QueryOutputModeSnapshot:      {constants.OutputFormatSnapshot},
	QueryOutputModeSnapshotShort: {OutputFormatPpSnapshotShort},
	QueryOutputModeTable:         {constants.OutputFormatTable},
	QueryOutputModeNone:          {constants.OutputFormatNone},
}

type DashboardOutputMode enumflag.Flag

const (
	DashboardOutputModeSnapshot DashboardOutputMode = iota
	DashboardOutputModeSnapshotShort
	DashboardOutputModeNone
)

var DashboardOutputModeIds = map[DashboardOutputMode][]string{
	DashboardOutputModeSnapshot:      {constants.OutputFormatSnapshot},
	DashboardOutputModeSnapshotShort: {OutputFormatPpSnapshotShort},
	DashboardOutputModeNone:          {constants.OutputFormatNone},
}

type CheckOutputMode enumflag.Flag

type DetectionOutputMode enumflag.Flag

const (
	CheckOutputModeText CheckOutputMode = iota
	CheckOutputModeBrief
	CheckOutputModeCsv
	CheckOutputModeHTM
	CheckOutputModeJSON
	CheckOutputModeMd
	CheckOutputModeTest
	CheckOutputModeSnapshot
	CheckOutputModeSnapshotShort
	CheckOutputModeNone
)

const (
	DetectionOutputModeText DetectionOutputMode = iota
	DetectionOutputModeJSON
)

var CheckOutputModeIds = map[CheckOutputMode][]string{
	CheckOutputModeText:          {constants.OutputFormatText},
	CheckOutputModeBrief:         {constants.OutputFormatBrief},
	CheckOutputModeCsv:           {constants.OutputFormatCSV},
	CheckOutputModeHTM:           {constants.OutputFormatHTML},
	CheckOutputModeJSON:          {constants.OutputFormatJSON},
	CheckOutputModeMd:            {constants.OutputFormatMD},
	CheckOutputModeTest:          {constants.OutputFormatText},
	CheckOutputModeSnapshot:      {constants.OutputFormatSnapshot},
	CheckOutputModeSnapshotShort: {OutputFormatPpSnapshotShort},
	CheckOutputModeNone:          {constants.OutputFormatNone},
}

var DetectionOutputModeIds = map[DetectionOutputMode][]string{
	DetectionOutputModeText: {constants.OutputFormatText},
	DetectionOutputModeJSON: {constants.OutputFormatJSON},
}
