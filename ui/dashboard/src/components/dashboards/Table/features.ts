import {
  Column,
  columnFilteringFeature,
  columnVisibilityFeature,
  createFilteredRowModel,
  createSortedRowModel,
  reSplitAlphaNumeric,
  Row,
  rowSortingFeature,
  SortFn,
  sortFn_alphanumeric,
  sortFn_basic,
  sortFn_datetime,
  sortFn_text,
  tableFeatures,
} from "@tanstack/react-table";
import { KeyValuePairs } from "../common/types";

// The TanStack Table features the table panel uses. Since v9 a table only has
// the APIs of the features registered here (v8 bundled every feature).
export const tablePanelFeatures = tableFeatures({
  columnFilteringFeature,
  columnVisibilityFeature,
  rowSortingFeature,
  filteredRowModel: createFilteredRowModel(),
  sortedRowModel: createSortedRowModel(),
});

export type TablePanelFeatures = typeof tablePanelFeatures;

export type TablePanelColumn = Column<TablePanelFeatures, KeyValuePairs>;

export type TablePanelRow = Row<TablePanelFeatures, KeyValuePairs>;

type TablePanelSortFn = SortFn<TablePanelFeatures, KeyValuePairs>;

// v8's "auto" sort picked a column's sort function from the filtered rows
// AFTER the tenth, so a table of ten rows or fewer compared values as-is
// (text starting with a digit first). v9 samples the first ten rows instead,
// which reorders those tables; this keeps the order v8 gave.
const pickV8AutoSortFn = (
  rows: TablePanelRow[],
  columnId: string,
): TablePanelSortFn => {
  let isString = false;
  for (const row of rows.slice(10)) {
    const value = row.getValue(columnId);
    if (Object.prototype.toString.call(value) === "[object Date]") {
      return sortFn_datetime;
    }
    if (typeof value === "string") {
      isString = true;
      if (value.split(reSplitAlphaNumeric).length > 1) {
        return sortFn_alphanumeric;
      }
    }
  }
  return isString ? sortFn_text : sortFn_basic;
};

const v8AutoSortFns = new WeakMap<
  TablePanelRow[],
  Record<string, TablePanelSortFn>
>();

export const sortFn_v8Auto: TablePanelSortFn = (rowA, rowB, columnId) => {
  const rows = rowA.table.getFilteredRowModel().flatRows;
  let byColumn = v8AutoSortFns.get(rows);
  if (!byColumn) {
    byColumn = {};
    v8AutoSortFns.set(rows, byColumn);
  }
  if (!byColumn[columnId]) {
    byColumn[columnId] = pickV8AutoSortFn(rows, columnId);
  }
  return byColumn[columnId](rowA, rowB, columnId);
};
