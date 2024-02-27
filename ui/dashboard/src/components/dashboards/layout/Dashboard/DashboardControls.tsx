import CheckFilterConfig from "@powerpipe/components/dashboards/check/CheckFilterConfig";
import CheckGroupingConfig from "@powerpipe/components/dashboards/check/CheckGroupingConfig";
import DashboardControlsSummary from "@powerpipe/components/DashboardControlsSummary";
import useCheckGroupingConfig from "@powerpipe/hooks/useCheckGroupingConfig";
import useCheckFilterConfig from "@powerpipe/hooks/useCheckFilterConfig";
import useDashboardSearchPathPrefix from "@powerpipe/hooks/useDashboardSearchPathPrefix";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useState } from "react";

const DashboardControls = () => {
  const { dashboard } = useDashboard();
  const [expandControls, setExpandControls] = useState(false);
  const searchPathPrefix = useDashboardSearchPathPrefix();
  const filterConfig = useCheckFilterConfig();
  const groupingConfig = useCheckGroupingConfig();

  const toggleControls = () => setExpandControls((e) => !e);

  const isBenchmark =
    dashboard?.children && dashboard.children[0].panel_type === "benchmark";

  return (
    <div className="h-full bg-dashboard-panel print:hidden">
      <div className="flex items-center justify-between p-4 space-x-4">
        <DashboardControlsSummary
          searchPathPrefix={searchPathPrefix}
          filterConfig={filterConfig}
          groupingConfig={groupingConfig}
          toggleControls={toggleControls}
        />
      </div>
      {expandControls && (
        <div className="border-t border-divide divide-y divide-divide">
          {!!dashboard && (
            <>
              {isBenchmark && (
                <>
                  <div className="p-4 space-y-3">
                    <span className="font-semibold">Filters</span>
                    <CheckFilterConfig onClose={toggleControls} />
                  </div>
                  <div className="p-4 space-y-3">
                    <span className="font-semibold">Grouping</span>
                    <CheckGroupingConfig onClose={toggleControls} />
                  </div>
                </>
              )}
              {/*<div className="p-4 space-y-3">*/}
              {/*  <span className="font-semibold">Search Path</span>*/}
              {/*  <SearchPathConfig onClose={toggleControls} />*/}
              {/*</div>*/}
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default DashboardControls;
