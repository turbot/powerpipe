import ControlDimension from "../grouping/Benchmark/ControlDimension";
import isEmpty from "lodash/isEmpty";
import isObject from "lodash/isObject";
import useDeepCompareEffect from "use-deep-compare-effect";
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
import { classNames } from "@powerpipe/utils/styles";
import { injectSearchPathPrefix } from "@powerpipe/utils/url";
import { memo, useEffect, useMemo, useState } from "react";
import { getComponent, registerComponent } from "../index";
import { PanelDefinition } from "@powerpipe/types";
import { RowRenderResult } from "../common/types";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useSortBy, useTable } from "react-table";
import Icon from "@powerpipe/components/Icon";

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
  showTitle?: boolean;
  handleAddFilter: (column: string, value: any) => void;
  handleRemoveFilter: (filter: { column: string; value: any }) => void;
};

const CellValue = ({
  column,
  rowIndex,
  rowTemplateData,
  value,
  showTitle = false,
  handleAddFilter,
  handleRemoveFilter,
}: CellValueProps) => {
  const ExternalLink = getComponent("external_link");
  const { searchPathPrefix } = useDashboard();
  const [href, setHref] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

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
      <div className="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
        <button
          onClick={() => handleAddFilter(column.name, value)}
          className="text-blue-500 hover:text-blue-700 focus:outline-none"
          title="Add Filter"
        >
          <Icon className="h-5 w-5" icon="add_circle"/>
        </button>
        <button
          onClick={() => handleRemoveFilter({ column: column.name, value })}
          className="text-red-500 hover:text-red-700 focus:outline-none"
          title="Exclude Filter"
        >
          <Icon className="h-5 w-5" icon="do_not_disturb_on"/>
        </button>
      </div>
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
  };

const TableView = ({
  rowData,
  columns,
  hiddenColumns,
  hasTopBorder = false,
}) => {
  const [activeFilters, setActiveFilters] = useState([]);
  const [excludedFilters, setExcludedFilters] = useState([]);
  const { ready: templateRenderReady, renderTemplates } = useTemplateRender();
  const [rowTemplateData, setRowTemplateData] = useState<RowRenderResult[]>([]);

  const handleAddFilter = (column, value) => {
    const newFilter = { column, value };
    if (!activeFilters.some(f => f.column === column && f.value === value)) {
      setActiveFilters([...activeFilters, newFilter]);
      setExcludedFilters(excludedFilters.filter(f => f.column !== column || f.value !== value));
    }
  };

  const handleRemoveFilter = (filter) => {
    if (!excludedFilters.some(f => f.column === filter.column && f.value === filter.value)) {
      setExcludedFilters([...excludedFilters, filter]);
      setActiveFilters(activeFilters.filter(f => f.column !== filter.column || f.value !== filter.value));
    }
  };

  const filteredData = useMemo(() => {
    let filtered = rowData;

    if (activeFilters.length > 0) {
      filtered = filtered.filter(row => {
        return activeFilters.every(filter => row[filter.column] === filter.value);
      });
    }

    if (excludedFilters.length > 0) {
      filtered = filtered.filter(row => {
        return excludedFilters.every(filter => row[filter.column] !== filter.value);
      });
    }

    return filtered;
  }, [rowData, activeFilters, excludedFilters]);

  const { getTableProps, getTableBodyProps, headerGroups, prepareRow, rows } =
    useTable(
      { columns, data: filteredData, initialState: { hiddenColumns } },
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

  const capitalize = (str) => str.charAt(0).toUpperCase() + str.slice(1);

  return (
    <div className="overflow-x-auto">
       {(activeFilters.length > 0 || excludedFilters.length > 0) && (
        <div className="mb-4 p-2 black-shade-2  rounded shadow-sm flex flex-wrap gap-2">
          {[...activeFilters, ...excludedFilters].map((filter, index) => {
            const isActive = activeFilters.some(f => f.column === filter.column && f.value === filter.value);
            return (
              <div
                key={index}
                className="flex items-center bg-black-scale-2 text-black-scale-8 px-3 py-1 rounded-full shadow-sm"
              >
                <span className="mr-2">
                  {`${isActive ? '+ | ' : '- | '}${capitalize(filter.column)}: ${filter.value}`}
                </span>
                <button
                  onClick={() => {
                    if (isActive) {
                      setActiveFilters(activeFilters.filter(f => f.column !== filter.column || f.value !== filter.value));
                    } else {
                      setExcludedFilters(excludedFilters.filter(f => f.column !== filter.column || f.value !== filter.value));
                    }
                  }}
                  className="text-black-scale-6 hover:text-black-scale-8 focus:outline-none"
                  title="Remove Filter"
                >
                  &times;
                </button>
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
                        "py-3 text-left text-sm font-medium tracking-wider whitespace-nowrap pl-4",
                        isNumericCol(column.data_type) ? "text-right" : "text-left",
                      )}
                    >
                      {capitalize(column.render("Header"))}
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
                className="px-4 py-4 align-top content-center text-sm italic whitespace-normal"
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
              <tr key={key} {...otherRowProps}>
                {row.cells.map((cell) => {
                  const { key, ...otherCellProps } = cell.getCellProps();
                  return (
                    <td
                      key={key}
                      {...otherCellProps}
                      className={classNames(
                        "px-4 py-4 align-top content-center text-sm",
                        isNumericCol(cell.column.data_type) ? "text-right" : "",
                        cell.column.wrap === "all"
                          ? "break-keep"
                          : "whitespace-normal",
                      )}
                    >
                      <MemoCellValue
                        column={cell.column}
                        rowIndex={index}
                        rowTemplateData={rowTemplateData}
                        value={cell.value}
                        handleAddFilter={handleAddFilter}
                        handleRemoveFilter={handleRemoveFilter}
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

const TableViewWrapper = (props: TableProps) => {
  const { columns, hiddenColumns } = useMemo(
    () => getColumns(props.data ? props.data.columns : [], props.properties),
    [props.data, props.properties],
  );
  const rowData = useMemo(
    () => getData(columns, props.data ? props.data.rows : []),
    [columns, props.data],
  );

  return props.data ? (
    <TableView
      rowData={rowData}
      columns={columns}
      hiddenColumns={hiddenColumns}
      hasTopBorder={!!props.title}
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

export { TableView };
