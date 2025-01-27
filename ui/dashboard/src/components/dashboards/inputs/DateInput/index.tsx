import D from "@powerpipe/components/dashboards/inputs/D";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";

const DateInput = (props: InputProps) => {
  return (
    <D name={props.name} display_type="date" properties={props.properties} />
  );
};

const definition: IInput = {
  type: "date",
  component: DateInput,
};

registerInputComponent(definition.type, definition);

export default definition;
