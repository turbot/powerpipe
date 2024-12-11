import {
  CheckDisplayGroup,
  Filter,
} from "@powerpipe/components/dashboards/grouping/common";
import { classNames } from "@powerpipe/utils/styles";
import { Noop } from "@powerpipe/types/func";
import { ReactNode } from "react";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/FilterEditor";

type DashboardControlsSummaryProps = {
  filterConfig: Filter;
  groupingConfig: CheckDisplayGroup[];
  onClose: Noop;
};

type DashboardFilterButtonCountProps = {
  count: number;
};

type DashboardFilterButtonProps = {
  children: ReactNode;
  className?: string;
  onClick: Noop;
};

const DashboardFilterButtonCount = ({
  count,
}: DashboardFilterButtonCountProps) => {
  if (!count) {
    return null;
  }

  return (
    <span className="bg-info bg-opacity-20 text-info text-sm px-1.5 py-0.5 rounded-md">
      {count}
    </span>
  );
};

const DashboardFilterButton = ({
  children,
  className,
  onClick,
}: DashboardFilterButtonProps) => (
  <button
    className={classNames(
      "border border-black-scale-3 px-2.5 py-1.5 whitespace-nowrap rounded-md cursor-pointer hover:bg-dashboard",
      className,
    )}
    onClick={onClick}
  >
    {children}
  </button>
);

const DashboardFilterControlButton = ({ filterConfig, toggleControls }) => {
  const filterCount = filterConfig?.expressions?.length
    ? filterConfig.expressions.filter(validateFilter).length
    : 0;
  // <div className="flex items-center space-x-3 shrink-0">
  //   <Icon className="h-5 w-5" icon="filter_list" />
  //   {filterConfig.operator === "and" &&
  //     !!filterConfig.expressions &&
  //     filterConfig.expressions.length > 0 && (
  //       <div className="space-x-2">{filtersToText(filterConfig)}</div>
  //     )}
  //   {filterConfig.operator === "and" &&
  //     (!filterConfig.expressions ||
  //       filterConfig.expressions.length === 0) && (
  //       <span className="text-foreground-lighter">No filters</span>
  //     )}
  //   {!showEditor && (
  //     <Icon
  //       className="h-5 w-5 cursor-pointer shrink-0"
  //       icon="edit_square"
  //       onClick={() => setShowEditor(true)}
  //       title="Edit filter"
  //     />
  //   )}
  // </div>

  return (
    <DashboardFilterButton className="block space-x-1" onClick={toggleControls}>
      <DashboardFilterButtonCount count={filterCount} /> <span>Filters</span>
    </DashboardFilterButton>
  );
};

const DashboardGroupingControlButton = ({ toggleControls }) => {
  return (
    <DashboardFilterButton onClick={toggleControls}>
      Grouping
    </DashboardFilterButton>
  );
};

const DashboardControlsSummary = ({
  filterConfig,
  groupingConfig,
  onClose,
}: DashboardControlsSummaryProps) => {
  return (
    <>
      <div className="grow flex items-center justify-end space-x-4">
        <DashboardFilterControlButton
          filterConfig={filterConfig}
          toggleControls={onClose}
        />
        <DashboardGroupingControlButton
          groupingConfig={groupingConfig}
          toggleControls={onClose}
        />
      </div>
    </>
  );
};

export default DashboardControlsSummary;
