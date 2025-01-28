import DashboardListDisplayModeSelect from "@powerpipe/components/DashboardListDisplayModeSelect";
import DashboardTagGroupSelect from "@powerpipe/components/DashboardTagGroupSelect";
import Icon from "@powerpipe/components/Icon";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import { forwardRef } from "react";
import { Popover } from "@headlessui/react";

const PopoverButton = forwardRef((props, ref) => {
  return (
    // @ts-ignore
    <NeutralButton
      ref={ref}
      className="inline-flex items-center space-x-2 h-full"
      //disabled={!enabled}
      title="Manage display options"
      {...props}
    >
      <Icon
        className="inline-block text-foreground-lighter w-5 h-5"
        icon="dashboard"
      />
      <span>Display</span>
    </NeutralButton>
  );
});

const Label = ({ children }) => (
  <span className="block text-xs uppercase text-foreground-lighter">
    {children}
  </span>
);

const DashboardListOptionsButton = () => {
  return (
    <Popover className="relative">
      <Popover.Button as={PopoverButton} />
      <Popover.Panel className="absolute left-1/2 z-10 mt-1 flex w-screen max-w-max -translate-x-1/2 px-4">
        {({ close }) => (
          <div className="max-w-sm flex-auto overflow-hidden rounded-md bg-dashboard border border-divide shadow-lg ring-1 ring-gray-900/5">
            <div className="divide-y divide-divide">
              <div className="p-3 space-y-2">
                <Label>Group By</Label>
                <DashboardTagGroupSelect onClose={close} />
              </div>
              <div className="p-3 space-y-2">
                <Label>Dashboards</Label>
                <DashboardListDisplayModeSelect onClose={close} />
              </div>
            </div>
          </div>
        )}
      </Popover.Panel>
    </Popover>
  );
};

export default DashboardListOptionsButton;
