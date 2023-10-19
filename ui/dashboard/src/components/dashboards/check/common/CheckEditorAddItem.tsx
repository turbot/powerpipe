import Icon from "../../../Icon";
import { classNames } from "../../../../utils/styles";
import { ThemeNames } from "../../../../hooks/useTheme";
import { useDashboard } from "../../../../hooks/useDashboard";

interface CheckEditorAddItemProps {
  label: string;
  onClick: () => void;
}

const CheckEditorAddItem = ({ label, onClick }: CheckEditorAddItemProps) => {
  const {
    themeContext: { theme },
  } = useDashboard();
  return (
    <div className="flex items-center space-x-3">
      <Icon className="block h-5 w-5 invisible" icon="drag_indicator" />
      <span
        className={classNames(
          "flex items-center cursor-pointer min-h-[38px] pl-[8px] border rounded-[4px] grow",
          theme.name === ThemeNames.STEAMPIPE_DARK
            ? "border-[#444] hover:border-[#b3b3b3]"
            : "border-[#d3d3d3] hover:border-[#b3b3b3]",
        )}
        onClick={onClick}
      >
        <span className="block text-foreground-lighter">{label}</span>
      </span>
      <Icon className="block h-5 w-5 invisible" icon="trash" />
    </div>
  );
};

export default CheckEditorAddItem;
