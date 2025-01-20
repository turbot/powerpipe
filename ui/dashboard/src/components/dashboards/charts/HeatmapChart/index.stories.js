import HeatmapChart from "./index";
import { PanelStoryDecorator } from "@powerpipe/utils/storybook";
import {
  LargeDailySingleTimeSeriesDefaults,
  LargeHourlySingleTimeSeriesDefaults,
} from "../Chart/index.stories";

const story = {
  title: "Charts/Heatmap",
  component: HeatmapChart.component,
};

export default story;

const Template = (args) => (
  <PanelStoryDecorator
    definition={{ ...args, display_type: "heatmap" }}
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

export const DailyData = Template.bind({});
DailyData.storyName = "Daily Data";
DailyData.args = LargeDailySingleTimeSeriesDefaults;

export const HourlyData = Template.bind({});
HourlyData.storyName = "Hourly Data";
HourlyData.args = LargeHourlySingleTimeSeriesDefaults;
