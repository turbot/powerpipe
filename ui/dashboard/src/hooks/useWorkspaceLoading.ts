import { useEffect, useState } from "react";

// How often the server is asked whether its workspace is loading. The reloading message shows within about one
// interval of a reload starting.
export const WORKSPACE_POLL_INTERVAL_MS = 500;

// A request that has not answered by now counts as "not loading", so that a message never stays up on a server that
// has stopped answering. The server answers this request during a reload, by design.
export const WORKSPACE_REQUEST_TIMEOUT_MS = 3000;

const getWorkspaceUrl = () =>
  // In development the browser is at http://localhost:3000, so assume the server is on its default port. The server
  // sends no CORS headers, so there the request fails and the message never shows.
  process.env.NODE_ENV === "development"
    ? "http://localhost:9033/api/v0/workspace"
    : "/api/v0/workspace";

// Whether the server is reloading its workspace, as reported at GET /api/v0/workspace. Asks at once, then every
// WORKSPACE_POLL_INTERVAL_MS while the page is visible, so that a page opened during a reload sees it. Anything but a
// clear "loading" from the server (an error, a timeout, an answer that is not JSON) is "not loading".
const useWorkspaceLoading = (enabled = true): boolean => {
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!enabled || typeof fetch !== "function") {
      setLoading(false);
      return;
    }

    let stopped = false;
    let inFlight = false;
    let timer: ReturnType<typeof setTimeout> | undefined;

    const schedule = () => {
      if (stopped) {
        return;
      }
      clearTimeout(timer);
      timer = setTimeout(poll, WORKSPACE_POLL_INTERVAL_MS);
    };

    const poll = async () => {
      clearTimeout(timer);
      if (stopped || inFlight) {
        return;
      }
      if (document.hidden) {
        // visibilitychange asks again when the page is shown
        schedule();
        return;
      }
      inFlight = true;
      const controller = new AbortController();
      const abortTimer = setTimeout(
        () => controller.abort(),
        WORKSPACE_REQUEST_TIMEOUT_MS,
      );
      let isLoading = false;
      try {
        const response = await fetch(getWorkspaceUrl(), {
          cache: "no-store",
          signal: controller.signal,
        });
        if (response.ok) {
          const body = await response.json();
          isLoading = body?.loading === true;
        }
      } catch (error) {
        isLoading = false;
      } finally {
        clearTimeout(abortTimer);
        inFlight = false;
      }
      if (stopped) {
        return;
      }
      setLoading(isLoading);
      schedule();
    };

    const onVisibilityChange = () => {
      if (!document.hidden) {
        poll();
      }
    };

    document.addEventListener("visibilitychange", onVisibilityChange);
    poll();

    return () => {
      stopped = true;
      clearTimeout(timer);
      document.removeEventListener("visibilitychange", onVisibilityChange);
    };
  }, [enabled]);

  return enabled && loading;
};

export default useWorkspaceLoading;
