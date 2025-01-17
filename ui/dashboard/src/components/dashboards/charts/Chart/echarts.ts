import * as echarts from "echarts/core";
import {
  BarChart,
  GraphChart,
  HeatmapChart,
  LineChart,
  PieChart,
  SankeyChart,
  TreeChart,
} from "echarts/charts";
import { CanvasRenderer } from "echarts/renderers";
import {
  CalendarComponent,
  DatasetComponent,
  GridComponent,
  LegendComponent,
  TitleComponent,
  TooltipComponent,
  MarkLineComponent,
  VisualMapComponent,
} from "echarts/components";
import { LabelLayout } from "echarts/features";

echarts.use([
  BarChart,
  CalendarComponent,
  CanvasRenderer,
  DatasetComponent,
  GraphChart,
  GridComponent,
  HeatmapChart,
  LabelLayout,
  LegendComponent,
  LineChart,
  MarkLineComponent,
  PieChart,
  SankeyChart,
  TitleComponent,
  TooltipComponent,
  TreeChart,
  VisualMapComponent,
]);

export { echarts };
