import CheckPanel from "../CheckPanel";
import sortBy from "lodash/sortBy";
import {
  CheckGroupNodeStates,
  GroupingActions,
  useBenchmarkGrouping,
} from "@powerpipe/hooks/useBenchmarkGrouping";
import { CheckNode } from "../common";
import { registerComponent } from "@powerpipe/components/dashboards";
import { useCallback, useEffect, useState } from "react";

type CheckGroupingProps = {
  node: CheckNode;
};

const CheckGrouping = ({ node }: CheckGroupingProps) => {
  const { dispatch, nodeStates } = useBenchmarkGrouping();
  const [restoreNodeStates, setRestoreNodeStates] =
    useState<CheckGroupNodeStates | null>(null);

  const expand = useCallback(() => {
    setRestoreNodeStates(nodeStates);
    dispatch({ type: GroupingActions.EXPAND_ALL_NODES });
  }, [dispatch, nodeStates]);

  const restore = useCallback(() => {
    if (restoreNodeStates) {
      dispatch({
        type: GroupingActions.UPDATE_NODES,
        nodes: restoreNodeStates,
      });
    }
  }, [dispatch, restoreNodeStates]);

  useEffect(() => {
    window.onbeforeprint = expand;
    window.onafterprint = restore;
  }, [expand, restore]);

  return (
    <div className="space-y-4 md:space-y-6 col-span-12">
      {sortBy(node.children, "sort")?.map((child) => (
        <CheckPanel key={child.name} depth={1} node={child} />
      ))}
    </div>
  );
};

registerComponent("check_grouping", CheckGrouping);

export default CheckGrouping;
