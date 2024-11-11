import CheckGrouping from ".";
import { PanelStoryDecorator } from "@powerpipe/utils/storybook";
import { GroupingContext } from "@powerpipe/hooks/useBenchmarkGrouping";

const story = {
  title: "Benchmarks/Check Grouping",
  component: CheckGrouping,
};

export default story;

const Template = (args) => (
  <GroupingContext.Provider
    value={{
      firstChildSummaries: [{ alarm: 10, ok: 10, info: 5, skip: 5, error: 2 }],
      groupingsConfig: [
        { type: "benchmark" },
        { name: "control" },
        { name: "result" },
      ],
      nodeStates: {},
    }}
  >
    <PanelStoryDecorator definition={args} panelType="check_grouping" />
  </GroupingContext.Provider>
);

export const LoadingBenchmark = Template.bind({});
LoadingBenchmark.args = {
  depth: 1,
  node: {
    name: "root",
    children: [
      {
        sort: "1",
        name: "benchmark",
        title: "My Benchmark",
        status: "loading",
        summary: { alarm: 10, ok: 10, info: 5, skip: 5, error: 2 },
        severity_summary: {},
      },
    ],
  },
};
