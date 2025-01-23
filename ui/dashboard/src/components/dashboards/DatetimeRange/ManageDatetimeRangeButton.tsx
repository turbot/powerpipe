import DatetimeRangeConfig from "@powerpipe/components/dashboards/DatetimeRange/DatetimeRangeConfig";
import Icon from "@powerpipe/components/Icon";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import { forwardRef } from "react";
import { Popover } from "@headlessui/react";
import { useDashboardDatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";

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
        {!range.relative && <span>{range.to?.format("")}</span>}
        {range.relative && <span>{range.relative}</span>}
      </>
    </NeutralButton>
  );
});

const ManageDatetimeRangeButton = () => {
  return (
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
  );
};

export default ManageDatetimeRangeButton;
