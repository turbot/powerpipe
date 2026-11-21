import LoadingIndicator from "@powerpipe/components/dashboards/LoadingIndicator";
import useWorkspaceLoading from "@powerpipe/hooks/useWorkspaceLoading";
import { DashboardDataModeLive } from "@powerpipe/types";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useRef } from "react";

const WorkspaceReloadingMessage = () => (
  <div
    className="fixed top-3 left-1/2 -translate-x-1/2 z-50 flex items-center space-x-2 rounded-md border border-divide bg-dashboard-panel px-4 py-2 text-sm text-foreground shadow-md print:hidden"
    role="status"
  >
    <LoadingIndicator className="h-4 w-4 shrink-0" />
    <span>Reloading the workspace…</span>
  </div>
);

// Shows a message while the server reloads its workspace, after a change to a mod or variable file or to the installed
// mods. A workspace error that arrives during the reload replaces it. An error that was already on screen when the
// reload began does not hide it, nor does clearing it: it is the reload the person started to fix it. A snapshot has no workspace to reload.
const WorkspaceReloadingMessageWrapper = () => {
  const { dataMode, error } = useDashboardState();
  const loading = useWorkspaceLoading(dataMode === DashboardDataModeLive);
  const errorWhenLoadingStarted = useRef(error);
  const wasLoading = useRef(false);
  if (loading && !wasLoading.current) {
    errorWhenLoadingStarted.current = error;
  }
  wasLoading.current = loading;

  if (!loading || (error && error !== errorWhenLoadingStarted.current)) {
    return null;
  }
  return <WorkspaceReloadingMessage />;
};

export default WorkspaceReloadingMessageWrapper;
