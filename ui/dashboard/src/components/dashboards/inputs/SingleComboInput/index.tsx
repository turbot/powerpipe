import ComboInput from "@powerpipe/components/dashboards/inputs/ComboInput";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";

const SingleComboInput = (props: InputProps) => {
  return <ComboInput {...props} />;
};

const definition: IInput = {
  type: "combo",
  component: SingleComboInput,
};

registerInputComponent(definition.type, definition);

export default definition;
