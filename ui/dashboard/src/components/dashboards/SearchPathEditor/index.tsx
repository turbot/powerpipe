import CheckEditorAddItem from "components/dashboards/check/common/CheckEditorAddItem";
import Icon from "components/Icon";
import { classNames } from "utils/styles";
import { Reorder, useDragControls } from "framer-motion";
import { useCallback, useState } from "react";

type SearchPathEditorProps = {
  originalSearchPath: string[];
  isValid: { value: boolean; reason: string };
  onCancel: () => void;
  onSave: (newValue: string[]) => void;
};

type SearchPathEditorItemProps = {
  searchPath: string[];
  item: string;
  index: number;
  remove: (index: number) => void;
  update: (index: number, item: string) => void;
};

const SearchPathEditorItem = ({
  searchPath,
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
          searchPath.length > 1
            ? "text-foreground-light hover:text-steampipe-red cursor-pointer"
            : "text-foreground-lightest",
        )}
        onClick={searchPath.length > 1 ? () => remove(index) : undefined}
        title={
          searchPath.length > 1
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
  isValid,
  originalSearchPath,
  onCancel,
  onSave,
}: SearchPathEditorProps) => {
  const [innerSearchPath, setInnerSearchPath] = useState(originalSearchPath);

  const remove = useCallback(
    (index: number) => {
      const removed = [
        ...innerSearchPath.slice(0, index),
        ...innerSearchPath.slice(index + 1),
      ];
      setInnerSearchPath(removed);
    },
    [innerSearchPath, setInnerSearchPath],
  );

  const update = useCallback(
    (index: number, updatedItem: string) => {
      const updated = [
        ...innerSearchPath.slice(0, index),
        updatedItem,
        ...innerSearchPath.slice(index + 1),
      ];
      setInnerSearchPath(updated);
    },
    [innerSearchPath, setInnerSearchPath],
  );

  return (
    <div className="flex flex-col space-y-4">
      <Reorder.Group
        axis="y"
        values={innerSearchPath}
        onReorder={setInnerSearchPath}
        as="div"
        className="flex flex-col space-y-4"
      >
        {innerSearchPath.map((connection, idx) => (
          <SearchPathEditorItem
            key={connection}
            searchPath={innerSearchPath}
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
        onAdd={() => setInnerSearchPath([...innerSearchPath, ""])}
        onClear={() => onSave([])}
        onCancel={onCancel}
        onSave={() => onSave(innerSearchPath)}
      />
    </div>
  );
};

export default SearchPathEditor;
