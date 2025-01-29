import CreatableSelect from "react-select/creatable";
import useSelectInputStyles from "@powerpipe/components/dashboards/inputs/common/useSelectInputStyles";
import useSelectInputValues from "@powerpipe/components/dashboards/inputs/common/useSelectInputValues";
import { DashboardDataModeLive } from "@powerpipe/types";
import {
  InputProps,
  SelectOption,
} from "@powerpipe/components/dashboards/inputs/types";
import {
  MultiValueLabelWithTags,
  OptionWithTags,
  SingleValueWithTags,
} from "@powerpipe/components/dashboards/inputs/common/Common";
import { useDashboardInputs } from "@powerpipe/hooks/useDashboardInputs";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect, useState } from "react";

type SelectInputProps = InputProps & {
  multi?: boolean;
  name: string;
};

const getValueForState = (multi, option) => {
  if (multi) {
    // @ts-ignore
    return option.map((v) => v.value).join(",");
  } else {
    return option.value;
  }
};

const findOptionsForUrlValue = (
  options,
  multi,
  urlValue,
): SelectOption | SelectOption[] => {
  // If we can't find any of the options in the data, we accept it, as this is a
  // combo box and the user can enter anything they like.
  if (multi) {
    const matchingOptions: SelectOption[] = [];
    for (const urlValuePart of urlValue) {
      const existingOption = options.find(
        (option) => option.value === urlValuePart,
      );
      if (existingOption) {
        matchingOptions.push(existingOption);
      } else {
        matchingOptions.push({
          label: urlValuePart,
          value: urlValuePart,
        } as SelectOption);
      }
    }
    return matchingOptions;
  } else {
    const existingOption = options.find((option) => option.value === urlValue);
    if (existingOption) {
      return existingOption;
    } else {
      return {
        label: urlValue,
        value: urlValue,
      } as SelectOption;
    }
  }
};

const ComboInput = ({
  data,
  multi,
  name,
  properties,
  status,
}: SelectInputProps) => {
  const { dataMode } = useDashboardState();
  const { inputs, updateInput, deleteInput } = useDashboardInputs();
  const [initialisedFromState, setInitialisedFromState] = useState(false);
  const [value, setValue] = useState<SelectOption | SelectOption[] | null>(
    null,
  );

  // Get the options for the select
  const options = useSelectInputValues(properties.options, data, status);

  const stateValue = inputs[name];

  // Bind the selected option to the reducer state
  useEffect(() => {
    // If we haven't got the data we need yet...
    if (
      // This property is only present in workspaces >=v0.16.x
      (status !== undefined && status !== "complete") ||
      !options ||
      options.length === 0
    ) {
      return;
    }

    // If this is first load, and we have a value from state, initialise it
    if (!initialisedFromState && stateValue) {
      const parsedUrlValue = multi ? stateValue.split(",") : stateValue;
      const foundOptions = findOptionsForUrlValue(
        options,
        multi,
        parsedUrlValue,
      );
      setValue(foundOptions);
      setInitialisedFromState(true);
    } else if (!initialisedFromState && !stateValue && properties.placeholder) {
      setInitialisedFromState(true);
    } else if (
      !initialisedFromState &&
      !stateValue &&
      !properties.placeholder
    ) {
      setInitialisedFromState(true);
      const newValue = multi ? [options[0]] : options[0];
      setValue(newValue);
      updateInput(name, getValueForState(multi, newValue), false);
    } else if (initialisedFromState && stateValue) {
      const parsedUrlValue = multi ? stateValue.split(",") : stateValue;
      const foundOptions = findOptionsForUrlValue(
        options,
        multi,
        parsedUrlValue,
      );
      setValue(foundOptions);
    } else if (initialisedFromState && !stateValue) {
      if (properties.placeholder) {
        setValue(null);
      } else {
        const newValue = multi ? [options[0]] : options[0];
        setValue(newValue);
        updateInput(name, getValueForState(multi, newValue), false);
      }
    }
  }, [
    initialisedFromState,
    multi,
    name,
    options,
    properties.placeholder,
    stateValue,
    status,
  ]);

  const updateValue = (newValue) => {
    setValue(newValue);
    if (!newValue || newValue.length === 0) {
      deleteInput(name, true);
    } else {
      updateInput(name, getValueForState(multi, newValue), true);
    }
  };

  const styles = useSelectInputStyles();

  if (!styles) {
    return null;
  }

  return (
    <form>
      {properties && properties.label && (
        <label
          className="block mb-1 text-sm"
          id={`${name}.label`}
          htmlFor={`${name}.input`}
        >
          {properties.label}
        </label>
      )}
      <CreatableSelect
        aria-labelledby={`${name}.input`}
        className="basic-single"
        classNamePrefix="select"
        components={{
          // @ts-ignore
          MultiValueLabel: MultiValueLabelWithTags,
          // @ts-ignore
          Option: OptionWithTags,
          // @ts-ignore
          SingleValue: SingleValueWithTags,
        }}
        createOptionPosition="first"
        formatCreateLabel={(inputValue) => `Use "${inputValue}"`}
        // @ts-ignore as this element definitely exists
        menuPortalTarget={document.getElementById("portals")}
        inputId={`${name}.input`}
        isDisabled={
          (!properties.options && !data) || dataMode !== DashboardDataModeLive
        }
        isLoading={!properties.options && !data}
        isClearable={!!properties.placeholder}
        isRtl={false}
        isSearchable
        isMulti={multi}
        // menuIsOpen
        name={name}
        // @ts-ignore
        onChange={updateValue}
        options={options}
        placeholder={
          properties && properties.placeholder ? properties.placeholder : null
        }
        styles={styles}
        value={value}
      />
    </form>
  );
};

export default ComboInput;
