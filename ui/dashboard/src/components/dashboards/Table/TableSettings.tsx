import Modal from "@powerpipe/components/Modal";
import SearchInput from "@powerpipe/components/SearchInput";
import sortBy from "lodash/sortBy";
import useTableConfig from "@powerpipe/hooks/useTableConfig";
import { AsyncNoop } from "@powerpipe/types/func";
import { classNames } from "@powerpipe/utils/styles";
import { Column, Table } from "@tanstack/react-table";
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import { useEffect, useMemo, useState } from "react";

type TableSettingsColumnsViewType = "all" | "visible" | "hidden";

const TableSettingsColumns = ({
  name,
  table,
}: {
  name: string;
  table: Table<KeyValuePairs>;
}) => {
  const { update } = useTableConfig(name);
  const [search, setSearch] = useState("");
  const [view, setView] = useState<TableSettingsColumnsViewType>("all");
  const { allLeafColumns, columnMap } = useMemo(() => {
    const allLeafColumns = table.getAllLeafColumns();
    const columnMap: KeyValuePairs<{
      column: Column<KeyValuePairs>;
      index: number;
    }> = allLeafColumns.reduce((acc, column, currentIndex) => {
      acc[column.id] = { column, index: currentIndex };
      return acc;
    }, {});
    return { allLeafColumns, columnMap };
  }, []);
  const [columnVisibility, setColumnVisibility] = useState<
    {
      id: string;
      visible: boolean;
    }[]
  >(
    allLeafColumns.map((column) => ({
      id: column.id,
      visible: column.getIsVisible(),
    })),
  );
  const filteredColumns = columnVisibility
    .filter(({ id }) => !search || `${id}`.match(search.trim()))
    .filter(({ visible }) => {
      if (view === "hidden") {
        return !visible;
      } else if (view === "visible") {
        return visible;
      } else {
        return true;
      }
    });

  useEffect(() => {
    const visibleColumns = columnVisibility
      .filter(({ visible }) => visible)
      .map((c) => c.id);
    update(
      visibleColumns.length
        ? {
            display_columns: sortBy(visibleColumns, (c) => columnMap[c].index),
          }
        : null,
    );
  }, [columnVisibility]);

  const ColumnRender = ({ column }) => (
    <label
      className="flex items-center space-x-2 cursor-pointer"
      key={column.id}
    >
      <input
        className="inline-block focus:outline-none focus:ring-0"
        type="checkbox"
        checked={column.visible}
        onChange={(e) => {
          const origColumnInfo = columnMap[column.id];
          origColumnInfo.column.getToggleVisibilityHandler()(e);
          setColumnVisibility((previous) => [
            ...previous.slice(0, origColumnInfo.index),
            { id: column.id, visible: !column.visible },
            ...previous.slice(origColumnInfo.index + 1),
          ]);
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
      {/*<span className="font-semibold">Visible Columns</span>*/}
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

const TableSettings = ({
  name,
  table,
  show,
  onClose,
}: {
  name: string;
  table: Table<KeyValuePairs>;
  show: boolean;
  onClose: AsyncNoop;
}) => {
  if (!show) {
    return null;
  }

  return (
    <Modal allowClickAway onClose={onClose} title="Select table columns">
      <TableSettingsColumns name={name} table={table} />
    </Modal>
  );
};

export default TableSettings;

// http://localhost:3000/snapshot/aws_detections.benchmark.mitre_v151.20241210T145435.pps?where=%7B%22aws_detections.benchmark.mitre_v151%22%3A%7B%22operator%22%3A%22and%22%2C%22expressions%22%3A%5B%7B%22operator%22%3A%22equal%22%2C%22value%22%3A%22aws_detections.benchmark.mitre_v151_ta0001%22%2C%22type%22%3A%22benchmark%22%7D%5D%7D%7D&grouping=%7B%22aws_detections.benchmark.mitre_v151%22%3A%5B%7B%22type%22%3A%22detection%22%7D%2C%7B%22type%22%3A%22result%22%7D%5D%7D&table=%7B%22aws_detections.detection.cloudtrail_logs_detect_iam_root_console_logins%22%3A%7B%22display_columns%22%3A%5B%22timestamp%22%2C%22operation%22%2C%22resource%22%2C%22account_id%22%2C%22region%22%2C%22tp_id%22%5D%7D%7D&input.detection_range=%7B%22from%22%3A%222024-12-09T15%3A33%3A45.126Z%22%2C%22to%22%3Anull%2C%22relative%22%3A%221d%22%7D
//
