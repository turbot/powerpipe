import FilterAndGroupSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/FilterAndGroupSidePanel";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";

const DashboardSidePanel = () => {
  const { selectedSidePanel } = useDashboardPanelDetail();

  if (!selectedSidePanel) {
    return null;
  }

  if (selectedSidePanel.viewType === "group_and_filter") {
    return <FilterAndGroupSidePanel sidePanel={selectedSidePanel} />;
  }

  if (selectedSidePanel.viewType === "table_row") {
    return null;
  }

  return null;
};

export default DashboardSidePanel;
