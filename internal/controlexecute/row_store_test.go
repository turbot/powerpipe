package controlexecute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/turbot/pipe-fittings/v2/steampipeconfig"
	"github.com/turbot/powerpipe/internal/controlstatus"
	"github.com/turbot/powerpipe/internal/dashboardtypes"
)

func newTestRowStore(t *testing.T) *rowStore {
	t.Helper()
	store, err := newRowStore()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.close() })
	return store
}

// mixedRows has every status, out of status order, with awkward values
func mixedRows(n int) []*ResultRow {
	statuses := []string{"ok", "alarm", "skip", "error", "info"}
	var rows []*ResultRow
	for i := range n {
		row := &ResultRow{
			Reason:   fmt.Sprintf("row %d: <b> & \"quoted\"\nnew line é 日本  ", i),
			Resource: fmt.Sprintf("arn:aws:s3:::bucket-%d", i),
			Status:   statuses[i%len(statuses)],
		}
		if i%3 != 0 {
			row.Dimensions = append(row.Dimensions, Dimension{Key: "region", Value: "us-east-1", SqlType: "TEXT"})
		}
		if i%4 == 1 {
			row.Dimensions = append(row.Dimensions, Dimension{Key: "account_id", Value: fmt.Sprint(i), SqlType: "BIGINT"})
		}
		if i == 7 {
			row.Dimensions = append(row.Dimensions, Dimension{Key: "empty", Value: "", SqlType: "TEXT"})
		}
		rows = append(rows, row)
	}
	return rows
}

// newRunWithRows adds the rows to a new control run as a query would, keeping them on disk if store is set
func newRunWithRows(name string, rows []*ResultRow, store *rowStore) *ControlRun {
	run := &ControlRun{
		FullName: name, Title: "Title <" + name + ">", NodeType: "control",
		Properties: map[string]any{"sql": "select 1", "severity": "high"},
		Summary:    &controlstatus.StatusSummary{}, RunStatus: dashboardtypes.RunComplete,
		rowMap: make(map[string]ResultRows),
	}
	if store != nil {
		run.storedRows = newStoredRows(store)
	}
	for _, row := range rows {
		copied := *row
		copied.Run = run
		run.addResultRow(&copied)
	}
	run.createdOrderedResultRows()
	run.Data = newControlRunData(run)
	return run
}

func collectRows(t *testing.T, run *ControlRun) []ResultRow {
	t.Helper()
	var res []ResultRow
	if err := run.eachResultRow(func(row *ResultRow) error {
		if row.Run != run {
			t.Fatalf("row has the wrong run")
		}
		res = append(res, *row)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return res
}

func TestStoredRows_SameAsMemory(t *testing.T) {
	store := newTestRowStore(t)
	// 3,000 rows of ~150 bytes is several segments for each status
	rows := mixedRows(3000)
	inMemory := newRunWithRows("mod.control.memory", rows, nil)
	onDisk := newRunWithRows("mod.control.disk", rows, store)

	if len(onDisk.Rows) != 0 {
		t.Fatalf("rows on disk should not be in Rows, got %d", len(onDisk.Rows))
	}
	if len(onDisk.storedRows.segments["ok"]) < 2 {
		t.Fatalf("expected several segments, got %d", len(onDisk.storedRows.segments["ok"]))
	}
	want := collectRows(t, inMemory)
	got := collectRows(t, onDisk)
	if len(got) != len(rows) || !reflect.DeepEqual(got, rowsWithRun(want, onDisk)) {
		t.Fatalf("rows from disk differ from rows in memory")
	}
	if !reflect.DeepEqual(onDisk.DimensionKeys, inMemory.DimensionKeys) || len(inMemory.DimensionKeys) != 3 {
		t.Fatalf("dimension keys differ: disk %v, memory %v", onDisk.DimensionKeys, inMemory.DimensionKeys)
	}
	if *onDisk.Summary != *inMemory.Summary {
		t.Fatalf("summaries differ: disk %v, memory %v", onDisk.Summary, inMemory.Summary)
	}
}

func rowsWithRun(rows []ResultRow, run *ControlRun) []ResultRow {
	res := make([]ResultRow, len(rows))
	for i, row := range rows {
		row.Run = run
		res[i] = row
	}
	return res
}

// the JSON of a control run with rows on disk is the same as with rows in memory
func TestStoredRows_SameJSON(t *testing.T) {
	store := newTestRowStore(t)
	rows := mixedRows(500)
	inMemory := newRunWithRows("mod.control.c", rows, nil)
	onDisk := newRunWithRows("mod.control.c", rows, store)
	// the columns come from a map, so use the same order in both
	onDisk.Data.columns = inMemory.Data.columns

	wantJson, err := json.Marshal(inMemory)
	if err != nil {
		t.Fatal(err)
	}
	gotJson, err := json.Marshal(onDisk)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotJson, wantJson) {
		t.Fatalf("marshalled run differs\n got: %s\nwant: %s", gotJson, wantJson)
	}

	snapshotOf := func(run *ControlRun) *steampipeconfig.SteampipeSnapshot {
		return &steampipeconfig.SteampipeSnapshot{SchemaVersion: "20240607", Panels: map[string]steampipeconfig.SnapshotPanel{run.FullName: run}}
	}
	for _, indent := range []bool{false, true} {
		want, err := legacyStrippedJson(snapshotOf(inMemory), indent)
		if err != nil {
			t.Fatal(err)
		}
		got, err := snapshotOf(onDisk).AsStrippedJson(indent)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("indent=%v: snapshot with rows on disk differs from legacy", indent)
		}
	}
}

