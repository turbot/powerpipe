import ComboInput from "@powerpipe/components/dashboards/inputs/ComboInput";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";

const MultiComboInput = (props: InputProps) => {
  return <ComboInput {...props} multi />;
};

const definition: IInput = {
  type: "multicombo",
  component: MultiComboInput,
};

registerInputComponent(definition.type, definition);

export default definition;
