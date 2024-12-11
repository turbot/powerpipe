import { ClearIcon, SubmitIcon } from "@powerpipe/constants/icons";
import { DashboardDataModeLive } from "@powerpipe/types";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { useDashboardInputs } from "@powerpipe/hooks/useDashboardInputs";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect, useState } from "react";

const TextInput = (props: InputProps) => {
  const { dataMode } = useDashboardState();
  const { inputs, updateInput, deleteInput } = useDashboardInputs();
  const stateValue = inputs[props.name];
  const [value, setValue] = useState<string>(() => {
    return stateValue || "";
  });
  const [isDirty, setIsDirty] = useState<boolean>(false);

  const updateValue = (e) => {
    setValue(e.target.value);
    setIsDirty(true);
  };

  const submit = () => {
    setIsDirty(false);
    if (value) {
      updateInput(props.name, value, !!stateValue);
    } else {
      deleteInput(props.name, !!stateValue);
    }
  };

  const clear = () => {
    setValue("");
    setIsDirty(false);
    deleteInput(props.name, true);
  };

  useEffect(() => {
    setValue(stateValue || "");
    setIsDirty(false);
  }, [stateValue]);

  const readOnly = dataMode !== DashboardDataModeLive;

  return (
    <div>
      {props.properties.label && (
        <label htmlFor={props.name} className="block mb-1">
          {props.properties.label}
        </label>
      )}
      <div className="relative">
        <input
          type="text"
          name={props.name}
          id={props.name}
          className="flex-1 block w-full bg-dashboard-panel rounded-md border border-black-scale-3 pr-8 overflow-x-auto text-sm md:text-base disabled:bg-black-scale-1 focus:ring-0"
          onChange={updateValue}
          onKeyPress={(e) => {
            if (e.key !== "Enter") {
              return;
            }
            submit();
          }}
          placeholder={props.properties.placeholder}
          readOnly={readOnly}
          value={value}
        />
        {value && isDirty && !readOnly && (
          <div
            className="absolute inset-y-0 right-0 pr-3 flex items-center cursor-pointer text-foreground-light"
            onClick={submit}
            title="Submit"
          >
            <SubmitIcon className="h-4 w-4" />
          </div>
        )}
        {value && !isDirty && !readOnly && (
          <div
            className="absolute inset-y-0 right-0 pr-3 flex items-center cursor-pointer text-foreground-light"
            onClick={clear}
            title="Clear"
          >
            <ClearIcon className="h-4 w-4" />
          </div>
        )}
      </div>
    </div>
  );
};

const definition: IInput = {
  type: "text",
  component: TextInput,
};

registerInputComponent(definition.type, definition);

export default definition;
