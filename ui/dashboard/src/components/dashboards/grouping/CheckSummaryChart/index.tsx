import IntegerDisplay from "../../../IntegerDisplay";
import usePrevious from "@powerpipe/hooks/usePrevious";
import { CheckNodeStatus, CheckSummary } from "../common";
import { classNames } from "@powerpipe/utils/styles";
import { ReactNode, useEffect, useRef, useState } from "react";

type ProgressBarGroupProps = {
  children: ReactNode;
  className?: string;
};

type ProgressBarProps = {
  className?: string;
  width: number;
};

type CheckSummaryChartProps = {
  status: CheckNodeStatus;
  summary: CheckSummary;
  firstChildSummaries: CheckSummary[];
};

type AlertProgressBarGroupTotalProps = {
  className?: string;
  summary: CheckSummary;
};

type NonAlertProgressBarGroupTotalProps = {
  className?: string;
  summary: CheckSummary;
};

type ProgressBarGroupTotalProps = {
  className?: string;
  total: number;
};

const getWidth = (x: number, y: number) => {
  const percent = (x / (x + y)) * 100;
  return percent >= 0.5 ? Math.round(percent) : 1;
};

export const ProgressBarGroupTotal = ({
  className,
  total,
}: ProgressBarGroupTotalProps) => (
  <span className={classNames(className, "text-right text-sm font-semibold")}>
    {total > 0 ? <IntegerDisplay num={total} withTitle={false} /> : "0"}
  </span>
);

const getSummaryTitle = (summary: CheckSummary): string => {
  const titleParts: string[] = [];
  if (summary.error) {
    titleParts.push(`Error: ${summary.error.toLocaleString()}`);
  }
  if (summary.alarm) {
    titleParts.push(`Alarm: ${summary.alarm.toLocaleString()}`);
  }
  if (summary.ok) {
    titleParts.push(`OK: ${summary.ok.toLocaleString()}`);
  }
  if (summary.info) {
    titleParts.push(`Info: ${summary.info.toLocaleString()}`);
  }
  if (summary.skip) {
    titleParts.push(`Skipped: ${summary.skip.toLocaleString()}`);
  }
  if (titleParts.length === 0) {
    return "";
  }
  return titleParts.join(`
`);
};

const AlertProgressBarGroupTotal = ({
  className,
  summary,
}: AlertProgressBarGroupTotalProps) => {
  const alertTotal = summary.error + summary.alarm;
  const newClassName = classNames(
    className,
    alertTotal > 0 ? "text-alert" : "text-foreground-lightest",
  );
  return <ProgressBarGroupTotal className={newClassName} total={alertTotal} />;
};

const NonAlertProgressBarGroupTotal = ({
  className,
  summary,
}: NonAlertProgressBarGroupTotalProps) => {
  const nonAlertTotal = summary.ok + summary.info + summary.skip;
  let textClassName: string;
  if (nonAlertTotal === 0) {
    textClassName = "text-foreground-lightest";
  } else if (summary.skip > summary.info && summary.skip > summary.ok) {
    textClassName = "text-black-scale-5";
  } else if (summary.info > summary.ok && summary.info >= summary.skip) {
    textClassName = "text-info";
  } else {
    textClassName = "text-ok";
  }

  const newClassName = classNames(className, textClassName);
  return (
    <ProgressBarGroupTotal className={newClassName} total={nonAlertTotal} />
  );
};

export const ProgressBarGroup = ({
  children,
  className,
}: ProgressBarGroupProps) => (
  <div className={classNames("flex h-3 items-center", className)}>
    {children}
  </div>
);

export const ProgressBar = ({ className, width }: ProgressBarProps) => {
  if (!width) {
    return null;
  }

  return (
    <div
      className={classNames("flex-shrink-0 h-3", className)}
      style={{ display: "inline-block", width: `${width}px` }}
    />
  );
};

const getSegmentPixelWidth = (
  value: number,
  divisor: number,
  containerWidth: number,
) => {
  if (!value || !divisor || !containerWidth) {
    return 0;
  }

  const percent = value / divisor;
  return Math.max(Math.round(percent * containerWidth), 1);
};

