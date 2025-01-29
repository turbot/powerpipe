package dashboardexecute

import (
	"sort"
	"sync"
	"time"

	"github.com/turbot/powerpipe/internal/controlexecute"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/schema"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/controlstatus"
)

const RootResultGroup_Name = "root_result_group"

// DetectionBenchmarkDisplay is a struct representing a grouping of control results
// It may correspond to a Benchmark, or some other arbitrary grouping
// TODO ultimately just use benchmark

type SummaryProvider interface {
	SetSummary()
	GetSummary() *DetectionBenchmarkSummary
}

type DetectionBenchmarkDisplay struct {
	GroupId       string            `json:"name" csv:"group_id"`
	Title         string            `json:"title,omitempty" csv:"title"`
	Description   string            `json:"description,omitempty" csv:"description"`
	Tags          map[string]string `json:"tags,omitempty"`
	Documentation string            `json:"documentation,omitempty"`
	Display       string            `json:"display,omitempty"`
	Type          string            `json:"type,omitempty"`

	// the overall summary of the group
	Summary *DetectionBenchmarkSummary `json:"summary"`
	// child result groups
	Groups []*DetectionBenchmarkDisplay `json:"groups"`
	// child runs
	DetectionRuns []*DetectionRun                        `json:"detection_runs"` // list of children stored as controlexecute.ExecutionTreeNode
	Children      []controlexecute.ExecutionTreeNode     `json:"-"`
	Severity      map[string]controlstatus.StatusSummary `json:"-"`
	// "benchmark"
	NodeType string `json:"panel_type"`
	// the control tree item associated with this group(i.e. a mod/benchmark)
	GroupItem modconfig.ModTreeItem      `json:"-"`
	Parent    *DetectionBenchmarkDisplay `json:"-"`
	Duration  time.Duration              `json:"-"`

	// a list of distinct dimension keys from descendant controls
	DimensionKeys []string `json:"-"`

	//childrenComplete   uint32
	//executionStartTime time.Time
	// lock to prevent multiple control_runs updating this
	updateLock *sync.Mutex
}

type DetectionBenchmarkSummary struct {
	Count int `json:"count"`
}

func NewDetectionBenchmarkSummary() *DetectionBenchmarkSummary {
	return &DetectionBenchmarkSummary{}
}

// NewRootBenchmarkDisplay creates a DetectionBenchmarkDisplay to act as the root node of a control execution tree
func NewRootBenchmarkDisplay(rootItem modconfig.ModTreeItem) (*DetectionBenchmarkDisplay, error) {
	root := &DetectionBenchmarkDisplay{
		GroupId:    RootResultGroup_Name,
		Groups:     []*DetectionBenchmarkDisplay{},
		Tags:       make(map[string]string),
		Summary:    NewDetectionBenchmarkSummary(),
		Severity:   make(map[string]controlstatus.StatusSummary),
		updateLock: new(sync.Mutex),
		NodeType:   schema.BlockTypeBenchmark,
		Title:      rootItem.GetTitle(),
	}

	return root, nil
}

// NewDetectionBenchmarkDisplay creates a result group from a ModTreeItem
func NewDetectionBenchmarkDisplay(benchmarkRun *DetectionBenchmarkRun, parent *DetectionBenchmarkDisplay) (*DetectionBenchmarkDisplay, error) {
	group := &DetectionBenchmarkDisplay{
		GroupId:     benchmarkRun.Name,
		Title:       benchmarkRun.GetTitle(),
		Description: benchmarkRun.resource.GetDescription(),
		Tags:        benchmarkRun.resource.GetTags(),
		GroupItem:   benchmarkRun.resource,
		Parent:      parent,
		Groups:      []*DetectionBenchmarkDisplay{},
		Summary:     NewDetectionBenchmarkSummary(),
		Severity:    make(map[string]controlstatus.StatusSummary),
		updateLock:  new(sync.Mutex),
		NodeType:    schema.BlockTypeBenchmark,
	}

	//add child groups for children which are benchmarks
	for _, c := range benchmarkRun.GetChildren() {
		switch child := c.(type) {
		case *DetectionBenchmarkRun:
			// create a result group for this item
			benchmarkGroup, err := NewDetectionBenchmarkDisplay(child, group)
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

// SetSummary populates the summary (which just counts the number of detections)
func (r *DetectionBenchmarkDisplay) SetSummary() {

	for _, c := range r.Children {
		if summaryProvider, ok := c.(SummaryProvider); ok {
			summaryProvider.SetSummary()
			r.Summary.Count = summaryProvider.GetSummary().Count
		}
	}
	for _, run := range r.DetectionRuns {
		r.Summary.Count = len(run.Data.Rows)
	}
}

func (r *DetectionBenchmarkDisplay) GetSummary() *DetectionBenchmarkSummary {
	return r.Summary
}

func (r *DetectionBenchmarkDisplay) AllTagKeys() []string {
	var tags []string
	for k := range r.Tags {
		tags = append(tags, k)
	}
	for _, child := range r.Groups {
		tags = append(tags, child.AllTagKeys()...)
	}
	for _, run := range r.DetectionRuns {
		for k := range run.resource.GetTags() {
			tags = append(tags, k)
		}
	}
	tags = helpers.StringSliceDistinct(tags)
	sort.Strings(tags)
	return tags
}

// GetGroupByName finds an immediate child DetectionBenchmarkDisplay with a specific name
func (r *DetectionBenchmarkDisplay) GetGroupByName(name string) *DetectionBenchmarkDisplay {
	for _, group := range r.Groups {
		if group.GroupId == name {
			return group
		}
	}
	return nil
}

// GetChildGroupByName finds a nested child DetectionBenchmarkDisplay with a specific name
func (r *DetectionBenchmarkDisplay) GetChildGroupByName(name string) *DetectionBenchmarkDisplay {
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

// GetRunByName finds a child ControlRun with a specific control name
func (r *DetectionBenchmarkDisplay) GetRunByName(name string) *DetectionRun {
	for _, run := range r.DetectionRuns {
		if run.resource.Name() == name {
			return run
		}
	}
	return nil
}

func (r *DetectionBenchmarkDisplay) RunCount() int {
	count := len(r.DetectionRuns)
	for _, g := range r.Groups {
		count += g.RunCount()
	}
	return count
}

// IsSnapshotPanel implements SnapshotPanel
func (*DetectionBenchmarkDisplay) IsSnapshotPanel() {}

// IsExecutionTreeNode implements ExecutionTreeNode
func (*DetectionBenchmarkDisplay) IsExecutionTreeNode() {}

// GetChildren implements ExecutionTreeNode
func (r *DetectionBenchmarkDisplay) GetChildren() []controlexecute.ExecutionTreeNode {
	return r.Children
}

// GetName implements ExecutionTreeNode
func (r *DetectionBenchmarkDisplay) GetName() string { return r.GroupId }

// AsTreeNode implements ExecutionTreeNode
func (r *DetectionBenchmarkDisplay) AsTreeNode() *steampipeconfig.SnapshotTreeNode {
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

// AddResultGroup adds result group into our list, and also add a tree node into our child list
func (r *DetectionBenchmarkDisplay) AddResultGroup(group *DetectionBenchmarkDisplay) {
	r.Groups = append(r.Groups, group)
	r.Children = append(r.Children, group)
}

// AddDetection add run into our list, and also add a tree node into our child list
func (r *DetectionBenchmarkDisplay) AddDetection(detectionRun *DetectionRun) {
	r.DetectionRuns = append(r.DetectionRuns, detectionRun)
	r.Children = append(r.Children, detectionRun)
}
