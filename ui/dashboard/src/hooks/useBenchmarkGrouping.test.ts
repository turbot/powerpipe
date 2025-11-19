import { includeResult } from "./useBenchmarkGrouping";
import { Filter } from "@powerpipe/components/dashboards/grouping/common";

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
    expect(includeResult({ tags: {} }, filter)).toBe(true);
  });

  it("includes missing tag for not_in", () => {
    const filter = makeFilter("not_in", "deprecated", ["true", "false"]);
    expect(includeResult({ tags: {} }, filter)).toBe(true);
  });

  it("excludes present disallowed value for not_equal", () => {
    const filter = makeFilter("not_equal", "deprecated", "true");
    expect(includeResult({ tags: { deprecated: "true" } }, filter)).toBe(false);
  });

  it("includes present allowed value for not_equal", () => {
    const filter = makeFilter("not_equal", "deprecated", "true");
    expect(includeResult({ tags: { deprecated: "false" } }, filter)).toBe(true);
  });

  it("requires match for equal", () => {
    const filter = makeFilter("equal", "deprecated", "true");
    expect(includeResult({ tags: {} }, filter)).toBe(false);
    expect(includeResult({ tags: { deprecated: "true" } }, filter)).toBe(true);
  });
});
