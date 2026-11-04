package controlexecute

import (
	"encoding/json"
	"testing"

	"github.com/turbot/pipe-fittings/v2/queryresult"
	"github.com/turbot/powerpipe/internal/controlstatus"
)

func testRows() ResultRows {
	return ResultRows{
		{Reason: "bucket <b> is public & open", Resource: "arn:aws:s3:::b1", Status: "alarm",
			Dimensions: []Dimension{{Key: "region", Value: "us-east-1", SqlType: "TEXT"}, {Key: "account_id", Value: "123", SqlType: "TEXT"}}},
		{Reason: "ok", Resource: "arn:aws:s3:::b2", Status: "ok",
			Dimensions: []Dimension{{Key: "region", Value: "eu-west-1", SqlType: "TEXT"}}},
	}
}

// the JSON of ControlRunData must be the JSON of the LeafData it replaces
func TestControlRunData_MarshalsAsLeafData(t *testing.T) {
	run := &ControlRun{Rows: testRows()}
	schema := map[string]*queryresult.ColumnDef{"region": {Name: "region", DataType: "TEXT"}}

	got, err := json.Marshal(&ControlRunData{run: run, dimensionsSchema: schema})
	if err != nil {
		t.Fatal(err)
	}
	want, err := json.Marshal(run.Rows.ToLeafData(schema))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestControlRun_DataJSON(t *testing.T) {
	// a run whose results were never gathered (e.g. the query failed) has null data
	notRun := &ControlRun{FullName: "mod.control.c1", Summary: &controlstatus.StatusSummary{}}
	b, err := json.Marshal(notRun)
	if err != nil {
		t.Fatal(err)
	}
	var res map[string]any
	if err := json.Unmarshal(b, &res); err != nil {
		t.Fatal(err)
	}
	if v, ok := res["data"]; !ok || v != nil {
		t.Fatalf("expected data to be null, got %v (present=%v)", v, ok)
	}

	// a completed run marshals its rows, with dimensions flattened into each row
	run := &ControlRun{FullName: "mod.control.c1", Summary: &controlstatus.StatusSummary{}, Rows: testRows()}
	run.Data = &ControlRunData{run: run, dimensionsSchema: run.getDimensionSchema()}
	b, err = json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	var parsed struct {
		Data struct {
			Columns []map[string]string `json:"columns"`
			Rows    []map[string]string `json:"rows"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatal(err)
	}
	if len(parsed.Data.Rows) != 2 || len(parsed.Data.Columns) != 5 {
		t.Fatalf("expected 2 rows and 5 columns, got %d rows and %d columns", len(parsed.Data.Rows), len(parsed.Data.Columns))
	}
	first := parsed.Data.Rows[0]
	if first["reason"] != "bucket <b> is public & open" || first["region"] != "us-east-1" || first["account_id"] != "123" || first["status"] != "alarm" {
		t.Fatalf("unexpected first row %v", first)
	}
	if _, ok := parsed.Data.Rows[1]["account_id"]; ok {
		t.Fatalf("row without an account_id dimension should not have one: %v", parsed.Data.Rows[1])
	}
}

func TestExecutionTree_ControlRunInstances(t *testing.T) {
	g1, g2 := &ResultGroup{GroupId: "g1"}, &ResultGroup{GroupId: "g2"}
	shared := &ControlRun{FullName: "mod.control.shared", Rows: testRows(), Parents: []*ResultGroup{g1, g2}}
	single := &ControlRun{FullName: "mod.control.single", Rows: testRows()[:1], Parents: []*ResultGroup{g2}}
	tree := &ExecutionTree{ControlRuns: map[string]*ControlRun{shared.FullName: shared, single.FullName: single}}

	instances := tree.ControlRunInstances()
	if len(instances) != 3 {
		t.Fatalf("expected one instance per parent (3), got %d", len(instances))
	}
	// count instances by control and group (ControlRuns is a map, so the order is not fixed)
	counts := map[string]int{}
	for _, inst := range instances {
		counts[inst.FullName+"/"+inst.Group.GroupId] = len(inst.Rows)
		for _, row := range inst.Rows {
			if row.ControlRun.Group != inst.Group {
				t.Fatalf("row instance has the wrong group")
			}
		}
	}
	want := map[string]int{"mod.control.shared/g1": 2, "mod.control.shared/g2": 2, "mod.control.single/g2": 1}
	for k, n := range want {
		if counts[k] != n {
			t.Fatalf("expected %d rows for %s, got %v", n, k, counts)
		}
	}
	if (&ExecutionTree{}).ControlRunInstances() != nil {
		t.Fatalf("expected no instances for an empty tree")
	}
}
