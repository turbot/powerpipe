import DatetimeRangePicker from "@powerpipe/components/dashboards/DatetimeRange/DatetimeRangePicker";
import { DashboardDataModeLive } from "@powerpipe/types";
import { useDashboardDatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect, useState } from "react";
import { parseDate } from "@powerpipe/utils/date";

const ManageDatetimeRangeButton = () => {
  const { range, setRange } = useDashboardDatetimeRange();
  const { metadata, dashboard, dashboardsMetadata, dataMode } =
    useDashboardState();
  const [show, setShow] = useState(false);

  useEffect(() => {
    if (!metadata && !dashboardsMetadata && !dashboard) {
      return;
    }
    if (dashboard && dashboard.name in dashboardsMetadata) {
      setShow(!!dashboardsMetadata[dashboard.name]?.supports_time_range);
    } else {
      setShow(!!metadata?.supports_time_range);
    }
  }, [metadata?.supports_time_range, dashboard, dashboardsMetadata]);

  const handleChange = (
    from: string,
    to?: string | null | undefined,
    relative?: string | null | undefined,
  ) => {
    if (!from) {
      return;
    }

    if (range.from === from && range.to === to && range.relative === relative) {
      return;
    }

    setRange({ from, to, relative });
  };

  return show ? (
    <DatetimeRangePicker
      from={parseDate(range.from)}
      to={range.to ? parseDate(range.to) : null}
      relative={range.relative}
      disabled={dataMode !== DashboardDataModeLive}
      onChange={handleChange}
    />
  ) : null;
};

export default ManageDatetimeRangeButton;
