import { DashboardDataModeLive } from "@powerpipe/types";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { useDashboardInputs } from "@powerpipe/hooks/useDashboardInputs";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect, useState } from "react";

const DateInput = (props: InputProps) => {
  const { dataMode } = useDashboardState();
  const { inputs, updateInput } = useDashboardInputs();
  const stateValue = inputs[props.name];
  const [value, setValue] = useState<string>(stateValue || "");

  useEffect(() => {
    if (!value || value === stateValue) {
      return;
    }
    updateInput(props.name, value, !!stateValue);
  }, [value, stateValue]);

  useEffect(() => {
    setValue(stateValue || "");
  }, [stateValue]);

  return (
    <div>
      {props.properties.label && (
        <label htmlFor={props.name} className="block mb-1">
          {props.properties.label}
        </label>
      )}
      <input
        type="date"
        name={props.name}
        id={props.name}
        className="flex-1 block bg-dashboard-panel rounded-md border border-black-scale-3 overflow-x-auto text-sm md:text-base disabled:bg-black-scale-1 focus:ring-0"
        onChange={(e) => setValue(e.target.value)}
        placeholder={props.properties.placeholder}
        readOnly={dataMode !== DashboardDataModeLive}
        value={value}
      />
    </div>
  );
};

const definition: IInput = {
  type: "date",
  component: DateInput,
};

registerInputComponent(definition.type, definition);

export default definition;
