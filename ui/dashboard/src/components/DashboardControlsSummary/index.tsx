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
    className="border border-black-scale-3 px-2 py-1 rounded-md cursor-pointer hover:bg-dashboard"
    onClick={onClick}
  >
    {children}
  </div>
);

const DashboardSearchPathPrefixControlButton = ({
  searchPathPrefix,
  toggleControls,
}) => {
  // <div className="flex items-center space-x-3 shrink-0">
  //   <Icon className="h-5 w-5 shrink-0" icon="list" />
  //   <div className="space-x-0.5 truncate">
  //     {modifiedSearchPath.length > 0 &&
  //       modifiedSearchPath
  //         .map<ReactNode>((item, i) => (
  //           <span key={`${item}-${i}`} className="font-medium">
  //             {item}
  //           </span>
  //         ))
  //         .reduce((prev, curr, idx) => [
  //           prev,
  //           <span key={idx} className="text-foreground-lighter">
  //             ,
  //           </span>,
  //           curr,
  //         ])}
  //     {modifiedSearchPath.length === 0 && (
  //       <span className="text-foreground-lighter">
  //         No search path prefix set
  //       </span>
  //     )}
  //   </div>
  //   <Icon
  //     className="h-5 w-5 cursor-pointer shrink-0"
  //     icon="edit_square"
  //     onClick={() => setShowEditor(true)}
  //     title="Edit search path prefix"
  //   />
  // </div>;

  return (
    <DashboardFilterButton onClick={toggleControls}>
      Search Path
    </DashboardFilterButton>
  );
};

const DashboardFilterControlButton = ({ filterConfig, toggleControls }) => {
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
    <DashboardFilterButton onClick={toggleControls}>
      Filters
    </DashboardFilterButton>
  );
};

const DashboardGroupingControlButton = ({ groupingConfig, toggleControls }) => {
  // <div className="flex items-center space-x-3 shrink-0">
  //   <Icon className="h-5 w-5" icon="workspaces" />
  //   {groupingConfig
  //     .map<ReactNode>((item) => (
  //       <CheckGroupingTitleLabel
  //         key={`${item.type}${!!item.value ? `-${item.value}` : ""}`}
  //         item={item}
  //       />
  //     ))
  //     .reduce((prev, curr, idx) => [
  //       prev,
  //       <Icon key={idx} className="h-4 w-4" icon="arrow-long-right" />,
  //       curr,
  //     ])}
  //   {!showEditor && (
  //     <Icon
  //       className="h-5 w-5 cursor-pointer shrink-0"
  //       icon="edit_square"
  //       onClick={() => setShowEditor(true)}
  //       title="Edit grouping"
  //     />
  //   )}
  // </div>

  return (
    <DashboardFilterButton onClick={toggleControls}>
      Grouping
    </DashboardFilterButton>
  );
};

const DashboardControlsSummary = ({
  searchPathPrefix,
  filterConfig,
  groupingConfig,
  toggleControls,
}: DashboardControlsSummaryProps) => {
  return (
    <>
      <div className="grow flex items-center justify-end space-x-4">
        <DashboardSearchPathPrefixControlButton
          searchPathPrefix={searchPathPrefix}
          toggleControls={toggleControls}
        />
        <DashboardFilterControlButton
          filterConfig={filterConfig}
          toggleControls={toggleControls}
        />
        <DashboardGroupingControlButton
          groupingConfig={groupingConfig}
          toggleControls={toggleControls}
        />
      </div>
    </>
  );
};

export default DashboardControlsSummary;
