import FilterAndGroupSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/FilterAndGroupSidePanel";
import TableRowSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/TableRowSidePanel";
import TableSettingsSidePanel from "@powerpipe/components/dashboards/layout/DashboardSidePanel/TableSettingsSidePanel";
import { classNames } from "@powerpipe/utils/styles";
import { SidePanelInfo } from "@powerpipe/hooks/useDashboardPanelDetail";
import { useBreakpoint } from "@powerpipe/hooks/useBreakpoint";

const DashboardSidePanel = ({
  sidePanel,
}: {
  sidePanel: SidePanelInfo | null;
}) => {
  const { minBreakpoint } = useBreakpoint();
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
      {(sidePanel.panel.panel_type === "benchmark" ||
        sidePanel.panel.panel_type === "control" ||
        sidePanel.panel.panel_type === "detection") && (
        <FilterAndGroupSidePanel panelName={sidePanel.panel.name} />
      )}
      {sidePanel.panel.panel_type === "table" &&
        sidePanel.context.mode === "row" && (
          <TableRowSidePanel
            data={sidePanel.panel.data}
            requestedColumnName={sidePanel.context.requestedColumnName}
            rowIndex={sidePanel.context.rowIndex}
          />
        )}
      {sidePanel.panel.panel_type === "table" &&
        sidePanel.context.mode === "settings" && (
          <TableSettingsSidePanel
            panelName={sidePanel.panel.name}
            leafColumns={sidePanel.context.leafColumns}
          />
        )}
    </div>
  );
};

export default DashboardSidePanel;
