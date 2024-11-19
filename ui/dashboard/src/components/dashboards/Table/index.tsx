import ControlDimension from "../grouping/Benchmark/ControlDimension";
import Icon from "@powerpipe/components/Icon";
import isEmpty from "lodash/isEmpty";
import isObject from "lodash/isObject";
import useDeepCompareEffect from "use-deep-compare-effect";
import useGroupingFilterConfig from "@powerpipe/hooks/useGroupingFilterConfig";
import useTemplateRender from "@powerpipe/hooks/useTemplateRender";
import {
  AlarmIcon,
  InfoIcon,
  OKIcon,
  SkipIcon,
  UnknownIcon,
  ErrorIcon,
  SortAscendingIcon,
  SortDescendingIcon,
} from "@powerpipe/constants/icons";
import {
  BasePrimitiveProps,
  ExecutablePrimitiveProps,
  isNumericCol,
  LeafNodeDataColumn,
  LeafNodeDataRow,
} from "../common";
import { CheckFilter } from "@powerpipe/components/dashboards/grouping/common";
import { classNames } from "@powerpipe/utils/styles";
import {
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { getComponent, registerComponent } from "../index";
import { injectSearchPathPrefix } from "@powerpipe/utils/url";
import { memo, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { PanelDefinition } from "@powerpipe/types";
import { KeyValuePairs, RowRenderResult } from "../common/types";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useSearchParams } from "react-router-dom";
import { useVirtualizer } from "@tanstack/react-virtual";

const ExternalLink = getComponent("external_link");

export type TableColumnDisplay = "all" | "none";
export type TableColumnWrap = "all" | "none";

type TableColumnInfo = {
  header: string;
  title: string;
  accessorKey: string;
  name: string;
  data_type: string;
  display?: "all" | "none";
  wrap: TableColumnWrap;
  href_template?: string;
  sortType?: any;
};

const getColumns = (
  cols: LeafNodeDataColumn[],
  properties?: TableProperties,
): {
  columns: TableColumnInfo[];
  columnVisibility: {
    [key: string]: boolean;
  };
} => {
  if (!cols || cols.length === 0) {
    return { columns: [], columnVisibility: {} };
  }

  const columnVisibility: {
    [key: string]: boolean;
  } = {};
  const columns: TableColumnInfo[] = cols.map((col) => {
    let colHref: string | null = null;
    let colWrap: TableColumnWrap = "none";
    if (
      properties &&
      properties.columns &&
      properties.columns[col.original_name || col.name]
    ) {
      const c = properties.columns[col.original_name || col.name];
      if (c.display === "none") {
        columnVisibility[col.name] = false;
      }
      if (c.wrap) {
        colWrap = c.wrap as TableColumnWrap;
      }
      if (c.href) {
        colHref = c.href;
      }
    }

    const colInfo: TableColumnInfo = {
      header: col.original_name || col.name,
      title: col.original_name || col.name,
      accessorKey: col.name,
      name: col.name,
      data_type: col.data_type,
      wrap: colWrap,
      sortType: col.data_type === "BOOL" ? "basic" : "alphanumeric",
    };
    if (colHref) {
      colInfo.href_template = colHref;
    }
    return colInfo;
  });
  return { columns, columnVisibility };
};

const getData = (columns: TableColumnInfo[], rows: LeafNodeDataRow[]) => {
  if (!columns || columns.length === 0) {
    return [];
  }

  if (!rows || rows.length === 0) {
    return [];
  }
  return rows;
};

type CellValueProps = {
  column: TableColumnInfo;
  rowIndex: number;
  rowTemplateData: RowRenderResult[];
  value: any;
  showTitle?: boolean;
  addFilter?: (
    operator: "equal" | "not_equal",
    key: string,
    value: any,
    context?: string,
  ) => void;
  filterEnabled?: boolean;
  context?: string;
};

const CellValue = ({
  column,
  rowIndex,
  rowTemplateData,
  value,
  showTitle = false,
  addFilter,
  filterEnabled = false,
  context = "",
}: CellValueProps) => {
  const { searchPathPrefix } = useDashboard();
  const [href, setHref] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const renderedTemplateObj = rowTemplateData[rowIndex];

    if (!renderedTemplateObj) {
      setHref(() => null);
      setError(() => null);
      return;
    }
    const renderedTemplateForColumn = renderedTemplateObj[column.name];

    if (!renderedTemplateForColumn) {
      setHref(null);
      setError(null);
      return;
    }
    if (renderedTemplateForColumn.result) {
      const withSearchPathPrefix = injectSearchPathPrefix(
        renderedTemplateForColumn.result,
        searchPathPrefix,
      );
      setHref(withSearchPathPrefix);
      setError(null);
    } else if (renderedTemplateForColumn.error) {
      setHref(null);
      setError(renderedTemplateForColumn.error);
    }
  }, [column, rowIndex, rowTemplateData, searchPathPrefix]);

  let cellContent;
  if (value === null || value === undefined) {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className="link-highlight"
        title={showTitle ? `${column.title}=null` : undefined}
      >
        <>null</>
      </ExternalLink>
    ) : (
      <span
        className="text-foreground-lightest"
        title={showTitle ? `${column.title}=null` : undefined}
      >
        <>null</>
      </span>
    );
  }

  const dataType = column.data_type.toLowerCase();
  if (dataType === "control_status") {
    switch (value) {
      case "alarm":
        cellContent = (
          <span title="Status = Alarm">
            <AlarmIcon className="text-alert w-5 h-5" />
          </span>
        );
        break;
      case "error":
        cellContent = (
          <span title="Status = Error">
            <AlarmIcon className="text-alert w-5 h-5" />
          </span>
        );
        break;
      case "ok":
        cellContent = (
          <span title="Status = OK">
            <OKIcon className="text-ok w-5 h-5" />
          </span>
        );
        break;
      case "info":
        cellContent = (
          <span title="Status = Info">
            <InfoIcon className="text-info w-5 h-5" />
          </span>
        );
        break;
      case "skip":
        cellContent = (
          <span title="Status = Skipped">
            <SkipIcon className="text-skip w-5 h-5" />
          </span>
        );
        break;
      default:
        cellContent = (
          <span title="Status = Unknown">
            <UnknownIcon className="text-foreground-light w-5 h-5" />
          </span>
        );
    }
  } else if (dataType === "control_dimensions") {
    cellContent = (
      <div className="space-x-2">
        {(value || []).map((dimension) => (
          <ControlDimension
            key={dimension.key}
            dimensionKey={dimension.key}
            dimensionValue={dimension.value}
          />
        ))}
      </div>
    );
  } else if (dataType === "bool") {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className="link-highlight"
        title={showTitle ? `${column.title}=${value.toString()}` : undefined}
      >
        <>{value.toString()}</>
      </ExternalLink>
    ) : (
      <span
        className={classNames(value ? null : "text-foreground-light")}
        title={showTitle ? `${column.title}=${value.toString()}` : undefined}
      >
        <>{value.toString()}</>
      </span>
    );
  } else if (dataType === "jsonb" || isObject(value)) {
    const asJsonString = JSON.stringify(value, null, 2);
    cellContent = href ? (
      <ExternalLink
        to={href}
        className="link-highlight"
        title={showTitle ? `${column.title}=${asJsonString}` : undefined}
      >
        <>{asJsonString}</>
      </ExternalLink>
    ) : (
      <span title={showTitle ? `${column.title}=${asJsonString}` : undefined}>
        {asJsonString}
      </span>
    );
  } else if (dataType === "text") {
    if (!!value.match && value.match("^https?://")) {
      cellContent = (
        <ExternalLink
          className="link-highlight tabular-nums"
          to={value}
          title={showTitle ? `${column.title}=${value}` : undefined}
        >
          {value}
        </ExternalLink>
      );
    }
    const mdMatch =
      !!value.match && value.match("^\\[(.*)\\]\\((https?://.*)\\)$");
    if (mdMatch) {
      cellContent = (
        <ExternalLink
          className="tabular-nums"
          to={mdMatch[2]}
          title={showTitle ? `${column.title}=${value}` : undefined}
        >
          {mdMatch[1]}
        </ExternalLink>
      );
    }
  } else if (dataType === "timestamp" || dataType === "timestamptz") {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className="link-highlight tabular-nums"
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value}
      </ExternalLink>
    ) : (
      <span
        className="tabular-nums"
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value}
      </span>
    );
  } else if (isNumericCol(dataType)) {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className="link-highlight tabular-nums"
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value.toLocaleString()}
      </ExternalLink>
    ) : (
      <span
        className="tabular-nums"
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value.toLocaleString()}
      </span>
    );
  }
  if (!cellContent) {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className="link-highlight tabular-nums"
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value}
      </ExternalLink>
    ) : (
      <span
        className="tabular-nums"
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value}
      </span>
    );
  }

  return error ? (
    <span className="flex items-center space-x-2" title={error}>
      {cellContent} <ErrorIcon className="inline h-4 w-4 text-alert" />
    </span>
  ) : (
    <div className="flex items-center space-x-2 group">
      {cellContent}
      {/*{filterEnabled && addFilter && (*/}
      {/*  <div className="flex items-center space-x-1 opacity-0 group-hover:opacity-100">*/}
      {/*    <button*/}
      {/*      onClick={() => addFilter("equal", column.name, value, context)}*/}
      {/*      className="text-black-scale-7 hover:text-black-scale-8 focus:outline-none"*/}
      {/*      title="Add value to include filter"*/}
      {/*    >*/}
      {/*      <Icon className="h-5 w-5" icon="add_circle" />*/}
      {/*    </button>*/}
      {/*    <button*/}
      {/*      onClick={() => addFilter("not_equal", column.name, value, context)}*/}
      {/*      className="text-black-scale-7 hover:text-black-scale-8 focus:outline-none"*/}
      {/*      title="Add value to exclude filter"*/}
      {/*    >*/}
      {/*      <Icon className="h-5 w-5" icon="do_not_disturb_on" />*/}
      {/*    </button>*/}
      {/*  </div>*/}
      {/*)}*/}
    </div>
  );
};

