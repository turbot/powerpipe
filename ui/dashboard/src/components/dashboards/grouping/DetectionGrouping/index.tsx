import DetectionPanel from "@powerpipe/components/dashboards/grouping/DetectionPanel";
import sortBy from "lodash/sortBy";
import {
  CheckGroupNodeStates,
  GroupingActions,
  useDetectionGrouping,
} from "@powerpipe/hooks/useDetectionGrouping";
import { DetectionNode } from "../common";
import { useCallback, useEffect, useState } from "react";

type DetectionGroupingProps = {
  node: DetectionNode;
};

const DetectionGrouping = ({ node }: DetectionGroupingProps) => {
  const { dispatch, nodeStates } = useDetectionGrouping();
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
        <DetectionPanel key={child.name} depth={1} node={child} />
      ))}
    </div>
  );
};

export default DetectionGrouping;
