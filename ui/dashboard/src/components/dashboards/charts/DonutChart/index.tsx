import Chart from "@powerpipe/components/dashboards/charts/Chart";
import {
  ChartProps,
  IChart,
} from "@powerpipe/components/dashboards/charts/types";
import { registerChartComponent } from "@powerpipe/components/dashboards/charts";

const DonutChart = (props: ChartProps) => {
  return <Chart {...props} />;
};

const definition: IChart = {
  type: "donut",
  component: DonutChart,
};

registerChartComponent(definition.type, definition);

export default definition;