const MemoCellValue = memo(CellValue);

type TableColumnOptions = {
  display?: TableColumnDisplay;
  href?: string;
  wrap?: TableColumnWrap;
};

type TableColumns = {
  [column: string]: TableColumnOptions;
};

type TableType = "table" | "line" | null;

export type TableProperties = {
  columns?: TableColumns;
};

export type TableProps = PanelDefinition &
  BasePrimitiveProps &
  ExecutablePrimitiveProps & {
    display_type?: TableType;
    properties?: TableProperties;
    filterEnabled?: boolean;
    context?: string;
  };

const useTableFilters = () => {
  const urlFilters = useGroupingFilterConfig();
  const [searchParams, setSearchParams] = useSearchParams();
  const expressions = urlFilters.expressions;
  const filters: CheckFilter[] = [];

  for (const expression of expressions || []) {
    if (
      expression.operator === "equal" &&
      expression.type === "dimension" &&
      !!expression.key &&
      !!expression.value
    ) {
      filters.push(expression);
    } else if (
      expression.operator === "not_equal" &&
      expression.type === "dimension" &&
      !!expression.key &&
      !!expression.value
    ) {
      filters.push(expression);
    }
  }

  const addFilter = useCallback(
    (
      operator: "equal" | "not_equal",
      key: string,
      value: any,
      context?: string,
    ) => {
      const index = urlFilters.expressions?.findIndex(
        (e) =>
          e.type === "dimension" &&
          e.key === key &&
          e.value === value &&
          e.context === context,
      );
      let newFilters =
        index !== undefined && index > -1
          ? [
              ...urlFilters.expressions?.slice(0, index),
              ...urlFilters.expressions?.slice(index + 1),
            ]
          : urlFilters.expressions || [];
      if (
        newFilters.length === 1 &&
        newFilters[0].operator === "equal" &&
        !newFilters[0].type
      ) {
        newFilters = [
          {
            operator,
            value,
            type: "dimension",
            key,
            title: value,
            context,
          },
        ];
      } else {
        newFilters.push({
          operator,
          value,
          type: "dimension",
          key,
          title: value,
          context,
        });
      }
      urlFilters.expressions = newFilters;
      searchParams.set("where", JSON.stringify(urlFilters));
      setSearchParams(searchParams);
    },
    [urlFilters, searchParams, setSearchParams],
  );

  const removeFilter = useCallback(
    (key: string, value: any, context: string) => {
      const index = urlFilters.expressions?.findIndex(
        (e) =>
          e.type === "dimension" &&
          e.key === key &&
          e.value === value &&
          e.context === context,
      );
      const newFilters =
        index !== undefined
          ? [
              ...urlFilters.expressions?.slice(0, index),
              ...urlFilters.expressions?.slice(index + 1),
            ]
          : urlFilters.expressions || [];
      if (newFilters.length === 0) {
        urlFilters.expressions = [{ operator: "equal" }];
      } else {
        urlFilters.expressions = newFilters;
      }
      searchParams.set("where", JSON.stringify(urlFilters));
      setSearchParams(searchParams);
    },
    [urlFilters, searchParams, setSearchParams],
  );

  return {
    filters,
    addFilter,
    removeFilter,
  };
};

