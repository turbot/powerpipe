import LineChart from "./index";
import { PanelStoryDecorator } from "@powerpipe/utils/storybook";
import {
  MultiTimeSeriesDefaults,
  SingleTimeSeriesDefaults,
  MultiTimeSeriesCrosstabDefaults,
} from "../Chart/index.stories";

const story = {
  title: "Charts/Line",
  component: LineChart.component,
};

export default story;

const Template = (args) => (
  <PanelStoryDecorator
    definition={{ ...args, display_type: "line" }}
    panelType="chart"
  />
);

export const Loading = Template.bind({});
Loading.args = {
  loading: true,
};

export const Error = Template.bind({});
Error.args = {
  loading: false,
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

export const MultiSeries = Template.bind({});
MultiSeries.storyName = "Multi-Series";
MultiSeries.args = {
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

// Line charts cannot stack, use Area for that
export const MultiTimeSeries = Template.bind({});
MultiTimeSeries.storyName = "Multiple Time Series";
MultiTimeSeries.args = MultiTimeSeriesDefaults;

export const MultiTimeSeriesCrosstab = Template.bind({});
MultiTimeSeriesCrosstab.storyName = "Multiple Time Series (crosstabbed)";
MultiTimeSeriesCrosstab.args = MultiTimeSeriesCrosstabDefaults;

export const SingleSeriesDiff = Template.bind({});
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