const CheckSummaryChart = ({
  status,
  summary,
  firstChildSummaries,
}: CheckSummaryChartProps) => {
  const [, setVersion] = useState(0);
  const alertsContainerRef = useRef<HTMLDivElement>(null);
  const nonAlertsContainerRef = useRef<HTMLDivElement>(null);
  const previousContainers = usePrevious<{
    alerts: HTMLDivElement | null;
    nonAlerts: HTMLDivElement | null;
  }>({
    alerts: alertsContainerRef.current,
    nonAlerts: nonAlertsContainerRef.current,
  });

  useEffect(() => {
    if (
      (!previousContainers &&
        (alertsContainerRef.current || nonAlertsContainerRef.current)) ||
      (previousContainers &&
        (previousContainers.alerts !== alertsContainerRef.current ||
          previousContainers.nonAlerts !== nonAlertsContainerRef.current))
    ) {
      // Trigger a re-render when the container refs change
      // This is necessary to recalculate widths based on new container sizes
      setVersion((v) => v + 1);
    }
  }, [previousContainers]);

  let maxAlerts = 0;
  let maxNonAlerts = 0;
  for (const firstChildSummary of firstChildSummaries) {
    const currentMaxAlerts = firstChildSummary.error + firstChildSummary.alarm;
    const currentMaxNonAlerts =
      firstChildSummary.ok + firstChildSummary.info + firstChildSummary.skip;
    if (currentMaxAlerts > maxAlerts) {
      maxAlerts = currentMaxAlerts;
    }
    if (currentMaxNonAlerts > maxNonAlerts) {
      maxNonAlerts = currentMaxNonAlerts;
    }
  }
  const alertsWidth = getWidth(maxAlerts, maxNonAlerts) * 0.9;
  const nonAlertsWidth = getWidth(maxNonAlerts, maxAlerts) * 0.9;

  const calculateWidths = () => {
    if (!alertsContainerRef.current || !nonAlertsContainerRef.current) {
      return {
        renderDivider: false,
        alarm: 0,
        error: 0,
        ok: 0,
        info: 0,
        skip: 0,
      };
    }

    const alertsContainerWidth = alertsContainerRef.current?.clientWidth ?? 0;
    const nonAlertsContainerWidth =
      nonAlertsContainerRef.current?.clientWidth ?? 0;

    const rawAlarm = getSegmentPixelWidth(
      summary.alarm,
      maxAlerts,
      alertsContainerWidth,
    );
    const rawError = getSegmentPixelWidth(
      summary.error,
      maxAlerts,
      alertsContainerWidth,
    );
    const rawOk = getSegmentPixelWidth(
      summary.ok,
      maxNonAlerts,
      nonAlertsContainerWidth,
    );
    const rawInfo = getSegmentPixelWidth(
      summary.info,
      maxNonAlerts,
      nonAlertsContainerWidth,
    );
    const rawSkip = getSegmentPixelWidth(
      summary.skip,
      maxNonAlerts,
      nonAlertsContainerWidth,
    );

    return {
      alertsContainerWidth,
      nonAlertsContainerWidth,
      renderDivider: true,
      alarm: rawAlarm,
      error: rawError,
      ok: rawOk,
      info: rawInfo,
      skip: rawSkip,
    };
  };

  const widths = calculateWidths();

  return (
    <div className="flex items-center" title={getSummaryTitle(summary)}>
      <div
        ref={alertsContainerRef}
        className="my-auto px-0"
        style={{ width: `${alertsWidth}%` }}
      >
        <ProgressBarGroup className="flex-row-reverse">
          <ProgressBar
            className={classNames(
              status === "running" ? "summary-chart-alarm-animate" : "bg-alert",
            )}
            width={widths.alarm}
          />
          <ProgressBar
            className={classNames(
              "border-alert",

              status === "running"
                ? "border summary-chart-error-animate"
                : widths.error < 2
                  ? "border-l"
                  : "border",
            )}
            width={widths.error}
          />
          <AlertProgressBarGroupTotal className="mr-2" summary={summary} />
        </ProgressBarGroup>
      </div>
      {widths.renderDivider && (
        <div
          className={classNames(
            "h-6 w-px bg-black-scale-4",
            status === "running" ? "subtle-ping" : null,
          )}
        />
      )}
      <div
        ref={nonAlertsContainerRef}
        className="my-auto px-0"
        style={{ width: `${nonAlertsWidth}%` }}
      >
        <ProgressBarGroup>
          <ProgressBar
            className={classNames(
              status === "running" ? "summary-chart-ok-animate" : "bg-ok",
            )}
            width={widths.ok}
          />
          <ProgressBar
            className={classNames(
              status === "running" ? "summary-chart-info-animate" : "bg-info",
            )}
            width={widths.info}
          />
          <ProgressBar
            className={classNames(
              status === "running" ? "summary-chart-skip-animate" : "bg-skip",
            )}
            width={widths.skip}
          />
          <NonAlertProgressBarGroupTotal className="ml-2" summary={summary} />
        </ProgressBarGroup>
      </div>
    </div>
  );
};

export default CheckSummaryChart;
