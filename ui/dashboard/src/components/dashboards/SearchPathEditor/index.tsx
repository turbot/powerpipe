import CheckEditorAddItem from "@powerpipe/components/dashboards/check/common/CheckEditorAddItem";
import Icon from "@powerpipe/components/Icon";
import { classNames } from "@powerpipe/utils/styles";
import { Reorder, useDragControls } from "framer-motion";
import { useCallback, useEffect, useMemo, useState } from "react";
import Select from "react-select";
import {
  MultiValueLabelWithTags,
  OptionWithTags,
  SingleValueWithTags,
} from "@powerpipe/components/dashboards/inputs/common/Common";
import { SelectOption } from "@powerpipe/components/dashboards/inputs/types";
import useSelectInputStyles from "@powerpipe/components/dashboards/inputs/common/useSelectInputStyles";

type SearchPathEditorProps = {
  availableConnections: string[];
  searchPathPrefix: string[];
  onCancel: () => void;
  onApply: (newValue: string[]) => void;
  onSave: (newValue: string[]) => void;
};

interface SearchPathEditorItemSelectProps {
  availableConnections: string[];
  searchPathPrefix: string[];
  item: string;
  index: number;
  update: (index: number, item: string) => void;
}

interface SearchPathEditorItemProps {
  availableConnections: string[];
  searchPathPrefix: string[];
  item: string;
  index: number;
  remove?: (index: number) => void;
  update: (index: number, item: string) => void;
}

const SearchPathEditorItemSelect = ({
  availableConnections,
  searchPathPrefix,
  item,
  index,
  update,
}: SearchPathEditorItemSelectProps) => {
  const connections = useMemo(() => {
    const c = (availableConnections || [])
      .filter(
        (c) => c === item || !(searchPathPrefix || []).find((s) => s === c),
      )
      .map((c) => ({ value: c, label: c }));
    c.sort((x, y) => {
      if (x.value < y.value) {
        return -1;
      } else if (x.value > y.value) {
        return 1;
      } else {
        return 0;
      }
    });
    return c;
  }, [availableConnections, item, searchPathPrefix]);

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
      onChange={(t) => update(index, (t as SelectOption).value)}
      options={connections}
      inputId={`${index}.input`}
      placeholder="Select a connection…"
      styles={styles}
      value={connections.find((t) => t.value === item)}
    />
  );
};

const SearchPathEditorItem = ({
  availableConnections,
  searchPathPrefix,
  index,
  item,
  remove,
  update,
}: SearchPathEditorItemProps) => {
  const dragControls = useDragControls();

  return (
    <Reorder.Item
      as="div"
      id={`${item}-${index}`}
      className="flex space-x-3 items-center"
      dragControls={dragControls}
      dragListener={false}
      value={item}
    >
      <div className="cursor-grab" onPointerDown={(e) => dragControls.start(e)}>
        <Icon className="h-5 w-5" icon="drag_indicator" />
      </div>
      <div className="grow">
        <SearchPathEditorItemSelect
          availableConnections={availableConnections}
          searchPathPrefix={searchPathPrefix}
          index={index}
          item={item}
          // @ts-ignore
          update={update}
        />
      </div>
      <span
        className="text-foreground-light hover:text-steampipe-red cursor-pointer"
        onClick={remove ? () => remove(index) : undefined}
        title="Remove"
      >
        <Icon className="h-5 w-5" icon="trash" />
      </span>
    </Reorder.Item>
  );
};

const SearchPathEditor = ({
  availableConnections,
  searchPathPrefix,
  onCancel,
  onApply,
  onSave,
}: SearchPathEditorProps) => {
  const [innerSearchPathPrefix, setInnerSearchPathPrefix] =
    useState<string[]>(searchPathPrefix);

  const [isValid, setIsValid] = useState({ value: false, reason: "" });

  useEffect(() => {
    const isValid = innerSearchPathPrefix.every((c) => !!c);
    setIsValid({
      value: isValid,
      reason: !isValid ? "Search path contains empty connection" : "",
    });
  }, [innerSearchPathPrefix, setIsValid]);

  const remove = useCallback(
    (index: number) => {
      setInnerSearchPathPrefix((existing) => [
        ...existing.slice(0, index),
        ...existing.slice(index + 1),
      ]);
    },
    [setInnerSearchPathPrefix],
  );

  const update = useCallback(
    (index: number, updatedItem: string) => {
      setInnerSearchPathPrefix((existing) => [
        ...existing.slice(0, index),
        updatedItem,
        ...existing.slice(index + 1),
      ]);
    },
    [setInnerSearchPathPrefix],
  );

  return (
    <div className="flex flex-col space-y-4">
      <Reorder.Group
        axis="y"
        values={innerSearchPathPrefix}
        onReorder={setInnerSearchPathPrefix}
        as="div"
        className="flex flex-col space-y-4"
      >
        {innerSearchPathPrefix.map((connection, idx) => (
          <SearchPathEditorItem
            key={connection}
            availableConnections={availableConnections}
            searchPathPrefix={innerSearchPathPrefix}
            item={connection}
            index={idx}
            remove={remove}
            update={update}
          />
        ))}
        {!innerSearchPathPrefix.length && (
          <SearchPathEditorItem
            availableConnections={availableConnections}
            searchPathPrefix={innerSearchPathPrefix}
            item=""
            index={innerSearchPathPrefix.length}
            update={update}
          />
        )}
      </Reorder.Group>
      <CheckEditorAddItem
        isValid={isValid}
        onAdd={() => setInnerSearchPathPrefix((existing) => [...existing, ""])}
        onClear={() => {
          setInnerSearchPathPrefix([]);
          onSave([]);
        }}
        onCancel={onCancel}
        onApply={() => onApply(innerSearchPathPrefix)}
        onSave={() => onSave(innerSearchPathPrefix)}
      />
    </div>
  );
};

export default SearchPathEditor;
