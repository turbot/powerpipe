import Icon from "@powerpipe/components/Icon";
import { CheckResult } from "@powerpipe/components/dashboards/grouping/common";
import { TableRowItem } from "@powerpipe/components/dashboards/layout/DashboardSidePanel/TableRowSidePanel";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";

const ControlSidePanel = ({ result }: { result: CheckResult }) => {
  const { closeSidePanel } = useDashboardPanelDetail();

  return (
    <div className="flex flex-col h-full md:min-w-[300px] md:max-w-[800px]">
      <div className="flex flex-col">
        <div className="flex items-center justify-between p-4">
          <h3>{result.control.title || result.control.name}</h3>
          <Icon
            className="w-5 h-5 text-foreground cursor-pointer hover:text-foreground-light shrink-0"
            icon="close"
            onClick={closeSidePanel}
            title="Close"
          />
        </div>
      </div>
      <div className="flex-1 h-full max-h-full overflow-y-auto divide-y divide-divide">
        <TableRowItem
          dataType="control_status"
          name="status"
          value={result.status}
        />
        <TableRowItem dataType="text" name="reason" value={result.reason} />
        <TableRowItem dataType="text" name="resource" value={result.resource} />
        {result.severity && (
          <TableRowItem
            dataType="text"
            name="severity"
            value={result.severity}
          />
        )}
        {result.dimensions?.map((d) => (
          <TableRowItem
            key={d.key}
            dataType="text"
            name={d.key}
            value={d.value}
          />
        ))}
        {Object.entries(result.tags || {}).map(([key, value]) => (
          <TableRowItem key={key} dataType="text" name={key} value={value} />
        ))}
      </div>
    </div>
  );
};

export default ControlSidePanel;
