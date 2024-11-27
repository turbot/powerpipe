import EditorAddItem from "../common/EditorAddItem";
import CreatableSelect from "react-select/creatable";
import Icon from "@powerpipe/components/Icon";
import Select, { MultiValue, SingleValue } from "react-select";
import useDeepCompareEffect from "use-deep-compare-effect";
import useSelectInputStyles from "../../inputs/common/useSelectInputStyles";
import { classNames } from "@powerpipe/utils/styles";
import { DashboardPanelType } from "@powerpipe/types";
import {
  DisplayGroupType,
  Filter,
  FilterOperator,
  FilterType,
} from "../common";
import { filterKeysSorter, filterTypeMap } from "@powerpipe/utils/filterEditor";
import {
  MultiValueLabelWithTags,
  OptionWithTags,
  SingleValueWithTags,
} from "../../inputs/common/Common";
import { Reorder, useDragControls } from "framer-motion";
import { SelectOption } from "../../inputs/types";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useDashboardControls } from "../../layout/Dashboard/DashboardControlsProvider";

type FilterEditorProps = {
  filter: Filter;
  panelType: DashboardPanelType;
  onApply: (toSave: Filter) => void;
};

type FilterEditorItemProps = {
  filter: Filter;
  item: Filter;
  panelType: DashboardPanelType;
  index: number;
  remove: (index: number) => void;
  update: (index: number, item: Filter) => void;
};

type FilterTypeSelectProps = {
  className?: string;
  filter: Filter;
  index: number;
  item: Filter;
  panelType: DashboardPanelType;
  type: FilterType;
  onChange: (type: DisplayGroupType) => void;
};

type FilterValueSelectProps = {
  className?: string;
  index: number;
  item: Filter;
  type: FilterType;
  onChange: (
    // value:
    //   | { value: string; title?: string }
    //   | { value: string; title?: string }[],
    value: string | string[],
  ) => void;
  value: string | string[] | undefined;
};

const validateFilter = (filter: Filter): boolean => {
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

  if (filter.operator === "equal" || filter.operator === "not_equal") {
    const valueExists = !!filter.value?.trim();
    const typeExists = !!filter.type?.trim();
    const keyExists = !!filter.key?.trim();
    if (!valueExists) return false;
    if (!typeExists) return false;
    if (keyExists) {
      return typeExists && valueExists;
    }
    return true;
  }

  if (filter.operator === "in" || filter.operator === "not_in") {
    const valueExists = !!filter.value?.length;
    const typeExists = !!filter.type?.trim();
    const keyExists = !!filter.key?.trim();
    if (!valueExists) return false;
    if (!typeExists) return false;
    if (keyExists) {
      return typeExists && valueExists;
    }
    return true;
  }

  return false;
};

const isValidFilterTypeForPanel = (
  type: string,
  panelType: DashboardPanelType,
) => {
  return !(
    type === "status" &&
    (panelType === "detection_benchmark" || panelType === "detection")
  );
};

const FilterTypeSelect = ({
  className,
  filter,
  item,
  panelType,
  type,
  onChange,
}: FilterTypeSelectProps) => {
  // const [currentType, setCurrentType] = useState<FilterType>(type);
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
        .sort(filterKeysSorter),
    [filterValues],
  );

  // useDeepCompareEffect(() => {
  //   if (currentType) {
  //     console.log("Type calling update");
  //     update(index, {
  //       ...item,
  //       value: item.value,
  //       type: currentType?.includes("|")
  //         ? (currentType?.split("|")[0] as FilterType)
  //         : currentType,
  //       key: currentType?.includes("|")
  //         ? currentType?.split("|")[1]
  //         : undefined,
  //     });
  //   }
  // }, [currentType, index, item]);

  const types = useMemo(() => {
    // @ts-ignore
    const existingTypes = filter.expressions
      ?.map((c) => c.type?.toString())
      .filter((t) => !!t);

    return allFilters.filter((t) => {
      return (
        (t.value === type ||
          t.value === "dimension" ||
          t.value === "control_tag" ||
          t.value === "detection_tag" ||
          // @ts-ignore
          !existingTypes.includes(t?.value)) &&
        isValidFilterTypeForPanel(t.value, panelType)
      );
    });
  }, [allFilters, filter, panelType, type]);

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
      onChange={(t) => {
        const v = (t as SelectOption).value;
        onChange(v as DisplayGroupType);
      }}
      options={types}
      inputId={`${type}.input`}
      placeholder="Select a filter…"
      // @ts-ignore
      styles={{
        ...styles,
        menu: (provided) => ({
          ...styles?.menu(provided),
          width: "275px",
        }),
      }}
      value={types
        .reduce((acc, curr) => {
          if (curr?.options) {
            return acc.concat(...curr.options);
          }
          return acc.concat(curr);
        }, [])
        .find((t) => t.value === type + (item?.key ? `|${item.key}` : ""))}
    />
  );
};

