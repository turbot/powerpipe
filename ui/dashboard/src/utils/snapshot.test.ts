import {
  DashboardExecutionEventWithSchema,
  PanelDefinition,
} from "@powerpipe/types";
import {
  EXECUTION_SCHEMA_VERSION_20220614,
  EXECUTION_SCHEMA_VERSION_20220929,
} from "@powerpipe/constants/versions";
import {
  groupingToSnapshotMetadata,
  stripSnapshotDataForExport,
} from "./snapshot";
import {
  CheckDisplayGroup,
  CheckDisplayGroupType,
} from "@powerpipe/components/dashboards/grouping/common";

describe("snapshot utils", () => {
  describe("stripSnapshotDataForExport", () => {
    test("Schema 20220614", () => {
      const inputSnapshot: DashboardExecutionEventWithSchema = {
        schema_version: EXECUTION_SCHEMA_VERSION_20220614,
        execution_id: "0x140029247e0",
        dashboard_node: {
          name: "aws_insights.dashboard.aws_iam_user_dashboard",
        },
        layout: {
          name: "aws_insights.dashboard.aws_iam_user_dashboard",
          panel_type: "dashboard",
          children: [
            {
              name: "aws_insights.container.dashboard_aws_iam_user_dashboard_anonymous_container_0",
              panel_type: "container",
            },
          ],
        },
        panels: {
          "aws_insights.dashboard.aws_iam_user_dashboard": {
            name: "aws_insights.dashboard.aws_iam_user_dashboard",
            documentation: "# Some documentation",
            sql: "select something from somewhere",
            source_definition: 'some { hcl: "values" }',
            properties: {
              search_path: ["some_schema"],
              search_path_prefix: ["some_prefix"],
              sql: "select something from somewhere",
            },
          },
          "aws_insights.container.dashboard_aws_iam_user_dashboard_anonymous_container_0":
            {
              name: "aws_insights.container.dashboard_aws_iam_user_dashboard_anonymous_container_0",
              documentation: "# Some documentation",
              sql: "select something from somewhere",
              source_definition: 'some { hcl: "values" }',
              properties: {
                search_path: ["some_schema"],
                search_path_prefix: ["some_prefix"],
                sql: "select something from somewhere",
              },
            },
        },
        inputs: {
          "input.foo": "bar",
        },
        variables: {
          foo: "bar",
        },
        search_path: ["some_schema"],
        search_path_prefix: ["some_prefix"],
        start_time: "2022-10-27T14:43:57.79514+01:00",
        end_time: "2022-10-27T14:43:58.045925+01:00",
      };

      const migratedEvent = stripSnapshotDataForExport(inputSnapshot);

      const expectedPanels = {};

      for (const [name, panel] of Object.entries(inputSnapshot.panels)) {
        const { documentation, sql, source_definition, ...rest } =
          panel as PanelDefinition;
        expectedPanels[name] = { ...rest, properties: {} };
      }

      const expectedEvent = {
        schema_version: inputSnapshot.schema_version,
        dashboard_node: inputSnapshot.dashboard_node,
        execution_id: inputSnapshot.execution_id,
        layout: inputSnapshot.layout,
        panels: expectedPanels,
        inputs: inputSnapshot.inputs,
        variables: inputSnapshot.variables,
        start_time: inputSnapshot.start_time,
        end_time: inputSnapshot.end_time,
      };

      expect(migratedEvent).toEqual(expectedEvent);
    });

    test("Schema 20220929", () => {
      const inputSnapshot: DashboardExecutionEventWithSchema = {
        schema_version: EXECUTION_SCHEMA_VERSION_20220929,
        layout: {
          name: "aws_insights.dashboard.aws_iam_user_dashboard",
          panel_type: "dashboard",
          children: [
            {
              name: "aws_insights.container.dashboard_aws_iam_user_dashboard_anonymous_container_0",
              panel_type: "container",
            },
          ],
        },
        panels: {
          "aws_insights.dashboard.aws_iam_user_dashboard": {
            name: "aws_insights.dashboard.aws_iam_user_dashboard",
            documentation: "# Some documentation",
            sql: "select something from somewhere",
            source_definition: 'some { hcl: "values" }',
            properties: {
              search_path: ["some_schema"],
              search_path_prefix: ["some_prefix"],
              sql: "select something from somewhere",
            },
          },
          "aws_insights.container.dashboard_aws_iam_user_dashboard_anonymous_container_0":
            {
              name: "aws_insights.container.dashboard_aws_iam_user_dashboard_anonymous_container_0",
              documentation: "# Some documentation",
              sql: "select something from somewhere",
              source_definition: 'some { hcl: "values" }',
              properties: {
                search_path: ["some_schema"],
                search_path_prefix: ["some_prefix"],
                sql: "select something from somewhere",
              },
            },
        },
        inputs: {
          "input.foo": "bar",
        },
        variables: {
          foo: "bar",
        },
        search_path: ["some_schema"],
        search_path_prefix: ["some_prefix"],
        start_time: "2022-10-27T14:43:57.79514+01:00",
        end_time: "2022-10-27T14:43:58.045925+01:00",
      };

      const migratedEvent = stripSnapshotDataForExport(inputSnapshot);

      const expectedPanels = {};

      for (const [name, panel] of Object.entries(inputSnapshot.panels)) {
        const { documentation, sql, source_definition, ...rest } =
          panel as PanelDefinition;
        expectedPanels[name] = { ...rest, properties: {} };
      }

      const expectedEvent = {
        schema_version: inputSnapshot.schema_version,
        layout: inputSnapshot.layout,
        panels: expectedPanels,
        inputs: inputSnapshot.inputs,
        variables: inputSnapshot.variables,
        start_time: inputSnapshot.start_time,
        end_time: inputSnapshot.end_time,
      };

      expect(migratedEvent).toEqual(expectedEvent);
    });

    test("Unsupported schema", () => {
      const inputSnapshot: DashboardExecutionEventWithSchema = {
        // @ts-ignore
        schema_version: "20221010",
      };

      expect(() => stripSnapshotDataForExport(inputSnapshot)).toThrow(
        `Unsupported dashboard event schema ${inputSnapshot.schema_version}`,
      );
    });
  });

  describe("groupingToSnapshotMetadata", () => {
    it("should return an empty array for empty input", () => {
      const input: CheckDisplayGroup[] = [];
      expect(groupingToSnapshotMetadata(input)).toEqual([]);
    });

    const types: CheckDisplayGroupType[] = [
      "status",
      "reason",
      "resource",
      "severity",
      "benchmark",
      "control",
      "result",
    ];
    types.forEach((type) => {
      it(`should handle type ${type} without value`, () => {
        const input: CheckDisplayGroup[] = [{ type }];
        const expected = [{ type }];
        expect(groupingToSnapshotMetadata(input)).toEqual(expected);
      });
    });

    it("should omit dimension with no value", () => {
      const input: CheckDisplayGroup[] = [{ type: "dimension" }];
      const expected = [];
      expect(groupingToSnapshotMetadata(input)).toEqual(expected);
    });

    it("should handle dimension with value", () => {
      const input: CheckDisplayGroup[] = [
        { type: "dimension", value: "region" },
      ];
      const expected = [{ type: "dimension", value: "region" }];
      expect(groupingToSnapshotMetadata(input)).toEqual(expected);
    });

    it("should omit tag with no value", () => {
      const input: CheckDisplayGroup[] = [{ type: "control_tag" }];
      const expected = [];
      expect(groupingToSnapshotMetadata(input)).toEqual(expected);
    });

    it("should handle tag with value", () => {
      const input: CheckDisplayGroup[] = [
        { type: "control_tag", value: "category" },
      ];
      const expected = [{ type: "control_tag", value: "category" }];
      expect(groupingToSnapshotMetadata(input)).toEqual(expected);
    });

    it("should handle multiple groupings", () => {
      const input: CheckDisplayGroup[] = [
        { type: "status" },
        { type: "dimension", value: "region" },
      ];
      const expected = [
        { type: "status" },
        { type: "dimension", value: "region" },
      ];
      expect(groupingToSnapshotMetadata(input)).toEqual(expected);
    });

    it("should handle null input gracefully", () => {
      expect(groupingToSnapshotMetadata(null)).toEqual([]);
    });

    it("should handle undefined input gracefully", () => {
      expect(groupingToSnapshotMetadata(undefined)).toEqual([]);
    });
  });
});
