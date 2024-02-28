import Icon from "@powerpipe/components/Icon";

interface CheckEditorAddItemProps {
  addLabel?: string;
  applyLabel?: string;
  isDirty: boolean;
  isValid: { value: boolean; reason: string };
  onAdd: () => void;
  onClear?: () => void;
  onApply: () => void;
}

const CheckEditorAddItem = ({
  addLabel = "Add",
  applyLabel = "Apply",
  isDirty = false,
  isValid,
  onAdd,
  onClear,
  onApply,
}: CheckEditorAddItemProps) => (
  <div className="flex items-center justify-between space-x-3">
    <div className="flex items-center space-x-3">
      <div className="flex items-center">
        <Icon
          className="block h-5 w-5 cursor-pointer hover:text-foreground-light"
          icon="add"
          onClick={onAdd}
          title={addLabel}
        />
      </div>
    </div>
    <div className="flex items-center justify-end space-x-2">
      <button
        type="button"
        className="rounded-md bg-dashboard-panel border border-black-scale-3 px-2.5 py-1.5 text-sm font-semibold text-foreground"
        onClick={onClear}
      >
        Reset
      </button>
      <button
        type="button"
        className="rounded-md bg-ok px-2.5 py-1.5 text-sm font-semibold text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
        disabled={!isValid.value || !isDirty}
        onClick={onApply}
        title={isValid.reason}
      >
        {applyLabel}
      </button>
      {/*<button*/}
      {/*  type="button"*/}
      {/*  className="rounded-md bg-ok px-2.5 py-1.5 text-sm font-semibold text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"*/}
      {/*  disabled={!isValid.value}*/}
      {/*  onClick={onSave}*/}
      {/*  title={isValid.reason}*/}
      {/*>*/}
      {/*  Save*/}
      {/*</button>*/}
      <Icon className="block h-5 w-5 invisible" icon="trash" />
    </div>
  </div>
);

export default CheckEditorAddItem;
