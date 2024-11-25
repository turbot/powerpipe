package dashboardexecute

import (
	"github.com/turbot/powerpipe/internal/controlexecute"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/controlstatus"
)

const RootResultGroup_Name = "root_result_group"

// ResultGroup_SNAP is a struct representing a grouping of control results
// It may correspond to a Benchmark, or some other arbitrary grouping
// TODO ultimately just use benchmark
type ResultGroup_SNAP struct {
	GroupId       string            `json:"name" csv:"group_id"`
	Title         string            `json:"title,omitempty" csv:"title"`
	Description   string            `json:"description,omitempty" csv:"description"`
	Tags          map[string]string `json:"tags,omitempty"`
	Documentation string            `json:"documentation,omitempty"`
	Display       string            `json:"display,omitempty"`
	Type          string            `json:"type,omitempty"`

	// the overall summary of the group
	Summary *GroupSummary_ `json:"summary"`
	// child result groups
	Groups []*ResultGroup_SNAP `json:"-"`
	// child control runs
	// TODO K combine and use interface?? make tree and groups??
	ControlRuns   []*controlexecute.ControlRun `json:"-"`
	DetectionRuns []*DetectionRun              `json:"-"`
	// list of children stored as controlexecute.ExecutionTreeNode
	Children []controlexecute.ExecutionTreeNode     `json:"-"`
	Severity map[string]controlstatus.StatusSummary `json:"-"`
	// "benchmark"
	NodeType string `json:"panel_type"`
	// the control tree item associated with this group(i.e. a mod/benchmark)
	GroupItem modconfig.ModTreeItem `json:"-"`
	Parent    *ResultGroup_SNAP     `json:"-"`
	Duration  time.Duration         `json:"-"`

	// a list of distinct dimension keys from descendant controls
	DimensionKeys []string `json:"-"`

	childrenComplete   uint32
	executionStartTime time.Time
	// lock to prevent multiple control_runs updating this
	updateLock *sync.Mutex
}

type GroupSummary_ struct {
	Status   controlstatus.StatusSummary            `json:"status"`
	Severity map[string]controlstatus.StatusSummary `json:"-"`
}

func NewGroupSummary_() *GroupSummary_ {
	return &GroupSummary_{Severity: make(map[string]controlstatus.StatusSummary)}
}

// NewRootResultGroup_ creates a ResultGroup_SNAP to act as the root node of a control execution tree
func NewRootResultGroup_(rootItem modconfig.ModTreeItem) (*ResultGroup_SNAP, error) {
	root := &ResultGroup_SNAP{
		GroupId:    RootResultGroup_Name,
		Groups:     []*ResultGroup_SNAP{},
		Tags:       make(map[string]string),
		Summary:    NewGroupSummary_(),
		Severity:   make(map[string]controlstatus.StatusSummary),
		updateLock: new(sync.Mutex),
		NodeType:   schema.BlockTypeBenchmark,
		Title:      rootItem.GetTitle(),
	}

	return root, nil
}

// NewResultGroup_ creates a result group from a ModTreeItem
func NewResultGroup_(benchmarkRun *BenchmarkRun, parent *ResultGroup_SNAP) (*ResultGroup_SNAP, error) {
	group := &ResultGroup_SNAP{
		GroupId:     benchmarkRun.Name,
		Title:       benchmarkRun.GetTitle(),
		Description: benchmarkRun.resource.GetDescription(),
		Tags:        benchmarkRun.resource.GetTags(),
		GroupItem:   benchmarkRun.resource,
		Parent:      parent,
		Groups:      []*ResultGroup_SNAP{},
		Summary:     NewGroupSummary_(),
		Severity:    make(map[string]controlstatus.StatusSummary),
		updateLock:  new(sync.Mutex),
		NodeType:    schema.BlockTypeBenchmark,
	}

	//add child groups for children which are benchmarks
	for _, c := range benchmarkRun.GetChildren() {
		switch child := c.(type) {
		case *BenchmarkRun:
			// create a result group for this item
			benchmarkGroup, err := NewResultGroup_(child, group)
			if err != nil {
				return nil, err
			}
			// create a new result group with 'group' as the parent
			group.AddResultGroup(benchmarkGroup)

		case *DetectionRun:
			group.AddDetection(child)

		}
	}

	return group, nil
}

func (r *ResultGroup_SNAP) AllTagKeys() []string {
	var tags []string
	for k := range r.Tags {
		tags = append(tags, k)
	}
	for _, child := range r.Groups {
		tags = append(tags, child.AllTagKeys()...)
	}
	for _, run := range r.ControlRuns {
		for k := range run.Control.Tags {
			tags = append(tags, k)
		}
	}
	tags = helpers.StringSliceDistinct(tags)
	sort.Strings(tags)
	return tags
}

// GetGroupByName finds an immediate child ResultGroup_SNAP with a specific name
func (r *ResultGroup_SNAP) GetGroupByName(name string) *ResultGroup_SNAP {
	for _, group := range r.Groups {
		if group.GroupId == name {
			return group
		}
	}
	return nil
}

