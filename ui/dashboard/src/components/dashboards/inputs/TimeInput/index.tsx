import Date from "@powerpipe/components/dashboards/inputs/Date";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";

const TimeInput = (props: InputProps) => {
  return (
    <Date name={props.name} display_type="time" properties={props.properties} />
  );
};

const definition: IInput = {
  type: "time",
  component: TimeInput,
};

registerInputComponent(definition.type, definition);

export default definition;
