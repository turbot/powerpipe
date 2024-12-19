import CodeBlock from "@powerpipe/components/CodeBlock";
import CopyToClipboard from "@powerpipe/components/CopyToClipboard";
import {
  LeafNodeData,
  LeafNodeDataColumn,
} from "@powerpipe/components/dashboards/common";
import { parseDate } from "@powerpipe/utils/date";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";
import { useEffect, useMemo, useState } from "react";
import SearchInput from "@powerpipe/components/SearchInput";
import Icon from "@powerpipe/components/Icon";

const getNumericValue = (value) => {
  if (
    !value ||
    value.NaN ||
    value.Exp === null ||
    value.Status === null ||
    value.InfinityModifier === null
  ) {
    return Number.NaN.toString();
  }
  if (value.InfinityModifier === 1) {
    return Number.POSITIVE_INFINITY.toString();
  }
  if (value.InfinityModifier === -1) {
    return Number.NEGATIVE_INFINITY.toString();
  }

  const parts: string[] = [];
  if (value.Int === null) {
    parts.push("0");
  } else if (value.Int !== undefined) {
    parts.push(value.Int.toString());
  } else {
    parts.push(value.toString());
  }
  if (value.Exp !== undefined && value.Exp !== null) {
    parts.push("e");
    parts.push(parseInt(value.Exp, 10).toString());
  }
  return parseFloat(parts.join("")).toString();
};

const renderValue = (name: string, dataType: string, value: any) => {
  switch (dataType.toLowerCase()) {
    case "text":
    case "varchar":
      return <CodeBlock language="yaml">{value}</CodeBlock>;
    case "timestamptz":
      return (
        <CodeBlock language="yaml">
          {parseDate(value)?.format() || ""}
        </CodeBlock>
      );
    case "jsonb":
    case "varchar[]":
      return (
        <CodeBlock language="json" style={{ fontSize: "12px" }}>
          {JSON.stringify(value, null, 2)}
        </CodeBlock>
      );
    case "numeric":
    case "bigint": {
      if (name === "timestamp") {
        return (
          <CodeBlock language="yaml">
            {parseDate(value)?.format() || ""}
          </CodeBlock>
        );
      }
      return <CodeBlock language="json">{getNumericValue(value)}</CodeBlock>;
    }
    default:
      return (
        <CodeBlock language="json">{JSON.stringify(value, null, 2)}</CodeBlock>
      );
  }
};

const TableRowItem = ({ dataType, name, value }) => {
  const [showOptions, setShowOptions] = useState(false);
  return (
    <div
      key={name}
      id={name}
      className="p-4 space-y-1"
      onMouseEnter={() => setShowOptions(true)}
      onMouseLeave={() => setShowOptions(false)}
    >
      <div className="flex icon-spacer items-center text-sm">
        <span className="block font-light tracking-wider text-table-head">
          {name}
        </span>
        {showOptions && (
          <>
            <CopyToClipboard data={JSON.stringify(value, null, 2)} />
          </>
        )}
      </div>
      <div>
        {value === null && (
          <span className="text-foreground-lightest font-mono">null</span>
        )}
        {value !== null && renderValue(name, dataType, value)}
      </div>
    </div>
  );
};

const TableRowSidePanel = ({
  data,
  requestedColumnName,
  rowIndex,
}: {
  data: LeafNodeData | undefined;
  requestedColumnName?: string;
  rowIndex: number | undefined;
}) => {
  const { closeSidePanel } = useDashboardPanelDetail();
  const [search, setSearch] = useState("");

  // if (!data || !data.columns || !data.rows || rowIndex === undefined) {
  //   return null;
  // }

  const { columns } = data;
  const row = data.rows[rowIndex];

  const orderedRow: { column: LeafNodeDataColumn; value: any }[] = [];
  for (const column of columns) {
    orderedRow.push({ column, value: row[column.name] });
  }

  const filteredObj = useMemo(() => {
    if (!search) {
      return orderedRow;
    }

    const searchParts = search.trim().toLowerCase().split(" ");
    const filtered: { column: LeafNodeDataColumn; value: any }[] = [];
    for (const item of orderedRow) {
      const dataType = item.column.data_type.toLowerCase();
      if (
        searchParts.every((searchPart) => {
          if (item.column.name.toLowerCase().indexOf(searchPart) >= 0) {
            return true;
          } else if (search === "null" && item.value === null) {
            return true;
          } else if (
            (dataType === "jsonb" ||
              dataType === "varchar[]" ||
              dataType.startsWith("struct")) &&
            item.value &&
            JSON.stringify(item.value).toLowerCase().indexOf(searchPart) >= 0
          ) {
            return true;
          } else if (
            item.value &&
            item.value.toString().toLowerCase().indexOf(searchPart) >= 0
          ) {
            return true;
          } else {
            return false;
          }
        })
      ) {
        filtered.push(item);
      }
    }
    return filtered;
  }, [orderedRow, search]);

  useEffect(() => {
    if (!requestedColumnName) {
      return;
    }

    const element = document.getElementById(requestedColumnName);
    if (element) {
      element.scrollIntoView();
    }
  }, [requestedColumnName]);

  return (
    <>
      <div className="flex items-center justify-between p-4 min-w-[300px]">
        <h3>Row</h3>
        <Icon
          className="w-5 h-5 text-foreground cursor-pointer hover:text-foreground-light shrink-0"
          icon="close"
          onClick={closeSidePanel}
          title="Close"
        />
      </div>
      <div className="flex flex-col w-full pt-3">
        <div className="px-4 pt-3">
          <SearchInput
            placeholder="Search row"
            setValue={setSearch}
            value={search}
          />
        </div>
        <div className="flex-1 overflow-auto divide-y divide-divide space-y-3">
          {filteredObj.map((item) => (
            <TableRowItem
              key={item.column.name}
              dataType={item.column.data_type}
              name={item.column.name}
              value={item.value}
            />
          ))}
        </div>
      </div>
    </>
  );
};

export default TableRowSidePanel;
