import Icon from "@powerpipe/components/Icon";
import TableSettings from "@powerpipe/components/dashboards/Table/TableSettings";
import { Column, RowData } from "@tanstack/react-table";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";

const TableSettingsSidePanel = ({
  panelName,
  leafColumns,
}: {
  panelName: string;
  leafColumns: Column<RowData>[];
}) => {
  const { closeSidePanel } = useDashboardPanelDetail();
  return (
    <div className="flex flex-col h-full md:min-w-[300px] md:max-w-[800px]">
      <div className="flex flex-col">
        <div className="flex items-center justify-between p-4">
          <h3>Visible Columns</h3>
          <Icon
            className="w-5 h-5 text-foreground cursor-pointer hover:text-foreground-light shrink-0"
            icon="close"
            onClick={closeSidePanel}
            title="Close"
          />
        </div>
      </div>
      <TableSettings name={panelName} leafColumns={leafColumns} />
    </div>
  );
};

export default TableSettingsSidePanel;
