import Chart from "@powerpipe/components/dashboards/charts/Chart";
import {
  ChartProps,
  IChart,
} from "@powerpipe/components/dashboards/charts/types";
import { registerChartComponent } from "@powerpipe/components/dashboards/charts";

const ColumnChart = (props: ChartProps) => {
  return <Chart {...props} />;
};

const definition: IChart = {
  type: "column",
  component: ColumnChart,
};

registerChartComponent(definition.type, definition);

export default definition;
