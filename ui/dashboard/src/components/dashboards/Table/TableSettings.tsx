import Icon from "@powerpipe/components/Icon";
import SearchInput from "@powerpipe/components/SearchInput";
import { classNames } from "@powerpipe/utils/styles";
import { Column, Table } from "@tanstack/react-table";
import { createPortal } from "react-dom";
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import { Popover } from "@headlessui/react";
import { ThemeProvider, ThemeWrapper } from "@powerpipe/hooks/useTheme";
import { useEffect, useState } from "react";
import { usePopper } from "react-popper";

type TableSettingsColumnsViewType = "all" | "visible" | "hidden";

const TableSettingsColumns = ({ table }: { table: Table<KeyValuePairs> }) => {
  const [search, setSearch] = useState("");
  const [view, setView] = useState<TableSettingsColumnsViewType>("all");
  const [filteredColumns, setFilteredColumns] = useState<
    Column<KeyValuePairs>[]
  >(() =>
    table
      .getAllLeafColumns()
      .filter((column) => !search || `${column.id}`.match(search.trim())),
  );

  useEffect(() => {
    setFilteredColumns(() =>
      table
        .getAllLeafColumns()
        .filter((column) => !search || `${column.id}`.match(search.trim())),
    );
  }, [search]);

  const ColumnRender = ({ column }) => (
    <label
      className="flex items-center space-x-2 cursor-pointer"
      key={column.id}
    >
      <input
        className="inline-block focus:outline-none focus:ring-0"
        {...{
          type: "checkbox",
          checked: column.getIsVisible(),
          onChange: column.getToggleVisibilityHandler(),
        }}
      />{" "}
      <span className="inline-block truncate" title={column.id}>
        {column.id}
      </span>
    </label>
  );

  const ViewSelector = ({
    label,
    view,
    setView,
  }: {
    label: TableSettingsColumnsViewType;
    view: TableSettingsColumnsViewType;
    setView: (view: TableSettingsColumnsViewType) => void;
  }) => (
    <span
      className={classNames(
        "capitalize",
        label === view
          ? "font-semibold underline"
          : "font-light cursor-pointer",
      )}
      onClick={label !== view ? () => setView(label) : undefined}
    >
      {label}
    </span>
  );

  return (
    <div className="space-y-3">
      <span className="font-semibold">Visible Columns</span>
      <SearchInput
        placeholder="Search columns..."
        value={search}
        setValue={setSearch}
      />
      <div className="flex space-x-4">
        <ViewSelector label="all" view={view} setView={setView} />
        <ViewSelector label="visible" view={view} setView={setView} />
        <ViewSelector label="hidden" view={view} setView={setView} />
      </div>
      <div className="max-h-64 overflow-x-auto pl-px">
        {filteredColumns
          .filter((column) => {
            if (view === "hidden") {
              return !column.getIsVisible();
            } else if (view === "visible") {
              return column.getIsVisible();
            } else {
              return true;
            }
          })
          .map((column) => (
            <ColumnRender key={column.id} column={column} />
          ))}
      </div>
    </div>
  );
};

const TableSettings = ({ table }: { table: Table<KeyValuePairs> }) => {
  const [popperElement, setPopperElement] = useState(null);
  const [referenceElement, setReferenceElement] = useState(null);
  const { styles, attributes } = usePopper(referenceElement, popperElement, {
    placement: "bottom-end",
  });

  return (
    <Popover className="relative">
      <Popover.Button ref={setReferenceElement} as="div" className="">
        <Icon icon="settings" className="h-4 w-4 cursor-pointer" />
      </Popover.Button>
      <Popover.Panel className="absolute z-10 pt-px">
        {createPortal(
          <ThemeProvider>
            <ThemeWrapper>
              <div
                // @ts-ignore
                ref={setPopperElement}
                style={{ ...styles.popper }}
                {...attributes.popper}
              >
                <div className="border border-dashboard-panel rounded-md bg-dashboard mt-1 p-3 space-y-3 min-w-60 max-w-96">
                  <TableSettingsColumns table={table} />
                </div>
              </div>
            </ThemeWrapper>
          </ThemeProvider>,
          // @ts-ignore as this element definitely exists
          document.getElementById("portals"),
        )}
        {/*{({ close }) => (*/}
        {/*  <div className="border border-dashboard-panel rounded-md bg-dashboard p-3 space-y-3 min-w-60 max-w-96">*/}
        {/*    <TableSettingsColumns table={table} />*/}
        {/*  </div>*/}
        {/*)}*/}
      </Popover.Panel>
    </Popover>
  );
};

export default TableSettings;
