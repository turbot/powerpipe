import CheckEditorAddItem from "../common/CheckEditorAddItem";
import CreatableSelect from "react-select/creatable";
import Icon from "components/Icon";
import Select from "react-select";
import useDeepCompareEffect from "use-deep-compare-effect";
import useSelectInputStyles from "../../inputs/common/useSelectInputStyles";
import { CheckFilter, CheckFilterType } from "../common";
import { classNames } from "utils/styles";
import {
  MultiValueLabelWithTags,
  OptionWithTags,
  SingleValueWithTags,
} from "../../inputs/common/Common";
import { Reorder, useDragControls } from "framer-motion";
import { SelectOption } from "../../inputs/types";
import { useCallback, useMemo, useState } from "react";
import { useDashboardControls } from "../../layout/Dashboard/DashboardControlsProvider";

type CheckFilterEditorProps = {
  config: CheckFilter;
  isValid: { value: boolean; reason: string };
  onCancel: () => void;
  onSave: () => void;
  setConfig: (newValue: CheckFilter) => void;
};

type CheckFilterEditorItemProps = {
  config: CheckFilter;
  item: CheckFilter;
  index: number;
  remove: (index: number) => void;
  update: (index: number, item: CheckFilter) => void;
};

type CheckFilterTypeSelectProps = {
  config: CheckFilter;
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
  index: number;
  item: CheckFilter;
  type: CheckFilterType;
  update: (index: number, updatedItem: CheckFilter) => void;
  value: string | undefined;
};

const CheckFilterTypeSelect = ({
  config,
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
    const existingTypes = config.expressions
      ?.map((c) => c.type?.toString())
      .filter((t) => !!t);
    const allTypes: SelectOption[] = [
      { value: "benchmark", label: "Benchmark" },
      { value: "control", label: "Control" },
      { value: "dimension", label: "Dimension" },
      { value: "reason", label: "Reason" },
      { value: "resource", label: "Resource" },
      { value: "severity", label: "Severity" },
      { value: "status", label: "Status" },
      { value: "tag", label: "Tag" },
    ];
    return allTypes.filter(
      (t) =>
        t.value === type ||
        t.value === "dimension" ||
        t.value === "tag" ||
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
    return Object.keys(filterValues[type].key || {}).map((k) => ({
      value: k,
      label: k,
    }));
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
  index,
  item,
  type,
  value,
  update,
}: CheckFilterValueSelectProps) => {
  const [currentValue, setCurrentValue] = useState(value);
  const { context: filterValues } = useDashboardControls();

  useDeepCompareEffect(() => {
    update(index, {
      ...item,
      value: currentValue,
    });
  }, [currentValue, index, item]);

  const values = useMemo(() => {
    if (!type) {
      return [];
    }
    if (type === "status") {
      return (
        Object.entries(filterValues[type] || {})
          // @ts-ignore
          .filter(([, v]) => v > 0)
          .map(([k, v]) => ({
            value: k,
            label: k,
            tags: { occurrences: v },
          }))
      );
    }
    return Object.entries(filterValues[type].value || {}).map(([k, v]) => ({
      value: k,
      label: k,
      tags: { occurrences: v },
    }));
  }, [filterValues, type]);

  const styles = useSelectInputStyles();

  return (
    <CreatableSelect
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
      createOptionPosition="first"
      formatCreateLabel={(inputValue) => `Use "${inputValue}"`}
      // @ts-ignore as this element definitely exists
      menuPortalTarget={document.getElementById("portals")}
      onChange={(t) => setCurrentValue((t as SelectOption).value)}
      options={values}
      inputId={`${type}.input`}
      placeholder="Choose a value…"
      styles={styles}
      value={values.find((t) => t.value === value)}
    />
  );
};

const CheckFilterEditorItem = ({
  config,
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
      <div className="grow">
        <CheckFilterTypeSelect
          config={config}
          index={index}
          item={item}
          // @ts-ignore
          type={item.type}
          update={update}
        />
      </div>
      {(item.type === "dimension" || item.type === "tag") && (
        <>
          <span>=</span>
          <div className="grow">
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
      <div className="grow">
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
          // @ts-ignore
          config.expressions?.length > 0
            ? "text-foreground-light hover:text-steampipe-red cursor-pointer"
            : "text-foreground-lightest",
        )}
        // @ts-ignore
        onClick={() => remove(index)}
        title="Remove"
      >
        <Icon className="h-5 w-5" icon="trash" />
      </span>
    </Reorder.Item>
  );
};

const CheckFilterEditor = ({
  config,
  setConfig,
  isValid,
  onCancel,
  onSave,
}: CheckFilterEditorProps) => {
  const remove = useCallback(
    (index: number) => {
      const newConfig: CheckFilter = {
        ...config,
        expressions: [
          ...(config.expressions?.slice(0, index) || []),
          ...(config.expressions?.slice(index + 1) || []),
        ],
      };
      setConfig(newConfig);
    },
    [config, setConfig],
  );

  const update = useCallback(
    (index: number, updatedItem: CheckFilter) => {
      const newConfig: CheckFilter = {
        ...config,
        expressions: [
          ...(config.expressions?.slice(0, index) || []),
          updatedItem,
          ...(config.expressions?.slice(index + 1) || []),
        ],
      };
      setConfig(newConfig);
    },
    [config, setConfig],
  );

  return (
    <div className="flex flex-col space-y-4">
      {(config.expressions?.length || 0) > 0 && (
        <Reorder.Group
          axis="y"
          values={config.expressions || []}
          onReorder={(a) => {
            if (!!config) {
              const newConfig = {
                ...config,
                expressions: a,
              };
              setConfig(newConfig);
            }
          }}
          as="div"
          className="flex flex-col space-y-4"
        >
          {config.expressions?.map((c: CheckFilter, idx: number) => (
            <CheckFilterEditorItem
              key={`${c.type}-${c.value}`}
              config={config}
              item={c}
              index={idx}
              remove={remove}
              update={update}
            />
          ))}
        </Reorder.Group>
      )}
      <CheckEditorAddItem
        addLabel="Add filter"
        clearLabel="Clear filter"
        isValid={isValid}
        onAdd={() =>
          setConfig({
            ...config,
            expressions: [
              // @ts-ignore
              ...(config.expressions || []),
              // @ts-ignore
              { operator: "equal", type: "" },
            ],
          })
        }
        onClear={() => {
          setConfig({
            ...config,
            expressions: [],
          });
          onSave();
        }}
        onCancel={onCancel}
        onSave={onSave}
      />
    </div>
  );
};

export default CheckFilterEditor;
