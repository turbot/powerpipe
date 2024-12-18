import FilterAndGroupSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/FilterAndGroupSidePanel";
import TableRowSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/TableRowSidePanel";
import { SidePanelInfo } from "@powerpipe/hooks/useDashboardPanelDetail";

const DashboardSidePanel = ({ sidePanel }: { sidePanel: SidePanelInfo }) => (
  <div className="h-full w-full bg-dashboard-panel divide-y divide-divide print:hidden overflow-y-auto">
    {(sidePanel.panel.panel_type === "benchmark" ||
      sidePanel.panel.panel_type === "control" ||
      sidePanel.panel.panel_type === "detection") && (
      <FilterAndGroupSidePanel panelName={sidePanel.panel.name} />
    )}
    {sidePanel.panel.panel_type === "table" && (
      <TableRowSidePanel
        data={sidePanel.panel.data}
        requestedColumnName={sidePanel.context.requestedColumnName}
        rowIndex={sidePanel.context.rowIndex}
      />
    )}
  </div>
);

export default DashboardSidePanel;
