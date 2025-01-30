import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/FilterEditor";

const CustomizeViewSummary = ({ panelName }: { panelName: string }) => {
  const { filter: filterConfig } = useFilterConfig(panelName);

  const filterCount = filterConfig?.expressions?.length
    ? filterConfig.expressions.filter(validateFilter).length
    : 0;

  return (
    <span>
      Filter & Group
      {!!filterCount ? ": On" : null}
    </span>
  );
};

export default CustomizeViewSummary;
