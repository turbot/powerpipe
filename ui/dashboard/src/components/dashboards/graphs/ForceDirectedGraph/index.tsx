import Graph from "@powerpipe/components/dashboards/graphs/Graph";
import {
  GraphProps,
  IGraph,
} from "@powerpipe/components/dashboards/graphs/types";
import { registerGraphComponent } from "@powerpipe/components/dashboards/graphs";

const ForceDirectedGraph = (props: GraphProps) => <Graph {...props} />;

const definition: IGraph = {
  type: "graph",
  component: ForceDirectedGraph,
};

registerGraphComponent(definition.type, definition);

export default definition;
