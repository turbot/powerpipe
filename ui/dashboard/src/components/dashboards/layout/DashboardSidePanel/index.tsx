import FilterAndGroupSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/FilterAndGroupSidePanel";
import TableRowSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/TableRowSidePanel";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";

const DashboardSidePanel = () => {
  const { selectedSidePanel } = useDashboardPanelDetail();

  if (!selectedSidePanel) {
    return null;
  }

  return (
    <div className="h-full bg-dashboard-panel divide-y divide-divide print:hidden overflow-y-auto">
      {selectedSidePanel.panel.panel_type === "benchmark" ||
        selectedSidePanel.panel.panel_type === "control" ||
        (selectedSidePanel.panel.panel_type === "detection" && (
          <FilterAndGroupSidePanel panelName={selectedSidePanel.panel.name} />
        ))}
      {selectedSidePanel.panel.panel_type === "table" && (
        <TableRowSidePanel
          data={selectedSidePanel.panel.data}
          rowIndex={selectedSidePanel.context.rowIndex}
        />
      )}
    </div>
  );
};

export default DashboardSidePanel;
