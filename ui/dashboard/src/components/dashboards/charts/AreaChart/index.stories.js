import AreaChart from "./index";
import { PanelStoryDecorator } from "@powerpipe/utils/storybook";
import {
  MultiTimeSeriesDefaults,
  MultiTimeSeriesGroupedDefaults,
  SingleTimeSeriesDefaults,
} from "../Chart/index.stories";

const story = {
  title: "Charts/Area",
  component: AreaChart.component,
};

export default story;

const Template = (args) => (
  <PanelStoryDecorator
    definition={{ ...args, display_type: "area" }}
    panelType="chart"
  />
);

export const Loading = Template.bind({});
Loading.args = {
  data: null,
};

export const Error = Template.bind({});
Error.args = {
  data: null,
  error: "Something went wrong!",
};

export const SingleSeries = Template.bind({});
SingleSeries.storyName = "Single Series";
SingleSeries.args = {
  data: {
    columns: [
      { name: "Type", data_type: "TEXT" },
      { name: "Count", data_type: "INT8" },
    ],
    rows: [
      { Type: "User", Count: 12 },
      { Type: "Policy", Count: 93 },
      { Type: "Role", Count: 48 },
    ],
  },
};

export const LargeSeries = Template.bind({});
LargeSeries.args = {
  data: {
    columns: [
      { name: "Region", data_type: "TEXT" },
      { name: "Total", data_type: "INT8" },
    ],
    rows: [
      { Region: "us-east-1", Total: 14 },
      { Region: "eu-central-1", Total: 6 },
      { Region: "ap-south-1", Total: 4 },
      { Region: "ap-southeast-1", Total: 3 },
      { Region: "ap-southeast-2", Total: 2 },
      { Region: "ca-central-1", Total: 2 },
      { Region: "eu-north-1", Total: 2 },
      { Region: "eu-west-1", Total: 1 },
      { Region: "eu-west-2", Total: 1 },
      { Region: "eu-west-3", Total: 1 },
      { Region: "sa-east-1", Total: 1 },
      { Region: "us-east-2", Total: 1 },
      { Region: "us-west-1", Total: 1 },
      { Region: "ap-northeast-1", Total: 1 },
      { Region: "us-west-2", Total: 1 },
      { Region: "ap-northeast-2", Total: 1 },
    ],
  },
};

export const MultiSeriesStacked = Template.bind({});
MultiSeriesStacked.storyName = "Multi-Series (stacked)";
MultiSeriesStacked.args = {
  data: {
    columns: [
      { name: "Country", data_type: "TEXT" },
      { name: "Men", data_type: "INT8" },
      { name: "Women", data_type: "INT8" },
      { name: "Children", data_type: "INT8" },
    ],
    rows: [
      { Country: "England", Men: 16000000, Women: 13000000, Children: 8000000 },
      { Country: "Scotland", Men: 8000000, Women: 7000000, Children: 3000000 },
      {
        Country: "Wales",
        Men: 5000000,
        Women: 3000000,
        Children: 2500000,
      },
      {
        Country: "Northern Ireland",
        Men: 3000000,
        Women: 2000000,
        Children: 1000000,
      },
    ],
  },
  properties: {
    grouping: "stack",
  },
};

export const MultiSeriesGrouped = Template.bind({});
MultiSeriesGrouped.storyName = "Multi-Series (grouped)";
MultiSeriesGrouped.args = {
  data: {
    columns: [
      { name: "Country", data_type: "TEXT" },
      { name: "Men", data_type: "INT8" },
      { name: "Women", data_type: "INT8" },
      { name: "Children", data_type: "INT8" },
    ],
    rows: [
      { Country: "England", Men: 16000000, Women: 13000000, Children: 8000000 },
      { Country: "Scotland", Men: 8000000, Women: 7000000, Children: 3000000 },
      {
        Country: "Wales",
        Men: 5000000,
        Women: 3000000,
        Children: 2500000,
      },
      {
        Country: "Northern Ireland",
        Men: 3000000,
        Women: 2000000,
        Children: 1000000,
      },
    ],
  },
  properties: {
    grouping: "compare",
  },
};

export const MultiSeriesOverrides = Template.bind({});
MultiSeriesOverrides.storyName = "Multi-Series with Series Overrides";
MultiSeriesOverrides.args = {
  data: {
    columns: [
      { name: "Country", data_type: "TEXT" },
      { name: "Men", data_type: "INT8" },
      { name: "Women", data_type: "INT8" },
      { name: "Children", data_type: "INT8" },
    ],
    rows: [
      { Country: "England", Men: 16000000, Women: 13000000, Children: 8000000 },
      { Country: "Scotland", Men: 8000000, Women: 7000000, Children: 3000000 },
      {
        Country: "Wales",
        Men: 5000000,
        Women: 3000000,
        Children: 2500000,
      },
      {
        Country: "Northern Ireland",
        Men: 3000000,
        Women: 2000000,
        Children: 1000000,
      },
    ],
  },
  properties: {
    series: {
      Children: {
        title: "Kids",
        color: "green",
      },
    },
  },
};

