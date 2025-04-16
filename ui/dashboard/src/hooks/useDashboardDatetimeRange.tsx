import dayjs from "dayjs";
import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useMemo
} from "react";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useSearchParams } from "react-router-dom";

export interface DatetimeRange {
  from: string;
  to?: string | null;
  relative?: string | null;
}

interface IDashboardDatetimeRangeContext {
  range: DatetimeRange;
  supportsTimeRange: boolean;
  setRange: (range: DatetimeRange) => void;
}

interface DashboardSearchPathProviderProps {
  children: ReactNode;
}

const DashboardDatetimeRangeContext =
  createContext<IDashboardDatetimeRangeContext | null>(null);

export const DashboardDatetimeRangeProvider = ({
                                                 children
                                               }: DashboardSearchPathProviderProps) => {
  const { metadata, dashboard, dashboardsMetadata } = useDashboardState();
  const [searchParams, setSearchParams] = useSearchParams();
  const rawDatetimeRange = searchParams.get("datetime_range");

  const serverSupportsTimeRange = !!metadata?.supports_time_range;
  const dashboardSupportsTimeRange = !!(
    dashboard && dashboardsMetadata[dashboard.name]?.supports_time_range
  );

  const datetimeRange = useMemo<DatetimeRange>(() => {
    if (!!rawDatetimeRange) {
      try {
        return JSON.parse(rawDatetimeRange);
      } catch (error) {
        console.error("Error parsing search path prefix", error);
        return {
          from: dayjs().subtract(7, "day").toISOString(),
          to: null,
          relative: "7d"
        };
      }
    } else {
      return {
        from: dayjs().subtract(7, "day").toISOString(),
        to: null,
        relative: "7d"
      };
    }
  }, [rawDatetimeRange]);

  const updateRange = (range: DatetimeRange) => {
    setSearchParams((previous) => {
      const newParams = new URLSearchParams(previous);
      newParams.set("datetime_range", JSON.stringify(range));
      return newParams;
    });
  };

  const initialiseRange = (range: DatetimeRange) => {
    updateRange(range);
  };

  const setRange = (range: DatetimeRange) => {
    updateRange(range);
  };

  useEffect(() => {
    if (
      (!serverSupportsTimeRange && !dashboardSupportsTimeRange) ||
      rawDatetimeRange
    ) {
      return;
    }
    initialiseRange({
      from: datetimeRange.from,
      to: datetimeRange.to,
      relative: datetimeRange.relative
    });
  }, [
    serverSupportsTimeRange,
    dashboardSupportsTimeRange,
    rawDatetimeRange,
    datetimeRange.from,
    datetimeRange.to,
    datetimeRange.relative
  ]);

  return (
    <DashboardDatetimeRangeContext.Provider
      value={{
        range: datetimeRange,
        supportsTimeRange:
          serverSupportsTimeRange ||
          dashboardSupportsTimeRange ||
          !!rawDatetimeRange,
        setRange
      }}
    >
      {children}
    </DashboardDatetimeRangeContext.Provider>
  );
};

export const useDashboardDatetimeRange = () => {
  const context = useContext(DashboardDatetimeRangeContext);
  if (!context) {
    throw new Error(
      "useDashboardDatetimeRange must be used within a DashboardDatetimeRangeContext"
    );
  }
  return context;
};
