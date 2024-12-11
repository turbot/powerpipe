import { DashboardActions } from "@powerpipe/types";
import { GlobalHotKeys } from "react-hotkeys";
import { noop } from "@powerpipe/utils/func";
import { useCallback, useEffect, useState } from "react";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

const DashboardGlobalHotKeys = () => {
  const { dispatch } = useDashboardState();

  const [hotKeysHandlers, setHotKeysHandlers] = useState({
    CLOSE_PANEL_DETAIL: noop,
  });

  const hotKeysMap = {
    CLOSE_PANEL_DETAIL: ["esc"],
  };

  const closePanelDetail = useCallback(() => {
    dispatch({
      type: DashboardActions.SELECT_PANEL,
      panel: null,
    });
  }, [dispatch]);

  useEffect(() => {
    setHotKeysHandlers({
      CLOSE_PANEL_DETAIL: closePanelDetail,
    });
  }, [closePanelDetail]);

  return (
    <GlobalHotKeys
      allowChanges
      keyMap={hotKeysMap}
      handlers={hotKeysHandlers}
    />
  );
};

export default DashboardGlobalHotKeys;
