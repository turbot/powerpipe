import DashboardListIncludeNestedSelect from "@powerpipe/components/DashboardListIncludeNestedSelect";
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
      title="View Options"
      {...props}
    >
      <Icon
        className="inline-block text-foreground-lighter w-5 h-5"
        icon="dashboard"
      />
      <span>View Options</span>
    </NeutralButton>
  );
});

const DashboardListOptionsButton = () => {
  return (
    <Popover className="hidden md:block relative">
      <Popover.Button as={PopoverButton} />
      <Popover.Panel className="absolute left-1/2 z-10 mt-1 flex w-screen max-w-max -translate-x-1/2 px-4">
        {({ close }) => (
          <div className="w-screen md:w-auto flex-auto overflow-hidden rounded-md bg-dashboard border border-divide shadow-lg ring-1 ring-gray-900/5">
            <div className="flex flex-col divide-y divide-divide">
              <DashboardTagGroupSelect onClose={close} />
              <DashboardListIncludeNestedSelect onClose={close} />
            </div>
          </div>
        )}
      </Popover.Panel>
    </Popover>
  );
};

export default DashboardListOptionsButton;
