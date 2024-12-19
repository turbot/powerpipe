import Children from "../Children";
import DashboardProgress from "./DashboardProgress";
import DashboardSidePanel from "../DashboardSidePanel";
import DashboardTitle from "@powerpipe/components/dashboards/titles/DashboardTitle";
import Grid from "../Grid";
import PanelDetail from "../PanelDetail";
import SnapshotRenderComplete from "@powerpipe/components/snapshot/SnapshotRenderComplete";
import usePageTitle from "@powerpipe/hooks/usePageTitle";
import { DashboardControlsProvider } from "./DashboardControlsProvider";
import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeLive,
  DashboardDefinition,
} from "@powerpipe/types";
import { registerComponent } from "@powerpipe/components/dashboards";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";
import { useDashboardSearch } from "@powerpipe/hooks/useDashboardSearch";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

type DashboardProps = {
  definition: DashboardDefinition;
  isRoot?: boolean;
  showPanelControls?: boolean;
  withPadding?: boolean;
};

type DashboardWrapperProps = {
  showPanelControls?: boolean;
};

const Dashboard = ({
  definition,
  isRoot = true,
  showPanelControls = true,
}: DashboardProps) => {
  const {
    components: { SnapshotHeader },
    dataMode,
  } = useDashboardState();
  const { selectedSidePanel } = useDashboardPanelDetail();
  const grid = (
    <Grid name={definition.name} width={isRoot ? 12 : definition.width}>
      {isRoot && !definition.artificial && (
        <DashboardTitle title={definition.title} />
      )}
      <Children
        children={definition.children}
        parentType="dashboard"
        showPanelControls={showPanelControls}
      />
    </Grid>
  );
  return (
    <DashboardControlsProvider>
      {dataMode === DashboardDataModeCLISnapshot && (
        <div className="p-4">
          <SnapshotHeader />
        </div>
      )}
      <div className="flex flex-col-reverse md:flex-row w-full h-full overflow-y-hidden">
        {isRoot ? (
          <div className="flex flex-col flex-1 h-full overflow-y-hidden">
            <DashboardProgress />
            <div className="h-full w-full overflow-y-auto p-4">{grid}</div>
          </div>
        ) : (
          <div className="w-full">{grid}</div>
        )}
        <DashboardSidePanel sidePanel={selectedSidePanel} />
      </div>
    </DashboardControlsProvider>
  );
};

const DashboardWrapper = ({
  showPanelControls = true,
}: DashboardWrapperProps) => {
  const { dashboard, dataMode, selectedDashboard } = useDashboardState();
  const { selectedPanel } = useDashboardPanelDetail();
  const { search } = useDashboardSearch();

  usePageTitle([
    selectedDashboard
      ? selectedDashboard.title || selectedDashboard.full_name
      : null,
    "Dashboards",
  ]);

  if (
    search.value ||
    !dashboard ||
    (!selectedDashboard && dataMode === DashboardDataModeLive)
  ) {
    return null;
  }

  if (selectedPanel) {
    return <PanelDetail definition={selectedPanel} />;
  }

  return (
    <>
      <Dashboard
        definition={dashboard}
        showPanelControls={showPanelControls}
        withPadding={true}
      />
      <SnapshotRenderComplete />
    </>
  );
};

registerComponent("dashboard", Dashboard);

export default DashboardWrapper;

export { Dashboard };
