import CheckFilterConfig from "../../check/CheckFilterConfig";
import CheckGroupingConfig from "../../check/CheckGroupingConfig";
import SearchPathConfig from "components/dashboards/SearchPathConfig";
import { DashboardDataModeCLISnapshot } from "types";
import { useDashboard } from "hooks/useDashboard";

const DashboardControls = () => {
  const {
    components: { SnapshotHeader },
    dataMode,
    dashboard,
  } = useDashboard();

  const isBenchmark =
    dashboard?.children && dashboard.children[0].panel_type === "benchmark";

  return (
    <div className="grid p-4 gap-6 grid-cols-2 bg-dashboard-panel print:hidden">
      <div className="col-span-2 grid grid-cols-2 gap-6">
        <div className="col-span-2 md:col-span-1 space-y-4">
          {!!dashboard && <SearchPathConfig />}
        </div>
      </div>
      {dataMode === DashboardDataModeCLISnapshot && <SnapshotHeader />}
      {isBenchmark && (
        <div className="col-span-2 grid grid-cols-2 gap-6">
          <div className="col-span-2 md:col-span-1 space-y-4">
            <CheckGroupingConfig />
          </div>
          <div className="col-span-2 md:col-span-1 space-y-4">
            <CheckFilterConfig />
          </div>
        </div>
      )}
    </div>
  );
};

export default DashboardControls;
