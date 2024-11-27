import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import {
  Filter,
  summaryCardFilterPath,
} from "@powerpipe/components/dashboards/grouping/common";
import { Link, useLocation } from "react-router-dom";
import { ReactNode } from "react";

const FilterCard = ({
  children,
  cardName,
  panelName,
  dimension,
  expressions,
}: {
  children: ReactNode;
  cardName: string;
  panelName: string;
  dimension: string;
  expressions: Filter[] | undefined;
}) => {
  const { allFilters } = useFilterConfig();
  const { pathname, search } = useLocation();

  const getSeverityCardMetric = () => {
    const parts = cardName.split(".");
    const type = parts[parts.length - 2];

    if (type !== dimension) {
      return null;
    }

    return parts[parts.length - 1];
  };

  const metric = getSeverityCardMetric();

  if (!metric) {
    return <>{children}</>;
  }

  return (
    <Link
      className="cursor-pointer"
      to={summaryCardFilterPath({
        allFilters,
        expressions,
        panelName,
        pathname,
        search,
        dimension,
        metric,
      })}
    >
      {children}
    </Link>
  );
};

export default FilterCard;
