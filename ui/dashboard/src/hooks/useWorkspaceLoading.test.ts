import useWorkspaceLoading, {
  WORKSPACE_POLL_INTERVAL_MS,
  WORKSPACE_REQUEST_TIMEOUT_MS,
} from "./useWorkspaceLoading";
import { act, renderHook } from "@testing-library/react";

// The real hook runs against a stand-in for the browser's fetch, which answers what GET /api/v0/workspace
// answers: the server's load state, whose `loading` the hook reads.

type Answer = { status?: number; body?: unknown } | "network-error" | "hang";

let answer: Answer;
let requests: string[];
const originalFetch = global.fetch;

const stubFetch = () => {
  global.fetch = jest.fn((url: string, init?: RequestInit) => {
    requests.push(url);
    if (answer === "network-error") {
      return Promise.reject(new TypeError("Failed to fetch"));
    }
    if (answer === "hang") {
      // never answers, but ends when the caller aborts, as fetch does
      return new Promise((_, reject) => {
        init?.signal?.addEventListener("abort", () =>
          reject(new DOMException("aborted", "AbortError")),
        );
      });
    }
    const status = answer.status ?? 200;
    return Promise.resolve({
      ok: status >= 200 && status < 300,
      status,
      json: () => Promise.resolve(answer.body),
    });
  }) as unknown as typeof fetch;
};

// lets the pending request, and the state change it causes, finish
const tick = async (ms = 0) => {
  await act(async () => {
    jest.advanceTimersByTime(ms);
  });
};

const loadState = (loading: boolean) => ({
  status: 200,
  body: { loading, load_count: 1 },
});

beforeEach(() => {
  jest.useFakeTimers();
  requests = [];
  answer = loadState(false);
  stubFetch();
});

afterEach(() => {
  jest.useRealTimers();
  global.fetch = originalFetch;
});

describe("useWorkspaceLoading", () => {
  test("reports nothing while no reload runs", async () => {
    const { result } = renderHook(() => useWorkspaceLoading());
    await tick();
    await tick(WORKSPACE_POLL_INTERVAL_MS * 3);

    expect(requests.length).toBeGreaterThan(1);
    expect(requests[0]).toBe("/api/v0/workspace");
    expect(result.current).toBe(false);
  });

  test("reports a reload that starts within a second, and its end", async () => {
    const { result } = renderHook(() => useWorkspaceLoading());
    await tick();
    expect(result.current).toBe(false);

    answer = loadState(true);
    await tick(WORKSPACE_POLL_INTERVAL_MS);
    expect(result.current).toBe(true);
    // the interval is inside the one second the issue allows
    expect(WORKSPACE_POLL_INTERVAL_MS).toBeLessThan(1000);

    answer = loadState(false);
    await tick(WORKSPACE_POLL_INTERVAL_MS);
    expect(result.current).toBe(false);
  });

  test("reports a reload that is already running when the page opens", async () => {
    answer = loadState(true);
    const { result } = renderHook(() => useWorkspaceLoading());
    // no wait for the first interval
    await tick();
    expect(result.current).toBe(true);
  });

  test("stops reporting when the reload ends in a failure", async () => {
    answer = loadState(true);
    const { result } = renderHook(() => useWorkspaceLoading());
    await tick();
    expect(result.current).toBe(true);

    // a failed load ends too: the server says loading false, with the error on its last load
    answer = {
      status: 200,
      body: { loading: false, load_count: 2, last_load: { error: "bad mod" } },
    };
    await tick(WORKSPACE_POLL_INTERVAL_MS);
    expect(result.current).toBe(false);
  });

  test("never sticks when the server stops answering", async () => {
    answer = loadState(true);
    const { result } = renderHook(() => useWorkspaceLoading());
    await tick();
    expect(result.current).toBe(true);

    answer = "network-error";
    await tick(WORKSPACE_POLL_INTERVAL_MS);
    expect(result.current).toBe(false);
  });

  test("never sticks on a request that hangs", async () => {
    answer = loadState(true);
    const { result } = renderHook(() => useWorkspaceLoading());
    await tick();
    expect(result.current).toBe(true);

    answer = "hang";
    await tick(WORKSPACE_POLL_INTERVAL_MS);
    await tick(WORKSPACE_REQUEST_TIMEOUT_MS);
    expect(result.current).toBe(false);
  });

  test("treats an error status as not loading", async () => {
    answer = loadState(true);
    const { result } = renderHook(() => useWorkspaceLoading());
    await tick();
    expect(result.current).toBe(true);

    answer = { status: 500, body: { loading: true } };
    await tick(WORKSPACE_POLL_INTERVAL_MS);
    expect(result.current).toBe(false);
  });

  test("treats an answer that is not JSON as not loading", async () => {
    answer = loadState(true);
    const { result } = renderHook(() => useWorkspaceLoading());
    await tick();
    expect(result.current).toBe(true);

    // a server with no such route answers with its page
    (global.fetch as jest.Mock).mockImplementation(() =>
      Promise.resolve({
        ok: true,
        status: 200,
        json: () => Promise.reject(new SyntaxError("Unexpected token <")),
      }),
    );
    await tick(WORKSPACE_POLL_INTERVAL_MS);
    expect(result.current).toBe(false);
  });

  test("asks nothing, and reports nothing, when disabled", async () => {
    answer = loadState(true);
    const { result } = renderHook(() => useWorkspaceLoading(false));
    await tick(WORKSPACE_POLL_INTERVAL_MS * 3);

    expect(requests).toHaveLength(0);
    expect(result.current).toBe(false);
  });

  test("stops asking when unmounted", async () => {
    const { unmount } = renderHook(() => useWorkspaceLoading());
    await tick();
    unmount();
    const before = requests.length;
    await tick(WORKSPACE_POLL_INTERVAL_MS * 3);

    expect(requests).toHaveLength(before);
  });
});
