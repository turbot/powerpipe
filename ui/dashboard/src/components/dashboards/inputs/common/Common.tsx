import { ColorGenerator } from "@powerpipe/utils/color";
import { components, OptionProps, SingleValueProps } from "react-select";
import isObject from "lodash/isObject";

const stringColorMap = {};
const colorGenerator = new ColorGenerator(24, 4);

const stringToColor = (str) => {
  if (stringColorMap[str]) {
    return stringColorMap[str];
  }
  const color = colorGenerator.nextColor().hex;
  stringColorMap[str] = color;
  return color;
};

const OptionTag = ({ tagKey, tagValue }) => (
  <span
    className="rounded-md text-xs"
    style={{ color: stringToColor(tagValue) }}
    title={`${tagKey} = ${tagValue}`}
  >
    {tagValue}
  </span>
);

const LabelTagWrapper = ({ label, tags, isDisabled }) => (
  <div className={`space-x-2 truncate ${isDisabled ? "text-gray-300" : ""}`}>
    {/*@ts-ignore*/}
    <span title={label}>{label}</span>
    {/*@ts-ignore*/}
    {Object.entries(tags || {}).map(([tagKey, tagValue]) => {
      if (isObject(tagValue)) {
        return Object.entries(tagValue || {}).map(([t, v]) => (
          <OptionTag key={t} tagKey={tagKey} tagValue={v} />
        ));
      }
      return <OptionTag key={tagKey} tagKey={tagKey} tagValue={tagValue} />;
    })}
  </div>
);

const OptionWithTags = (props: OptionProps) => (
  <components.Option {...props}>
    <LabelTagWrapper
      // @ts-ignore
      label={props.data.label}
      // @ts-ignore
      tags={props.data.tags}
      // @ts-ignore
      isDisabled={props.data.isDisabled}
    />
  </components.Option>
);

const SingleValueWithTags = ({ children, ...props }: SingleValueProps) => {
  return (
    <components.SingleValue {...props}>
      <LabelTagWrapper
        // @ts-ignore
        label={props.data.label}
        // @ts-ignore
        tags={props.data.tags}
        // @ts-ignore
        isDisabled={props.data.isDisabled}
      />
    </components.SingleValue>
  );
};

const MultiValueLabelWithTags = ({ children, ...props }: SingleValueProps) => {
  return (
    <components.MultiValueLabel {...props}>
      <LabelTagWrapper
        // @ts-ignore
        label={props.data.label}
        // @ts-ignore
        tags={props.data.tags}
        // @ts-ignore
        isDisabled={props.data.isDisabled}
      />
    </components.MultiValueLabel>
  );
};

export { MultiValueLabelWithTags, OptionWithTags, SingleValueWithTags };
