import useDashboardWebSocketEventHandler from "@powerpipe/hooks/useDashboardWebSocketEventHandler";
import WorkspaceReloadingMessage from "./index";
import {
  DashboardActions,
  DashboardDataMode,
  DashboardDataModeCLISnapshot,
  DashboardDataModeLive,
} from "@powerpipe/types";
import {
  DashboardStateProvider,
  useDashboardState,
} from "@powerpipe/hooks/useDashboardState";
import { act, render, screen } from "@testing-library/react";

// The real state, event handler and message run. The only stand-in is the WebSocket: the test hands the event
// handler the messages the server sends, as the WebSocket hook does.

const loadState = (loading: boolean, load_count = 2) => ({
  action: "workspace_load_state",
  load_state: { loading, load_count },
});

let send: (event: any) => void;

const Harness = () => {
  const { dispatch } = useDashboardState();
  const { eventHandler } = useDashboardWebSocketEventHandler(dispatch, {});
  send = (event) => act(() => eventHandler(event));
  return <WorkspaceReloadingMessage />;
};

const renderPage = (dataMode: DashboardDataMode = DashboardDataModeLive) =>
  render(
    <DashboardStateProvider
      analyticsContext={{
        setMetadata: () => {},
        setSelectedDashboard: () => {},
      }}
      breakpointContext={{}}
      dataMode={dataMode}
      rootPathname="/"
      versionMismatchCheck={false}
    >
      <Harness />
    </DashboardStateProvider>,
  );

const originalFetch = global.fetch;

beforeEach(() => {
  jest.useFakeTimers();
  global.fetch = jest.fn(() =>
    Promise.reject(new Error("the page must not poll the server")),
  ) as unknown as typeof fetch;
});

afterEach(() => {
  jest.useRealTimers();
  global.fetch = originalFetch;
});

describe("WorkspaceReloadingMessage", () => {
  test("shows nothing before the server says a reload is running", () => {
    renderPage();
    expect(screen.queryByRole("status")).toBeNull();
  });

  test("a page opened mid-reload shows the message from the first message the server sends", () => {
    renderPage();
    send(loadState(true));
    expect(screen.getByRole("status").textContent).toContain(
      "Reloading the workspace",
    );
  });

  test("the message goes when the server says the reload has ended", () => {
    renderPage();
    send(loadState(true));
    send(loadState(false, 3));
    expect(screen.queryByRole("status")).toBeNull();
  });

  test("a message that says no reload runs shows nothing", () => {
    renderPage();
    send(loadState(false, 1));
    expect(screen.queryByRole("status")).toBeNull();
  });

  test("a workspace error that arrives during the reload replaces the message", () => {
    renderPage();
    send(loadState(true));
    expect(screen.getByRole("status")).not.toBeNull();
    send({ action: DashboardActions.WORKSPACE_ERROR, error: "bad mod" });
    expect(screen.queryByRole("status")).toBeNull();
  });

  test("an error left on screen by an earlier failure does not hide the next reload", () => {
    renderPage();
    send({ action: DashboardActions.WORKSPACE_ERROR, error: "bad mod" });
    send(loadState(true));
    expect(screen.getByRole("status")).not.toBeNull();
  });

  test("clearing an earlier error during the reload keeps the message", () => {
    renderPage();
    send({ action: DashboardActions.WORKSPACE_ERROR, error: "bad mod" });
    send(loadState(true));
    send({ action: DashboardActions.WORKSPACE_ERROR, error: null });
    expect(screen.getByRole("status")).not.toBeNull();
  });

  test("a snapshot shows no message", () => {
    renderPage(DashboardDataModeCLISnapshot);
    send(loadState(true));
    expect(screen.queryByRole("status")).toBeNull();
  });

  test("no request is made, however long the page stays open, idle or mid-reload", () => {
    renderPage();
    send(loadState(false, 1));
    act(() => {
      jest.advanceTimersByTime(60 * 60 * 1000);
    });
    send(loadState(true));
    act(() => {
      jest.advanceTimersByTime(60 * 60 * 1000);
    });
    expect(global.fetch).not.toHaveBeenCalled();
  });
});
