import { DashboardActions } from "@powerpipe/types";
import { useEffect } from "react";

// A server that is gone cannot say its reload has ended, so the page stops saying the workspace is loading when the
// WebSocket is not open. The server tells the page its load state again when the WebSocket reconnects.
const useClearWorkspaceLoadingOnClose = (
  ready: boolean,
  dispatch: (action: any) => void,
) => {
  useEffect(() => {
    if (!ready) {
      dispatch({
        type: DashboardActions.WORKSPACE_LOAD_STATE,
        load_state: { loading: false },
      });
    }
  }, [ready, dispatch]);
};

export default useClearWorkspaceLoadingOnClose;
