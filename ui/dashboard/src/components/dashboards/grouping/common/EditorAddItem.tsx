import Icon from "@powerpipe/components/Icon";
import { getComponent } from "@powerpipe/components/dashboards";

interface EditorAddItemProps {
  addLabel?: string;
  applyLabel?: string;
  helpUrl?: string;
  isDirty: boolean;
  isValid: { value: boolean; reason: string };
  onAdd: () => void;
  onClear?: () => void;
  onApply: () => void;
}

const EditorAddItem = ({
  addLabel = "Add",
  applyLabel = "Apply",
  helpUrl,
  isDirty = false,
  isValid,
  onAdd,
  onClear,
  onApply,
}: EditorAddItemProps) => {
  const ExternalLink = getComponent("external_link");
  return (
    <div className="flex items-center justify-between space-x-3">
      <div className="flex items-center space-x-3">
        <div
          className="flex items-center space-x-3 cursor-pointer group"
          onClick={onAdd}
          title={addLabel}
        >
          <Icon
            className="block h-5 w-5 group-hover:text-foreground-light"
            icon="add"
          />
          <span className="group-hover:text-foreground-light">Add</span>
        </div>
      </div>
      <div className="flex items-center justify-end space-x-3">
        <span
          className="text-sm text-foreground-lighter cursor-pointer hover:text-link mr-2"
          onClick={onClear}
        >
          Reset
        </span>
        <button
          type="button"
          className="rounded-md bg-ok px-2.5 py-1.5 text-sm font-semibold text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={!isValid.value || !isDirty}
          onClick={onApply}
          title={isValid.reason}
        >
          {applyLabel}
        </button>
        {!!helpUrl ? (
          <ExternalLink
            to={helpUrl}
            className="block"
            title="Get help on this feature"
          >
            <Icon className="h-5 w-5 text-foreground-light" icon="help" />
          </ExternalLink>
        ) : (
          <Icon className="block h-5 w-5 invisible" icon="help_center" />
        )}
      </div>
    </div>
  );
};

export default EditorAddItem;
