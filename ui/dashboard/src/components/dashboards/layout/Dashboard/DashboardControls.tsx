import CheckFilterConfig from "@powerpipe/components/dashboards/check/CheckFilterConfig";
import CheckGroupingConfig from "@powerpipe/components/dashboards/check/CheckGroupingConfig";
import DashboardControlsSummary from "@powerpipe/components/DashboardControlsSummary";
import SearchPathConfig from "@powerpipe/components/dashboards/SearchPathConfig";
import useCheckGroupingConfig from "@powerpipe/hooks/useCheckGroupingConfig";
import useCheckFilterConfig from "@powerpipe/hooks/useCheckFilterConfig";
import useDashboardSearchPathPrefix from "@powerpipe/hooks/useDashboardSearchPathPrefix";
import { DashboardDataModeCLISnapshot } from "@powerpipe/types";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useState } from "react";

const DashboardControls = () => {
  const {
    components: { SnapshotHeader },
    dataMode,
    dashboard,
  } = useDashboard();
  const [expandControls, setExpandControls] = useState(false);
  const searchPathPrefix = useDashboardSearchPathPrefix();
  const filterConfig = useCheckFilterConfig();
  const groupingConfig = useCheckGroupingConfig();

  const toggleControls = () => setExpandControls((e) => !e);

  const isBenchmark =
    dashboard?.children && dashboard.children[0].panel_type === "benchmark";

  return (
    <>
      <div className="flex items-center justify-between w-full p-4 space-x-4 bg-dashboard-panel">
        {dataMode === DashboardDataModeCLISnapshot && <SnapshotHeader />}
        <DashboardControlsSummary
          searchPathPrefix={searchPathPrefix}
          filterConfig={filterConfig}
          groupingConfig={groupingConfig}
          toggleControls={toggleControls}
        />
      </div>
      {expandControls && (
        <div className="grid p-4 gap-6 grid-cols-3 bg-dashboard-panel print:hidden">
          <div className="col-span-3 md:col-span-1">
            {!!dashboard && <SearchPathConfig onClose={toggleControls} />}
          </div>
          {isBenchmark && (
            <>
              <div className="col-span-3 md:col-span-1">
                <CheckFilterConfig onClose={toggleControls} />
              </div>
              <div className="col-span-3 md:col-span-1">
                <CheckGroupingConfig onClose={toggleControls} />
              </div>
            </>
          )}
        </div>
      )}
    </>
  );
};

export default DashboardControls;
