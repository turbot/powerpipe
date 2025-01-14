import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/FilterEditor";
import { ReactNode } from "react";

const CustomizeViewSummary = ({
  panelName,
  title = "Filter & Group",
}: {
  panelName: string;
  title: ReactNode;
}) => {
  const { filter: filterConfig } = useFilterConfig(panelName);

  const filterCount = filterConfig?.expressions?.length
    ? filterConfig.expressions.filter(validateFilter).length
    : 0;

  return (
    <span>
      {title}
      {!!filterCount ? ": On" : null}
    </span>
  );
};

export default CustomizeViewSummary;
