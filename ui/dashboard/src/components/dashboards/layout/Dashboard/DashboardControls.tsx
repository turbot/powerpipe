import FilterConfig from "@powerpipe/components/dashboards/grouping/FilterConfig";
import GroupingConfig from "@powerpipe/components/dashboards/grouping/GroupingConfig";
import Icon from "@powerpipe/components/Icon";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";

const DashboardControls = ({ panelName }: { panelName: string | null }) => {
  const { closeFilterAndGroupPanel } = useDashboardPanelDetail();

  if (!panelName) {
    return null;
  }

  return (
    <div className="h-full bg-dashboard-panel divide-y divide-divide print:hidden overflow-y-auto">
      <div className="flex items-center justify-between p-4">
        <h3>Filter & Group</h3>
        <Icon
          className="w-5 h-5 text-foreground cursor-pointer hover:text-foreground-light shrink-0"
          icon="close"
          onClick={closeFilterAndGroupPanel}
          title="Close customize view"
        />
      </div>
      <div className="p-4 space-y-3">
        <span className="font-semibold">Filter</span>
        <FilterConfig panelName={panelName} />
      </div>
      <div className="p-4 space-y-3">
        <span className="font-semibold">Group</span>
        <GroupingConfig panelName={panelName} />
      </div>
    </div>
  );
};

export default DashboardControls;
