import { includeResult } from "./useBenchmarkGrouping";
import {
  CheckResult,
  CheckResultStatus,
  Filter,
} from "@powerpipe/components/dashboards/grouping/common";

const makeCheckResult = (tags: Record<string, string> = {}): CheckResult => ({
  dimensions: [],
  tags,
  control: {
    sort: "0",
    name: "control.test",
    title: "Control",
    type: "control",
    severity_summary: {},
    status: "complete",
    summary: { alarm: 0, ok: 0, info: 0, skip: 0, error: 0 },
  },
  benchmark_trunk: [],
  status: CheckResultStatus.ok,
  reason: "",
  resource: "",
  type: "result",
});

describe("includeResult control_tag semantics", () => {

  const makeFilter = (
    operator: "equal" | "not_equal" | "in" | "not_in",
    key: string,
    value: any,
  ): Filter => ({
    operator: "and",
    expressions: [
      {
        operator,
        type: "control_tag",
        key,
        value,
      },
    ],
  });

  it("includes missing tag for not_equal", () => {
    const filter = makeFilter("not_equal", "deprecated", "true");
    expect(includeResult(makeCheckResult({}), filter)).toBe(true);
  });

  it("includes missing tag for not_in", () => {
    const filter = makeFilter("not_in", "deprecated", ["true", "false"]);
    expect(includeResult(makeCheckResult({}), filter)).toBe(true);
  });

  it("excludes present disallowed value for not_equal", () => {
    const filter = makeFilter("not_equal", "deprecated", "true");
    expect(
      includeResult(makeCheckResult({ deprecated: "true" }), filter),
    ).toBe(false);
  });

  it("includes present allowed value for not_equal", () => {
    const filter = makeFilter("not_equal", "deprecated", "true");
    expect(
      includeResult(makeCheckResult({ deprecated: "false" }), filter),
    ).toBe(true);
  });

  it("requires match for equal", () => {
    const filter = makeFilter("equal", "deprecated", "true");
    expect(includeResult(makeCheckResult({}), filter)).toBe(false);
    expect(
      includeResult(makeCheckResult({ deprecated: "true" }), filter),
    ).toBe(true);
  });
});
