import FilterEditor from "@powerpipe/components/dashboards/grouping/FilterEditor";
import useFilterConfig from "@powerpipe/hooks/useFilterConfig";

const FilterConfig = ({ panelName }: { panelName: string }) => {
  const { filter: filterConfig, update } = useFilterConfig(panelName);
  return <FilterEditor filter={filterConfig} onApply={update} />;
};

export default FilterConfig;

const f1 = {
  view: {
    filter_by: {
      panel_1: {}, // some filter
      panel_2: {}, // some filter
    },
  },
};

const f2 = {
  view: {
    panel_1: {
      filter_by: {}, // some filter
    },
    panel_2: {
      filter_by: {}, // some filter
    },
  },
};
