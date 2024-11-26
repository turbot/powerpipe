import FilterEditor from "@powerpipe/components/dashboards/grouping/FilterEditor";
import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import { useDashboard } from "@powerpipe/hooks/useDashboard";

const FilterConfig = ({ panelName }: { panelName: string }) => {
  const { panelsMap } = useDashboard();
  const panel = panelsMap[panelName];
  const { filter: filterConfig, update } = useFilterConfig(panelName);
  return (
    <FilterEditor
      filter={filterConfig}
      panelType={panel?.panel_type}
      onApply={update}
    />
  );
};

export default FilterConfig;