// const TableView = ({
//   data,
//   columns,
//   columnVisibility,
//   hasTopBorder = false,
//   filterEnabled = false,
//   context = "",
// }) => {
//   const { filters, addFilter, removeFilter } = useTableFilters();
//   const { ready: templateRenderReady, renderTemplates } = useTemplateRender();
//   const [rowTemplateData, setRowTemplateData] = useState<RowRenderResult[]>([]);
//
//   const table = useReactTable({
//     data,
//     columns,
//     initialState: { columnVisibility },
//     getCoreRowModel: getCoreRowModel(),
//     getSortedRowModel: getSortedRowModel(),
//     debugTable: true,
//   });
//
//   const { rows } = table.getRowModel();
//
//   const visibleColumns = table.getVisibleLeafColumns();
//
//   //The virtualizers need to know the scrollable container element
//   const tableContainerRef = useRef<HTMLDivElement>(null);
//
//   //we are using a slightly different virtualization strategy for columns (compared to virtual rows) in order to support dynamic row heights
//   const columnVirtualizer = useVirtualizer({
//     count: visibleColumns.length,
//     estimateSize: (index) => 100, // visibleColumns[index].getSize(), //estimate width of each column for accurate scrollbar dragging
//     getScrollElement: () => tableContainerRef.current,
//     horizontal: true,
//     overscan: 3, //how many columns to render on each side off screen each way (adjust this for performance)
//   });
//
//   //dynamic row height virtualization - alternatively you could use a simpler fixed row height strategy without the need for `measureElement`
//   const rowVirtualizer = useVirtualizer({
//     count: rows.length,
//     estimateSize: () => 46.5, //estimate row height for accurate scrollbar dragging
//     getScrollElement: () => tableContainerRef.current,
//     //measure dynamic row height, except in firefox because it measures table border height incorrectly
//     measureElement:
//       typeof window !== "undefined" &&
//       navigator.userAgent.indexOf("Firefox") === -1
//         ? (element) => element?.getBoundingClientRect().height
//         : undefined,
//     overscan: 5,
//   });
//
//   const virtualColumns = columnVirtualizer.getVirtualItems();
//   const virtualRows = rowVirtualizer.getVirtualItems();
//
//   //different virtualization strategy for columns - instead of absolute and translateY, we add empty columns to the left and right
//   let virtualPaddingLeft: number | undefined;
//   let virtualPaddingRight: number | undefined;
//
//   if (columnVirtualizer && virtualColumns?.length) {
//     virtualPaddingLeft = virtualColumns[0]?.start ?? 0;
//     virtualPaddingRight =
//       columnVirtualizer.getTotalSize() -
//       (virtualColumns[virtualColumns.length - 1]?.end ?? 0);
//   }
//
//   //All important CSS styles are included as inline styles for this example. This is not recommended for your code.
//   return (
//     <div
//       className="overflow-auto relative min-h-[100px] h-[400px] max-h-[800px] w-full"
//       ref={tableContainerRef}
//     >
//       {/* Even though we're still using semantic table tags, we must use CSS grid and flexbox for dynamic row heights */}
//       <table
//         className={classNames(
//           "grid min-w-full divide-y divide-table-divide",
//           hasTopBorder ? "border-t border-divide" : null,
//         )}
//       >
//         <thead className="grid sticky top-0 z-[1] text-table-head border-b border-divide">
//           {table.getHeaderGroups().map((headerGroup) => (
//             <tr key={headerGroup.id} className="flex w-full">
//               {virtualPaddingLeft ? (
//                 //fake empty column to the left for virtualization scroll padding
//                 <th
//                   className={classNames("flex", `w-[${virtualPaddingLeft}]`)}
//                 />
//               ) : null}
//               {virtualColumns.map((vc) => {
//                 const header = headerGroup.headers[vc.index];
//                 return (
//                   <th
//                     key={header.id}
//                     className={classNames(
//                       "flex py-3 text-left text-sm font-normal tracking-wider whitespace-nowrap pl-4",
//                       //`w-[${header.getSize()}]`,
//                     )}
//                   >
//                     <div
//                       {...{
//                         className: header.column.getCanSort()
//                           ? "cursor-pointer select-none"
//                           : "",
//                         onClick: header.column.getToggleSortingHandler(),
//                       }}
//                     >
//                       {flexRender(
//                         header.column.columnDef.header,
//                         header.getContext(),
//                       )}
//                       {{
//                         asc: (
//                           <SortAscendingIcon
//                             className={classNames("inline-block h-4 w-4")}
//                           />
//                         ),
//                         desc: (
//                           <SortDescendingIcon className="inline-block h-4 w-4" />
//                         ),
//                       }[header.column.getIsSorted() as string] ?? null}
//                     </div>
//                   </th>
//                 );
//               })}
//               {virtualPaddingRight ? (
//                 //fake empty column to the right for virtualization scroll padding
//                 <th
//                   className={classNames("flex", `w-[${virtualPaddingRight}]`)}
//                 />
//               ) : null}
//             </tr>
//           ))}
//         </thead>
//         <tbody
//           className={classNames(
//             "grid relative divide-y divide-table-divide",
//             `h-[${rowVirtualizer.getTotalSize()}px]`, //tells scrollbar how big the table is
//           )}
//         >
//           {virtualRows.map((virtualRow, rowIndex) => {
//             const row = rows[virtualRow.index] as Row<KeyValuePairs>;
//             const visibleCells = row.getVisibleCells();
//
//             return (
//               <tr
//                 key={row.id}
//                 data-index={virtualRow.index} //needed for dynamic row height measurement
//                 ref={(node) => rowVirtualizer.measureElement(node)} //measure dynamic row height
//                 className="flex absolute w-full"
//                 style={{
//                   transform: `translateY(${virtualRow.start}px)`, //this should always be a `style` as it changes on scroll
//                 }}
//               >
//                 {virtualPaddingLeft ? (
//                   //fake empty column to the left for virtualization scroll padding
//                   <td style={{ display: "flex", width: virtualPaddingLeft }} />
//                 ) : null}
//                 {virtualColumns.map((vc) => {
//                   const cell = visibleCells[vc.index];
//                   return (
//                     <td
//                       key={cell.id}
//                       className={classNames(
//                         "flex px-4 py-4 align-top content-center text-sm",
//                         isNumericCol(cell.column.columnDef.data_type)
//                           ? "text-right"
//                           : "",
//                         cell.column.columnDef.wrap === "all"
//                           ? "break-keep"
//                           : "whitespace-nowrap",
//                       )}
//                       // style={{
//                       //   width: cell.column.getSize(),
//                       // }}
//                     >
//                       {/*{flexRender(*/}
//                       {/*  cell.column.columnDef.cell,*/}
//                       {/*  cell.getContext(),*/}
//                       {/*)}*/}
//                       <MemoCellValue
//                         column={cell.column.columnDef}
//                         rowIndex={rowIndex}
//                         rowTemplateData={rowTemplateData}
//                         value={cell.getValue()}
//                         addFilter={addFilter}
//                         filterEnabled={filterEnabled}
//                         context={context}
//                       />
//                     </td>
//                   );
//                 })}
//                 {virtualPaddingRight ? (
//                   //fake empty column to the right for virtualization scroll padding
//                   <td style={{ display: "flex", width: virtualPaddingRight }} />
//                 ) : null}
//               </tr>
//             );
//           })}
//         </tbody>
//       </table>
//     </div>
//   );
// };

