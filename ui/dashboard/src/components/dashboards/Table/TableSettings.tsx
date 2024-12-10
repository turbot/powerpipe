import Icon from "@powerpipe/components/Icon";
import SearchInput from "@powerpipe/components/SearchInput";
import { classNames } from "@powerpipe/utils/styles";
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import { Popover } from "@headlessui/react";
import { Table } from "@tanstack/react-table";
import { usePopper } from "react-popper";
import { useState } from "react";

type TableSettingsColumnsViewType = "all" | "visible" | "hidden";

const TableSettingsColumns = ({ table }: { table: Table<KeyValuePairs> }) => {
  const [search, setSearch] = useState("");
  const [view, setView] = useState<TableSettingsColumnsViewType>("all");
  const filteredColumns = table
    .getAllLeafColumns()
    .filter((column) => !search || `${column.id}`.match(search.trim()))
    .filter((column) => {
      if (view === "hidden") {
        return !column.getIsVisible();
      } else if (view === "visible") {
        return column.getIsVisible();
      } else {
        return true;
      }
    });

  const ColumnRender = ({ column }) => (
    <label
      className="flex items-center space-x-2 cursor-pointer"
      key={column.id}
    >
      <input
        className="inline-block focus:outline-none focus:ring-0"
        type="checkbox"
        checked={column.getIsVisible()}
        onChange={(e) => {
          e.stopPropagation();
          column.getToggleVisibilityHandler()(e);
        }}
      />
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
        {filteredColumns.map((column) => (
          <ColumnRender key={column.id} column={column} />
        ))}
        {!filteredColumns.length && <span>No columns</span>}
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
      {/*@ts-ignore*/}
      <Popover.Button ref={setReferenceElement} as="div">
        <Icon icon="data_table" className="h-4.5 w-4.5" />
      </Popover.Button>
      <Popover.Panel className="absolute z-10 pt-px">
        <div
          // @ts-ignore
          ref={setPopperElement}
          style={{ ...styles.popper }}
          {...attributes.popper}
        >
          <div
            onClick={(e) => e.stopPropagation()}
            className="border border-dashboard-panel rounded-md bg-dashboard mt-1 p-3 space-y-3 min-w-60 max-w-96"
          >
            <TableSettingsColumns table={table} />
          </div>
        </div>
      </Popover.Panel>
    </Popover>
  );
};

export default TableSettings;
