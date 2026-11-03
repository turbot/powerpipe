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
  test("sorts a column the way v8's auto sort did, both directions", () => {
    const { result } = renderPanelTable([]);
    expect(instanceIds(result.current)).toEqual(["i-10", "i-9", "i-100"]);

    const instanceColumn = result.current.getColumn("instance_id")!;
    expect(instanceColumn.getCanSort()).toBe(true);

    // A string column sorts ascending first, in natural (alphanumeric) order
    act(() => instanceColumn.getToggleSortingHandler()!({}));
    expect(result.current.getColumn("instance_id")!.getIsSorted()).toBe("asc");
    expect(instanceIds(result.current)).toEqual(["i-9", "i-10", "i-100"]);

    act(() => instanceColumn.toggleSorting(true));
    expect(instanceIds(result.current)).toEqual(["i-100", "i-10", "i-9"]);
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
