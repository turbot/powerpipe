import FilterConfig from "@powerpipe/components/dashboards/grouping/FilterConfig";
import GroupingConfig from "@powerpipe/components/dashboards/grouping/GroupingConfig";
import Icon from "@powerpipe/components/Icon";
import {
  SidePanelInfo,
  useDashboardPanelDetail,
} from "@powerpipe/hooks/useDashboardPanelDetail";

const FilterAndGroupSidePanel = ({
  sidePanel,
}: {
  sidePanel: SidePanelInfo;
}) => {
  const { closeSidePanel } = useDashboardPanelDetail();
  return (
    <div className="h-full bg-dashboard-panel divide-y divide-divide print:hidden overflow-y-auto">
      <div className="flex items-center justify-between p-4">
        <h3>Filter & Group</h3>
        <Icon
          className="w-5 h-5 text-foreground cursor-pointer hover:text-foreground-light shrink-0"
          icon="close"
          onClick={closeSidePanel}
          title="Close customize view"
        />
      </div>
      <div className="p-4 space-y-3">
        <span className="font-semibold">Filter</span>
        <FilterConfig panelName={sidePanel.name} />
      </div>
      <div className="p-4 space-y-3">
        <span className="font-semibold">Group</span>
        <GroupingConfig panelName={sidePanel.name} />
      </div>
    </div>
  );
};

export default FilterAndGroupSidePanel;
