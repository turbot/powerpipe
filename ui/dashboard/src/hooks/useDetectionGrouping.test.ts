import { includeDetectionResult } from "./useDetectionGrouping";
import { Filter } from "@powerpipe/components/dashboards/grouping/common";
import { PanelDefinition } from "@powerpipe/types";

describe("includeResult detection_tag semantics", () => {
  const panel = { name: "panel" } as PanelDefinition;

  const makeFilters = (
    operator: "equal" | "not_equal" | "in" | "not_in",
    key: string,
    value: any,
  ) => ({
    panel: {
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
    expect(includeDetectionResult({ tags: {} }, panel, filters)).toBe(true);
  });

  it("includes missing tag for not_in", () => {
    const filters = makeFilters("not_in", "deprecated", ["true", "false"]);
    expect(includeDetectionResult({ tags: {} }, panel, filters)).toBe(true);
  });

  it("excludes present disallowed value for not_equal", () => {
    const filters = makeFilters("not_equal", "deprecated", "true");
    expect(
      includeDetectionResult(
        { tags: { deprecated: "true" } },
        panel,
        filters,
      ),
    ).toBe(false);
  });

  it("includes present allowed value for not_equal", () => {
    const filters = makeFilters("not_equal", "deprecated", "true");
    expect(
      includeDetectionResult(
        { tags: { deprecated: "false" } },
        panel,
        filters,
      ),
    ).toBe(true);
  });

  it("requires match for equal", () => {
    const filters = makeFilters("equal", "deprecated", "true");
    expect(includeDetectionResult({ tags: {} }, panel, filters)).toBe(false);
    expect(
      includeDetectionResult(
        { tags: { deprecated: "true" } },
        panel,
        filters,
      ),
    ).toBe(true);
  });
});
