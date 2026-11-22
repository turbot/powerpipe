import { parseGroupings, usableGroups } from "./useGroupingConfig";
import {
  checkDisplayGroupTypes,
  detectionDisplayGroupTypes,
} from "@powerpipe/components/dashboards/grouping/common";

describe("parseGroupings", () => {
  const panel = "mod.benchmark.b";

  test.each([
    ["no parameter", null],
    ["an empty parameter", ""],
    ["text that is not JSON", "{not-valid-json"],
    ["a JSON string", '"a string"'],
    ["a JSON number", "5"],
    ["JSON null", "null"],
    ["a top-level array", '[{"type":"benchmark"}]'],
    ["a panel whose value is not a list", JSON.stringify({ [panel]: "x" })],
    ["a panel whose list is empty", JSON.stringify({ [panel]: [] })],
    ["a list of nulls", JSON.stringify({ [panel]: [null] })],
    ["a group with a null type", JSON.stringify({ [panel]: [{ type: null }] })],
    ["a group with a numeric type", JSON.stringify({ [panel]: [{ type: 5 }] })],
    ["a group that is a string", JSON.stringify({ [panel]: ["benchmark"] })],
  ])("returns no groupings for %s", (_, raw) => {
    expect(parseGroupings(raw)).toEqual({});
  });

  test("keeps well-formed groups and drops the rest", () => {
    const raw = JSON.stringify({
      [panel]: [null, { type: "control_tag", value: "service" }, { type: 1 }],
      other: "x",
    });
    expect(parseGroupings(raw)).toEqual({
      [panel]: [{ type: "control_tag", value: "service" }],
    });
  });

  test("keeps every panel's groupings", () => {
    const groupings = {
      [panel]: [{ type: "benchmark" }, { type: "result" }],
      "mod.benchmark.d": [{ type: "detection_tag", value: "stage" }],
    };
    expect(parseGroupings(JSON.stringify(groupings))).toEqual(groupings);
  });
});

describe("usableGroups", () => {
  test("keeps the types the control tree builder draws", () => {
    const groups = checkDisplayGroupTypes.map((type) => ({ type }));
    expect(usableGroups(groups, checkDisplayGroupTypes)).toEqual(groups);
  });

  test("keeps the types the detection tree builder draws", () => {
    const groups = detectionDisplayGroupTypes.map((type) => ({ type }));
    expect(usableGroups(groups, detectionDisplayGroupTypes)).toEqual(groups);
  });

  test.each([
    ["detection", checkDisplayGroupTypes],
    ["detection_tag", checkDisplayGroupTypes],
    ["banana", checkDisplayGroupTypes],
    ["control", detectionDisplayGroupTypes],
    ["control_tag", detectionDisplayGroupTypes],
    ["reason", detectionDisplayGroupTypes],
    ["resource", detectionDisplayGroupTypes],
    ["status", detectionDisplayGroupTypes],
    ["banana", detectionDisplayGroupTypes],
  ])("drops a %s group the panel cannot draw", (type, types) => {
    expect(
      usableGroups([{ type }, { type: "severity" }, { type: "result" }], types),
    ).toEqual([{ type: "severity" }, { type: "result" }]);
  });

  test("returns nothing when no groups are saved", () => {
    expect(usableGroups(undefined, checkDisplayGroupTypes)).toEqual([]);
  });
});
