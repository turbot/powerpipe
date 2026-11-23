import { addLevel, validateGrouping } from "./index";

describe("validateGrouping", () => {
  test.each([
    [
      "the control default",
      [{ type: "benchmark" }, { type: "control" }, { type: "result" }],
    ],
    [
      "the detection default",
      [{ type: "benchmark" }, { type: "detection" }, { type: "result" }],
    ],
    [
      "two different tags",
      [
        { type: "control_tag", value: "service" },
        { type: "control_tag", value: "category" },
        { type: "result" },
      ],
    ],
    ["no Result level", [{ type: "severity" }, { type: "status" }]],
  ])("accepts %s", (_, config) => {
    expect(validateGrouping(config)).toEqual({ value: true, reason: "" });
  });

  test.each([
    [
      "an untyped level",
      [{ type: "benchmark" }, { type: "" }, { type: "result" }],
      "Choose a type for every grouping",
    ],
    [
      "an untyped level after Result",
      [{ type: "benchmark" }, { type: "result" }, { type: "" }],
      "Choose a type for every grouping",
    ],
    [
      "a repeated level",
      [{ type: "benchmark" }, { type: "benchmark" }, { type: "result" }],
      "Each grouping can only be used once",
    ],
    [
      "a repeated tag",
      [
        { type: "control_tag", value: "service" },
        { type: "control_tag", value: "service" },
      ],
      "Each grouping can only be used once",
    ],
    [
      "Result before another level",
      [{ type: "result" }, { type: "benchmark" }],
      "Result must be the last grouping",
    ],
  ])("rejects %s, and says why", (_, config, reason) => {
    expect(validateGrouping(config)).toEqual({ value: false, reason });
  });
});

describe("addLevel", () => {
  const level = { type: "" };

  test("puts a new level before Result", () => {
    expect(
      addLevel(
        [{ type: "benchmark" }, { type: "control" }, { type: "result" }],
        level,
      ),
    ).toEqual([
      { type: "benchmark" },
      { type: "control" },
      level,
      { type: "result" },
    ]);
  });

  test("adds a new level at the end when there is no Result", () => {
    expect(addLevel([{ type: "benchmark" }], level)).toEqual([
      { type: "benchmark" },
      level,
    ]);
  });

  test("adds the first level to an empty grouping", () => {
    expect(addLevel([], level)).toEqual([level]);
  });

  test("a typed level added to a valid grouping keeps it valid", () => {
    const added = addLevel([{ type: "benchmark" }, { type: "result" }], {
      type: "severity",
    });
    expect(validateGrouping(added)).toEqual({ value: true, reason: "" });
  });
});