// const TableView2 = ({
//   data,
//   columns,
//   columnVisibility,
//   filterEnabled = false,
//   hasTopBorder = false,
//   context = "",
// }) => {
//   const table = useReactTable({
//     data,
//     columns,
//     initialState: { columnVisibility },
//     getCoreRowModel: getCoreRowModel(),
//     getSortedRowModel: getSortedRowModel(),
//     debugTable: true,
//   });
//
//   const { rows } = table.getRowModel();
//
//   const visibleColumns = table.getVisibleLeafColumns();
//
//   //The virtualizers need to know the scrollable container element
//   const tableContainerRef = useRef<HTMLDivElement>(null);
//
//   //we are using a slightly different virtualization strategy for columns (compared to virtual rows) in order to support dynamic row heights
//   const columnVirtualizer = useVirtualizer({
//     count: visibleColumns.length,
//     estimateSize: (index) => visibleColumns[index].getSize(), //estimate width of each column for accurate scrollbar dragging
//     getScrollElement: () => tableContainerRef.current,
//     horizontal: true,
//     overscan: 3, // how many columns to render on each side off screen each way (adjust this for performance)
//   });
//
//   //dynamic row height virtualization - alternatively you could use a simpler fixed row height strategy without the need for `measureElement`
//   const rowVirtualizer = useVirtualizer({
//     count: rows.length,
//     estimateSize: () => 46.5, //estimate row height for accurate scrollbar dragging
//     getScrollElement: () => tableContainerRef.current,
//     //measure dynamic row height, except in firefox because it measures table border height incorrectly
//     measureElement:
//       typeof window !== "undefined" &&
//       navigator.userAgent.indexOf("Firefox") === -1
//         ? (element) => element?.getBoundingClientRect().height
//         : undefined,
//     overscan: 5, // how many row to render on each side off screen each way (adjust this for performance)
//   });
//
//   const virtualColumns = columnVirtualizer.getVirtualItems();
//   const virtualRows = rowVirtualizer.getVirtualItems();
//
//   //different virtualization strategy for columns - instead of absolute and translateY, we add empty columns to the left and right
//   let virtualPaddingLeft: number | undefined;
//   let virtualPaddingRight: number | undefined;
//
//   if (columnVirtualizer && virtualColumns?.length) {
//     virtualPaddingLeft = virtualColumns[0]?.start ?? 0;
//     virtualPaddingRight =
//       columnVirtualizer.getTotalSize() -
//       (virtualColumns[virtualColumns.length - 1]?.end ?? 0);
//   }
//
//   console.log({ rows, data, visibleColumns, virtualColumns, virtualRows });
//
//   //All important CSS styles are included as inline styles for this example. This is not recommended for your code.
//   return (
//     <div
//       ref={tableContainerRef}
//       className="overflow-auto relative h-[800px] w-full"
//     >
//       {/* Even though we're still using sematic table tags, we must use CSS grid and flexbox for dynamic row heights */}
//       <table
//         className={classNames(
//           "grid divide-y divide-table-divide",
//           hasTopBorder ? "border-t border-divide" : null,
//         )}
//       >
//         <thead className="grid sticky top-0 z-[1] text-table-head border-b border-divide">
//           {table.getHeaderGroups().map((headerGroup) => (
//             <tr key={headerGroup.id} className="flex w-full">
//               {virtualPaddingLeft ? (
//                 //fake empty column to the left for virtualization scroll padding
//                 <th className={`flex w-[${virtualPaddingLeft}]`} />
//               ) : null}
//               {virtualColumns.map((vc) => {
//                 const header = headerGroup.headers[vc.index];
//                 return (
//                   <th
//                     key={header.id}
//                     className={classNames(
//                       "flex py-3 text-left text-sm font-normal tracking-wider whitespace-nowrap pl-4 w-auto",
//                       //`w-[${header.getSize()}]`,
//                     )}
//                   >
//                     <div
//                       className={classNames(
//                         header.column.getCanSort()
//                           ? "cursor-pointer select-none"
//                           : null,
//                       )}
//                       onClick={header.column.getToggleSortingHandler()}
//                     >
//                       {flexRender(
//                         header.column.columnDef.header,
//                         header.getContext(),
//                       )}
//                       {{
//                         asc: (
//                           <SortAscendingIcon
//                             className={classNames("inline-block h-4 w-4")}
//                           />
//                         ),
//                         desc: (
//                           <SortDescendingIcon className="inline-block h-4 w-4" />
//                         ),
//                       }[header.column.getIsSorted() as string] ?? null}
//                     </div>
//                   </th>
//                 );
//               })}
//               {virtualPaddingRight ? (
//                 //fake empty column to the right for virtualization scroll padding
//                 <th
//                   className={classNames("flex", `w-[${virtualPaddingRight}]`)}
//                 />
//               ) : null}
//             </tr>
//           ))}
//         </thead>
//         <tbody
//           className={classNames(
//             "grid relative",
//             `h-[${rowVirtualizer.getTotalSize()}px]`,
//           )}
//         >
//           {virtualRows.map((virtualRow) => {
//             const row = rows[virtualRow.index] as Row<KeyValuePairs>;
//             const visibleCells = row.getVisibleCells();
//
//             return (
//               <tr
//                 key={row.id}
//                 data-index={virtualRow.index} //needed for dynamic row height measurement
//                 ref={(node) => rowVirtualizer.measureElement(node)} //measure dynamic row height
//                 className={"flex absolute w-full"}
//                 style={{
//                   transform: `translateY(${virtualRow.start}px)`, //this should always be a `style` as it changes on scroll
//                 }}
//               >
//                 {virtualPaddingLeft ? (
//                   //fake empty column to the left for virtualization scroll padding
//                   <td
//                     className={classNames("flex", `w-[${virtualPaddingLeft}]`)}
//                   />
//                 ) : null}
//                 {virtualColumns.map((vc) => {
//                   const cell = visibleCells[vc.index];
//                   return (
//                     <td
//                       key={cell.id}
//                       className="flex px-4 py-4 align-top content-center text-sm w-auto"
//                       // style={{
//                       //   width: cell.column.getSize(),
//                       // }}
//                     >
//                       {flexRender(
//                         cell.column.columnDef.cell,
//                         cell.getContext(),
//                       )}
//                     </td>
//                   );
//                 })}
//                 {virtualPaddingRight ? (
//                   //fake empty column to the right for virtualization scroll padding
//                   <td
//                     className={classNames("flex", `w-[${virtualPaddingRight}]`)}
//                   />
//                 ) : null}
//               </tr>
//             );
//           })}
//         </tbody>
//       </table>
//     </div>
//   );
// };

