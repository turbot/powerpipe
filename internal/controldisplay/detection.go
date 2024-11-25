package controldisplay

import (
	"fmt"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"strings"

	"github.com/spf13/viper"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/powerpipe/internal/controlexecute"
)

type DetectionRenderer struct {
	run               *dashboardexecute.DetectionRun
	parent            *GroupRenderer
	maxFailedControls int
	maxTotalControls  int
	// screen width
	width          int
	colorGenerator *controlexecute.DimensionColorGenerator
	lastChild      bool
}

func NewDetectionRenderer(run *dashboardexecute.DetectionRun, parent *GroupRenderer) *DetectionRenderer {
	r := &DetectionRenderer{
		run:               run,
		parent:            parent,
		maxFailedControls: parent.maxFailedControls,
		maxTotalControls:  parent.maxTotalControls,
		colorGenerator:    parent.resultTree.DimensionColorGenerator,
		width:             parent.width,
	}
	r.lastChild = r.isLastChild()
	return r
}

// are we the last child of our parent?
// this affects the tree rendering
func (r DetectionRenderer) isLastChild() bool {
	if r.parent.group == nil || r.parent.group.GroupItem == nil {
		return true
	}
	siblings := r.parent.group.GroupItem.GetChildren()
	return r.run.Resource.Name() == siblings[len(siblings)-1].Name()
}

// get the indent inherited from our parent
// - this will depend on whether we are our parents last child
func (r DetectionRenderer) parentIndent() string {
	if r.lastChild {
		return r.parent.lastChildIndent()
	}
	return r.parent.childIndent()
}

// indent before first result
func (r DetectionRenderer) preResultIndent() string {
	// when we do not have any rows, do not add a '|' to indent
	if viper.GetBool(constants.ArgDryRun) || len(r.run.Data.Rows) == 0 {
		return r.parentIndent()
	}
	return r.parentIndent() + "| "
}

// indent before first result
func (r DetectionRenderer) resultIndent() string {
	return r.parentIndent()
}

// indent after last result
func (r DetectionRenderer) postResultIndent() string {
	return r.parentIndent()
}

func (r DetectionRenderer) Render() string {
	var controlStrings []string
	failedCount := len(r.run.Data.Rows)
	// use group heading renderer to render the control title and counts
	// TODO group heading renderer suitable for detection or need a special one?
	controlHeadingRenderer := NewGroupHeadingRenderer(typehelpers.SafeString(r.run.Resource.GetTitle()),
		failedCount,
		failedCount,
		r.maxFailedControls,
		r.maxTotalControls,
		r.width,
		r.parent.childGroupIndent())

	// set the severity on the heading renderer
	controlHeadingRenderer.severity = typehelpers.SafeString(r.run.Resource.Severity)

	// get formatted indents
	formattedPostResultIndent := fmt.Sprintf("%s", ControlColors.Indent(r.postResultIndent()))
	formattedPreResultIndent := fmt.Sprintf("%s", ControlColors.Indent(r.preResultIndent()))

	controlStrings = append(controlStrings,
		controlHeadingRenderer.Render(),
		// newline after control heading
		formattedPreResultIndent)

	// if the control is in error, render an error
	if r.run.GetError() != nil {
		errorRenderer := NewErrorRenderer(r.run.GetError(), r.width, r.parentIndent())
		controlStrings = append(controlStrings,
			errorRenderer.Render(),
			// newline after error
			formattedPostResultIndent)
	}

	// now render the results (if any)
	var resultStrings []string
	for _, row := range r.run.Data.Rows {
		// build dimension
		var dimensions = make([]controlexecute.Dimension, len(r.run.Data.Columns))
		for i, c := range r.run.Data.Columns {
			dimensions[i] = controlexecute.Dimension{
				Key:     c.Name,
				SqlType: c.DataType,
				Value:   typehelpers.ToString(row[c.Name]),
			}
		}
		resultRenderer := NewDetectionResultRenderer(
			dimensions,
			r.colorGenerator,
			r.width,
			r.resultIndent())
		// the result renderer may not render the result - in quiet mode only failures are rendered
		if resultString := resultRenderer.Render(); resultString != "" {
			resultStrings = append(resultStrings, resultString)
		}
	}

	// newline after results
	if len(resultStrings) > 0 {
		controlStrings = append(controlStrings, resultStrings...)
		if len(r.run.Data.Rows) > 0 || r.run.GetError() != nil {
			controlStrings = append(controlStrings, formattedPostResultIndent)
		}
	}

	return strings.Join(controlStrings, "\n")
}