// like rows in memory, the rows of a control run that did not complete are not shown
func TestStoredRows_NotComplete(t *testing.T) {
	store := newTestRowStore(t)
	run := &ControlRun{Summary: &controlstatus.StatusSummary{}, storedRows: newStoredRows(store)}
	for _, row := range mixedRows(3000) {
		run.addResultRow(row)
	}
	if len(run.storedRows.segments["ok"]) == 0 {
		t.Fatalf("expected some rows to be written to the store")
	}
	if rows := collectRows(t, run); len(rows) != 0 {
		t.Fatalf("expected no rows before the run completes, got %d", len(rows))
	}
	if run.Summary.Ok != 600 {
		t.Fatalf("expected the summary to count the rows, got %v", run.Summary)
	}
}

// control runs write to the same store at the same time
func TestStoredRows_ConcurrentRuns(t *testing.T) {
	store := newTestRowStore(t)
	runs := make([]*ControlRun, 8)
	var wg sync.WaitGroup
	for i := range runs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runs[i] = newRunWithRows(fmt.Sprintf("mod.control.c%d", i), mixedRows(1000+i), store)
		}()
	}
	wg.Wait()
	for i, run := range runs {
		rows := collectRows(t, run)
		if len(rows) != 1000+i {
			t.Fatalf("run %d: expected %d rows, got %d", i, 1000+i, len(rows))
		}
		for _, row := range rows {
			if row.Run != run {
				t.Fatalf("run %d has a row of another run", i)
			}
		}
	}
}

func TestStoredRows_ReadError(t *testing.T) {
	store := newTestRowStore(t)
	run := newRunWithRows("mod.control.c", mixedRows(10), store)
	if err := store.close(); err != nil {
		t.Fatal(err)
	}
	if err := run.EachRow(func(any) error { return nil }); err == nil {
		t.Fatal("expected an error reading rows from a closed store")
	}
	if _, err := json.Marshal(run); err == nil {
		t.Fatal("expected an error marshalling a run whose rows cannot be read")
	}
}

func TestStoredRows_CorruptSegment(t *testing.T) {
	store := newTestRowStore(t)
	run := newRunWithRows("mod.control.c", mixedRows(10), store)
	// overwrite the start of the file
	if _, err := store.file.WriteAt([]byte("not json"), 0); err != nil {
		t.Fatal(err)
	}
	if err := run.EachRow(func(any) error { return nil }); err == nil {
		t.Fatal("expected an error for a corrupt row")
	}
}

func TestNewRowStore_RemovesFileName(t *testing.T) {
	store := newTestRowStore(t)
	if _, err := os.Stat(store.file.Name()); !os.IsNotExist(err) {
		t.Fatalf("expected the row file name to be removed, got %v", err)
	}
	if _, err := store.write([]byte("row\n")); err != nil {
		t.Fatalf("expected the open file to be writable: %v", err)
	}
}

func TestExecutionTree_LoadRows(t *testing.T) {
	tree := &ExecutionTree{ControlRuns: map[string]*ControlRun{}}
	run := &ControlRun{FullName: "mod.control.c", Summary: &controlstatus.StatusSummary{}, rowMap: make(map[string]ResultRows)}
	tree.ControlRuns[run.FullName] = run
	if err := tree.StoreRowsOnDisk(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tree.Close() })
	rows := mixedRows(20)
	for _, row := range rows {
		row.Run = run
		run.addResultRow(row)
	}
	run.createdOrderedResultRows()
	tree.DimensionColorGenerator, _ = NewDimensionColorGenerator(4, 27)

	if len(run.Rows) != 0 {
		t.Fatalf("expected no rows in memory before LoadRows")
	}
	if err := tree.LoadRows(); err != nil {
		t.Fatal(err)
	}
	if len(run.Rows) != len(rows) || run.Rows[0].Status != "error" {
		t.Fatalf("expected %d rows in status order, got %d", len(rows), len(run.Rows))
	}
	if _, ok := tree.DimensionColorGenerator.Map["region"]["us-east-1"]; !ok {
		t.Fatalf("expected the dimension colors to include the loaded rows")
	}
	// a second load does nothing
	if err := tree.LoadRows(); err != nil || len(run.Rows) != len(rows) {
		t.Fatalf("expected a second LoadRows to keep %d rows, got %d (%v)", len(rows), len(run.Rows), err)
	}
}

// with rows in memory, StoreRowsOnDisk was not called and LoadRows and Close do nothing
func TestExecutionTree_RowsInMemory(t *testing.T) {
	run := newRunWithRows("mod.control.c", mixedRows(5), nil)
	tree := &ExecutionTree{ControlRuns: map[string]*ControlRun{run.FullName: run}}
	if err := tree.LoadRows(); err != nil || len(run.Rows) != 5 {
		t.Fatalf("expected the 5 rows to stay, got %d (%v)", len(run.Rows), err)
	}
	if err := tree.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestStoredRows_TruncatedSegment(t *testing.T) {
	store := newTestRowStore(t)
	run := newRunWithRows("mod.control.c", mixedRows(10), store)
	// lose the final new line of a segment
	run.storedRows.segments["ok"][0].length--
	if err := run.EachRow(func(any) error { return nil }); err == nil {
		t.Fatal("expected an error for a truncated row")
	}
}
