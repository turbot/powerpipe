export const WORKSPACE_POLL_INTERVAL_MS = 500;
export const WORKSPACE_REQUEST_TIMEOUT_MS = 3000;

// Not yet asking the server
const useWorkspaceLoading = (enabled = true): boolean => {
  return false && enabled;
};

export default useWorkspaceLoading;
