import useGroupingFilterConfig from "@powerpipe/hooks/useGroupingFilterConfig";
// import useCheckGroupingConfig from "@powerpipe/hooks/useCheckGroupingConfig";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/CheckFilterEditor";

const CustomizeViewSummary = () => {
  const filterConfig = useGroupingFilterConfig();
  // const groupingConfig = useCheckGroupingConfig();

  const filterCount = filterConfig?.expressions?.length
    ? filterConfig.expressions.filter(validateFilter).length
    : 0;

  return <span>Filter & Group{!!filterCount ? ": On" : null}</span>;
};

export default CustomizeViewSummary;
