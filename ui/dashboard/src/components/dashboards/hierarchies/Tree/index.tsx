import Hierarchy from "@powerpipe/components/dashboards/hierarchies/Hierarchy";
import {
  HierarchyProps,
  IHierarchy,
} from "@powerpipe/components/dashboards/hierarchies/types";
import { registerHierarchyComponent } from "@powerpipe/components/dashboards/hierarchies";

const Tree = (props: HierarchyProps) => <Hierarchy {...props} />;

const definition: IHierarchy = {
  type: "tree",
  component: Tree,
};

registerHierarchyComponent(definition.type, definition);

export default definition;
