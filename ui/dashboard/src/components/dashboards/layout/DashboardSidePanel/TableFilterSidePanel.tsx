import FilterConfig from "@powerpipe/components/dashboards/grouping/FilterConfig";
import Icon from "@powerpipe/components/Icon";
import { useDashboardControls } from "@powerpipe/components/dashboards/layout/Dashboard/DashboardControlsProvider";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect } from "react";

const TableFilterSidePanel = ({ panelName }: { panelName: string }) => {
  const { panelsMap } = useDashboardState();
  const { closeSidePanel } = useDashboardPanelDetail();
  const { setContext } = useDashboardControls();
  const panel = panelsMap[panelName];

  useEffect(() => {
    if (!panel.data) {
      return;
    }
    const filterValues = { dimension: { key: {}, value: {} } };
    for (const column of panel.data?.columns || []) {
      for (const row of panel.data?.rows || []) {
        if (!filterValues.dimension.key[column.name]) {
          filterValues.dimension.key[column.name] = {};
        }
        const rowValue = row[column.name];
        if (!filterValues.dimension.key[column.name][rowValue]) {
          filterValues.dimension.key[column.name][rowValue] = 1;
        } else {
          filterValues.dimension.key[column.name][rowValue] += 1;
        }

        if (!filterValues.dimension.value[rowValue]) {
          filterValues.dimension.value[rowValue] = {};
        }
        if (!filterValues.dimension.value[rowValue][column.name]) {
          filterValues.dimension.value[rowValue][column.name] = 1;
        } else {
          filterValues.dimension.value[rowValue][column.name] += 1;
        }
      }
    }
    setContext(() => filterValues);
  }, [panel.data]);

  return (
    <>
      <div className="flex items-center justify-between p-4 min-w-[520px] max-w-[1000px]">
        <h3>Filter</h3>
        <Icon
          className="w-5 h-5 text-foreground cursor-pointer hover:text-foreground-light shrink-0"
          icon="close"
          onClick={closeSidePanel}
          title="Close"
        />
      </div>
      <div className="w-full max-h-full border-t border-divide overflow-y-scroll">
        <div className="p-4 space-y-3">
          <FilterConfig panelName={panelName} />
        </div>
      </div>
    </>
  );
};

export default TableFilterSidePanel;
