import FilterConfig from "@powerpipe/components/dashboards/grouping/FilterConfig";
import GroupingConfig from "@powerpipe/components/dashboards/grouping/GroupingConfig";
import Icon from "@powerpipe/components/Icon";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";

const FilterAndGroupSidePanel = ({ panelName }: { panelName: string }) => {
  const { closeSidePanel } = useDashboardPanelDetail();
  return (
    <div className="h-full md:min-w-[500px] md:max-w-[1000px]">
      <div className="flex items-center justify-between p-4">
        <h3>Filter & Group</h3>
        <Icon
          className="w-5 h-5 text-foreground cursor-pointer hover:text-foreground-light shrink-0"
          icon="close"
          onClick={closeSidePanel}
          title="Close"
        />
      </div>
      <div className="flex-1 p-4 space-y-3">
        <span className="font-semibold">Filter</span>
        <FilterConfig panelName={panelName} />
      </div>
      <div className="flex-1 p-4 space-y-3">
        <span className="font-semibold">Group</span>
        <GroupingConfig panelName={panelName} />
      </div>
    </div>
  );
};

export default FilterAndGroupSidePanel;
