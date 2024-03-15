import { validateFilter } from "@powerpipe/components/dashboards/check/CheckFilterEditor"; // Replace with the actual module path

interface TestCase {
  name: string;
  input: any; // Input data for the test
  expected: boolean; // Expected result
}

const filterTestCases: TestCase[] = [
  {
    name: "valid filter with type and value",
    input: { operator: "equal", type: "status", value: "alarm" },
    expected: true,
  },
  {
    name: "valid filter with type, key and value",
    input: {
      operator: "equal",
      type: "control_tag",
      key: "environment",
      value: "production",
    },
    expected: true,
  },
  {
    name: "filter missing operator",
    input: { type: "resource", value: "name" },
    expected: false,
  },
  {
    name: "filter missing type",
    input: { operator: "equal", key: "name" },
    expected: false,
  },
  {
    name: "filter missing value",
    input: { operator: "equal", type: "control_tag", key: "environment" },
    expected: false,
  },
  {
    name: "filter missing both key and value",
    input: { operator: "equal", type: "resource" },
    expected: false,
  },
  {
    name: "valid AND filter with valid filters",
    input: {
      operator: "and",
      expressions: [
        { operator: "equal", type: "resource", value: "*mybucket*" },
        {
          operator: "equal",
          type: "control_tag",
          key: "environment",
          value: "production",
        },
        { operator: "equal", type: "dimension", key: "region", value: "us*" },
      ],
    },
    expected: true,
  },
  {
    name: "AND filter with an invalid filter",
    input: {
      and: [{ key: "name" }],
    },
    expected: false,
  },
];

// const orFilterTestCases: TestCase[] = [
//   {
//     name: "valid OR filter with valid filters",
//     input: {
//       or: [
//         { type: "resource", value: "*mybucket*" },
//         { type: "control_tag", key: "environment", value: "production" },
//       ],
//     },
//     expected: true,
//   },
//   {
//     name: "empty OR filter",
//     input: { or: [] },
//     expected: true,
//   },
//   {
//     name: "OR filter with an invalid filter",
//     input: {
//       or: [{ type: "resource", key: "name" }, { key: "name" }],
//     },
//     expected: false,
//   },
//   {
//     name: "OR filter with one invalid and one valid filter",
//     input: {
//       or: [{ type: "resource", value: "*mybucket*" }, { key: "name" }],
//     },
//     expected: false,
//   },
// ];

function runTestCases(
  testCases: TestCase[],
  validationFunction: (input: any) => boolean,
) {
  testCases.forEach((testCase) => {
    it(`should return ${testCase.expected} for ${testCase.name}`, () => {
      const result = validationFunction(testCase.input);
      expect(result).toBe(testCase.expected);
    });
  });
}

describe("Check Filter Validation", () => {
  describe("validateFilter", () => {
    runTestCases(filterTestCases, validateFilter);
  });
});
