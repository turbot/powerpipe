import Child from "@powerpipe/components/dashboards/layout/Child";
import { PanelDefinition } from "@powerpipe/types";
import { PanelDetailProps } from "@powerpipe/components/dashboards/layout/PanelDetail";

const PanelDetailPreview = ({
  definition: { children, name, panel_type, title, ...rest },
}: PanelDetailProps) => {
  const layoutDefinition = { children, name, panel_type };
  const panelDefinition = {
    name,
    panel_type,
    width: 12,
    ...rest,
  } as PanelDefinition;
  return (
    <Child
      layoutDefinition={layoutDefinition}
      panelDefinition={panelDefinition}
      parentType="dashboard"
      showPanelControls={false}
    />
  );
};

export default PanelDetailPreview;
