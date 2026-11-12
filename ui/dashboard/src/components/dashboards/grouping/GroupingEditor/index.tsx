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
  groups.map((c) => ({ ...c, __rowId: `row-${nextRowId++}` })) as DisplayGroup[];

const stripRowIds = (groups: DisplayGroup[]): DisplayGroup[] =>
  groups.map((c) => {
    const { __rowId, ...rest } = c as any;
    return rest;
  }) as DisplayGroup[];

const GroupingEditor = ({
  config,
  defaultConfig,
  onApply,
}: GroupingEditorProps) => {
  const [innerConfig, setInnerConfig] = useState<DisplayGroup[]>(config);
  const [isDirty, setIsDirty] = useState(false);
  const [isValid, setIsValid] = useState({ value: false, reason: "" });

  useEffect(() => {
    setInnerConfig(withRowIds(config) as any);
  }, [config, setInnerConfig]);

  useEffect(() => {
    let reason: string = "";

    // Every level must have a type. The switch below falls through to
    // `default: return true`, so an untyped row counted as valid and Apply
    // saved {"type":""} into the URL - a level that groups nothing.
    const untyped = innerConfig.some((c) => !c?.type);

    // ...and each level may only appear once. A repeat cannot subdivide
    // anything its twin has not already split, so it is silently inert.
    // control_tag/dimension are keyed with their value, so `domain` and
    // `label` remain distinct levels.
    const keys = innerConfig
      .filter((c) => !!c?.type)
      .map((c) => `${c.type}:${c.value ?? ""}`);
    const duplicated = keys.length !== new Set(keys).size;

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
    // Reported after the positional checks so the most specific message wins.
    if (untyped) {
      setIsValid({ value: false, reason: "Choose a type for every grouping" });
    } else if (duplicated) {
      setIsValid({ value: false, reason: "Each grouping can only be used once" });
    } else {
      setIsValid({ value: isValid, reason });
    }

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
            // Key by POSITION, not by content. A newly added row is {type: ""},
            // so a content key is `-undefined` - add two rows and both children
            // share one key. React's child map for the list is then ambiguous
            // and later updates cannot reliably unmount the right rows: after
            // Reset a discarded row survived on screen (and a second Reset could
            // not clear it) even though innerConfig, the saved grouping and the
            // rendered tree were all correct. Only a remount fixed it.
            //
            // Reorder identity is carried by Reorder.Item's `value={item}` prop,
            // not by the React key, so dragging is unaffected.
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
          // Insert BEFORE a trailing "result" rather than appending blindly.
          // "result" is the leaf level - there is nothing beneath a result to
          // subdivide - so the validator below requires it to be last. Appending
          // after it produced the one arrangement that is always invalid, which
          // greyed out Apply with the explanation hidden in a title tooltip.
          setInnerConfig((existing) => {
            const resultIndex = existing.findIndex((c) => c.type === "result");
            return resultIndex === -1
              ? [...existing, { type: "", __rowId: `row-${nextRowId++}` } as any]
              : [
                  ...existing.slice(0, resultIndex),
                  { type: "", __rowId: `row-${nextRowId++}` } as any,
                  ...existing.slice(resultIndex),
                ];
          })
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
