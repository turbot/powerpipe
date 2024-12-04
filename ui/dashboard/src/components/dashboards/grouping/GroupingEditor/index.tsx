import EditorAddItem from "@powerpipe/components/dashboards/grouping/common/EditorAddItem";
import Icon from "@powerpipe/components/Icon";
import Select from "react-select";
import useDeepCompareEffect from "use-deep-compare-effect";
import useSelectInputStyles from "../../inputs/common/useSelectInputStyles";
import { classNames } from "@powerpipe/utils/styles";
import { DisplayGroup, DisplayGroupType } from "../common";
import { filterKeysSorter, filterTypeMap } from "@powerpipe/utils/filterEditor";
import {
  MultiValueLabelWithTags,
  OptionWithTags,
  SingleValueWithTags,
} from "@powerpipe/components/dashboards/inputs/common/Common";
import { Reorder, useDragControls } from "framer-motion";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useDashboardControls } from "@powerpipe/components/dashboards/layout/Dashboard/DashboardControlsProvider";

type GroupingEditorProps = {
  config: DisplayGroup[];
  onApply: (newValue: DisplayGroup[]) => void;
};

type GroupingEditorItemProps = {
  config: DisplayGroup[];
  item: DisplayGroup;
  index: number;
  remove: (index: number) => void;
  update: (index: number, item: DisplayGroup) => void;
};

type GroupingTypeSelectProps = {
  index: number;
  item: DisplayGroup;
  type: DisplayGroupType;
  value?: string;
  update: (index: number, updatedItem: DisplayGroup) => void;
};

const GroupingTypeSelect = ({
  index,
  item,
  type,
  value,
  update,
}: GroupingTypeSelectProps) => {
  const [current, setCurrent] = useState<{
    type: DisplayGroupType;
    value?: string;
  }>({ type, value });

  useDeepCompareEffect(() => {
    update(index, {
      ...item,
      type: current.type,
      value: current.value,
    });
  }, [current, index, item]);

  const { context: filterValues } = useDashboardControls();

  const options = useMemo(
    () =>
      Object.entries(filterValues || {})
        .reduce((acc: any[], [key]): any[] => {
          if (filterValues[key]?.hasOwnProperty("key")) {
            let group: any = {
              type: key,
              dynamic: true,
              label: filterTypeMap[key],
              options: [],
            };
            for (let k in filterValues[key]?.key) {
              group.options.push({
                type: key,
                value: k,
                label: k,
              });
            }
            return acc.concat(group);
          }
          return acc.concat({
            type: key,
            label: filterTypeMap[key],
          });
        }, [])
        .concat({ label: "Result", type: "result", value: "result" })
        .sort(filterKeysSorter),
    [filterValues],
  );

  const styles = useSelectInputStyles();

  return (
    <Select
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
      // @ts-ignore as this element definitely exists
      menuPortalTarget={document.getElementById("portals")}
      onChange={setCurrent}
      options={options}
      inputId={`${type}.input`}
      placeholder="Select a group type…"
      // @ts-ignore
      styles={styles}
      value={options
        .reduce((acc, curr) => {
          if (curr?.options) {
            return acc.concat(...curr.options);
          }
          return acc.concat(curr);
        }, [])
        .find(
          (t) => t.value === type || (t.type === type && t.value === value),
        )}
    />
  );
};

const GroupingEditorItem = ({
  config,
  index,
  item,
  remove,
  update,
}: GroupingEditorItemProps) => {
  const dragControls = useDragControls();

  return (
    <Reorder.Item
      as="div"
      id={`${item.type}-${item.value}`}
      className="flex space-x-3 items-center"
      dragControls={dragControls}
      dragListener={false}
      value={item}
    >
      {/*<div className="flex space-x-3 items-center">*/}
      <div className="cursor-grab" onPointerDown={(e) => dragControls.start(e)}>
        <Icon className="h-5 w-5" icon="drag_indicator" />
      </div>
      <div className="grow">
        <GroupingTypeSelect
          index={index}
          item={item}
          type={item.type}
          value={item.value}
          update={update}
        />
      </div>
      <span
        className={classNames(
          config.length > 1
            ? "text-foreground-light hover:text-steampipe-red cursor-pointer"
            : "text-foreground-lightest",
        )}
        onClick={config.length > 1 ? () => remove(index) : undefined}
        title={
          config.length > 1
            ? "Remove"
            : "Grouping must contain at least one level"
        }
      >
        <Icon className="h-5 w-5" icon="trash" />
      </span>
    </Reorder.Item>
  );
};

const GroupingEditor = ({ config, onApply }: GroupingEditorProps) => {
  const [innerConfig, setInnerConfig] = useState<DisplayGroup[]>(config);
  const [isDirty, setIsDirty] = useState(false);
  const [isValid, setIsValid] = useState({ value: false, reason: "" });

  useEffect(() => {
    setInnerConfig(
      config.map((c) => ({
        ...c,
        type: c.type,
        value: c.value,
      })) as any,
    );
  }, [config, setInnerConfig]);

  useEffect(() => {
    let reason: string = "";
    const isValid = innerConfig.every((c, i) => {
      switch (c?.type) {
        case "benchmark":
        case "control":
        case "detection":
        case "reason":
        case "resource":
        case "severity":
        case "status":
          return !c.value;
        case "result":
          if (i !== innerConfig.length - 1) {
            reason = "Result must be the last grouping";
            return false;
          }
          return true;
        default:
          if (c?.type?.includes("|")) {
            return true;
          }
          return true;
      }
    });
    setIsValid({ value: isValid, reason });

    const removeEmpty = innerConfig.map((c) => {
      const noEmpty = {};
      for (const [k, v] of Object.entries(c)) {
        if (!v) {
          continue;
        }
        noEmpty[k] = v;
      }
      return noEmpty;
    });
    setIsDirty(JSON.stringify(config) !== JSON.stringify(removeEmpty));
  }, [config, innerConfig, setIsDirty, setIsValid]);

  const remove = useCallback(
    (index: number) =>
      setInnerConfig((existing) => [
        ...existing.slice(0, index),
        ...existing.slice(index + 1),
      ]),
    [setInnerConfig],
  );

  const update = useCallback(
    (index: number, updatedItem: DisplayGroup) =>
      setInnerConfig((existing) => [
        ...existing.slice(0, index),
        updatedItem,
        ...existing.slice(index + 1),
      ]),
    [setInnerConfig],
  );

  return (
    <div className="flex flex-col space-y-4">
      <Reorder.Group
        axis="y"
        values={innerConfig}
        onReorder={setInnerConfig}
        as="div"
        className="flex flex-col space-y-4"
      >
        {innerConfig.map((c, idx) => (
          <GroupingEditorItem
            key={`${c.type}-${c.value}`}
            config={innerConfig}
            item={c}
            index={idx}
            remove={remove}
            update={update}
          />
        ))}
      </Reorder.Group>
      <EditorAddItem
        isDirty={isDirty}
        isValid={isValid}
        // @ts-ignore
        onAdd={() => setInnerConfig((existing) => [...existing, { type: "" }])}
        onApply={() => onApply(innerConfig)}
        onClear={() => onApply([])}
        addLabel="Add grouping"
      />
    </div>
  );
};

export default GroupingEditor;
