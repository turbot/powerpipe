import CheckEditorAddItem from "../common/CheckEditorAddItem";
import CreatableSelect from "react-select/creatable";
import has from "lodash/has";
import Icon from "@powerpipe/components/Icon";
import Select from "react-select";
import useDeepCompareEffect from "use-deep-compare-effect";
import useSelectInputStyles from "../../inputs/common/useSelectInputStyles";
import { CheckFilter, CheckFilterType } from "../common";
import { classNames } from "@powerpipe/utils/styles";
import {
  MultiValueLabelWithTags,
  OptionWithTags,
  SingleValueWithTags,
} from "../../inputs/common/Common";
import { Reorder, useDragControls } from "framer-motion";
import { SelectOption } from "../../inputs/types";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useDashboardControls } from "../../layout/Dashboard/DashboardControlsProvider";

type CheckFilterEditorProps = {
  filter: CheckFilter;
  onApply: (newValue: CheckFilter) => void;
};

type CheckFilterEditorItemProps = {
  filter: CheckFilter;
  item: CheckFilter;
  index: number;
  remove: (index: number) => void;
  update: (index: number, item: CheckFilter) => void;
};

type CheckFilterTypeSelectProps = {
  className?: string;
  filter: CheckFilter;
  index: number;
  item: CheckFilter;
  type: CheckFilterType;
  update: (index: number, updatedItem: CheckFilter) => void;
};

type CheckFilterKeySelectProps = {
  index: number;
  item: CheckFilter;
  type: CheckFilterType;
  update: (index: number, updatedItem: CheckFilter) => void;
  filterKey: string | undefined;
};

type CheckFilterValueSelectProps = {
  className?: string;
  index: number;
  item: CheckFilter;
  type: CheckFilterType;
  update: (index: number, updatedItem: CheckFilter) => void;
  value: string | undefined;
};

const validateFilter = (filter: CheckFilter): boolean => {
  // Each and must have at least one expression
  if (
    filter.operator === "and" &&
    (!filter.expressions || filter.expressions.length < 1)
  ) {
    return false;
  }
  if (filter.operator === "and") {
    // @ts-ignore can't reach here if filter.expressions is not truthy
    return filter.expressions.every(validateFilter);
  }

  if (filter.operator === "equal") {
    return (
      !!filter.type && (filter.key !== undefined || filter.value !== undefined)
    );
  }

  return false;
};

const CheckFilterTypeSelect = ({
  className,
  filter,
  index,
  item,
  type,
  update,
}: CheckFilterTypeSelectProps) => {
  const [currentType, setCurrentType] = useState<CheckFilterType>(type);

  useDeepCompareEffect(() => {
    update(index, {
      ...item,
      type: currentType,
      value: "",
    });
  }, [currentType, index, item]);

  const types = useMemo(() => {
    // @ts-ignore
    const existingTypes = filter.expressions
      ?.map((c) => c.type?.toString())
      .filter((t) => !!t);
    const allTypes: SelectOption[] = [
      { value: "benchmark", label: "Benchmark" },
      { value: "control", label: "Control" },
      { value: "control_tag", label: "Control Tag" },
      { value: "dimension", label: "Dimension" },
      { value: "reason", label: "Reason" },
      { value: "resource", label: "Resource" },
      { value: "severity", label: "Severity" },
      { value: "status", label: "Status" },
    ];
    return allTypes.filter(
      (t) =>
        t.value === type ||
        t.value === "dimension" ||
        t.value === "control_tag" ||
        // @ts-ignore
        !existingTypes.includes(t.value),
    );
  }, [filter, type]);

  const styles = useSelectInputStyles();

  return (
    <Select
      className={classNames("basic-single", className)}
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
      onChange={(t) => setCurrentType((t as SelectOption).value)}
      options={types}
      inputId={`${type}.input`}
      placeholder="Select a filter…"
      styles={styles}
      value={types.find((t) => t.value === type)}
    />
  );
};

const CheckFilterKeySelect = ({
  index,
  item,
  type,
  filterKey,
  update,
}: CheckFilterKeySelectProps) => {
  const [currentKey, setCurrentKey] = useState(filterKey);
  const { context: filterValues } = useDashboardControls();

  useDeepCompareEffect(() => {
    update(index, {
      ...item,
      key: currentKey,
    });
  }, [currentKey, index, item]);

  const keys = useMemo(() => {
    return Object.keys(filterValues ? filterValues[type].key || {} : {}).map(
      (k) => ({
        value: k,
        label: k,
      }),
    );
  }, [filterValues, type]);

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
      onChange={(t) => setCurrentKey((t as SelectOption).value)}
      options={keys}
      inputId={`${type}.input`}
      placeholder={`Choose a ${type}…`}
      styles={styles}
      value={keys.find((t) => t.value === filterKey)}
    />
  );
};

