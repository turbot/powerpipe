import ControlDimension from "../grouping/Benchmark/ControlDimension";
import Icon from "@powerpipe/components/Icon";
import isEmpty from "lodash/isEmpty";
import isObject from "lodash/isObject";
import TableSettings from "@powerpipe/components/dashboards/Table/TableSettings";
import useDeepCompareEffect from "use-deep-compare-effect";
import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import useTableConfig from "@powerpipe/hooks/useTableConfig";
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
import { createPortal } from "react-dom";
import { Filter } from "@powerpipe/components/dashboards/grouping/common";
import {
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { formatDate, parseDate } from "@powerpipe/utils/date";
import { getComponent, registerComponent } from "../index";
import { injectSearchPathPrefix } from "@powerpipe/utils/url";
import { KeyValuePairs, RowRenderResult } from "../common/types";
import { memo, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { PanelDefinition } from "@powerpipe/types";
import { ThemeProvider, ThemeWrapper } from "@powerpipe/hooks/useTheme";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";
import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { usePanelControls } from "@powerpipe/hooks/usePanelControls";
import { usePopper } from "react-popper";
import { useSearchParams } from "react-router-dom";
import { useVirtualizer } from "@tanstack/react-virtual";
import useCopyToClipboard from "@powerpipe/hooks/useCopyToClipboard";
import { AsyncNoop } from "@powerpipe/types/func";

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
    if (properties?.columns?.[col.original_name || col.name]) {
      const c = properties.columns[col.original_name || col.name];

      // Column display always wins here, then we check if there are display_columns and whether the column is in that list
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

    // If we've got display columns set up and this column hasn't already had its default visibility set,
    // and it's not listed as a column to show, hide it by default
    if (
      !!properties?.display_columns?.length &&
      !properties?.display_columns.includes(col.name) &&
      !(col.name in columnVisibility)
    ) {
      columnVisibility[col.name] = false;
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
  panel: PanelDefinition;
  column: TableColumnInfo;
  rowIndex: number;
  rowTemplateData: RowRenderResult[];
  value: any;
  showTitle?: boolean;
  addFilter?: (
    operator: "equal" | "not_equal",
    key: string,
    value: any,
  ) => void;
  filterEnabled?: boolean;
  isScrolling?: boolean;
};

const CellValue = ({
  panel,
  column,
  rowIndex,
  rowTemplateData,
  value,
  addFilter,
  showTitle = false,
  filterEnabled = false,
  isScrolling = false,
}: CellValueProps) => {
  const baseClasses = "px-4 py-4";
  const { searchPathPrefix } = useDashboardSearchPath();
  const [href, setHref] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [referenceElement, setReferenceElement] = useState();
  const [showCellControls, setShowCellControls] = useState<boolean>(false);

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
    return href ? (
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
          <span className={baseClasses} title="Status = Alarm">
            <AlarmIcon className="text-alert w-5 h-5" />
          </span>
        );
        break;
      case "error":
        cellContent = (
          <span className={baseClasses} title="Status = Error">
            <AlarmIcon className="text-alert w-5 h-5" />
          </span>
        );
        break;
      case "ok":
        cellContent = (
          <span className={baseClasses} title="Status = OK">
            <OKIcon className="text-ok w-5 h-5" />
          </span>
        );
        break;
      case "info":
        cellContent = (
          <span className={baseClasses} title="Status = Info">
            <InfoIcon className="text-info w-5 h-5" />
          </span>
        );
        break;
      case "skip":
        cellContent = (
          <span className={baseClasses} title="Status = Skipped">
            <SkipIcon className="text-skip w-5 h-5" />
          </span>
        );
        break;
      default:
        cellContent = (
          <span className={baseClasses} title="Status = Unknown">
            <UnknownIcon className="text-foreground-light w-5 h-5" />
          </span>
        );
    }
  } else if (dataType === "control_dimensions") {
    cellContent = (
      <div className={classNames(baseClasses, "space-x-2")}>
        {(value || []).map((dimension) => (
          <ControlDimension
            key={dimension.key}
            dimensionKey={dimension.key}
            dimensionValue={dimension.value}
          />
        ))}
      </div>
    );
  } else if (dataType === "bool" || dataType === "boolean") {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames(baseClasses, "link-highlight")}
        title={showTitle ? `${column.title}=${value.toString()}` : undefined}
      >
        <>{value.toString()}</>
      </ExternalLink>
    ) : (
      <span
        className={classNames(
          baseClasses,
          value ? null : "text-foreground-light",
        )}
        title={showTitle ? `${column.title}=${value.toString()}` : undefined}
      >
        <>{value.toString()}</>
      </span>
    );
  } else if (
    dataType === "jsonb" ||
    dataType === "varchar[]" ||
    isObject(value)
  ) {
    const asJsonString = JSON.stringify(value, null, 2);
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames(baseClasses, "link-highlight")}
        title={showTitle ? `${column.title}=${asJsonString}` : undefined}
      >
        <>{asJsonString}</>
      </ExternalLink>
    ) : (
      <span
        className={baseClasses}
        title={showTitle ? `${column.title}=${asJsonString}` : undefined}
      >
        {asJsonString}
      </span>
    );
  } else if (dataType === "text" || dataType === "varchar") {
    if (!!value.match && value.match("^https?://")) {
      cellContent = (
        <ExternalLink
          className={classNames(baseClasses, "link-highlight tabular-nums")}
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
          className={classNames(baseClasses, "tabular-nums")}
          to={mdMatch[2]}
          title={showTitle ? `${column.title}=${value}` : undefined}
        >
          {mdMatch[1]}
        </ExternalLink>
      );
    } else {
      cellContent = href ? (
        <ExternalLink
          to={href}
          className={classNames(baseClasses, "link-highlight tabular-nums")}
          title={showTitle ? `${column.title}=${value}` : undefined}
        >
          {value}
        </ExternalLink>
      ) : (
        <span
          className={classNames(baseClasses, "tabular-nums")}
          title={showTitle ? `${column.title}=${value}` : undefined}
        >
          {value}
        </span>
      );
    }
  } else if (dataType === "date") {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames(baseClasses, "link-highlight tabular-nums")}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {formatDate(value)}
      </ExternalLink>
    ) : (
      <span
        className={classNames(baseClasses, "tabular-nums")}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {formatDate(value)}
      </span>
    );
  } else if (column.name === "timestamp" && dataType === "bigint") {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames(baseClasses, "link-highlight tabular-nums")}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {parseDate(value)?.format()}
      </ExternalLink>
    ) : (
      <span
        className={classNames(baseClasses, "tabular-nums")}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {parseDate(value)?.format()}
      </span>
    );
  } else if (dataType === "timestamp" || dataType === "timestamptz") {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames(baseClasses, "link-highlight tabular-nums")}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value}
      </ExternalLink>
    ) : (
      <span
        className={classNames(baseClasses, "tabular-nums")}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value}
      </span>
    );
  } else if (isNumericCol(dataType)) {
    cellContent = href ? (
      <ExternalLink
        to={href}
        className={classNames(baseClasses, "link-highlight tabular-nums")}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value.toLocaleString()}
      </ExternalLink>
    ) : (
      <span
        className={classNames(baseClasses, "tabular-nums")}
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
        className={classNames(baseClasses, "link-highlight tabular-nums")}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value}
      </ExternalLink>
    ) : (
      <span
        className={classNames(baseClasses, "tabular-nums")}
        title={showTitle ? `${column.title}=${value}` : undefined}
      >
        {value}
      </span>
    );
  }

  return error ? (
    <span
      className={classNames(baseClasses, "flex items-center space-x-2")}
      title={error}
    >
      {cellContent} <ErrorIcon className="inline h-4 w-4 text-alert" />
    </span>
  ) : isScrolling || !filterEnabled || !addFilter ? (
    cellContent
  ) : (
    <div
      ref={setReferenceElement}
      className="w-full"
      onMouseEnter={() => setShowCellControls(true)}
      onMouseLeave={() => setShowCellControls(false)}
    >
      {cellContent}
      {showCellControls && (
        <CellControls
          referenceElement={referenceElement}
          column={column}
          rowIndex={rowIndex}
          panel={panel}
          value={value}
          addFilter={addFilter}
        />
      )}
    </div>
  );
};

