import FilterConfig from "@powerpipe/components/dashboards/grouping/FilterConfig";
import GroupingConfig from "@powerpipe/components/dashboards/grouping/GroupingConfig";
import Icon from "@powerpipe/components/Icon";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";

const FilterAndGroupSidePanel = ({ panelName }: { panelName: string }) => {
  const { closeSidePanel } = useDashboardPanelDetail();
  return (
    <>
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
        <FilterConfig panelName={panelName} />
      </div>
      <div className="p-4 space-y-3">
        <span className="font-semibold">Group</span>
        <GroupingConfig panelName={panelName} />
      </div>
    </>
  );
};

export default FilterAndGroupSidePanel;
