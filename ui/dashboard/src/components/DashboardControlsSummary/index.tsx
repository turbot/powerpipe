import {
  CheckDisplayGroup,
  CheckFilter,
} from "components/dashboards/check/common";
import { Noop } from "types/func";
import { ReactNode } from "react";

type DashboardControlsSummaryProps = {
  searchPathPrefix: string[];
  filterConfig: CheckFilter;
  groupingConfig: CheckDisplayGroup[];
  toggleControls: Noop;
};

type DashboardFilterButtonProps = {
  children: ReactNode;
  onClick: Noop;
};

const DashboardFilterButton = ({
  children,
  onClick,
}: DashboardFilterButtonProps) => (
  <div
    className="border border-black-scale-4 px-2 py-1 rounded-md"
    onClick={onClick}
  >
    {children}
  </div>
);

const DashboardControlsSummary = ({
  toggleControls,
}: DashboardControlsSummaryProps) => {
  return (
    <>
      <div className="grow flex items-center justify-end space-x-4">
        <DashboardFilterButton onClick={toggleControls}>
          Search Path
        </DashboardFilterButton>
        <DashboardFilterButton onClick={toggleControls}>
          Filters
        </DashboardFilterButton>
        <DashboardFilterButton onClick={toggleControls}>
          Grouping
        </DashboardFilterButton>
      </div>
    </>
  );
};

export default DashboardControlsSummary;
