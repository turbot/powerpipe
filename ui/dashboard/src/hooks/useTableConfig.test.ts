import React from "react";
import { act, renderHook } from "@testing-library/react";
import { MemoryRouter, useSearchParams } from "react-router-dom";
import useTableConfig, { parseTableConfigs } from "./useTableConfig";

jest.mock("@powerpipe/hooks/useDashboardState", () => ({
  useDashboardState: () => ({
    panelsMap: { "mod.table.t": { name: "mod.table.t", panel_type: "table" } },
  }),
}));

describe("parseTableConfigs", () => {
  const panel = "mod.table.t";

  test.each([
    ["no parameter", null],
    ["an empty parameter", ""],
    ["text that is not JSON", "{bad"],
    ["a JSON string", '"a string"'],
    ["a JSON number", "5"],
    ["JSON null", "null"],
    ["a top-level array", '[{"display_columns":["a"]}]'],
    ["an entry that is a string", JSON.stringify({ [panel]: "x" })],
    ["an entry that is null", JSON.stringify({ [panel]: null })],
    ["an entry that is a list", JSON.stringify({ [panel]: ["a"] })],
    ["an entry with no columns", JSON.stringify({ [panel]: {} })],
    [
      "columns that are not a list",
      JSON.stringify({ [panel]: { display_columns: "a" } }),
    ],
    [
      "columns that are not strings",
      JSON.stringify({ [panel]: { display_columns: ["a", 1] } }),
    ],
  ])("returns no settings for %s", (_, raw) => {
    expect(parseTableConfigs(raw)).toEqual({});
  });

  test("keeps well-formed entries and drops the rest", () => {
    const raw = JSON.stringify({
      [panel]: { display_columns: ["a", "b"] },
      "mod.table.bad": { display_columns: 5 },
      "mod.table.other": "x",
    });
    expect(parseTableConfigs(raw)).toEqual({
      [panel]: { display_columns: ["a", "b"] },
    });
  });

  test("keeps every table's settings, including an empty column list", () => {
    const tables = {
      [panel]: { display_columns: ["a"] },
      "mod.table.u": { display_columns: [] },
    };
    expect(parseTableConfigs(JSON.stringify(tables))).toEqual(tables);
  });
});

describe("useTableConfig", () => {
  const wrapperFor =
    (url: string) =>
    ({ children }: { children: React.ReactNode }) =>
      React.createElement(MemoryRouter, { initialEntries: [url] }, children);

  const useBoth = (panelName: string) => {
    const config = useTableConfig(panelName);
    const [searchParams] = useSearchParams();
    return { ...config, tableParam: searchParams.get("table") };
  };

  test("a damaged parameter gives default columns and no settings", () => {
    const { result } = renderHook(() => useBoth("mod.table.t"), {
      wrapper: wrapperFor("/mod.dashboard.d?table=%7Bbad"),
    });
    expect(result.current.allTables).toEqual({});
    expect(result.current.table).toEqual({ display_columns: [] });
  });

  test("a valid entry still applies next to a malformed one", () => {
    const param = encodeURIComponent(
      JSON.stringify({
        "mod.table.t": { display_columns: ["a"] },
        "mod.table.bad": 7,
      }),
    );
    const { result } = renderHook(() => useBoth("mod.table.t"), {
      wrapper: wrapperFor(`/mod.dashboard.d?table=${param}`),
    });
    expect(result.current.table).toEqual({ display_columns: ["a"] });
  });

  test("changing columns over a damaged parameter writes a valid one", () => {
    const { result } = renderHook(() => useBoth("mod.table.t"), {
      wrapper: wrapperFor("/mod.dashboard.d?table=%7Bbad"),
    });
    act(() => result.current.update({ display_columns: ["x"] }));
    expect(JSON.parse(result.current.tableParam as string)).toEqual({
      "mod.table.t": { display_columns: ["x"] },
    });
  });

  test("changing columns keeps other tables' valid settings only", () => {
    const param = encodeURIComponent(
      JSON.stringify({
        "mod.table.other": { display_columns: ["o"] },
        "mod.table.bad": { display_columns: "nope" },
      }),
    );
    const { result } = renderHook(() => useBoth("mod.table.t"), {
      wrapper: wrapperFor(`/mod.dashboard.d?table=${param}`),
    });
    act(() => result.current.update({ display_columns: ["x"] }));
    expect(JSON.parse(result.current.tableParam as string)).toEqual({
      "mod.table.other": { display_columns: ["o"] },
      "mod.table.t": { display_columns: ["x"] },
    });
  });

  test("resetting the only table removes the parameter", () => {
    const param = encodeURIComponent(
      JSON.stringify({ "mod.table.t": { display_columns: ["a"] } }),
    );
    const { result } = renderHook(() => useBoth("mod.table.t"), {
      wrapper: wrapperFor(`/mod.dashboard.d?table=${param}`),
    });
    act(() => result.current.update(null));
    expect(result.current.tableParam).toBeNull();
  });
});