const FilterOperatorSelect = ({
  className,
  operator,
  onChange,
}: {
  className?: string;
  index: number;
  operator: FilterOperator;
  onChange: (operator: FilterOperator) => void;
}) => {
  const styles = useSelectInputStyles();

  // useDeepCompareEffect(() => {
  //   if (operator) {
  //     const currentOperator = item.operator;
  //     let newValue: string | string[] | undefined = item.value;
  //
  //     if (
  //       (currentOperator === "equal" || currentOperator === "not_equal") &&
  //       (operator === "in" || operator === "not_in") &&
  //       newValue !== undefined
  //     ) {
  //       newValue = [newValue];
  //     } else if (
  //       (currentOperator === "in" || currentOperator === "not_in") &&
  //       (operator === "equal" || operator === "not_equal") &&
  //       newValue !== undefined
  //     ) {
  //       newValue = newValue[0];
  //     }
  //
  //     console.log("Operator calling update", {
  //       currentOperator,
  //       operator,
  //       currentValue: item.value,
  //       value: newValue,
  //     });
  //
  //     update(index, {
  //       ...item,
  //       operator: operator,
  //       value: newValue,
  //     });
  //   }
  // }, [operator, index, item]);

  const operators = useMemo<{ value: FilterOperator; label: string }[]>(
    () => [
      { value: "equal", label: "=" },
      { value: "not_equal", label: "!=" },
      { value: "in", label: "in" },
      { value: "not_in", label: "!in" },
    ],
    [],
  );

  return (
    <Select
      className={classNames("basic-single", className)}
      classNamePrefix="select"
      components={{
        // @ts-ignore
        SingleValue: SingleValueWithTags,
      }}
      // @ts-ignore as this element definitely exists
      menuPortalTarget={document.getElementById("portals")}
      onChange={(t) => {
        const v = (t as SelectOption).value;
        onChange(v as FilterOperator);
      }}
      options={operators}
      inputId={`${operator}.input`}
      isClearable={false}
      // @ts-ignore
      styles={{
        ...styles,
        menu: (provided) => ({
          ...styles?.menu(provided),
        }),
      }}
      value={operators.find((o) => o.value === operator)}
    />
  );
};

const FilterValueSelect = ({
  className,
  item,
  type,
  value,
  onChange,
}: FilterValueSelectProps) => {
  // const [currentValue, setCurrentValue] = useState<
  //   | {
  //       value: any;
  //       title?: string;
  //     }
  //   | {
  //       value: any;
  //       title?: string;
  //     }[]
  // >(
  //   item.operator === "in" || item.operator === "not_in"
  //     ? []
  //     : { value, title: item.title },
  // );
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
    } else if (["control_tag", "detection_tag", "dimension"].includes(type)) {
      const keys = Object.entries(
        filterValues ? filterValues[type]?.key || {} : {},
      );
      return keys
        .filter(([k]) => k === item?.key)
        .flatMap(([, v]) => {
          const keys = Object.keys(v as any);
          return keys.map((key) => ({
            value: key,
            label: key,
            // @ts-ignore
            tags: { occurrences: v[key] },
          }));
        });
    } else if (
      type === "benchmark" ||
      type === "control" ||
      type === "detection_benchmark" ||
      type === "detection"
    ) {
      return Object.entries(
        filterValues ? filterValues[type]?.value || {} : {},
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
      filterValues ? filterValues[type]?.value || {} : {},
    ).map(([k, v]) => {
      return {
        value: k,
        label: k,
        // @ts-ignore
        tags: { occurrences: v },
      };
    });
  }, [filterValues, item.key, type]);

  // useDeepCompareEffect(() => {
  //   console.log("Value calling update", {
  //     currentValue: item.value,
  //     value: currentValue,
  //   });
  //   update(index, {
  //     ...item,
  //     value:
  //       (item.operator === "in" || item.operator === "not_in") &&
  //       isArray(currentValue)
  //         ? (
  //             currentValue as {
  //               value: any;
  //               title?: string;
  //             }[]
  //           ).map((t) => t.value)
  //         : (item.operator === "in" || item.operator === "not_in") &&
  //             !isArray(currentValue)
  //           ? [currentValue]
  //           : { value: currentValue.value, title: currentValue.title },
  //   });
  // }, [currentValue, index, item]);

  const styles = useSelectInputStyles();

  const currentValue =
    item.operator === "in" || item.operator === "not_in"
      ? values.filter((v) => value?.includes(v.value))
      : values.find((v) => v.value === value);

  console.log("Current value", {
    currentValue,
    value,
    values,
    operator: item.operator,
  });

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
      onChange={(t) => {
        if (item.operator === "in" || item.operator === "not_in") {
          onChange(
            (
              (t as MultiValue<{
                value: any;
                title?: string;
              }>) || []
            ).map((t) => {
              // return {
              //   value: (t as SelectOption).value,
              //   title: (t as SelectOption).label as string,
              // };
              return t.value;
            }),
          );
          return;
        }
        // onChange({
        //   value: (t as SelectOption).value,
        //   title: (t as SelectOption).label as string,
        // });
        onChange(
          (
            t as SingleValue<{
              value: any;
              title?: string;
            }>
          ).value,
        );
      }}
      options={values}
      inputId={`${type}.input`}
      isMulti={item.operator === "in" || item.operator === "not_in"}
      placeholder="Choose a value…"
      // @ts-ignore
      styles={styles}
      value={currentValue}
    />
  );
};