const CheckFilterValueSelect = ({
  className,
  index,
  item,
  type,
  value,
  update,
}: CheckFilterValueSelectProps) => {
  const [currentValue, setCurrentValue] = useState<{
    value: any;
    title?: string;
  }>({ value, title: item.title });
  const { context: filterValues } = useDashboardControls();

  const values = useMemo(() => {
    if (!type) {
      return [];
    }
    if (type === "status") {
      return (
        Object.entries(filterValues ? filterValues[type] || {} : {})
          // @ts-ignore
          .filter(([, v]) => v > 0)
          .map(([k, v]) => ({
            value: k,
            label: k,
            tags: { occurrences: v },
          }))
      );
    } else if (type === "dimension") {
      return Object.entries(filterValues ? filterValues[type].value || {} : {})
        .filter(([, v]) => has(v, item.key as string))
        .map(([k, v]) => {
          return {
            value: k,
            label: k,
            // @ts-ignore
            tags: { occurrences: v[item.key] },
          };
        });
    } else if (type === "benchmark" || type === "control") {
      return Object.entries(
        filterValues ? filterValues[type].value || {} : {},
      ).map(([k, v]) => {
        return {
          value: k,
          // @ts-ignore
          label: v.title || k,
          // @ts-ignore
          tags: { occurrences: v.count },
        };
      });
    }
    return Object.entries(
      filterValues ? filterValues[type].value || {} : {},
    ).map(([k, v]) => {
      return {
        value: k,
        label: k,
        // @ts-ignore
        tags: { occurrences: v },
      };
    });
  }, [filterValues, item.key, type]);

  useDeepCompareEffect(() => {
    update(index, {
      ...item,
      value: currentValue.value,
      title: currentValue.title,
    });
  }, [currentValue, index, item]);

  const styles = useSelectInputStyles();

  return (
    <CreatableSelect
      className={classNames("basic-single", className)}
      classNamePrefix="select"
      components={{
        // @ts-ignore
        MultiValueLabel: MultiValueLabelWithTags,
        // @ts-ignore
        Option: OptionWithTags,
        // @ts-ignore
        SingleValue: SingleValueWithTags,
      }}
      createOptionPosition="first"
      formatCreateLabel={(inputValue) => `Use "${inputValue}"`}
      // @ts-ignore as this element definitely exists
      menuPortalTarget={document.getElementById("portals")}
      onChange={(t) =>
        setCurrentValue({
          value: (t as SelectOption).value,
          title: (t as SelectOption).label,
        })
      }
      options={values}
      inputId={`${type}.input`}
      placeholder="Choose a value…"
      styles={styles}
      value={values.find((t) => t.value === value)}
    />
  );
};

const CheckFilterEditorItem = ({
  filter,
  index,
  item,
  remove,
  update,
}: CheckFilterEditorItemProps) => {
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
      <div className="grow min-w-44 max-w-72">
        <CheckFilterTypeSelect
          filter={filter}
          index={index}
          item={item}
          // @ts-ignore
          type={item.type}
          update={update}
        />
      </div>
      {(item.type === "dimension" || item.type === "control_tag") && (
        <>
          <span>=</span>
          <div className="grow min-w-40 max-w-72">
            <CheckFilterKeySelect
              index={index}
              item={item}
              filterKey={item.key}
              type={item.type}
              update={update}
            />
          </div>
        </>
      )}
      <span>=</span>
      <div className="grow min-w-52 max-w-72">
        <CheckFilterValueSelect
          index={index}
          item={item}
          // @ts-ignore
          type={item.type}
          update={update}
          value={item.value}
        />
      </div>
      <span
        className={classNames(
          (filter.expressions?.length || 0) > 1
            ? "text-foreground-light hover:text-steampipe-red cursor-pointer"
            : "text-foreground-lightest",
        )}
        onClick={
          (filter.expressions?.length || 0) > 1
            ? () => remove(index)
            : undefined
        }
        title="Remove"
      >
        <Icon className="h-5 w-5" icon="trash" />
      </span>
    </Reorder.Item>
  );
};

const CheckFilterEditor = ({ filter, onApply }: CheckFilterEditorProps) => {
  const [innerFilter, setInnerFilter] = useState<CheckFilter>(filter);
  const [isDirty, setIsDirty] = useState(false);
  const [isValid, setIsValid] = useState({ value: false, reason: "" });

  useEffect(() => {
    if (!innerFilter) {
      setIsDirty(false);
      setIsValid({ value: true, reason: "" });
      return;
    }

    setIsValid({ value: validateFilter(innerFilter), reason: "" });
    setIsDirty(JSON.stringify(innerFilter) !== JSON.stringify(filter));
  }, [filter, innerFilter, setIsDirty, setIsValid]);

  const remove = useCallback(
    (index: number) => {
      setInnerFilter((existing) => ({
        ...existing,
        expressions: [
          ...(existing.expressions?.slice(0, index) || []),
          ...(existing.expressions?.slice(index + 1) || []),
        ],
      }));
    },
    [setInnerFilter],
  );

  const update = useCallback(
    (index: number, updatedItem: CheckFilter) => {
      setInnerFilter((existing) => ({
        ...existing,
        expressions: [
          ...(existing.expressions?.slice(0, index) || []),
          updatedItem,
          ...(existing.expressions?.slice(index + 1) || []),
        ],
      }));
    },
    [setInnerFilter],
  );

  return (
    <div className="flex flex-col space-y-4">
      <Reorder.Group
        axis="y"
        values={filter.expressions || []}
        onReorder={(a) => {
          if (!!innerFilter) {
            setInnerFilter((existing) => ({
              ...existing,
              expressions: a,
            }));
          }
        }}
        as="div"
        className="flex flex-col space-y-4"
      >
        {innerFilter.expressions?.map((c: CheckFilter, idx: number) => (
          <CheckFilterEditorItem
            key={`${c.type}-${c.value}`}
            filter={innerFilter}
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
        onAdd={() =>
          setInnerFilter((existing) => ({
            ...existing,
            expressions: [
              ...(existing.expressions || []),
              { operator: "equal" },
            ],
          }))
        }
        onClear={() => {
          const toSave: CheckFilter = {
            expressions: [{ operator: "equal" }],
            operator: "and",
          };
          setInnerFilter(toSave);
          onApply(toSave);
        }}
        onApply={() => onApply(innerFilter)}
        addLabel="Add filter"
      />
    </div>
  );
};

export default CheckFilterEditor;

export { validateFilter };
