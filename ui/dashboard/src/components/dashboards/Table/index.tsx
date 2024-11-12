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
  ColumnDiffState,
  ExecutablePrimitiveProps,
  isNumericCol,
  LeafNodeDataColumn,
  LeafNodeDataRow,
} from "../common";
import { CheckFilter } from "@powerpipe/components/dashboards/grouping/common";
import { classNames } from "@powerpipe/utils/styles";
import { injectSearchPathPrefix } from "@powerpipe/utils/url";
import { memo, useCallback, useEffect, useMemo, useState } from "react";
import { getComponent, registerComponent } from "../index";
import { PanelDefinition } from "@powerpipe/types";
import { RowRenderResult } from "../common/types";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useSearchParams } from "react-router-dom";
import { useSortBy, useTable } from "react-table";
import { ThemeNames } from "@powerpipe/hooks/useTheme";
import { tableRowDiffColumn } from "@powerpipe/utils/data";

export type TableColumnDisplay = "all" | "none";
export type TableColumnWrap = "all" | "none";

type TableColumnInfo = {
  Header: string;
  title: string;
  accessor: string;
  name: string;
  data_type: string;
  display?: "all" | "none";
  wrap: TableColumnWrap;
  href_template?: string;
  sortType?: any;
  diff: ColumnDiffState;
};

const getColumns = (
  cols: LeafNodeDataColumn[],
  properties?: TableProperties,
): { columns: TableColumnInfo[]; hiddenColumns: string[] } => {
  if (!cols || cols.length === 0) {
    return { columns: [], hiddenColumns: [] };
  }

  const hiddenColumns: string[] = [];
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
        hiddenColumns.push(col.name);
      }
      if (c.wrap) {
        colWrap = c.wrap as TableColumnWrap;
      }
      if (c.href) {
        colHref = c.href;
      }
    }

    const colInfo: TableColumnInfo = {
      Header: col.original_name || col.name,
      title: col.original_name || col.name,
      accessor: col.name,
      name: col.name,
      data_type: col.data_type,
      wrap: colWrap,
      sortType: col.data_type === "BOOL" ? "basic" : "alphanumeric",
      diff: col.__diff || "none",
    };
    if (colHref) {
      colInfo.href_template = colHref;
    }
    return colInfo;
  });
  return { columns, hiddenColumns };
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
  diffValue: any;
  showTitle?: boolean;
  addFilter?: (
    operator: "equal" | "not_equal",
    key: string,
    value: any,
    context?: string,
  ) => void;
  filterEnabled: boolean;
  context?: string;
};

