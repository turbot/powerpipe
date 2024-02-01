import Icon from "components/Icon";
import { classNames } from "utils/styles";
import { ThemeNames } from "hooks/useTheme";
import { useDashboard } from "hooks/useDashboard";

interface CheckEditorAddItemProps {
  addLabel: string;
  clearLabel?: string;
  isValid: { value: boolean; reason: string };
  onAdd: () => void;
  onClear?: () => void;
  onCancel: () => void;
  onSave: () => void;
}

const CheckEditorAddItem = ({
  addLabel,
  clearLabel,
  isValid,
  onAdd,
  onClear,
  onCancel,
  onSave,
}: CheckEditorAddItemProps) => {
  const {
    themeContext: { theme },
  } = useDashboard();
  return (
    <div className="flex items-center justify-between space-x-3">
      <div className="flex items-center space-x-3">
        <div className="flex items-center">
          <Icon className="block h-5 w-5 invisible" icon="drag_indicator" />
          <span
            className={classNames(
              "flex items-center text-link cursor-pointer min-h-[38px] pl-[8px] grow",
              theme.name === ThemeNames.STEAMPIPE_DARK
                ? "border-[#444] hover:border-[#b3b3b3]"
                : "border-[#d3d3d3] hover:border-[#b3b3b3]",
            )}
            onClick={onAdd}
          >
            <span className="block">{addLabel}</span>
          </span>
        </div>
        {!!onClear && (
          <div className="flex items-center">
            <Icon className="block h-5 w-5 invisible" icon="drag_indicator" />
            <span
              className={classNames(
                "flex items-center text-link cursor-pointer min-h-[38px] pl-[8px] grow",
                theme.name === ThemeNames.STEAMPIPE_DARK
                  ? "border-[#444] hover:border-[#b3b3b3]"
                  : "border-[#d3d3d3] hover:border-[#b3b3b3]",
              )}
              onClick={onClear}
            >
              <span className="block">{clearLabel}</span>
            </span>
          </div>
        )}
      </div>
      <div className="flex items-center justify-end space-x-2">
        <button
          type="button"
          className="rounded-md bg-dashboard-panel border border-gray-200 px-2.5 py-1.5 text-sm font-semibold text-foreground"
          onClick={onCancel}
        >
          Cancel
        </button>
        <button
          type="button"
          className="rounded-md bg-ok px-2.5 py-1.5 text-sm font-semibold text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={!isValid.value}
          onClick={onSave}
          title={isValid.reason}
        >
          Update
        </button>
        <Icon className="block h-5 w-5 invisible" icon="trash" />
      </div>
    </div>
  );
};

export default CheckEditorAddItem;
