// Ensure Table is loaded & registered first
import "@powerpipe/components/dashboards/Table";
import Panel from "../Panel";
import { getComponent } from "../../index";
import { PanelDetailProps } from "./index";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";

const Table = getComponent("table");

const PanelDetailData = ({ definition }: PanelDetailProps) => {
  const { panelOverrideData } = useDashboardPanelDetail();
  return (
    <Panel
      definition={definition}
      parentType="dashboard"
      showControls={false}
      forceBackground={true}
    >
      <Table
        name={`${definition}.table.detail`}
        panel_type="table"
        data={panelOverrideData || definition.data}
      />
    </Panel>
  );
};

export default PanelDetailData;
