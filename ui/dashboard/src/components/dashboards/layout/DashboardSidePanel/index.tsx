import ControlSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/ControlSidePanel";
import FilterAndGroupSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/FilterAndGroupSidePanel";
import TableFilterSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/TableFilterSidePanel";
import TableRowSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/TableRowSidePanel";
import TableSettingsSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/TableSettingsSidePanel";
import { classNames } from "@powerpipe/utils/styles";
import { SidePanelInfo } from "@powerpipe/hooks/useDashboardPanelDetail";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

const DashboardSidePanel = ({
  sidePanel,
}: {
  sidePanel: SidePanelInfo | null;
}) => {
  const {
    breakpointContext: { minBreakpoint },
  } = useDashboardState();
  const isDesktop = minBreakpoint("md");

  if (!sidePanel) {
    return null;
  }

  return (
    <div
      className={classNames(
        !isDesktop ? "w-full absolute" : null,
        "h-full bg-dashboard-panel overflow-y-hidden print:hidden pb-4",
      )}
    >
      {(sidePanel.panel.panel_type === "table" ||
        sidePanel.panel.display_type === "table") && (
        <>
          {sidePanel.context.mode === "filter" && (
            <TableFilterSidePanel panelName={sidePanel.panel.name} />
          )}
          {sidePanel.context.mode === "row" && (
            <TableRowSidePanel
              data={sidePanel.panel.data}
              requestedColumnName={sidePanel.context.requestedColumnName}
              rowIndex={sidePanel.context.rowIndex}
            />
          )}
          {sidePanel.context.mode === "settings" && (
            <TableSettingsSidePanel
              panelName={sidePanel.panel.name}
              leafColumns={sidePanel.context.leafColumns}
            />
          )}
        </>
      )}
      {(sidePanel.panel.panel_type === "benchmark" ||
        sidePanel.panel.panel_type === "detection") && (
        <FilterAndGroupSidePanel panelName={sidePanel.panel.name} />
      )}
      {sidePanel.panel.panel_type === "control" && (
        <ControlSidePanel result={sidePanel.context.result} />
      )}
    </div>
  );
};

export default DashboardSidePanel;