const FilterEditorItem = ({
  filter,
  index,
  item,
  panelType,
  remove,
  update,
}: FilterEditorItemProps) => {
  const dragControls = useDragControls();
  const [innerItem, setInnerItem] = useState<Filter>(item);

  useEffect(() => {
    console.log("Item updated, replicating to inner item", { item, innerItem });
    setInnerItem(() => item);
  }, [item]);

  const onItemChange = (updatedItem) => {
    console.log("Item changed", { item, updatedItem });
    update(index, updatedItem);
  };

  const onTypeChange = (type: DisplayGroupType) => {
    console.log("Type changed", { innerItem, type });

    const currentOperator = innerItem.operator;
    const currentType = innerItem.type;
    const newItem = {
      ...innerItem,
      value: currentOperator === "in" || currentOperator === "not_in" ? [] : "",
      type: currentType?.includes("|")
        ? (currentType?.split("|")[0] as FilterType)
        : currentType,
      key: currentType?.includes("|") ? currentType?.split("|")[1] : undefined,
    };

    onItemChange(newItem);
  };

  const onOperatorChange = (operator: FilterOperator) => {
    console.log("Operator changed", { innerItem, operator });

    if (!operator) {
      return;
    }

    const currentOperator = innerItem.operator;
    let newValue: string | string[] | undefined = innerItem.value;

    if (
      (currentOperator === "equal" || currentOperator === "not_equal") &&
      (operator === "in" || operator === "not_in") &&
      newValue !== undefined
    ) {
      newValue = [newValue];
    } else if (
      (currentOperator === "in" || currentOperator === "not_in") &&
      (operator === "equal" || operator === "not_equal") &&
      newValue !== undefined
    ) {
      newValue = newValue[0];
    }

    const newItem = {
      ...innerItem,
      operator: operator,
      value: newValue,
    };

    onItemChange(newItem);
  };

  const onValueChange = (
    value:
      | { value: string; title?: string }
      | { value: string; title?: string }[],
  ) => {
    console.log("Value changed", { innerItem, value });

    const newItem = {
      ...innerItem,
      value,
    };

    onItemChange(newItem);
  };

  return (
    <Reorder.Item
      as="div"
      id={`${innerItem.type}-${innerItem.value}`}
      className="flex space-x-3 items-center"
      dragControls={dragControls}
      dragListener={false}
      value={innerItem}
    >
      {/*<div className="flex space-x-3 items-center">*/}
      <div className="cursor-grab" onPointerDown={(e) => dragControls.start(e)}>
        <Icon className="h-5 w-5" icon="drag_indicator" />
      </div>
      <div className="grow min-w-44 max-w-72">
        <FilterTypeSelect
          filter={filter}
          index={index}
          item={innerItem}
          panelType={panelType}
          // @ts-ignore
          type={innerItem.type}
          onChange={onTypeChange}
        />
      </div>
      <div>
        <FilterOperatorSelect
          index={index}
          operator={innerItem.operator}
          onChange={onOperatorChange}
        />
      </div>
      <div className="grow min-w-52 max-w-72">
        <FilterValueSelect
          index={index}
          item={innerItem}
          // @ts-ignore
          type={innerItem.type}
          value={innerItem.value}
          onChange={onValueChange}
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

const FilterEditor = ({ filter, panelType, onApply }: FilterEditorProps) => {
  const [innerFilter, setInnerFilter] = useState<Filter>(filter);
  const [isDirty, setIsDirty] = useState(false);
  const [isValid, setIsValid] = useState({ value: false, reason: "" });

  useEffect(() => {
    if (!innerFilter) {
      setIsDirty(() => false);
      setIsValid(() => ({ value: true, reason: "" }));
      return;
    }

    setIsValid({ value: validateFilter(innerFilter), reason: "" });
    setIsDirty(JSON.stringify(innerFilter) !== JSON.stringify(filter));
  }, [filter, innerFilter]);

  useDeepCompareEffect(() => {
    setInnerFilter(filter);
  }, [filter]);

  const remove = useCallback((index: number) => {
    setInnerFilter((existing) => ({
      ...existing,
      expressions: [
        ...(existing.expressions?.slice(0, index) || []),
        ...(existing.expressions?.slice(index + 1) || []),
      ],
    }));
  }, []);

  const update = useCallback((index: number, updatedItem: Filter) => {
    console.log("update", { index, updatedItem });
    setInnerFilter((existing) => ({
      ...existing,
      expressions: [
        ...(existing.expressions?.slice(0, index) || []),
        updatedItem,
        ...(existing.expressions?.slice(index + 1) || []),
      ],
    }));
  }, []);

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
        {innerFilter.expressions?.map((c: Filter, idx: number) => (
          <FilterEditorItem
            key={`${c.type}-${c.value}`}
            filter={innerFilter}
            item={c}
            panelType={panelType}
            index={idx}
            remove={remove}
            update={update}
          />
        ))}
      </Reorder.Group>
      <EditorAddItem
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
          const toSave: Filter = {
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

export default FilterEditor;

export { validateFilter };
