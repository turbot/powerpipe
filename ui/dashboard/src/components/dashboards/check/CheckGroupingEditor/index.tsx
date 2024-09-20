import CheckEditorAddItem from "../common/CheckEditorAddItem";
import Icon from "@powerpipe/components/Icon";
import Select from "react-select";
import useDeepCompareEffect from "use-deep-compare-effect";
import useSelectInputStyles from "../../inputs/common/useSelectInputStyles";
import { CheckDisplayGroup, CheckDisplayGroupType } from "../common";
import { classNames } from "@powerpipe/utils/styles";
import {
  MultiValueLabelWithTags,
  OptionWithTags,
  SingleValueWithTags,
} from "@powerpipe/components/dashboards/inputs/common/Common";
import { Reorder, useDragControls } from "framer-motion";
import { SelectOption } from "@powerpipe/components/dashboards/inputs/types";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useDashboardControls } from "@powerpipe/components/dashboards/layout/Dashboard/DashboardControlsProvider";
import { filterKeysSorter, filterTypeMap } from "@powerpipe/utils/filterEditor";

type CheckGroupingEditorProps = {
  config: CheckDisplayGroup[];
  onCancel: () => void;
  onApply: (newValue: CheckDisplayGroup[]) => void;
};

type CheckGroupingEditorItemProps = {
  config: CheckDisplayGroup[];
  item: CheckDisplayGroup;
  index: number;
  remove: (index: number) => void;
  update: (index: number, item: CheckDisplayGroup) => void;
};

type CheckGroupingTypeSelectProps = {
  config: CheckDisplayGroup[];
  index: number;
  item: CheckDisplayGroup;
  type: CheckDisplayGroupType;
  update: (index: number, updatedItem: CheckDisplayGroup) => void;
};

const CheckGroupingTypeSelect = ({
  config,
  index,
  item,
  type,
  update,
}: CheckGroupingTypeSelectProps) => {
  const [currentType, setCurrentType] = useState<CheckDisplayGroupType>(type);

  useDeepCompareEffect(() => {
    update(index, {
      ...item,
      type: currentType,
    });
  }, [currentType, index, item]);

  const { context: filterValues } = useDashboardControls();

  const allFilters = useMemo(
    () =>
      Object.entries(filterValues || {})
        .reduce((acc: any[], [key]): any[] => {
          if (filterValues[key]?.hasOwnProperty("key")) {
            let group: any = {
              label: filterTypeMap[key],
              options: [],
            };
            for (let k in filterValues[key]?.key) {
              group.options.push({
                value: `${key}|${k}`,
                label: k,
              });
            }
            return acc.concat(group);
          }
          return acc.concat({ value: key, label: filterTypeMap[key] });
        }, [])
        .concat({ label: "Result", value: "result" })
        .sort(filterKeysSorter),
    [filterValues],
  );

  const types = useMemo(() => {
    const existingTypes = config.map((c) => c.type.toString());
    // const allTypes: SelectOption[] = [
    //   { value: "benchmark", label: "Benchmark" },
    //   { value: "control", label: "Control" },
    //   { value: "control_tag", label: "Control Tag" },
    //   { value: "dimension", label: "Dimension" },
    //   { value: "reason", label: "Reason" },
    //   { value: "resource", label: "Resource" },
    //   { value: "result", label: "Result" },
    //   { value: "severity", label: "Severity" },
    //   { value: "status", label: "Status" },
    // ];
    return allFilters.filter(
      (t) =>
        t.value === type ||
        t.value === "dimension" ||
        t.value === "control_tag" ||
        // @ts-ignore
        !existingTypes.includes(t.value),
    );
  }, [config, type]);

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
      onChange={(t) =>
        setCurrentType((t as SelectOption).value as CheckDisplayGroupType)
      }
      options={types}
      inputId={`${type}.input`}
      placeholder="Select a group type…"
      // @ts-ignore
      styles={styles}
      value={types
        .reduce((acc, curr) => {
          if (curr?.options) {
            return acc.concat(...curr.options);
          }
          return acc.concat(curr);
        }, [])
        .find((t) => t.value === type)}
    />
  );
};

const CheckGroupingEditorItem = ({
  config,
  index,
  item,
  remove,
  update,
}: CheckGroupingEditorItemProps) => {
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
        <CheckGroupingTypeSelect
          config={config}
          index={index}
          item={item}
          type={item.type}
          update={update}
        />
      </div>
      {/* {(item.type === "dimension" || item.type === "control_tag") && (
        <>
          <span>=</span>
          <div className="grow">
            <CheckGroupingValueSelect
              index={index}
              item={item}
              type={item.type}
              update={update}
              value={item.value}
            />
          </div>
        </>
      )} */}
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

const CheckGroupingEditor = ({ config, onApply }: CheckGroupingEditorProps) => {
  const [innerConfig, setInnerConfig] = useState<CheckDisplayGroup[]>(config);
  const [isDirty, setIsDirty] = useState(false);
  const [isValid, setIsValid] = useState({ value: false, reason: "" });

  useEffect(() => {
    setInnerConfig(
      config.map((c) => ({
        ...c,
        type: c?.value ? `${c.type}|${c.value}` : c?.type,
      })) as any,
    );
  }, [config, setInnerConfig]);

  useEffect(() => {
    let reason: string = "";
    const isValid = innerConfig.every((c, i) => {
      switch (c?.type) {
        case "benchmark":
        case "control":
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
    (index: number, updatedItem: CheckDisplayGroup) =>
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
          <CheckGroupingEditorItem
            key={`${c.type}-${c.value}`}
            config={innerConfig}
            item={c}
            index={idx}
            remove={remove}
            update={update}
          />
        ))}
      </Reorder.Group>
      <CheckEditorAddItem
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

export default CheckGroupingEditor;
