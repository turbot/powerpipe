import DatetimeRangePicker from "@powerpipe/components/dashboards/DatetimeRange/DatetimeRangePicker";
import { DashboardDataModeLive } from "@powerpipe/types";
import { useDashboardDatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

const DatetimeRangeConfig = ({ onClose }) => {
  const { dataMode } = useDashboardState();
  const { range, setRange } = useDashboardDatetimeRange();

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
    onClose();
  };

  return (
    <DatetimeRangePicker
      from={range.from}
      to={range.to}
      relative={range.relative}
      disabled={dataMode !== DashboardDataModeLive}
      onChange={handleChange}
    />
  );
};

export default DatetimeRangeConfig;
