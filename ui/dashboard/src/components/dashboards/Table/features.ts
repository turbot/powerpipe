import {
  Column,
  columnFilteringFeature,
  columnVisibilityFeature,
  createFilteredRowModel,
  createSortedRowModel,
  Row,
  rowSortingFeature,
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
  // The sort functions a column's "auto" sort picks from. v9 only picks
  // registered ones, so register all four v8 chose between.
  sortFns: {
    alphanumeric: sortFn_alphanumeric,
    basic: sortFn_basic,
    datetime: sortFn_datetime,
    text: sortFn_text,
  },
});

export type TablePanelFeatures = typeof tablePanelFeatures;

export type TablePanelColumn = Column<TablePanelFeatures, KeyValuePairs>;

export type TablePanelRow = Row<TablePanelFeatures, KeyValuePairs>;
