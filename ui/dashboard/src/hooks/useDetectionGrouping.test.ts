import { includeDetectionResult } from "./useDetectionGrouping";
import {
  CheckResultStatus,
  DetectionResult,
} from "@powerpipe/components/dashboards/grouping/common";
import { PanelDefinition } from "@powerpipe/types";

const makeDetectionResult = (
  tags: Record<string, string> = {},
): DetectionResult => ({
  rows: [],
  columns: [],
  dimension_columns: [],
  dimensions: [],
  tags,
  detection: {
    sort: "0",
    name: "detection.test",
    title: "Detection",
    type: "detection",
    severity_summary: {},
    status: "complete",
    summary: { total: 0, error: 0 },
  },
  benchmark_trunk: [],
  status: CheckResultStatus.ok,
  reason: "",
  resource: "",
  type: "result",
});

describe("includeResult detection_tag semantics", () => {
  const panel = { name: "panel" } as PanelDefinition;

  const makeFilters = (
    operator: "equal" | "not_equal" | "in" | "not_in",
    key: string,
    value: any,
  ) => ({
    [panel.name]: {
      operator: "and",
      expressions: [
        {
          operator,
          type: "detection_tag",
          key,
          value,
        },
      ],
    },
  });

  it("includes missing tag for not_equal", () => {
    const filters = makeFilters("not_equal", "deprecated", "true");
    expect(includeDetectionResult(makeDetectionResult({}), panel, filters)).toBe(
      true,
    );
  });

  it("includes missing tag for not_in", () => {
    const filters = makeFilters("not_in", "deprecated", ["true", "false"]);
    expect(includeDetectionResult(makeDetectionResult({}), panel, filters)).toBe(
      true,
    );
  });

  it("excludes present disallowed value for not_equal", () => {
    const filters = makeFilters("not_equal", "deprecated", "true");
    expect(
      includeDetectionResult(
        makeDetectionResult({ deprecated: "true" }),
        panel,
        filters,
      ),
    ).toBe(false);
  });

  it("includes present allowed value for not_equal", () => {
    const filters = makeFilters("not_equal", "deprecated", "true");
    expect(
      includeDetectionResult(
        makeDetectionResult({ deprecated: "false" }),
        panel,
        filters,
      ),
    ).toBe(true);
  });

  it("requires match for equal", () => {
    const filters = makeFilters("equal", "deprecated", "true");
    expect(includeDetectionResult(makeDetectionResult({}), panel, filters)).toBe(
      false,
    );
    expect(
      includeDetectionResult(
        makeDetectionResult({ deprecated: "true" }),
        panel,
        filters,
      ),
    ).toBe(true);
  });
});
