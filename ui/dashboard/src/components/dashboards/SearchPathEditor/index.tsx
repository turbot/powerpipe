import CheckEditorAddItem from "@powerpipe/components/dashboards/check/common/CheckEditorAddItem";
import Icon from "@powerpipe/components/Icon";
import { classNames } from "@powerpipe/utils/styles";
import { Reorder, useDragControls } from "framer-motion";
import { useCallback, useEffect, useState } from "react";

type SearchPathEditorProps = {
  searchPathPrefix: string[];
  onCancel: () => void;
  onApply: (newValue: string[]) => void;
  onSave: (newValue: string[]) => void;
};

type SearchPathEditorItemProps = {
  searchPathPrefix: string[];
  item: string;
  index: number;
  remove: (index: number) => void;
  update: (index: number, item: string) => void;
};

const SearchPathEditorItem = ({
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
        <input
          type="text"
          name={`search-path-${index}`}
          id={`search-path-${index}`}
          className="flex-1 block w-full bg-dashboard-panel rounded-md border border-black-scale-3 pr-8 overflow-x-auto text-sm md:text-base disabled:bg-black-scale-1 focus:ring-0"
          onChange={(e) => update(index, e.target.value)}
          placeholder="Enter a connection"
          value={item}
        />
      </div>
      <span
        className={classNames(
          searchPathPrefix.length > 1
            ? "text-foreground-light hover:text-steampipe-red cursor-pointer"
            : "text-foreground-lightest",
        )}
        onClick={searchPathPrefix.length > 1 ? () => remove(index) : undefined}
        title={
          searchPathPrefix.length > 1
            ? "Remove"
            : "Search path must contain at least one entry"
        }
      >
        <Icon className="h-5 w-5" icon="trash" />
      </span>
    </Reorder.Item>
  );
};

const SearchPathEditor = ({
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
            searchPathPrefix={innerSearchPathPrefix}
            item={connection}
            index={idx}
            remove={remove}
            update={update}
          />
        ))}
      </Reorder.Group>
      <CheckEditorAddItem
        addLabel="Add connection"
        clearLabel="Reset"
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
