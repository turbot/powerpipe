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
        status: "complete",
        summary: { alarm: 10, ok: 10, info: 5, skip: 5, error: 2 },
        summary_diff: { alarm: 2, ok: -4, info: 0, skip: 0, error: 2 },
        severity_summary: {},
      },
    ],
  },
};