const TableViewVirtualizedRows = ({
  data,
  columns,
  columnVisibility,
  hasTopBorder = false,
  filterEnabled = false,
  context = "",
}) => {
  const { filters, addFilter, removeFilter } = useTableFilters();
  const { ready: templateRenderReady, renderTemplates } = useTemplateRender();
  const [rowTemplateData, setRowTemplateData] = useState<RowRenderResult[]>([]);
  const parentRef = useRef<HTMLDivElement>(null);

  const table = useReactTable<KeyValuePairs>({
    data,
    columns,
    initialState: { columnVisibility },
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    //debugTable: true,
  });

  const { rows } = table.getRowModel();

  const virtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 46.5,
    overscan: 10,
  });

  const virtualizedRows = virtualizer.getVirtualItems();

  useEffect(() => {
    if (!templateRenderReady || columns.length === 0 || rows.length === 0) {
      setRowTemplateData([]);
      return;
    }

    const doRender = async () => {
      const templates = Object.fromEntries(
        columns
          .filter((col) => col.display !== "none" && !!col.href_template)
          .map((col) => [col.name, col.href_template as string]),
      );
      if (isEmpty(templates)) {
        setRowTemplateData([]);
        return;
      }
      const data = virtualizedRows.map((virtualRow) => {
        const row = rows[virtualRow.index];
        return row.original;
      });
      const renderedResults = await renderTemplates(templates, data);
      setRowTemplateData(renderedResults || []);
    };

    doRender();
  }, [columns, renderTemplates, rows, virtualizedRows, templateRenderReady]);

  return (
    <div>
      {filterEnabled &&
        filters.filter((f) => f.context === context).length > 0 && (
          <div className="p-4 pb-4 rounded shadow-sm flex flex-wrap gap-2">
            {filters.map((filter) => {
              return (
                <div
                  key={`${filter.operator}:${filter.key}:${filter.value}`}
                  className="flex items-center bg-black-scale-2 px-3 py-1 rounded-md space-x-2"
                >
                  <Icon
                    className="w-4 h-4"
                    icon={
                      filter.operator === "equal"
                        ? "add_circle"
                        : "do_not_disturb_on"
                    }
                  />
                  <span>{`${filter.key}: ${filter.value}`}</span>
                  <span
                    onClick={() =>
                      removeFilter(filter.key, filter.value, filter.context)
                    }
                    className="cursor-pointer text-black-scale-6 hover:text-black-scale-8 focus:outline-none"
                  >
                    <Icon className="w-4 h-4" icon="close" />
                  </span>
                </div>
              );
            })}
          </div>
        )}
      <div
        ref={parentRef}
        className="relative overflow-auto min-h-[46.5px] max-h-[800px]"
      >
        <div className={`h-[${virtualizer.getTotalSize()}px}]`}>
          <table
            className={classNames(
              "w-full divide-y divide-table-divide",
              hasTopBorder ? "border-t border-divide" : null,
            )}
          >
            <thead className="text-table-head border-b border-divide">
              {table.getHeaderGroups().map((headerGroup) => (
                <tr key={headerGroup.id}>
                  {headerGroup.headers.map((header) => {
                    return (
                      <th
                        key={header.id}
                        colSpan={header.colSpan}
                        scope="col"
                        className={classNames(
                          "py-3 text-left text-sm font-normal tracking-wider whitespace-nowrap pl-4",
                        )}
                        //style={{ width: header.getSize() }}
                      >
                        {header.isPlaceholder ? null : (
                          <div
                            {...{
                              className: header.column.getCanSort()
                                ? "cursor-pointer select-none"
                                : "",
                              onClick: header.column.getToggleSortingHandler(),
                            }}
                          >
                            {flexRender(
                              header.column.columnDef.header,
                              header.getContext(),
                            )}
                            {{
                              asc: (
                                <SortAscendingIcon
                                  className={classNames("inline-block h-4 w-4")}
                                />
                              ),
                              desc: (
                                <SortDescendingIcon className="inline-block h-4 w-4" />
                              ),
                            }[header.column.getIsSorted() as string] ?? null}
                          </div>
                        )}
                      </th>
                    );
                  })}
                </tr>
              ))}
            </thead>
            <tbody className="divide-y divide-table-divide">
              {rows.length === 0 && (
                <tr>
                  <td
                    className="px-4 py-4 align-top content-center text-sm italic whitespace-nowrap"
                    colSpan={columns.length}
                  >
                    No results
                  </td>
                </tr>
              )}
              {virtualizer.getVirtualItems().map((virtualRow, index) => {
                const row = rows[virtualRow.index];
                return (
                  <tr
                    key={row.id}
                    style={{
                      height: `${virtualRow.size}px`,
                      transform: `translateY(${
                        virtualRow.start - index * virtualRow.size
                      }px)`,
                    }}
                  >
                    {row.getVisibleCells().map((cell) => {
                      return (
                        <td
                          key={cell.id}
                          className={classNames(
                            "px-4 py-4 align-top content-center text-sm max-w-[500px] overflow-x-hidden",
                            isNumericCol(cell.column.columnDef.data_type)
                              ? "text-right"
                              : "",
                            cell.column.columnDef.wrap === "all"
                              ? "break-keep"
                              : "whitespace-nowrap",
                          )}
                        >
                          {/*{flexRender(*/}
                          {/*  cell.column.columnDef.cell,*/}
                          {/*  cell.getContext(),*/}
                          {/*)}*/}
                          <CellValue
                            column={cell.column.columnDef}
                            rowIndex={index}
                            rowTemplateData={rowTemplateData}
                            value={cell.getValue()}
                            addFilter={addFilter}
                            filterEnabled={filterEnabled}
                            context={context}
                          />
                        </td>
                      );
                    })}
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

const TableView3 = ({
  data,
  columns,
  columnVisibility,
  hasTopBorder = false,
  filterEnabled = false,
  context = "",
}) => {
  const { filters, addFilter, removeFilter } = useTableFilters();
  const { ready: templateRenderReady, renderTemplates } = useTemplateRender();
  const [rowTemplateData, setRowTemplateData] = useState<RowRenderResult[]>([]);

  // const filteredData = useMemo(() => {
  //   let filtered = rowData;
  //
  //   if (activeFilters.length > 0 || excludedFilters.length > 0) {
  //     filtered = filtered.filter((row) => {
  //       return (
  //         activeFilters.every(
  //           (filter) => row[filter.column] === filter.value,
  //         ) &&
  //         excludedFilters.every((filter) => row[filter.column] !== filter.value)
  //       );
  //     });
  //   }
  //
  //   return filtered;
  // }, [rowData, activeFilters, excludedFilters]);

  const table = useReactTable({
    columns,
    data,
    initialState: { columnVisibility },
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  const rows = table.getRowModel().rows;

  useDeepCompareEffect(() => {
    if (!templateRenderReady || columns.length === 0 || rows.length === 0) {
      setRowTemplateData([]);
      return;
    }

    const doRender = async () => {
      const templates = Object.fromEntries(
        columns
          .filter((col) => col.display !== "none" && !!col.href_template)
          .map((col) => [col.name, col.href_template as string]),
      );
      if (isEmpty(templates)) {
        setRowTemplateData([]);
        return;
      }
      const data = rows.map((row) => row.original);
      const renderedResults = await renderTemplates(templates, data);
      setRowTemplateData(renderedResults || []);
    };

    doRender();
  }, [columns, renderTemplates, rows, templateRenderReady]);

  return (
    <div>
      {filterEnabled &&
        filters.filter((f) => f.context === context).length > 0 && (
          <div className="p-4 pb-4 rounded shadow-sm flex flex-wrap gap-2">
            {filters.map((filter) => {
              return (
                <div
                  key={`${filter.operator}:${filter.key}:${filter.value}`}
                  className="flex items-center bg-black-scale-2 px-3 py-1 rounded-md space-x-2"
                >
                  <Icon
                    className="w-4 h-4"
                    icon={
                      filter.operator === "equal"
                        ? "add_circle"
                        : "do_not_disturb_on"
                    }
                  />
                  <span>{`${filter.key}: ${filter.value}`}</span>
                  <span
                    onClick={() =>
                      removeFilter(filter.key, filter.value, filter.context)
                    }
                    className="cursor-pointer text-black-scale-6 hover:text-black-scale-8 focus:outline-none"
                  >
                    <Icon className="w-4 h-4" icon="close" />
                  </span>
                </div>
              );
            })}
          </div>
        )}

      <table
        className={classNames(
          "min-w-full divide-y divide-table-divide overflow-hidden",
          hasTopBorder ? "border-t border-divide" : null,
        )}
      >
        <thead className="text-table-head border-b border-divide">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th
                  key={header.id}
                  scope="col"
                  className={classNames(
                    "py-3 text-left text-sm font-normal tracking-wider whitespace-nowrap pl-4",
                  )}
                >
                  {flexRender(
                    header.column.columnDef.header,
                    header.getContext(),
                  )}
                  {{
                    asc: (
                      <SortAscendingIcon
                        className={classNames("inline-block h-4 w-4")}
                      />
                    ),
                    desc: (
                      <SortDescendingIcon className="inline-block h-4 w-4" />
                    ),
                  }[header.column.getIsSorted() as string] ?? null}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody className="divide-y divide-table-divide">
          {rows.length === 0 && (
            <tr>
              <td
                className="px-4 py-4 align-top content-center text-sm italic whitespace-nowrap"
                colSpan={columns.length}
              >
                No results
              </td>
            </tr>
          )}
          {rows.map((row, index) => {
            return (
              <tr key={row.id}>
                {row.getVisibleCells().map((cell) => {
                  // console.log(cell.column.columnDef);
                  return (
                    <td
                      key={cell.id}
                      className={classNames(
                        "px-4 py-4 align-top content-center text-sm",
                        isNumericCol(cell.column.columnDef.data_type)
                          ? "text-right"
                          : "",
                        cell.column.columnDef.wrap === "all"
                          ? "break-keep"
                          : "whitespace-nowrap",
                      )}
                    >
                      {/*{flexRender(*/}
                      {/*  cell.column.columnDef.cell,*/}
                      {/*  cell.getContext(),*/}
                      {/*)}*/}
                      <MemoCellValue
                        column={cell.column.columnDef}
                        rowIndex={index}
                        rowTemplateData={rowTemplateData}
                        value={cell.getValue()}
                        addFilter={addFilter}
                        filterEnabled={filterEnabled}
                        context={context}
                      />
                    </td>
                  );
                })}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};

// TODO retain full width on mobile, no padding
const TableViewWrapper = (props: TableProps) => {
  const { columns, columnVisibility } = useMemo(
    () => getColumns(props.data ? props.data.columns : [], props.properties),
    [props.data, props.properties],
  );
  const rowData = useMemo(
    () => getData(columns, props.data ? props.data.rows : []),
    [columns, props.data],
  );

  // State for managing column visibility
  // const [visibleColumns, setVisibleColumns] = useState(
  //   columns.map((col) => ({
  //     ...col,
  //     visible: !hiddenColumns.includes(col.name),
  //   })),
  // );

  // Handler to toggle column visibility
  // const toggleColumnVisibility = (columnName) => {
  //   setVisibleColumns((prevColumns) =>
  //     prevColumns.map((col) =>
  //       col.name === columnName ? { ...col, visible: !col.visible } : col,
  //     ),
  //   );
  // };

  // Filter columns based on visibility state
  // const filteredColumns = useMemo(
  //   () => visibleColumns.filter((col) => col.visible),
  //   [visibleColumns],
  // );

  // Render column selection UI
  // const renderColumnSelector = () => (
  //   <div className="p-2 border-b mb-4">
  //     <label className="block font-bold mb-2">Select Columns to Display:</label>
  //     <div className="flex flex-wrap gap-2">
  //       {visibleColumns.map((col) => (
  //         <div key={col.name} className="flex items-center">
  //           <input
  //             type="checkbox"
  //             checked={col.visible}
  //             onChange={() => toggleColumnVisibility(col.name)}
  //             id={`toggle-${col.name}`}
  //           />
  //           <label htmlFor={`toggle-${col.name}`} className="ml-2">
  //             {col.title}
  //           </label>
  //         </div>
  //       ))}
  //     </div>
  //   </div>
  // );

  return props.data ? (
    <TableViewVirtualizedRows
      data={rowData}
      columns={columns} // Use filtered columns for the table
      columnVisibility={columnVisibility}
      hasTopBorder={!!props.title}
      filterEnabled={props.filterEnabled}
      context={props.context}
    />
  ) : null;
};

const LineView = (props: TableProps) => {
  const { ready: templateRenderReady, renderTemplates } = useTemplateRender();
  const [columns, setColumns] = useState<TableColumnInfo[]>([]);
  const [rows, setRows] = useState<LeafNodeDataRow[]>([]);
  const [rowTemplateData, setRowTemplateData] = useState<RowRenderResult[]>([]);

  useEffect(() => {
    if (!props.data || !props.data.columns || !props.data.rows) {
      setColumns([]);
      setRows([]);
      return;
    }
    const newColumns: TableColumnInfo[] = [];
    props.data.columns.forEach((col) => {
      const columnOverrides =
        props.properties?.columns &&
        props.properties.columns[col.original_name || col.name];
      const newColDef: TableColumnInfo = {
        ...col,
        header: col.original_name || col.name,
        title: col.original_name || col.name,
        accessorKey: col.name,
        display: columnOverrides?.display ? columnOverrides.display : "all",
        wrap: columnOverrides?.wrap ? columnOverrides.wrap : "none",
        href_template: columnOverrides?.href,
      };
      newColumns.push(newColDef);
    });

    setColumns(newColumns);
    setRows(props.data.rows);
  }, [props.data, props.properties]);

  useDeepCompareEffect(() => {
    if (!templateRenderReady || columns.length === 0 || rows.length === 0) {
      setRowTemplateData([]);
      return;
    }

    const doRender = async () => {
      const templates = Object.fromEntries(
        columns
          .filter((col) => col.display !== "none" && !!col.href_template)
          .map((col) => [col.name, col.href_template as string]),
      );
      if (isEmpty(templates)) {
        setRowTemplateData([]);
        return;
      }
      const renderedResults = await renderTemplates(templates, rows);
      setRowTemplateData(renderedResults);
    };

    doRender();
  }, [columns, renderTemplates, rows, templateRenderReady]);

  if (columns.length === 0 || rows.length === 0) {
    return null;
  }

  return (
    <div className="px-4 py-3 space-y-4">
      {rows.map((row, rowIndex) => {
        return (
          <div key={rowIndex} className="space-y-2">
            {columns.map((col) => {
              if (col.display === "none") {
                return null;
              }
              return (
                <div key={`${col.name}-${rowIndex}`}>
                  <span className="block text-sm text-table-head truncate">
                    {col.title}
                  </span>
                  <span
                    className={classNames(
                      "block",
                      col.wrap === "all" ? "break-keep" : "truncate",
                    )}
                  >
                    <MemoCellValue
                      column={col}
                      rowIndex={rowIndex}
                      rowTemplateData={rowTemplateData}
                      value={row[col.name]}
                      showTitle
                    />
                  </span>
                </div>
              );
            })}
          </div>
        );
      })}
    </div>
  );
};

const Table = (props: TableProps) => {
  if (props.display_type === "line") {
    return <LineView {...props} />;
  }
  return <TableViewWrapper {...props} />;
};

registerComponent("table", Table);

export default Table;

export { TableViewWrapper };
