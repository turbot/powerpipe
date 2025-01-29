// Ensure Table is loaded & registered first
import "@powerpipe/components/dashboards/Table";
import {
  DashboardLayoutNode,
  DashboardPanelType,
  PanelDefinition,
} from "@powerpipe/types";
import { getComponent } from "@powerpipe/components/dashboards";
import { getNodeAndEdgeDataFormat } from "@powerpipe/components/dashboards/common/useNodeAndEdgeData";
import { NodeAndEdgeProperties } from "@powerpipe/components/dashboards/common/types";

type ChildProps = {
  layoutDefinition: DashboardLayoutNode;
  panelDefinition: PanelDefinition;
  parentType: DashboardPanelType;
  showPanelControls?: boolean;
};

const Child = ({
  layoutDefinition,
  panelDefinition,
  parentType,
  showPanelControls = true,
}: ChildProps) => {
  const Panel = getComponent("panel");
  switch (layoutDefinition.panel_type) {
    case "benchmark":
      if (panelDefinition.benchmark_type === "detection") {
        const DetectionBenchmark = getComponent("detection_benchmark");
        return (
          <DetectionBenchmark
            definition={panelDefinition}
            benchmarkChildren={layoutDefinition.children}
            showControls={showPanelControls}
          />
        );
      } else {
        const Benchmark = getComponent("benchmark");
        return (
          <Benchmark
            definition={panelDefinition}
            benchmarkChildren={layoutDefinition.children}
            showControls={showPanelControls}
          />
        );
      }
    case "control":
      const Benchmark = getComponent("benchmark");
      return (
        <Benchmark
          definition={panelDefinition}
          showControls={showPanelControls}
        />
      );
    case "detection":
      const DetectionBenchmark = getComponent("detection_benchmark");
      return (
        <DetectionBenchmark
          definition={panelDefinition}
          showControls={showPanelControls}
        />
      );
    case "card":
      const Card = getComponent("card");
      return (
        <Panel
          definition={panelDefinition}
          parentType={parentType}
          showControls={showPanelControls}
          showPanelStatus={false}
        >
          <Card {...panelDefinition} />
        </Panel>
      );
    case "chart":
      const Chart = getComponent("chart");
      return (
        <Panel
          definition={panelDefinition}
          parentType={parentType}
          showControls={showPanelControls}
        >
          <Chart {...panelDefinition} />
        </Panel>
      );
    case "container":
      const Container = getComponent("container");
      return <Container layoutDefinition={layoutDefinition} />;
    case "dashboard":
      const Dashboard = getComponent("dashboard");
      return <Dashboard definition={panelDefinition} isRoot={false} />;
    case "error":
      const ErrorPanel = getComponent("error");
      return (
        <Panel
          definition={panelDefinition}
          parentType={parentType}
          showControls={showPanelControls}
        >
          <ErrorPanel {...panelDefinition} />
        </Panel>
      );
    case "flow": {
      const Flow = getComponent("flow");
      const format = getNodeAndEdgeDataFormat(
        panelDefinition.properties as NodeAndEdgeProperties,
      );
      return (
        <Panel
          definition={panelDefinition}
          parentType={parentType}
          showControls={showPanelControls}
          showPanelStatus={
            format === "LEGACY" ||
            panelDefinition.status === "cancelled" ||
            panelDefinition.status === "error"
          }
          // Node and edge format will show error info on the panel information
          showPanelError={format === "LEGACY"}
        >
          <Flow {...panelDefinition} />
        </Panel>
      );
    }
    case "graph": {
      const Graph = getComponent("graph");
      const format = getNodeAndEdgeDataFormat(
        panelDefinition.properties as NodeAndEdgeProperties,
      );
      return (
        <Panel
          definition={panelDefinition}
          parentType={parentType}
          showControls={showPanelControls}
          showPanelStatus={
            format === "LEGACY" ||
            panelDefinition.status === "cancelled" ||
            panelDefinition.status === "error"
          }
          // Node and edge format will show error info on the panel information
          showPanelError={format === "LEGACY"}
        >
          <Graph {...panelDefinition} />
        </Panel>
      );
    }
    case "hierarchy": {
      const Hierarchy = getComponent("hierarchy");
      const format = getNodeAndEdgeDataFormat(
        panelDefinition.properties as NodeAndEdgeProperties,
      );
      return (
        <Panel
          definition={panelDefinition}
          parentType={parentType}
          showControls={showPanelControls}
          showPanelStatus={
            format === "LEGACY" ||
            panelDefinition.status === "cancelled" ||
            panelDefinition.status === "error"
          }
          // Node and edge format will show error info on the panel information
          showPanelError={format === "LEGACY"}
        >
          <Hierarchy {...panelDefinition} />
        </Panel>
      );
    }
    case "image":
      const Image = getComponent("image");
      return (
        <Panel
          definition={panelDefinition}
          parentType={parentType}
          showControls={showPanelControls}
        >
          <Image {...panelDefinition} />
        </Panel>
      );
    case "input":
      const Input = getComponent("input");
      return (
        <Panel
          definition={panelDefinition}
          parentType={parentType}
          showControls={
            showPanelControls &&
            (panelDefinition.title || panelDefinition.display_type === "table")
          }
          showPanelStatus={false}
        >
          <Input {...panelDefinition} />
        </Panel>
      );
    case "table":
      const Table = getComponent("table");
      return (
        <Panel
          definition={panelDefinition}
          parentType={parentType}
          showControls={showPanelControls}
        >
          <Table {...panelDefinition} />
        </Panel>
      );
    case "text":
      const Text = getComponent("text");
      return (
        <Panel
          definition={panelDefinition}
          parentType={parentType}
          showControls={false}
        >
          <Text {...panelDefinition} />
        </Panel>
      );
    default:
      return null;
  }
};

export default Child;
