package controlexecute

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/turbot/pipe-fittings/v2/steampipeconfig"
	"github.com/turbot/powerpipe/internal/controlstatus"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
)

func testRows() ResultRows {
	return ResultRows{
		{Reason: "bucket <b> is public & open", Resource: "arn:aws:s3:::b1", Status: "alarm",
			Dimensions: []Dimension{{Key: "region", Value: "us-east-1", SqlType: "TEXT"}, {Key: "account_id", Value: "123", SqlType: "TEXT"}}},
		{Reason: "ok", Resource: "arn:aws:s3:::b2", Status: "ok",
			Dimensions: []Dimension{{Key: "region", Value: "eu-west-1", SqlType: "TEXT"}}},
	}
}

func TestControlRunData_MarshalJSON(t *testing.T) {
	run := &ControlRun{Rows: testRows()[1:]}

	got, err := json.Marshal(newControlRunData(run))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"columns":[{"name":"reason","data_type":"TEXT"},{"name":"resource","data_type":"TEXT"},{"name":"status","data_type":"TEXT"},{"name":"region","data_type":"TEXT"}],` +
		`"rows":[{"reason":"ok","region":"eu-west-1","resource":"arn:aws:s3:::b2","status":"ok"}]}`
	if string(got) != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
	if run.DimensionKeys[0] != "region" {
		t.Fatalf("expected the dimension keys to be set, got %v", run.DimensionKeys)
	}
}

// legacyStrippedJson is how a snapshot was exported before control rows were written one at a time
func legacyStrippedJson(s *steampipeconfig.SteampipeSnapshot, indent bool) ([]byte, error) {
	res, err := s.AsCloudSnapshot()
	if err != nil {
		return nil, err
	}
	if err := steampipeconfig.StripSnapshot(res); err != nil {
		return nil, err
	}
	if indent {
		return json.MarshalIndent(res, "", "  ")
	}
	return json.Marshal(res)
}

func TestControlRun_SnapshotRowsMatchLegacy(t *testing.T) {
	newRun := func(name string, rows ResultRows) *ControlRun {
		run := &ControlRun{
			FullName: name, Title: "Title <" + name + ">", NodeType: "control", Tags: map[string]string{"service": "AWS/S3"},
			Properties: map[string]any{"sql": "select 1", "severity": "high"},
			Summary:    &controlstatus.StatusSummary{Alarm: 1, Ok: 1}, RunStatus: dashboardtypes.RunComplete, Rows: rows,
		}
		run.Data = newControlRunData(run)
		return run
	}
	notRun := &ControlRun{FullName: "mod.control.not_run", NodeType: "control", RunStatus: dashboardtypes.RunError, RunErrorString: "query failed"}
	snapshot := &steampipeconfig.SteampipeSnapshot{
		SchemaVersion: "20240607",
		Panels: map[string]steampipeconfig.SnapshotPanel{
			"mod.control.rows":    newRun("mod.control.rows", testRows()),
			"mod.control.no_rows": newRun("mod.control.no_rows", ResultRows{}),
			"mod.control.not_run": notRun,
		},
	}
	for _, indent := range []bool{false, true} {
		want, err := legacyStrippedJson(snapshot, indent)
		if err != nil {
			t.Fatal(err)
		}
		got, err := snapshot.AsStrippedJson(indent)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("indent=%v: streamed snapshot differs from legacy\n got: %s\nwant: %s", indent, got, want)
		}
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
	run.Data = newControlRunData(run)
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
