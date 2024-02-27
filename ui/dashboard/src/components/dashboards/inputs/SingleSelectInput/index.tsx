import SelectInput from "@powerpipe/components/dashboards/inputs/SelectInput";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";

const SingleSelectInput = (props: InputProps) => {
  return <SelectInput {...props} />;
};

const definition: IInput = {
  type: "select",
  component: SingleSelectInput,
};

registerInputComponent(definition.type, definition);

export default definition;