const CellControls = ({
  referenceElement,
  panel,
  rowIndex,
  column,
  value,
  addFilter,
}) => {
  const { selectSidePanel } = useDashboardPanelDetail();
  const { setShowPanelControls } = usePanelControls();
  const [popperElement, setPopperElement] = useState(null);
  const offset = useMemo(() => {
    return {
      name: "offset",
      options: {
        offset: [14, -1],
      },
    };
  }, []);
  const { styles, attributes } = usePopper(referenceElement, popperElement, {
    modifiers: [offset],
    placement: "bottom-start",
  });
  const { copy, copySuccess } = useCopyToClipboard();

  return (
    <>
      {createPortal(
        <ThemeProvider>
          <ThemeWrapper>
            <div
              // @ts-ignore
              ref={setPopperElement}
              style={{ ...styles.popper }}
              {...attributes.popper}
            >
              <div className="flex items-center space-x-1">
                <CellControl
                  iconClassName={copySuccess ? "text-ok" : undefined}
                  icon={
                    copySuccess
                      ? "materialsymbols-solid:content_copy"
                      : "content_copy"
                  }
                  title="Copy value"
                  onClick={!copySuccess ? async () => copy(value) : undefined}
                />
                <CellControl
                  icon="filter_alt"
                  title="Filter by this value"
                  onClick={() => addFilter("equal", column.name, value)}
                />
                <CellControl
                  // className="h-4 w-4"
                  icon="close"
                  title="Exclude value from results"
                  onClick={() => addFilter("not_equal", column.name, value)}
                />
                <CellControl
                  icon="split_scene"
                  title="View row"
                  onClick={async () => {
                    selectSidePanel({
                      panel,
                      context: {
                        requestedColumnName: column.name,
                        rowIndex,
                      },
                    });
                    setShowPanelControls(false);
                  }}
                />
              </div>
            </div>
          </ThemeWrapper>
        </ThemeProvider>,
        // @ts-ignore as this element definitely exists
        document.getElementById("portals"),
      )}
    </>
  );
};

