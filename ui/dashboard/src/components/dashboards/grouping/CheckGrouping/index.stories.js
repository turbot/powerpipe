import CheckGrouping from ".";
import { PanelStoryDecorator } from "@powerpipe/utils/storybook";
import { GroupingContext } from "@powerpipe/hooks/useBenchmarkGrouping";

const story = {
  title: "Benchmarks/Check Grouping",
  component: CheckGrouping,
};

export default story;

const Template = (args) => {
  const { firstChildSummaries, ...rest } = args;
  return (
    <GroupingContext.Provider
      value={{
        firstChildSummaries,
        groupingsConfig: [
          { type: "benchmark" },
          { name: "control" },
          { name: "result" },
        ],
        nodeStates: {},
      }}
    >
      <PanelStoryDecorator definition={rest} panelType="check_grouping" />
    </GroupingContext.Provider>
  );
};

export const LoadingBenchmark = Template.bind({});
LoadingBenchmark.args = {
  depth: 1,
  node: {
    status: "running",
    name: "root",
    children: [
      {
        sort: "1",
        name: "aws_compliance.benchmark.cis_v400_1",
        title: "1 Identity and Access Management",
        status: "running",
        summary: { alarm: 10, ok: 10, info: 5, skip: 5, error: 2 },
        severity_summary: {},
      },
    ],
  },
};

export const DiffBenchmark = Template.bind({});
DiffBenchmark.args = {
  firstChildSummaries: [{ alarm: 10, ok: 10, info: 5, skip: 5, error: 2 }],
  depth: 1,
  node: {
    name: "root",
    children: [
      {
        sort: "1",
        name: "aws_compliance.benchmark.cis_v400_1",
        title: "1 Identity and Access Management",
        panel_type: "benchmark",
        status: "complete",
        summary: {
          alarm: 10,
          alarm_diff: 12,
          ok: 10,
          ok_diff: 6,
          info: 5,
          info_diff: 5,
          skip: 5,
          skip_diff: 5,
          error: 2,
          error_diff: 4,
          __diff: "updated",
        },
        severity_summary: {},
      },
    ],
  },
};