import Icon from "@powerpipe/components/Icon";
import useDashboardSearchPathPrefix from "@powerpipe/hooks/useDashboardSearchPathPrefix";
import useCheckFilterConfig from "@powerpipe/hooks/useCheckFilterConfig";
import useCheckGroupingConfig from "@powerpipe/hooks/useCheckGroupingConfig";
import { validateFilter } from "@powerpipe/components/dashboards/check/CheckFilterEditor";

type CustomizeViewFilterButtonCountProps = {
  count: number;
};

const CustomizeViewFilterButtonCount = ({
  count,
}: CustomizeViewFilterButtonCountProps) => {
  if (!count) {
    return null;
  }

  return (
    <span className="bg-info bg-opacity-20 text-info text-sm px-1.5 py-0.5 rounded-md">
      {count}
    </span>
  );
};

const CustomizeViewSummary = () => {
  const searchPathPrefix = useDashboardSearchPathPrefix();
  const filterConfig = useCheckFilterConfig();
  const groupingConfig = useCheckGroupingConfig();

  const filterCount = filterConfig?.expressions?.length
    ? filterConfig.expressions.filter(validateFilter).length
    : 0;

  return (
    <div className="flex items-center space-x-2">
      {!!filterCount ? (
        <>
          <span>Filters</span>
          <CustomizeViewFilterButtonCount count={filterCount} />
        </>
      ) : (
        <>
          <span>Customize</span>
          <Icon className="w-4.5 h-4.5" icon="design_services" />
        </>
      )}
    </div>
  );
};

export default CustomizeViewSummary;
