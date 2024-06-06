import ControlDimension from "../check/Benchmark/ControlDimension";
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
} from "@powerpipe/constants/icons";
import {
  BasePrimitiveProps,
  ExecutablePrimitiveProps,
  isNumericCol,
  LeafNodeDataColumn,
  LeafNodeDataRow,
} from "../common";
import { classNames } from "@powerpipe/utils/styles";
import {
  ErrorIcon,
  SortAscendingIcon,
  SortDescendingIcon,
} from "@powerpipe/constants/icons";
import { memo, useEffect, useMemo, useState } from "react";
import { getComponent, registerComponent } from "../index";
import { PanelDefinition } from "@powerpipe/types";
import { RowRenderResult } from "../common/types";
import { useSortBy, useTable } from "react-table";

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
      name: col.original_name || col.name,
      data_type: col.data_type,
      wrap: colWrap,
      // Boolean data types do not sort under the default alphanumeric sorting logic of react-table
      // On the next column type that needs specialising we'll move this out into a function / hook
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
};

const CellValue = ({
  column,
  rowIndex,
  rowTemplateData,
  value,
  showTitle = false,
}: CellValueProps) => {
  const ExternalLink = getComponent("external_link");
  const [href, setHref] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Calculate a link for this cell
  useEffect(() => {
    const renderedTemplateObj = rowTemplateData[rowIndex];

    if (!renderedTemplateObj) {
      setHref(null);
      setError(null);
      return;
    }
    const renderedTemplateForColumn = renderedTemplateObj[column.accessor];
    if (!renderedTemplateForColumn) {
      setHref(null);
      setError(null);
      return;
    }
    if (renderedTemplateForColumn.result) {
      setHref(renderedTemplateForColumn.result);
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
    // True should be
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
  // Fallback is just show it as a string
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
    cellContent
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
  const { ready: templateRenderReady, renderTemplates } = useTemplateRender();
  const [rowTemplateData, setRowTemplateData] = useState<RowRenderResult[]>([]);

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
          .map((col) => [col.accessor, col.href_template as string]),
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
    <>
      <table
        {...getTableProps()}
        className={classNames(
          "min-w-full divide-y divide-table-divide overflow-hidden",
          hasTopBorder ? "border-t border-divide" : null,
        )}
      >
        <thead className="text-table-head border-b border-divide">
          {headerGroups.map((headerGroup) => (
            <tr {...headerGroup.getHeaderGroupProps()}>
              {headerGroup.headers.map((column) => (
                <th
                  {...column.getHeaderProps(column.getSortByToggleProps())}
                  scope="col"
                  className={classNames(
                    "py-3 text-left text-sm font-normal tracking-wider whitespace-nowrap pl-4",
                    isNumericCol(column.data_type) ? "text-right" : null,
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
              ))}
            </tr>
          ))}
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
            return (
              <tr {...row.getRowProps()}>
                {row.cells.map((cell) => (
                  <td
                    {...cell.getCellProps()}
                    className={classNames(
                      "px-4 py-4 align-top content-center text-sm",
                      isNumericCol(cell.column.data_type) ? "text-right" : "",
                      cell.column.wrap === "all"
                        ? "break-keep"
                        : "whitespace-nowrap",
                    )}
                  >
                    <MemoCellValue
                      column={cell.column}
                      rowIndex={index}
                      rowTemplateData={rowTemplateData}
                      value={cell.value}
                    />
                  </td>
                ))}
              </tr>
            );
          })}
        </tbody>
      </table>
    </>
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
          .map((col) => [col.accessor, col.href_template as string]),
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
                <div key={`${col.accessor}-${rowIndex}`}>
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
                      value={row[col.accessor]}
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
