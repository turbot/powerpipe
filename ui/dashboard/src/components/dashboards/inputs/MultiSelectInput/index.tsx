import SelectInput from "@powerpipe/components/dashboards/inputs/SelectInput";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";

const MultiSelectInput = (props: InputProps) => {
  return <SelectInput {...props} multi />;
};

const definition: IInput = {
  type: "multiselect",
  component: MultiSelectInput,
};

registerInputComponent(definition.type, definition);

export default definition;