const CellControl = ({
  iconClassName,
  icon,
  title,
  onClick,
}: {
  iconClassName?: string;
  icon: string;
  title: string;
  onClick: AsyncNoop | undefined;
}) => {
  return (
    <div
      onClick={onClick}
      className={classNames(
        "text-table-head hover:text-foreground",
        onClick ? "cursor-pointer" : null,
      )}
      title={title}
    >
      <Icon className={classNames(iconClassName, "h-4 w-4")} icon={icon} />
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
  display_columns?: string[];
  columns?: TableColumns;
};

export type TableProps = PanelDefinition &
  BasePrimitiveProps &
  ExecutablePrimitiveProps & {
    display_type?: TableType;
    properties?: TableProperties;
    filterEnabled?: boolean;
  };

const useTableFilters = (panelName: string) => {
  const { allFilters, filter: urlFilters } = useFilterConfig(panelName);
  const [searchParams, setSearchParams] = useSearchParams();
  const expressions = urlFilters.expressions;
  const filters: Filter[] = [];

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
    } else if (
      expression.operator === "in" &&
      expression.type === "dimension" &&
      !!expression.key &&
      !!expression.value
    ) {
      filters.push(expression);
    } else if (
      expression.operator === "not_in" &&
      expression.type === "dimension" &&
      !!expression.key &&
      !!expression.value
    ) {
      filters.push(expression);
    }
  }

  const addFilter = useCallback(
    (operator: "equal" | "not_equal", key: string, value: any) => {
      const newUrlFilters = { ...urlFilters };
      const expressions = [...(newUrlFilters.expressions || [])];
      const index = expressions.findIndex(
        (e) => e.type === "dimension" && e.key === key && e.value === value,
      );
      let newFilters =
        index !== undefined && index > -1
          ? [...expressions.slice(0, index), ...expressions.slice(index + 1)]
          : expressions || [];
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
          },
        ];
      } else {
        newFilters.push({
          operator,
          value,
          type: "dimension",
          key,
          title: value,
        });
      }
      newUrlFilters.expressions = newFilters;
      const newPanelFilters = {
        ...allFilters,
        [panelName]: newUrlFilters,
      };
      searchParams.set("where", JSON.stringify(newPanelFilters));
      setSearchParams(searchParams);
    },
    [urlFilters, searchParams, setSearchParams],
  );

  const removeFilter = useCallback(
    (key: string, value: any) => {
      const newUrlFilters = { ...urlFilters };
      let expressions = [...(newUrlFilters.expressions || [])];
      const index = expressions.findIndex(
        (e) => e.type === "dimension" && e.key === key && e.value === value,
      );
      let newFilters =
        index !== undefined
          ? [...expressions.slice(0, index), ...expressions.slice(index + 1)]
          : expressions;
      if (newFilters.length === 0) {
        newFilters = [{ operator: "equal" }];
      }
      newUrlFilters.expressions = newFilters;
      const newPanelFilters = {
        ...allFilters,
        [panelName]: newUrlFilters,
      };
      searchParams.set("where", JSON.stringify(newPanelFilters));
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

const useDisableHoverOnScroll = (scrollElement: HTMLDivElement | null) => {
  const isScrolling = useRef<boolean>(false);
  const scrollTimeout = useRef<NodeJS.Timeout | undefined>(undefined);

  const handleScroll = () => {
    if (!isScrolling.current) {
      isScrolling.current = true;
    }

    clearTimeout(scrollTimeout.current);
    scrollTimeout.current = setTimeout(() => {
      isScrolling.current = false;
    }, 200); // Wait for 200ms after scrolling stops
  };

  useEffect(() => {
    if (!scrollElement) {
      return;
    }

    scrollElement.addEventListener("scroll", handleScroll, { passive: true });

    return () => {
      scrollElement.removeEventListener("scroll", handleScroll);
      isScrolling.current = false;
      clearTimeout(scrollTimeout.current);
    };
  }, [scrollElement]);

  return isScrolling.current;
};

const TableViewVirtualizedRows = ({
  panel,
  data,
  columns,
  columnVisibility,
  hasTopBorder = false,
  filterEnabled = false,
}) => {
  const { filters, addFilter, removeFilter } = useTableFilters(panel.name);
  const { ready: templateRenderReady, renderTemplates } = useTemplateRender();
  const [rowTemplateData, setRowTemplateData] = useState<RowRenderResult[]>([]);
  const parentRef = useRef<HTMLDivElement>(null);
  const isScrolling = useDisableHoverOnScroll(parentRef.current);

  const table = useReactTable<KeyValuePairs>({
    data,
    columns,
    initialState: { columnVisibility },
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  const { customControls, setCustomControls } = usePanelControls();

  const [showColumnSettingsModal, setShowColumnSettingsModal] = useState(false);

  useDeepCompareEffect(() => {
    const tableColumnChooser = customControls.find(
      (c) => c.key === "table-select-columns",
    );
    if (tableColumnChooser) {
      return;
    }
    setCustomControls([
      ...customControls,
      {
        key: "table-select-columns",
        title: "Select table columns",
        icon: "add_column_right",
        action: async () => {
          setShowColumnSettingsModal(true);
        },
      },
    ]);
  }, [customControls, setCustomControls]);

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
    <>
      <div className="flex flex-col w-full overflow-hidden">
        {filterEnabled && !!filters.length && (
          <div className="flex flex-wrap gap-2 w-full p-4">
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
                        ? "filter_alt"
                        : "filter_alt_off"
                    }
                  />
                  <span>{`${filter.key}: ${filter.value}`}</span>
                  <span
                    onClick={() => removeFilter(filter.key, filter.value)}
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
                            "py-3 text-left font-normal tracking-wider whitespace-nowrap pl-4",
                          )}
                          //style={{ width: header.getSize() }}
                        >
                          {header.isPlaceholder ? null : (
                            <div
                              {...{
                                className: header.column.getCanSort()
                                  ? "cursor-pointer select-none"
                                  : "",
                                onClick:
                                  header.column.getToggleSortingHandler(),
                              }}
                            >
                              {flexRender(
                                header.column.columnDef.header,
                                header.getContext(),
                              )}
                              {{
                                asc: (
                                  <SortAscendingIcon
                                    className={classNames(
                                      "inline-block h-4 w-4",
                                    )}
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
                      className="px-4 py-4 align-top content-center italic whitespace-nowrap"
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
                              "align-top content-center max-w-[500px] overflow-x-hidden",
                              isNumericCol(cell.column.columnDef.data_type)
                                ? "text-right"
                                : "",
                              cell.column.columnDef.wrap === "all"
                                ? "break-keep"
                                : "whitespace-nowrap",
                            )}
                          >
                            <CellValue
                              panel={panel}
                              column={cell.column.columnDef}
                              rowIndex={index}
                              rowTemplateData={rowTemplateData}
                              value={cell.getValue()}
                              filterEnabled={filterEnabled}
                              isScrolling={isScrolling}
                              addFilter={addFilter}
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
      <TableSettings
        name={panel.name}
        table={table}
        show={showColumnSettingsModal}
        onClose={async () => setShowColumnSettingsModal(false)}
      />
    </>
  );
};

// TODO retain full width on mobile, no padding
const TableViewWrapper = (props: TableProps) => {
  const {
    table: { display_columns },
  } = useTableConfig(props.name);

  const { columns, columnVisibility } = useMemo(
    () =>
      getColumns(props.data ? props.data.columns : [], {
        ...props.properties,
        display_columns,
      }),
    [props.data, props.properties, display_columns],
  );
  const rowData = useMemo(
    () => getData(columns, props.data ? props.data.rows : []),
    [columns, props.data],
  );

  return props.data ? (
    <TableViewVirtualizedRows
      panel={props}
      data={rowData}
      columns={columns} // Use filtered columns for the table
      columnVisibility={columnVisibility}
      hasTopBorder={!!props.title}
      filterEnabled={props.filterEnabled}
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

      // If we've got display columns set up, it doesn't have a column override,
      // and it's not listed as a column to show, hide it by default
      if (
        !!props.properties?.display_columns?.length &&
        !props.properties?.display_columns.includes(col.name) &&
        !columnOverrides?.display
      ) {
        newColDef.display = "none";
      }

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
                  <span className="block text-table-head truncate">
                    {col.title}
                  </span>
                  <span
                    className={classNames(
                      "block",
                      col.wrap === "all" ? "break-keep" : "truncate",
                    )}
                  >
                    <MemoCellValue
                      panel={props}
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
