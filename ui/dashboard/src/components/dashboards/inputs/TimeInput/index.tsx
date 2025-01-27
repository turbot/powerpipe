import D from "@powerpipe/components/dashboards/inputs/D";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";

const TimeInput = (props: InputProps) => {
  return (
    <D name={props.name} display_type="time" properties={props.properties} />
  );
};

const definition: IInput = {
  type: "time",
  component: TimeInput,
};

registerInputComponent(definition.type, definition);

export default definition;
