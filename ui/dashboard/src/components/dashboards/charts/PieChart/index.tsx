import Chart from "@powerpipe/components/dashboards/charts/Chart";
import {
  ChartProps,
  IChart,
} from "@powerpipe/components/dashboards/charts/types";
import { registerChartComponent } from "@powerpipe/components/dashboards/charts";

const PieChart = (props: ChartProps) => {
  return <Chart {...props} />;
};

const definition: IChart = {
  type: "pie",
  component: PieChart,
};

registerChartComponent(definition.type, definition);

export default definition;
