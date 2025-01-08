import Icon from "@powerpipe/components/Icon";
import SearchInput from "@powerpipe/components/SearchInput";
import sortBy from "lodash/sortBy";
import useTableConfig from "@powerpipe/hooks/useTableConfig";
import { AsyncNoop } from "@powerpipe/types/func";
import { Column, RowData } from "@tanstack/react-table";
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import { useEffect, useMemo, useState } from "react";

const ShortcutSelector = ({
  icon,
  label,
  action,
}: {
  icon: string;
  label: string;
  action: AsyncNoop;
}) => (
  <div
    className="flex items-center space-x-1 cursor-pointer hover:text-foreground-light"
    onClick={action}
  >
    <Icon icon={icon} className="h-5 w-5" />
    <span>{label}</span>
  </div>
);

const TableSettingsColumns = ({
  name,
  leafColumns,
}: {
  name: string;
  leafColumns: Column<RowData>[];
}) => {
  const { update } = useTableConfig(name);
  const [search, setSearch] = useState("");
  const columnMap = useMemo(() => {
    const columnMap: KeyValuePairs<{
      column: Column<KeyValuePairs>;
      index: number;
    }> = leafColumns.reduce((acc, column, currentIndex) => {
      acc[column.id] = { column, index: currentIndex };
      return acc;
    }, {});
    return columnMap;
  }, []);
  const [columnVisibility, setColumnVisibility] = useState<
    {
      id: string;
      visible: boolean;
    }[]
  >(
    leafColumns.map((column) => ({
      id: column.id,
      visible: column.getIsVisible(),
    })),
  );
  const filteredColumns = columnVisibility.filter(
    ({ id }) => !search || `${id}`.match(search.trim()),
  );

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

  const selectAllColumns = async () => {
    for (const column of leafColumns) {
      column.toggleVisibility(true);
    }
    setColumnVisibility(
      leafColumns.map((column) => ({
        id: column.id,
        visible: true,
      })),
    );
  };

  const clearAllColumns = async () => {
    for (const column of leafColumns) {
      column.toggleVisibility(false);
    }
    setColumnVisibility(
      leafColumns.map((column) => ({
        id: column.id,
        visible: false,
      })),
    );
  };

  const ColumnRender = ({ column }) => (
    <label
      className="flex px-4 items-center space-x-2 cursor-pointer"
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

  return (
    <>
      <div className="px-4">
        <SearchInput
          placeholder="Search row"
          setValue={setSearch}
          value={search}
        />
      </div>
      <div className="flex items-center p-4 py-3 space-x-4">
        <ShortcutSelector
          icon="select_all"
          label="Select all"
          action={selectAllColumns}
        />
        <ShortcutSelector
          icon="clear_all"
          label="Clear"
          action={clearAllColumns}
        />
      </div>
      <div className="flex-1 h-full max-h-full overflow-y-auto">
        {filteredColumns.map((column) => (
          <ColumnRender key={column.id} column={column} />
        ))}
        {!filteredColumns.length && <span>No columns</span>}
      </div>
    </>
  );
};

export default TableSettingsColumns;