export const SingleSeriesLegend = Template.bind({});
SingleSeriesLegend.storyName = "Single Series with Legend";
SingleSeriesLegend.args = {
  data: {
    columns: [
      { name: "Type", data_type: "TEXT" },
      { name: "Count", data_type: "INT8" },
    ],
    rows: [
      { Type: "User", Count: 12 },
      { Type: "Policy", Count: 93 },
      { Type: "Role", Count: 48 },
    ],
  },
  properties: {
    legend: {
      display: "all",
    },
  },
};

export const SingleSeriesLegendPosition = Template.bind({});
SingleSeriesLegendPosition.storyName = "Single Series With Legend At Bottom";
SingleSeriesLegendPosition.args = {
  data: {
    columns: [
      { name: "Type", data_type: "TEXT" },
      { name: "Count", data_type: "INT8" },
    ],
    rows: [
      { Type: "User", Count: 12 },
      { Type: "Policy", Count: 93 },
      { Type: "Role", Count: 48 },
    ],
  },
  properties: {
    legend: {
      display: "all",
      position: "bottom",
    },
  },
};

export const SingleSeriesXAxisTitle = Template.bind({});
SingleSeriesXAxisTitle.storyName = "Single Series with X Axis Title";
SingleSeriesXAxisTitle.args = {
  data: {
    columns: [
      { name: "Type", data_type: "TEXT" },
      { name: "Count", data_type: "INT8" },
    ],
    rows: [
      { Type: "User", Count: 12 },
      { Type: "Policy", Count: 93 },
      { Type: "Role", Count: 48 },
    ],
  },
  properties: {
    axes: {
      x: {
        title: {
          display: "all",
          value: "I am a the X Axis title",
        },
      },
    },
  },
};

export const SingleSeriesXAxisNoLabels = Template.bind({});
SingleSeriesXAxisNoLabels.storyName = "Single Series with no X Axis Labels";
SingleSeriesXAxisNoLabels.args = {
  data: {
    columns: [
      { name: "Type", data_type: "TEXT" },
      { name: "Count", data_type: "INT8" },
    ],
    rows: [
      { Type: "User", Count: 12 },
      { Type: "Policy", Count: 93 },
      { Type: "Role", Count: 48 },
    ],
  },
  properties: {
    axes: {
      x: {
        labels: {
          display: "none",
        },
      },
    },
  },
};

export const SingleSeriesYAxisNoLabels = Template.bind({});
SingleSeriesYAxisNoLabels.storyName = "Single Series with no Y Axis Labels";
SingleSeriesYAxisNoLabels.args = {
  data: {
    columns: [
      { name: "Type", data_type: "TEXT" },
      { name: "Count", data_type: "INT8" },
    ],
    rows: [
      { Type: "User", Count: 12 },
      { Type: "Policy", Count: 93 },
      { Type: "Role", Count: 48 },
    ],
  },
  properties: {
    axes: {
      y: {
        labels: {
          display: "none",
        },
      },
    },
  },
};

export const TimeSeries = Template.bind({});
TimeSeries.storyName = "Single Time Series";
TimeSeries.args = SingleTimeSeriesDefaults;

export const MultiTimeSeries = Template.bind({});
MultiTimeSeries.storyName = "Multiple Time Series with Overrides (stacked)";
MultiTimeSeries.args = MultiTimeSeriesDefaults;

export const MultiTimeSeriesGrouped = Template.bind({});
MultiTimeSeriesGrouped.storyName = "Multiple Time Series (grouped)";
MultiTimeSeriesGrouped.args = MultiTimeSeriesGroupedDefaults;

export const SingleSeriesDiff = Template.bind({});
SingleSeriesDiff.storyName = "Single Series Diff";
SingleSeriesDiff.args = {
  data: {
    columns: [
      { name: "Type", data_type: "TEXT" },
      { name: "__diff", data_type: "INT8" },
      { name: "Count_diff", data_type: "INT8" },
      { name: "Count", data_type: "INT8" },
    ],
    rows: [
      { Type: "User", Count: 12, Count_diff: 10, _diff: "updated" },
      { Type: "Policy", Count: 93, Count_diff: 100, _diff: "updated" },
      { Type: "Role", Count: 48, Count_diff: 50, _diff: "updated" },
    ],
  },
};

export const SingleTimeSeriesDiff = Template.bind({});
SingleTimeSeriesDiff.storyName = "Single Time Series Diff";
SingleTimeSeriesDiff.args = {
  data: {
    columns: [
      { name: "time", data_type: "TIMESTAMP" },
      { name: "Count", data_type: "INT8" },
      { name: "Count_diff", data_type: "INT8" },
      { name: "__diff", data_type: "TEXT" },
    ],
    rows: [
      {
        time: "2023-01-01T00:00:00.000",
        Count: 20,
        Count_diff: 18,
        _diff: "updated",
      },
      {
        time: "2023-02-01T00:00:00.000",
        Count: 32,
        Count_diff: 40,
        _diff: "updated",
      },
      {
        time: "2023-04-01T00:00:00.000",
        Count: -15,
        Count_diff: -25,
        _diff: "updated",
      },
      {
        time: "2023-05-01T00:00:00.000",
        Count: 18,
        Count_diff: 20,
        _diff: "updated",
      },
      {
        time: "2023-06-01T00:00:00.000",
        Count: -9,
        Count_diff: -4,
        _diff: "updated",
      },
      {
        time: "2023-12-01T00:00:00.000",
        Count: 3,
        Count_diff: 7,
        _diff: "updated",
      },
    ],
  },
};
