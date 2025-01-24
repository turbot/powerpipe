import DatetimeRangeConfig from "@powerpipe/components/dashboards/DatetimeRange/DatetimeRangeConfig";
import Icon from "@powerpipe/components/Icon";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import { forwardRef, useEffect, useState } from "react";
import { parseDate } from "@powerpipe/utils/date";
import { Popover } from "@headlessui/react";
import { useDashboardDatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

const PopoverButton = forwardRef((props, ref) => {
  const { range } = useDashboardDatetimeRange();

  return (
    // @ts-ignore
    <NeutralButton
      ref={ref}
      className="inline-flex items-center space-x-2 h-full"
      {...props}
    >
      <>
        <Icon
          className="inline-block text-foreground-lighter w-5 h-5"
          icon="calendar_month"
        />
        {(!range.relative || range.relative === "custom") && (
          <span>
            From: {parseDate(range.from)?.format("")}{" "}
            {range.to ? <>To: {parseDate(range.to)?.format("")}</> : ""}
          </span>
        )}
        {range.relative && range.relative !== "custom" && (
          <span>{range.relative}</span>
        )}
      </>
    </NeutralButton>
  );
});

const ManageDatetimeRangeButton = () => {
  const { metadata, dashboard, dashboardsMetadata } = useDashboardState();
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

  return show ? (
    <Popover className="hidden md:block relative">
      <Popover.Button as={PopoverButton} />
      <Popover.Panel className="absolute left-1/2 z-10 mt-4 flex w-screen max-w-max -translate-x-1/2 px-4">
        {({ close }) => (
          <div className="w-screen max-w-md flex-auto overflow-hidden rounded-md bg-dashboard border border-divide shadow-lg ring-1 ring-gray-900/5 p-4">
            <DatetimeRangeConfig onClose={close} />
          </div>
        )}
      </Popover.Panel>
    </Popover>
  ) : null;
};

export default ManageDatetimeRangeButton;
