package controldisplay

import (
	"fmt"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"log/slog"
	"strings"

	"github.com/turbot/powerpipe/internal/controlexecute"
	"github.com/turbot/powerpipe/internal/resources"
)

type DetectionGroupRenderer struct {
	group *dashboardexecute.DetectionBenchmarkDisplay
	// screen width
	width      int
	resultTree *dashboardexecute.DetectionBenchmarkDisplayTree
	lastChild  bool
	parent     *DetectionGroupRenderer
}

func NewDetectionGroupRenderer(group *dashboardexecute.DetectionBenchmarkDisplay, parent *DetectionGroupRenderer, resultTree *dashboardexecute.DetectionBenchmarkDisplayTree, width int) *DetectionGroupRenderer {
	r := &DetectionGroupRenderer{
		group:      group,
		parent:     parent,
		resultTree: resultTree,

		width: width,
	}
	r.lastChild = r.isLastChild(group)
	return r
}

// are we the last child of our parent?
// this affects the tree rendering
func (r DetectionGroupRenderer) isLastChild(group *dashboardexecute.DetectionBenchmarkDisplay) bool {
	if group.Parent == nil || group.Parent.GroupItem == nil {
		return true
	}
	siblings := group.Parent.GroupItem.GetChildren()
	// get the name of the last sibling which has controls (or is a control)
	var finalSiblingName string
	for _, s := range siblings {
		if b, ok := s.(*resources.ControlBenchmark); ok {
			// find the result group for this benchmark and see if it has controls
			resultGroup := r.resultTree.Root.GetChildGroupByName(b.Name())
			// if the result group has not controls, we will not find it in the result tree
			if resultGroup == nil || resultGroup.RunCount() == 0 {
				continue
			}
		}
		// store the name of this sibling
		finalSiblingName = s.Name()
	}

	res := group.GroupItem.Name() == finalSiblingName

	return res
}

// the indent for blank lines
// same as for (not last) children
func (r DetectionGroupRenderer) blankLineIndent() string {
	return r.childIndent()
}

// the indent for group heading
func (r DetectionGroupRenderer) headingIndent() string {
	// if this is the first displayed node, no indent
	if r.parent == nil || r.parent.group.GroupId == controlexecute.RootResultGroupName {
		return ""
	}
	// as our parent for the indent for a group
	i := r.parent.childGroupIndent()
	return i
}

// the indent for child groups/controls (which are not the final child)
// include the tree '|'
func (r DetectionGroupRenderer) childIndent() string {
	return r.parentIndent() + "| "
}

// the indent for the FINAL child groups/controls
// just a space
func (r DetectionGroupRenderer) lastChildIndent() string {
	return r.parentIndent() + "  "
}

// the indent for child groups - our parent indent with the group expander "+ "
func (r DetectionGroupRenderer) childGroupIndent() string {
	return r.parentIndent() + "+ "
}

// get the indent inherited from our parent
// - this will depend on whether we are our parents last child
func (r DetectionGroupRenderer) parentIndent() string {
	if r.parent == nil || r.parent.group.GroupId == controlexecute.RootResultGroupName {
		return ""
	}
	if r.lastChild {
		return r.parent.lastChildIndent()
	}
	return r.parent.childIndent()
}

func (r DetectionGroupRenderer) Render() string {
	if r.width <= 0 {
		// this should never happen, since the minimum width is set by the formatter
		slog.Warn("DetectionGroupRenderer.Render unexpected negative width", "width", r.width)

		return ""
	}

	if r.group.GroupId == controlexecute.RootResultGroupName {
		return r.renderRootResultGroup()
	}

	groupHeadingRenderer := NewDetectionGroupHeadingRenderer(
		r.group.Title,
		r.group.Summary.Count,
		r.width,
		r.headingIndent())

	// render this group header
	tableStrings := append([]string{},
		groupHeadingRenderer.Render(),
		// newline after group
		fmt.Sprintf("%s", ControlColors.Indent(r.blankLineIndent())))

	// now render the group children, in the order they are specified
	childStrings := r.renderChildren()
	tableStrings = append(tableStrings, childStrings...)
	return strings.Join(tableStrings, "\n")
}

// for root result group, there will either be one or more groups, or one or more control runs
// there will be no order specified so just loop through them
func (r DetectionGroupRenderer) renderRootResultGroup() string {
	var resultStrings = make([]string, len(r.group.Groups)+len(r.group.DetectionRuns))
	for i, group := range r.group.Groups {
		groupRenderer := NewDetectionGroupRenderer(group, &r, r.resultTree, r.width)
		resultStrings[i] = groupRenderer.Render()
	}
	for i, run := range r.group.DetectionRuns {
		controlRenderer := NewDetectionRenderer(run, &r)
		resultStrings[i] = controlRenderer.Render()
	}
	return strings.Join(resultStrings, "\n")
}

// render the children of this group, in the order they are specified in the hcl
func (r DetectionGroupRenderer) renderChildren() []string {
	children := r.group.GroupItem.GetChildren()
	var childStrings []string

	for _, child := range children {
		if detection, ok := child.(*resources.Detection); ok {
			// get Result group with a matching name
			if run := r.group.GetRunByName(detection.Name()); run != nil {
				detectionRender := NewDetectionRenderer(run, &r)
				childStrings = append(childStrings, detectionRender.Render())
			}
		} else {
			if childGroup := r.group.GetGroupByName(child.Name()); childGroup != nil {
				groupRenderer := NewDetectionGroupRenderer(childGroup, &r, r.resultTree, r.width)
				childStrings = append(childStrings, groupRenderer.Render())
			}
		}
	}

	return childStrings
}
