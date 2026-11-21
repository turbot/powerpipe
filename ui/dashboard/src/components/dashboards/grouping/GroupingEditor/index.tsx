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
  defaultConfig: DisplayGroup[];
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

// Rows need a stable identity of their own.
//
// Keying by content (`${type}-${value}`) collides: a newly added row is
// {type: ""}, so adding two produces two children with the key "-undefined".
// React's child map for the list is then ambiguous and later updates strand
// rows on screen - a Reset-discarded row survived, and a second Reset could not
// clear it, even though state, saved grouping and the rendered tree were right.
//
// Keying by index is worse: the react-select at a given position keeps its
// displayed value, so inserting rows made the new ones render as "Result".
//
// So each row carries an id, generated when the row is created and preserved
// across edits and reordering. It is stripped before saving - onApply writes
// straight to the URL, and this is presentation state, not part of the grouping.
let nextRowId = 0;

const withRowIds = (groups: DisplayGroup[]): DisplayGroup[] =>
  groups.map((c) => ({
    ...c,
    __rowId: `row-${nextRowId++}`,
  })) as DisplayGroup[];

const stripRowIds = (groups: DisplayGroup[]): DisplayGroup[] =>
  groups.map((c) => {
    const { __rowId, ...rest } = c as any;
    return rest;
  }) as DisplayGroup[];

// Report an untyped or repeated level before the order rule: an untyped level
// groups nothing and a repeated one splits nothing its twin has not, so either
// would be saved as a level that does nothing. control_tag and dimension levels
// count with their value, so two different tags are two levels.
export const validateGrouping = (
  config: DisplayGroup[],
): { value: boolean; reason: string } => {
  if (config.some((c) => !c?.type)) {
    return { value: false, reason: "Choose a type for every grouping" };
  }
  const keys = config.map((c) => `${c.type}:${c.value ?? ""}`);
  if (keys.length !== new Set(keys).size) {
    return { value: false, reason: "Each grouping can only be used once" };
  }

  let reason = "";
  const value = config.every((c, i) => {
    switch (c.type) {
      case "benchmark":
      case "control":
      case "detection":
      case "reason":
      case "resource":
      case "severity":
      case "status":
        return !c.value;
      case "result":
        if (i !== config.length - 1) {
          reason = "Result must be the last grouping";
          return false;
        }
        return true;
      default:
        return true;
    }
  });
  return { value, reason };
};

// A new level goes before a Result level, which must stay last.
export const addLevel = (
  config: DisplayGroup[],
  level: DisplayGroup,
): DisplayGroup[] => {
  const resultIndex = config.findIndex((c) => c.type === "result");
  return resultIndex === -1
    ? [...config, level]
    : [...config.slice(0, resultIndex), level, ...config.slice(resultIndex)];
};

const GroupingEditor = ({
  config,
  defaultConfig,
  onApply,
}: GroupingEditorProps) => {
  const [innerConfig, setInnerConfig] = useState<DisplayGroup[]>(() =>
    withRowIds(config),
  );
  const [isDirty, setIsDirty] = useState(false);
  const [isValid, setIsValid] = useState({ value: false, reason: "" });

  useEffect(() => {
    setInnerConfig(withRowIds(config) as any);
  }, [config, setInnerConfig]);

  useEffect(() => {
    setIsValid(validateGrouping(innerConfig));

    const removeEmpty = stripRowIds(innerConfig).map((c) => {
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
        // keep the row's id: the child rebuilds the item and would drop it
        { ...updatedItem, __rowId: (existing[index] as any)?.__rowId } as any,
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
            // The row's own id, not its content or position: see withRowIds.
            key={(c as any).__rowId ?? idx}
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
        onAdd={() =>
          setInnerConfig((existing) =>
            addLevel(existing, {
              type: "",
              __rowId: `row-${nextRowId++}`,
            } as DisplayGroup),
          )
        }
        onApply={() => onApply(stripRowIds(innerConfig))}
        onClear={() => {
          // Show the DEFAULT immediately, then clear the saved grouping.
          //
          // Two traps here, both hit in earlier attempts:
          //
          //  - onApply([]) alone does nothing visible. It removes this panel's
          //    saved entry so useGroupingConfig falls back to the default, but
          //    when the panel is already on the default the search params do
          //    not change, so nothing re-renders and the useEffect that syncs
          //    innerConfig never fires.
          //  - Setting innerConfig to `config` is worse than useless: `config`
          //    is the grouping being discarded, so the rows being removed were
          //    briefly re-asserted and could survive the reset.
          //
          // defaultConfig is what onApply([]) will resolve to, computed by
          // useGroupingConfig - which knows the control vs detection default,
          // so this stays correct for detection benchmarks too.
          setInnerConfig(withRowIds(defaultConfig) as any);
          onApply([]);
        }}
        addLabel="Add grouping"
      />
    </div>
  );
};

export default GroupingEditor;
