import ErrorPanel from "@powerpipe/components/dashboards/Error";
import merge from "lodash/merge";
import useChartThemeColors from "@powerpipe/hooks/useChartThemeColors";
import useNodeAndEdgeData from "@powerpipe/components/dashboards/common/useNodeAndEdgeData";
import {
  buildNodesAndEdges,
  buildSankeyDataInputs,
  LeafNodeData,
  toEChartsType,
} from "@powerpipe/components/dashboards/common";
import { Chart } from "@powerpipe/components/dashboards/charts/Chart";
import {
  FlowProperties,
  FlowProps,
  FlowType,
} from "@powerpipe/components/dashboards/flows/types";
import { getFlowComponent } from "..";
import { NodesAndEdges } from "@powerpipe/components/dashboards/common/types";
import { registerComponent } from "@powerpipe/components/dashboards";
import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { useDashboardTheme } from "@powerpipe/hooks/useDashboardTheme";

const getCommonBaseOptions = () => ({
  animation: false,
  tooltip: {
    trigger: "item",
    triggerOn: "mousemove",
  },
});

const getCommonBaseOptionsForFlowType = (type: FlowType = "sankey") => {
  switch (type) {
    case "sankey":
      return {};
    default:
      return {};
  }
};

const getSeriesForFlowType = (
  type: FlowType = "sankey",
  data: LeafNodeData | undefined,
  properties: FlowProperties | undefined,
  nodesAndEdges: NodesAndEdges,
  themeColors,
) => {
  if (!data) {
    return {};
  }
  const series: any[] = [];
  const seriesLength = 1;
  for (let seriesIndex = 0; seriesIndex < seriesLength; seriesIndex++) {
    switch (type) {
      case "sankey": {
        const { data: sankeyData, links } = buildSankeyDataInputs(
          nodesAndEdges,
          themeColors,
        );
        series.push({
          type: toEChartsType(type),
          layout: "none",
          draggable: true,
          label: { color: themeColors.foreground, formatter: "{b}" },
          emphasis: {
            focus: "adjacency",
            blurScope: "coordinateSystem",
          },
          lineStyle: {
            color: "source",
            curveness: 0.5,
          },
          data: sankeyData,
          links,
          tooltip: {
            formatter: "{b}",
          },
        });
        break;
      }
    }
  }

  return { series };
};

const getOptionOverridesForFlowType = (
  type: FlowType = "sankey",
  properties: FlowProperties | undefined,
) => {
  if (!properties) {
    return {};
  }

  return {};
};

const buildFlowOptions = (props: FlowProps, themeColors) => {
  const nodesAndEdges = buildNodesAndEdges(
    props.categories,
    props.data,
    props.properties,
    themeColors,
  );

  return merge(
    getCommonBaseOptions(),
    getCommonBaseOptionsForFlowType(props.display_type),
    getSeriesForFlowType(
      props.display_type,
      props.data,
      props.properties,
      nodesAndEdges,
      themeColors,
    ),
    getOptionOverridesForFlowType(props.display_type, props.properties),
  );
};

const FlowWrapper = (props: FlowProps) => {
  const themeColors = useChartThemeColors();
  const { searchPathPrefix } = useDashboardSearchPath();
  const { wrapperRef } = useDashboardTheme();
  const nodeAndEdgeData = useNodeAndEdgeData(
    props.data,
    props.properties,
    props.status,
  );

  if (!wrapperRef) {
    return null;
  }

  if (
    !nodeAndEdgeData ||
    !nodeAndEdgeData.data ||
    !nodeAndEdgeData.data.rows ||
    nodeAndEdgeData.data.rows.length === 0
  ) {
    return null;
  }

  return (
    <Chart
      options={buildFlowOptions(
        {
          ...props,
          categories: nodeAndEdgeData.categories,
          data: nodeAndEdgeData.data,
          properties: nodeAndEdgeData.properties,
        },
        themeColors,
      )}
      searchPathPrefix={searchPathPrefix}
      type={props.display_type || "sankey"}
    />
  );
};

const renderFlow = (definition: FlowProps) => {
  // We default to sankey diagram if not specified
  const { display_type = "sankey" } = definition;

  const flow = getFlowComponent(display_type);

  if (!flow) {
    return <ErrorPanel error={`Unknown flow type ${display_type}`} />;
  }

  const Component = flow.component;
  return <Component {...definition} />;
};

const RenderFlow = (props: FlowProps) => {
  return renderFlow(props);
};

registerComponent("flow", RenderFlow);

export default FlowWrapper;
