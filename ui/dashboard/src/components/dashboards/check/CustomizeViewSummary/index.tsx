import useDashboardSearchPathPrefix from "@powerpipe/hooks/useDashboardSearchPathPrefix";
import useCheckFilterConfig from "@powerpipe/hooks/useCheckFilterConfig";
import useCheckGroupingConfig from "@powerpipe/hooks/useCheckGroupingConfig";

const CustomizeViewSummary = () => {
  const searchPathPrefix = useDashboardSearchPathPrefix();
  const filterConfig = useCheckFilterConfig();
  const groupingConfig = useCheckGroupingConfig();

  return null;
};

export default CustomizeViewSummary;
