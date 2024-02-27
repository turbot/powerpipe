import Flow from "@powerpipe/components/dashboards/flows/Flow";
import { FlowProps, IFlow } from "@powerpipe/components/dashboards/flows/types";
import { registerFlowComponent } from "@powerpipe/components/dashboards/flows";

const Sankey = (props: FlowProps) => <Flow {...props} />;

const definition: IFlow = {
  type: "sankey",
  component: Sankey,
};

registerFlowComponent(definition.type, definition);

export default definition;