// GetChildGroupByName finds a nested child ResultGroup_SNAP with a specific name
func (r *ResultGroup_SNAP) GetChildGroupByName(name string) *ResultGroup_SNAP {
	for _, group := range r.Groups {
		if group.GroupId == name {
			return group
		}
		if child := group.GetChildGroupByName(name); child != nil {
			return child
		}
	}
	return nil
}

// GetControlRunByName finds a child ControlRun with a specific control name
func (r *ResultGroup_SNAP) GetControlRunByName(name string) *controlexecute.ControlRun {
	for _, run := range r.ControlRuns {
		if run.Control.Name() == name {
			return run
		}
	}
	return nil
}

func (r *ResultGroup_SNAP) ControlRunCount() int {
	count := len(r.ControlRuns)
	for _, g := range r.Groups {
		count += g.ControlRunCount()
	}
	return count
}

// IsSnapshotPanel implements SnapshotPanel
func (*ResultGroup_SNAP) IsSnapshotPanel() {}

// IsExecutionTreeNode implements ExecutionTreeNode
func (*ResultGroup_SNAP) IsExecutionTreeNode() {}

// GetChildren implements ExecutionTreeNode
func (r *ResultGroup_SNAP) GetChildren() []controlexecute.ExecutionTreeNode { return r.Children }

// GetName implements ExecutionTreeNode
func (r *ResultGroup_SNAP) GetName() string { return r.GroupId }

// AsTreeNode implements ExecutionTreeNode
func (r *ResultGroup_SNAP) AsTreeNode() *steampipeconfig.SnapshotTreeNode {
	res := &steampipeconfig.SnapshotTreeNode{
		Name:     r.GroupId,
		Children: make([]*steampipeconfig.SnapshotTreeNode, len(r.Children)),
		NodeType: r.NodeType,
	}
	for i, c := range r.Children {
		res.Children[i] = c.AsTreeNode()
	}
	return res
}

// add result group into our list, and also add a tree node into our child list
func (r *ResultGroup_SNAP) AddResultGroup(group *ResultGroup_SNAP) {
	r.Groups = append(r.Groups, group)
	r.Children = append(r.Children, group)
}

// add control into our list, and also add a tree node into our child list
func (r *ResultGroup_SNAP) AddControl(controlRun *controlexecute.ControlRun) {
	r.ControlRuns = append(r.ControlRuns, controlRun)
	r.Children = append(r.Children, controlRun)
}
func (r *ResultGroup_SNAP) AddDetection(detectionRun *DetectionRun) {
	r.DetectionRuns = append(r.DetectionRuns, detectionRun)
	r.Children = append(r.Children, detectionRun)
}

func (r *ResultGroup_SNAP) addDimensionKeys(keys ...string) {
	r.updateLock.Lock()
	defer r.updateLock.Unlock()
	r.DimensionKeys = append(r.DimensionKeys, keys...)
	if r.Parent != nil {
		r.Parent.addDimensionKeys(keys...)
	}
	r.DimensionKeys = helpers.StringSliceDistinct(r.DimensionKeys)
	sort.Strings(r.DimensionKeys)
}

// onChildDone is a callback that gets called from the children of this result group when they are done
func (r *ResultGroup_SNAP) onChildDone() {
	newCount := atomic.AddUint32(&r.childrenComplete, 1)
	totalCount := uint32(len(r.ControlRuns) + len(r.Groups)) //nolint:gosec // will not overflow
	if newCount < totalCount {
		// all children haven't finished execution yet
		return
	}

	// all children are done
	r.Duration = time.Since(r.executionStartTime)
	if r.Parent != nil {
		r.Parent.onChildDone()
	}
}

func (r *ResultGroup_SNAP) updateSummary(summary *controlstatus.StatusSummary) {
	r.updateLock.Lock()
	defer r.updateLock.Unlock()

	r.Summary.Status.Skip += summary.Skip
	r.Summary.Status.Alarm += summary.Alarm
	r.Summary.Status.Info += summary.Info
	r.Summary.Status.Ok += summary.Ok
	r.Summary.Status.Error += summary.Error

	if r.Parent != nil {
		r.Parent.updateSummary(summary)
	}
}

func (r *ResultGroup_SNAP) updateSeverityCounts(severity string, summary *controlstatus.StatusSummary) {
	r.updateLock.Lock()
	defer r.updateLock.Unlock()

	val, exists := r.Severity[severity]
	if !exists {
		val = controlstatus.StatusSummary{}
	}
	val.Alarm += summary.Alarm
	val.Error += summary.Error
	val.Info += summary.Info
	val.Ok += summary.Ok
	val.Skip += summary.Skip

	r.Summary.Severity[severity] = val
	if r.Parent != nil {
		r.Parent.updateSeverityCounts(severity, summary)
	}
}
