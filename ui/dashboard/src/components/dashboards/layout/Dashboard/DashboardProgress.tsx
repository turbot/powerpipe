import { DashboardDataModeLive } from "@powerpipe/types";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

const DashboardProgress = () => {
  const { dataMode, progress, state } = useDashboardState();

  // We only show a progress indicator in live mode
  if (dataMode !== DashboardDataModeLive) {
    return null;
  }

  return (
    <div className="w-full h-[4px] bg-dashboard print:hidden">
      {state === "running" ? (
        <div
          className="h-full bg-black-scale-3"
          style={{ width: `${progress}%` }}
        />
      ) : null}
    </div>
  );
};

export default DashboardProgress;
