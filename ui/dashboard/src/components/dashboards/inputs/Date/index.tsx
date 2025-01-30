import { DashboardDataModeLive } from "@powerpipe/types";
import {
  InputProperties,
  InputType,
} from "@powerpipe/components/dashboards/inputs/types";
import { useDashboardInputs } from "@powerpipe/hooks/useDashboardInputs";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect, useState } from "react";

// const DateComp

const Date = ({
  name,
  display_type,
  properties,
}: {
  name: string;
  display_type: InputType;
  properties: InputProperties;
}) => {
  const { dataMode } = useDashboardState();
  const { inputs, updateInput } = useDashboardInputs();
  const stateValue = inputs[name];
  const [value, setValue] = useState<string>(stateValue || "");

  useEffect(() => {
    if (!value || value === stateValue) {
      return;
    }
    updateInput(name, value, !!stateValue);
  }, [value, stateValue]);

  useEffect(() => {
    setValue(stateValue || "");
  }, [stateValue]);

  return (
    <div>
      {properties?.label && (
        <label htmlFor={name} className="block mb-1">
          {properties.label}
        </label>
      )}
      <input
        type={display_type === "datetime" ? "datetime-local" : display_type}
        name={name}
        id={name}
        className="flex-1 block bg-dashboard-panel rounded-md border border-black-scale-3 overflow-x-auto text-sm md:text-base disabled:bg-black-scale-1 focus:ring-0"
        onChange={(e) => setValue(e.target.value)}
        placeholder={properties?.placeholder}
        readOnly={dataMode !== DashboardDataModeLive}
        value={value}
      />
    </div>
  );
};

export default Date;
