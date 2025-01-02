import Children from "../Children";
import DashboardProgress from "./DashboardProgress";
import DashboardSidePanel from "../DashboardSidePanel";
import DashboardTitle from "@powerpipe/components/dashboards/titles/DashboardTitle";
import Grid from "../Grid";
import PanelDetail from "../PanelDetail";
import usePageTitle from "@powerpipe/hooks/usePageTitle";
import { classNames } from "@powerpipe/utils/styles";
import { DashboardControlsProvider } from "./DashboardControlsProvider";
import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeLive,
  DashboardDefinition,
} from "@powerpipe/types";
import { Fragment, ReactNode, useEffect, useRef, useState } from "react";
import { registerComponent } from "@powerpipe/components/dashboards";
import { useBreakpoint } from "@powerpipe/hooks/useBreakpoint";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";
import { useDashboardSearch } from "@powerpipe/hooks/useDashboardSearch";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { v4 as uuid } from "uuid";

type DashboardProps = {
  definition: DashboardDefinition;
  isRoot?: boolean;
  showPanelControls?: boolean;
  withPadding?: boolean;
};

type DashboardWrapperProps = {
  showPanelControls?: boolean;
};

type SplitPaneProps = {
  children: ReactNode[]; // Expect exactly 2 children
  split?: "vertical" | "horizontal"; // Orientation
  minRightPanelSize?: number; // Minimum size for the second pane
  maxRightPanelSize?: number; // Maximum size for the second pane
  defaultRightPanelSize?: number; // Initial size for the second pane
  onChange?: (size: number) => void; // Callback for size changes
};

const VerticalSplitPane = ({
  children,
  defaultRightPanelSize,
  minRightPanelSize,
  maxRightPanelSize,
  onChange,
}: SplitPaneProps) => {
  const { minBreakpoint } = useBreakpoint();
  const isDesktop = minBreakpoint("md");
  const [_, setRenderId] = useState(uuid());
  const [size, setSize] = useState(defaultRightPanelSize);
  const isDragging = useRef(false);
  const paneRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    setSize(() => defaultRightPanelSize);
  }, [defaultRightPanelSize]);

  const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging.current || !paneRef.current) return;

    const rect = paneRef.current.getBoundingClientRect();
    let newSize = rect.right - e.clientX;
    newSize = Math.max(minRightPanelSize, Math.min(newSize, maxRightPanelSize));

    setSize(newSize);
    if (onChange) onChange(newSize);
  };

  const handleMouseUp = () => {
    if (isDragging.current) {
      isDragging.current = false; // Stop dragging
      setRenderId(uuid());
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
    }
  };

  const handleMouseDown = () => {
    isDragging.current = true;
    setRenderId(uuid());
    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);
  };

  useEffect(() => {
    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, []);

  return (
    <div
      ref={paneRef}
      className={classNames(
        "flex flex-col-reverse md:flex-row w-full h-full overflow-y-hidden",
        isDragging.current ? "select-none" : "",
      )}
    >
      {children[0]}
      {!isDesktop && children[1]}
      {isDesktop && children[1] && (
        <div
          className={classNames(
            "border cursor-col-resize w-[3px]",
            isDragging.current
              ? "border-black-scale-4 bg-black-scale-4"
              : "border-dashboard bg-divide",
          )}
          onMouseDown={handleMouseDown}
        />
      )}
      {isDesktop && children[1] && (
        <div
          style={{
            width: `${size}px`,
          }}
        >
          {children[1]}
        </div>
      )}
    </div>
  );
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
        childPanels={definition.children}
        parentType="dashboard"
        showPanelControls={showPanelControls}
      />
    </Grid>
  );
  const selectedRightPanelType = selectedSidePanel?.panel?.panel_type;
  return (
    <DashboardControlsProvider>
      {dataMode === DashboardDataModeCLISnapshot && (
        <div className="p-4">
          <SnapshotHeader />
        </div>
      )}
      <VerticalSplitPane
        defaultRightPanelSize={selectedRightPanelType === "table" ? 400 : 600}
        minRightPanelSize={selectedRightPanelType === "table" ? 300 : 400}
        maxRightPanelSize={selectedRightPanelType === "table" ? 800 : 1000}
      >
        <Fragment key={definition.name}>
          {isRoot ? (
            <div className="flex flex-col overflow-y-hidden">
              <DashboardProgress />
              <div className="h-full w-full overflow-y-auto p-4">{grid}</div>
            </div>
          ) : (
            <div className="w-full">{grid}</div>
          )}
        </Fragment>
        {selectedSidePanel && (
          <DashboardSidePanel
            key={selectedSidePanel?.panel?.panel_type}
            sidePanel={selectedSidePanel}
          />
        )}
      </VerticalSplitPane>
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
    <Dashboard
      definition={dashboard}
      showPanelControls={showPanelControls}
      withPadding={true}
    />
  );
};

registerComponent("dashboard", Dashboard);

export default DashboardWrapper;

export { Dashboard };