const CellValue = ({
  column,
  rowIndex,
  rowTemplateData,
  value,
  diffValue,
  showTitle = false,
  addFilter,
  filterEnabled,
  context = "",
}: CellValueProps) => {
  const ExternalLink = getComponent("external_link");
  const { searchPathPrefix } = useDashboard();
  const [href, setHref] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const diffValueSpacer =
    diffValue || typeof diffValue === "boolean"
      ? "flex items-center space-x-2"
      : undefined;

  useEffect(() => {
    const renderedTemplateObj = rowTemplateData[rowIndex];

    if (!renderedTemplateObj) {
      setHref(null);
      setError(null);
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
  }, [column, rowIndex, rowTemplateData]);

  let cellContent;
  const dataType = column.data_type.toLowerCase();
  if (value === null || value === undefined) {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames("link-highlight", diffValueSpacer)}
        title={showTitle ? `${column.title}=null` : undefined}
      >
        <span>null</span>
        {diffValue && (
          <span className="text-foreground-lighter line-through">
            {JSON.stringify(diffValue, null, 2)}
          </span>
        )}
      </ExternalLink>
    ) : (
      <span
        className={classNames("text-foreground-lightest", diffValueSpacer)}
        title={showTitle ? `${column.title}=null` : undefined}
      >
        <span>null</span>
        {diffValue && (
          <span className="text-foreground-lighter line-through">
            {JSON.stringify(diffValue, null, 2)}
          </span>
        )}
      </span>
    );
  } else if (dataType === "control_status") {
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
        className={classNames("link-highlight", diffValueSpacer)}
        title={showTitle ? `${column.title}=${value?.toString()}` : undefined}
      >
        <span>{value?.toString()}</span>
        {!!diffValue && (
          <span className="text-foreground-lighter line-through">
            {diffValue.toString()}
          </span>
        )}
      </ExternalLink>
    ) : (
      <span
        className={classNames(
          "tabular-nums",
          value ? null : "text-foreground-light",
          diffValueSpacer,
        )}
        title={showTitle ? `${column.title}=${value?.toString()}` : undefined}
      >
        <span>{value?.toString()}</span>
        {(diffValue || diffValue?.toString() === "false") && (
          <span className="text-foreground-lighter line-through">
            {diffValue.toString()}
          </span>
        )}
      </span>
    );
  } else if (dataType === "jsonb" || isObject(value)) {
    const asJsonString = JSON.stringify(value, null, 2);
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames("link-highlight", diffValueSpacer)}
        title={showTitle ? `${column.title}=${asJsonString}` : undefined}
      >
        <span>{asJsonString}</span>
        {diffValue && (
          <span className="text-foreground-lighter line-through">
            {JSON.stringify(diffValue, null, 2)}
          </span>
        )}
      </ExternalLink>
    ) : (
      <span
        className={diffValueSpacer}
        title={showTitle ? `${column.title}=${asJsonString}` : undefined}
      >
        <span>{asJsonString}</span>
        {diffValue && (
          <span className="text-foreground-lighter line-through">
            {JSON.stringify(diffValue, null, 2)}
          </span>
        )}
      </span>
    );
  } else if (dataType === "text") {
    if (!!value.match && value.match("^https?://")) {
      cellContent = (
        <ExternalLink
          className={classNames("link-highlight tabular-nums", diffValueSpacer)}
          to={value}
          title={showTitle ? `${column.title}=${value}` : undefined}
        >
          <span>{value}</span>
          {diffValue && (
            <span className="text-foreground-lighter line-through">
              {diffValue.toString()}
            </span>
          )}
        </ExternalLink>
      );
    }
    const mdMatch =
      !!value.match && value.match("^\\[(.*)\\]\\((https?://.*)\\)$");
    if (mdMatch) {
      cellContent = (
        <ExternalLink
          className={classNames("link-highlight tabular-nums", diffValueSpacer)}
          to={mdMatch[2]}
          title={showTitle ? `${column.title}=${value}` : undefined}
        >
          <span>{mdMatch[1]}</span>
          {diffValue && (
            <span className="text-foreground-lighter line-through">
              {diffValue.toString()}
            </span>
          )}
        </ExternalLink>
      );
    }
  } else if (dataType === "timestamp" || dataType === "timestamptz") {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames("link-highlight tabular-nums", diffValueSpacer)}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        <span>{value}</span>
        {diffValue && (
          <span className="text-foreground-lighter line-through">
            {diffValue.toString()}
          </span>
        )}
      </ExternalLink>
    ) : (
      <span
        className={classNames("tabular-nums", diffValueSpacer)}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        <span>{value}</span>
        {diffValue && (
          <span className="text-foreground-lighter line-through">
            {diffValue.toString()}
          </span>
        )}
      </span>
    );
  } else if (isNumericCol(dataType)) {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames("link-highlight tabular-nums", diffValueSpacer)}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        <span>{value.toLocaleString()}</span>
        {diffValue && (
          <span className="text-foreground-lighter line-through">
            {diffValue.toLocaleString()}
          </span>
        )}
      </ExternalLink>
    ) : (
      <span
        className={classNames("tabular-nums", diffValueSpacer)}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        <span>{value.toLocaleString()}</span>
        {diffValue && (
          <span className="text-foreground-lighter line-through">
            {diffValue.toLocaleString()}
          </span>
        )}
      </span>
    );
  }
  if (!cellContent) {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames("link-highlight tabular-nums", diffValueSpacer)}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        <span>{value}</span>
        {diffValue && (
          <span className="text-foreground-lighter line-through">
            {diffValue.toString()}
          </span>
        )}
      </ExternalLink>
    ) : (
      <span
        className={classNames("tabular-nums", diffValueSpacer)}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        <span>{value}</span>
        {diffValue && (
          <span className="text-foreground-lighter line-through">
            {diffValue.toString()}
          </span>
        )}
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
      {filterEnabled && addFilter && (
        <div className="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
          <button
            onClick={() => addFilter("equal", column.name, value, context)}
            className="text-black-scale-7 hover:text-black-scale-8 focus:outline-none"
            title="Add value to include filter"
          >
            <Icon className="h-5 w-5" icon="add_circle" />
          </button>
          <button
            onClick={() => addFilter("not_equal", column.name, value, context)}
            className="text-black-scale-7 hover:text-black-scale-8 focus:outline-none"
            title="Add value to exclude filter"
          >
            <Icon className="h-5 w-5" icon="do_not_disturb_on" />
          </button>
        </div>
      )}
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

const TableView = ({
  rowData,
  columns,
  hiddenColumns,
  hasTopBorder = false,
  filterEnabled = false,
  context = "",
}) => {
  const {
    themeContext: { theme },
  } = useDashboard();
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

  const { getTableProps, getTableBodyProps, headerGroups, prepareRow, rows } =
    useTable(
      { columns, data: rowData, initialState: { hiddenColumns } },
      useSortBy,
    );

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
      const data = rows.map((row) => row.values);
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
        {...getTableProps()}
        className={classNames(
          "min-w-full divide-y divide-table-divide overflow-hidden",
          hasTopBorder ? "border-t border-divide" : null,
        )}
      >
        <thead className="text-table-head border-b border-divide">
          {headerGroups.map((headerGroup) => {
            const { key, ...otherHeaderGroupProps } =
              headerGroup.getHeaderGroupProps();
            return (
              <tr key={key} {...otherHeaderGroupProps}>
                {headerGroup.headers.map((column) => {
                  const { key, ...otherHeaderProps } = column.getHeaderProps(
                    column.getSortByToggleProps(),
                  );
                  return (
                    <th
                      key={key}
                      {...otherHeaderProps}
                      scope="col"
                      className={classNames(
                        "py-3 text-left text-sm font-normal tracking-wider whitespace-nowrap pl-4",
                      )}
                    >
                      {column.render("Header")}
                      {column.isSortedDesc ? (
                        <SortDescendingIcon className="inline-block h-4 w-4" />
                      ) : (
                        <SortAscendingIcon
                          className={classNames(
                            "inline-block h-4 w-4",
                            !column.isSorted ? "invisible" : null,
                          )}
                        />
                      )}
                    </th>
                  );
                })}
              </tr>
            );
          })}
        </thead>
        <tbody
          {...getTableBodyProps()}
          className="divide-y divide-table-divide"
        >
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
            prepareRow(row);
            const { key, ...otherRowProps } = row.getRowProps();
            return (
              <tr
                key={key}
                {...otherRowProps}
                className={classNames(
                  row.original?.__diff === "deleted" ? "bg-red-100" : null,
                  row.original?.__diff === "inserted" ? "bg-green-100" : null,
                )}
              >
                {row.cells.map((cell) => {
                  const { key, ...otherCellProps } = cell.getCellProps();
                  const diff = tableRowDiffColumn(
                    row.original,
                    cell.column.name,
                  );
                  return (
                    <td
                      key={key}
                      {...otherCellProps}
                      className={classNames(
                        "px-4 py-4 align-top content-center text-sm",
                        isNumericCol(cell.column.data_type) ? "text-right" : "",
                        cell.column.wrap === "all"
                          ? "break-keep"
                          : "whitespace-nowrap",
                        cell.column.diff === "deleted"
                          ? theme.name === ThemeNames.STEAMPIPE_DARK
                            ? "bg-red-900"
                            : "bg-red-100"
                          : null,
                        cell.column.diff === "inserted"
                          ? theme.name === ThemeNames.STEAMPIPE_DARK
                            ? "bg-green-900"
                            : "bg-green-100"
                          : null,
                        diff.hasDiffColumn
                          ? theme.name === ThemeNames.STEAMPIPE_DARK
                            ? "bg-amber-900"
                            : "bg-amber-100"
                          : null,
                      )}
                    >
                      <MemoCellValue
                        column={cell.column}
                        rowIndex={index}
                        rowTemplateData={rowTemplateData}
                        value={cell.value}
                        addFilter={addFilter}
                        filterEnabled={filterEnabled}
                        context={context}
                        diffValue={diff.diffValue}
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
  const { columns, hiddenColumns } = useMemo(
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
    <div>
      {/* Render column selector UI */}
      {/* {renderColumnSelector()}  */}
      <TableView
        rowData={rowData}
        columns={columns} // Use filtered columns for the table
        hiddenColumns={hiddenColumns}
        hasTopBorder={!!props.title}
        filterEnabled={props.filterEnabled}
        context={props.context}
      />
    </div>
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
        Header: col.original_name || col.name,
        title: col.original_name || col.name,
        accessor: col.name,
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

export { TableView, TableViewWrapper };
