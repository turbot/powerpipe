import { DashboardActions, PanelDefinition } from "@powerpipe/types";
import { useCallback } from "react";
import { useDashboardState } from "./useDashboardState";

const useSelectPanel = (definition: PanelDefinition) => {
  const { dispatch } = useDashboardState();
  const openPanelDetail = useCallback(
    async (e) => {
      e.stopPropagation();
      dispatch({
        type: DashboardActions.SELECT_PANEL,
        panel: definition,
      });
    },
    [dispatch, definition],
  );

  return { select: openPanelDetail };
};

export default useSelectPanel;
