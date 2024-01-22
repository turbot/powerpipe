import Icon from "../../../Icon";
import { classNames } from "../../../../utils/styles";
import { ThemeNames } from "../../../../hooks/useTheme";
import { useDashboard } from "../../../../hooks/useDashboard";

interface CheckEditorAddItemProps {
  label: string;
  isValid: { value: boolean; reason: string };
  onClick: () => void;
  onCancel: () => void;
  onSave: () => void;
}

const CheckEditorAddItem = ({
  label,
  isValid,
  onClick,
  onCancel,
  onSave,
}: CheckEditorAddItemProps) => {
  const {
    themeContext: { theme },
  } = useDashboard();
  return (
    <div className="flex items-center space-x-3">
      <Icon className="block h-5 w-5 invisible" icon="drag_indicator" />
      <span
        className={classNames(
          // "flex items-center cursor-pointer min-h-[38px] pl-[8px] border rounded-[4px] grow",
          "flex items-center text-link cursor-pointer min-h-[38px] pl-[8px] grow",
          theme.name === ThemeNames.STEAMPIPE_DARK
            ? "border-[#444] hover:border-[#b3b3b3]"
            : "border-[#d3d3d3] hover:border-[#b3b3b3]",
        )}
        onClick={onClick}
      >
        <span className="block">{label}</span>
      </span>
      <div className="flex items-center space-x-2 justify-end mr-8">
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
      </div>
      <Icon className="block h-5 w-5 invisible" icon="trash" />
    </div>
  );
};

export default CheckEditorAddItem;
