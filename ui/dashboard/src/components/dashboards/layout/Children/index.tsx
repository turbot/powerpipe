import Child from "../Child";
import {
  ContainerDefinition,
  DashboardPanelType,
  PanelDefinition,
} from "@powerpipe/types";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

type ChildrenProps = {
  children: ContainerDefinition[] | PanelDefinition[] | undefined;
  parentType: DashboardPanelType;
  showPanelControls?: boolean;
};

const Children = ({
  children = [],
  parentType,
  showPanelControls = true,
}: ChildrenProps) => {
  const { panelsMap } = useDashboardState();
  return (
    <>
      {children.map((child) => {
        const definition = panelsMap[child.name];
        if (!definition) {
          return null;
        }
        return (
          <Child
            key={definition.name}
            layoutDefinition={child}
            panelDefinition={definition}
            parentType={parentType}
            showPanelControls={showPanelControls}
          />
        );
      })}
    </>
  );
};

export default Children;
