import D from "@powerpipe/components/dashboards/inputs/D";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";

// const Foo = () => {
//   return <>Bar</>;
// };

const DatetimeInput = (props: InputProps) => {
  // return <Foo />;
  return (
    <D
      name={props.name}
      display_type="datetime"
      properties={props.properties}
    />
  );
};

const definition: IInput = {
  type: "datetime",
  component: DatetimeInput,
};

registerInputComponent(definition.type, definition);

export default definition;
