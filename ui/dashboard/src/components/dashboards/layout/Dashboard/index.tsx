import Children from "../Children";
import DashboardControls from "./DashboardControls";
import DashboardProgress from "./DashboardProgress";
import DashboardTitle from "@powerpipe/components/dashboards/titles/DashboardTitle";
import Grid from "../Grid";
import PanelDetail from "../PanelDetail";
import SnapshotRenderComplete from "@powerpipe/components/snapshot/SnapshotRenderComplete";
import { DashboardControlsProvider } from "./DashboardControlsProvider";
import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeLive,
  DashboardDefinition,
} from "@powerpipe/types";
import { registerComponent } from "@powerpipe/components/dashboards";
import { useDashboard } from "@powerpipe/hooks/useDashboard";

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
    showCustomizeBenchmarkPanel,
  } = useDashboard();
  const grid = (
    <Grid name={definition.name} width={isRoot ? 12 : definition.width}>
      {isRoot && <DashboardTitle title={definition.title} />}
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
      <div className="flex flex-col md:flex-row w-full md:w-auto h-full">
        <div className="w-full">
          {isRoot ? <DashboardProgress /> : null}
          {isRoot ? (
            <div className="h-full overflow-y-auto p-4">{grid}</div>
          ) : (
            grid
          )}
        </div>
        {showCustomizeBenchmarkPanel && <DashboardControls />}
      </div>
    </DashboardControlsProvider>
  );
};

const DashboardWrapper = ({
  showPanelControls = true,
}: DashboardWrapperProps) => {
  const { dashboard, dataMode, search, selectedDashboard, selectedPanel } =
    useDashboard();

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
