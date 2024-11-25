import FilterConfig from "@powerpipe/components/dashboards/grouping/FilterConfig";
import GroupingConfig from "@powerpipe/components/dashboards/grouping/GroupingConfig";
import Icon from "@powerpipe/components/Icon";
import { DashboardActions } from "@powerpipe/types";
import { useDashboard } from "@powerpipe/hooks/useDashboard";

const DashboardControls = ({ panelName }: { panelName: string }) => {
  const { dispatch } = useDashboard();
  const hideControls = () =>
    dispatch({ type: DashboardActions.HIDE_CUSTOMIZE_BENCHMARK_PANEL });

  return (
    <div className="h-full bg-dashboard-panel divide-y divide-divide print:hidden overflow-y-auto">
      <div className="flex items-center justify-between p-4">
        <h3>Filter & Group</h3>
        <Icon
          className="w-5 h-5 text-foreground cursor-pointer hover:text-foreground-light shrink-0"
          icon="close"
          onClick={hideControls}
          title="Close customize view"
        />
      </div>
      <div className="p-4 space-y-3">
        <span className="font-semibold">Filter</span>
        <FilterConfig panelName={panelName} />
      </div>
      <div className="p-4 space-y-3">
        <span className="font-semibold">Group</span>
        <GroupingConfig panelName={panelName} onClose={hideControls} />
      </div>
    </div>
  );
};

export default DashboardControls;
