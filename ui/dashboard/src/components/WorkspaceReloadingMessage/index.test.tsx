import WorkspaceReloadingMessage from "./index";
import { act, render, screen } from "@testing-library/react";
import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeLive,
} from "@powerpipe/types";

// The real message and hook run, over a fetch that answers a reload in progress. The dashboard state is the only
// stand-in: the message reads its data mode and error.

let mockDashboardState: { dataMode: string; error: unknown };

jest.mock("@powerpipe/hooks/useDashboardState", () => ({
  useDashboardState: () => mockDashboardState,
}));

const originalFetch = global.fetch;

const renderDuringReload = async () => {
  render(<WorkspaceReloadingMessage />);
  await act(async () => {
    jest.advanceTimersByTime(0);
  });
};

beforeEach(() => {
  jest.useFakeTimers();
  mockDashboardState = { dataMode: DashboardDataModeLive, error: null };
  global.fetch = jest.fn(() =>
    Promise.resolve({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ loading: true }),
    }),
  ) as unknown as typeof fetch;
});

afterEach(() => {
  jest.useRealTimers();
  global.fetch = originalFetch;
});

describe("WorkspaceReloadingMessage", () => {
  test("shows the message while the server reloads", async () => {
    await renderDuringReload();
    expect(screen.getByRole("status").textContent).toContain(
      "Reloading the workspace",
    );
  });

  test("shows nothing when no reload runs", async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ loading: false }),
      }),
    ) as unknown as typeof fetch;
    await renderDuringReload();
    expect(screen.queryByRole("status")).toBeNull();
  });

  test("a workspace error that arrives during the reload replaces the message", async () => {
    const { rerender } = render(<WorkspaceReloadingMessage />);
    await act(async () => {
      jest.advanceTimersByTime(0);
    });
    expect(screen.getByRole("status")).not.toBeNull();

    mockDashboardState = { dataMode: DashboardDataModeLive, error: "bad mod" };
    rerender(<WorkspaceReloadingMessage />);
    expect(screen.queryByRole("status")).toBeNull();
  });

  test("an error left on screen by an earlier failure does not hide the next reload", async () => {
    mockDashboardState = { dataMode: DashboardDataModeLive, error: "bad mod" };
    await renderDuringReload();
    expect(screen.getByRole("status")).not.toBeNull();
  });

  test("clearing an earlier error during the reload keeps the message", async () => {
    mockDashboardState = { dataMode: DashboardDataModeLive, error: "bad mod" };
    const { rerender } = render(<WorkspaceReloadingMessage />);
    await act(async () => {
      jest.advanceTimersByTime(0);
    });
    expect(screen.getByRole("status")).not.toBeNull();

    mockDashboardState = { dataMode: DashboardDataModeLive, error: null };
    rerender(<WorkspaceReloadingMessage />);
    expect(screen.getByRole("status")).not.toBeNull();
  });

  test("a snapshot shows no message and asks the server nothing", async () => {
    mockDashboardState = {
      dataMode: DashboardDataModeCLISnapshot,
      error: null,
    };
    await renderDuringReload();
    expect(screen.queryByRole("status")).toBeNull();
    expect(global.fetch).not.toHaveBeenCalled();
  });
});
