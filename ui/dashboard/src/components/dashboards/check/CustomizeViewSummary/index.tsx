import useCheckFilterConfig from "@powerpipe/hooks/useCheckFilterConfig";
import useCheckGroupingConfig from "@powerpipe/hooks/useCheckGroupingConfig";
import { validateFilter } from "@powerpipe/components/dashboards/check/CheckFilterEditor";

const CustomizeViewSummary = () => {
  const filterConfig = useCheckFilterConfig();
  const groupingConfig = useCheckGroupingConfig();

  const filterCount = filterConfig?.expressions?.length
    ? filterConfig.expressions.filter(validateFilter).length
    : 0;

  return <span>Filter & Group{!!filterCount ? ": On" : null}</span>;
};

export default CustomizeViewSummary;
