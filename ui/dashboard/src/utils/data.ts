import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import {
  LeafNodeData,
  LeafNodeDataColumn,
} from "@powerpipe/components/dashboards/common";

const hasData = (data: LeafNodeData | undefined) => {
  return (
    !!data &&
    data.columns &&
    data.rows &&
    data.columns.length > 0 &&
    data.rows.length > 0
  );
};

const getColumn = (
  columns: LeafNodeDataColumn[],
  name: string,
): LeafNodeDataColumn | undefined => {
  if (!columns || !name) {
    return undefined;
  }

  return columns.find((col) => col.name === name);
};

const tableRowDiffColumn = (row: KeyValuePairs, columnName: string) => {
  const match = Object.keys(row || {}).find((k) => {
    const match = matchDiffColumn(k);
    return !!match && match[1] === columnName;
  });
  return {
    hasDiffColumn: !!match,
    diffValue: match ? row[match] : undefined,
  };
};

const matchDiffColumn = (name: string) => /^(?!.*__)(.*?)_diff/.exec(name);

const isDiffColumn = (name: string) => !!matchDiffColumn(name);

const parseDiffColumn = (name: string) => {
  // const match = /^(.*)_diff(_[a-z\d]{4})?$/.exec(name);
  const match = matchDiffColumn(name);
  if (!match) {
    return { isDiff: false };
  }
  return {
    isDiff: true,
    pairedColumn: match[1],
  };
};

export {
  isDiffColumn,
  parseDiffColumn,
  getColumn,
  hasData,
  tableRowDiffColumn,
};
