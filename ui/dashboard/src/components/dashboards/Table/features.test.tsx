import { act, renderHook } from "@testing-library/react";
import { ColumnFilter, useTable } from "@tanstack/react-table";
import { getColumns, getData } from "./index";
import { KeyValuePairs } from "../common/types";
import { LeafNodeDataColumn, LeafNodeDataRow } from "../common";
import { Filter } from "@powerpipe/components/dashboards/grouping/common";
import { tablePanelFeatures, TablePanelFeatures } from "./features";

// The table panel's own columns and rows, built the way the panel builds them.
const dataColumns: LeafNodeDataColumn[] = [
  { name: "instance_id", data_type: "TEXT" },
  { name: "region", data_type: "TEXT" },
  { name: "cpu_count", data_type: "INT8" },
];

const dataRows: LeafNodeDataRow[] = [
  { instance_id: "i-10", region: "us-east-1", cpu_count: 4 },
  { instance_id: "i-9", region: "eu-west-2", cpu_count: 16 },
  { instance_id: "i-100", region: "us-east-1", cpu_count: 2 },
];

const renderPanelTable = (
  filters: Filter[],
  properties: KeyValuePairs = {},
) => {
  const { columns, columnVisibility } = getColumns(
    dataColumns,
    properties,
    filters,
    false,
  );
  const data = getData(columns, dataRows);
  const columnFilters: ColumnFilter[] = filters.map((filter) => ({
    id: filter.key as string,
    value: filter.value,
  }));
  return renderHook(() =>
    useTable<TablePanelFeatures, KeyValuePairs>({
      features: tablePanelFeatures,
      data,
      columns,
      initialState: { columnVisibility },
      state: { columnFilters },
    }),
  );
};

const instanceIds = (table) =>
  table.getRowModel().rows.map((row) => row.original.instance_id);

describe("table panel on TanStack Table v9", () => {
  test("sorts a short string column as plain strings, both directions", () => {
    const { result } = renderPanelTable([]);
    expect(instanceIds(result.current)).toEqual(["i-10", "i-9", "i-100"]);

    const instanceColumn = result.current.getColumn("instance_id")!;
    expect(instanceColumn.getCanSort()).toBe(true);

    act(() => instanceColumn.getToggleSortingHandler()!({}));
    expect(result.current.getColumn("instance_id")!.getIsSorted()).toBe("asc");
    expect(instanceIds(result.current)).toEqual(["i-10", "i-100", "i-9"]);

    act(() => instanceColumn.toggleSorting(true));
    expect(instanceIds(result.current)).toEqual(["i-9", "i-100", "i-10"]);
  });

  // Expected orders are TanStack Table v8.21.3's output for the same values.
  const sortedNames = (names: (string | null)[], desc: boolean) => {
    const { columns } = getColumns(
      [{ name: "name", data_type: "TEXT" }],
      {},
      [],
      false,
    );
    const { result } = renderHook(() =>
      useTable<TablePanelFeatures, KeyValuePairs>({
        features: tablePanelFeatures,
        data: getData(
          columns,
          names.map((name) => ({ name })),
        ),
        columns,
        state: { sorting: [{ id: "name", desc }] },
      }),
    );
    return result.current.getRowModel().rows.map((row) => row.original.name);
  };

  test("sorts text starting with a digit, mixed case, empty and null as v8 did", () => {
    const short = [
      "banana",
      "10 apples",
      "Apple",
      "2 pears",
      "zucchini",
      "",
      null,
      "apple",
    ];
    expect(sortedNames(short, false)).toEqual([
      null,
      "",
      "10 apples",
      "2 pears",
      "Apple",
      "apple",
      "banana",
      "zucchini",
    ]);
    expect(sortedNames(short, true)).toEqual([
      "zucchini",
      "banana",
      "apple",
      "Apple",
      "2 pears",
      "10 apples",
      "",
      null,
    ]);

    // Past ten rows v8 picked the sort from the later rows' values.
    const long = [...short, "cherry", "date", "item 2", "item 10"];
    expect(sortedNames(long, false)).toEqual([
      "",
      null,
      "Apple",
      "apple",
      "banana",
      "cherry",
      "date",
      "item 2",
      "item 10",
      "zucchini",
      "2 pears",
      "10 apples",
    ]);
    expect(sortedNames(long, true)).toEqual([
      "10 apples",
      "2 pears",
      "zucchini",
      "item 10",
      "item 2",
      "date",
      "cherry",
      "banana",
      "Apple",
      "apple",
      "",
      null,
    ]);
  });

  test("applies the panel's column filters and shows no rows when none match", () => {
    const matching = renderPanelTable([
      { operator: "equal", key: "region", value: "us-east-1" },
    ]);
    expect(instanceIds(matching.result.current)).toEqual(["i-10", "i-100"]);

    const excluding = renderPanelTable([
      { operator: "not_equal", key: "region", value: "us-east-1" },
    ]);
    expect(instanceIds(excluding.result.current)).toEqual(["i-9"]);

    const none = renderPanelTable([
      { operator: "equal", key: "region", value: "ap-south-1" },
    ]);
    expect(none.result.current.getRowModel().rows).toHaveLength(0);
  });

  test("hides columns the panel hides and shows them again when toggled", () => {
    const { result } = renderPanelTable([], {
      columns: { cpu_count: { display: "none" } },
    });
    const visibleCellIds = () =>
      result.current
        .getRowModel()
        .rows[0].getVisibleCells()
        .map((cell) => cell.column.id);
    expect(visibleCellIds()).toEqual(["instance_id", "region"]);

    const cpuColumn = result.current
      .getAllLeafColumns()
      .find((column) => column.id === "cpu_count")!;
    expect(cpuColumn.getIsVisible()).toBe(false);

    act(() => cpuColumn.toggleVisibility(true));
    expect(visibleCellIds()).toEqual(["instance_id", "region", "cpu_count"]);
  });
});
