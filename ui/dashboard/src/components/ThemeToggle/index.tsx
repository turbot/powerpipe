import DashboardIcon from "../dashboards/common/DashboardIcon";
import { classNames } from "@powerpipe/utils/styles";
import { ThemeNames } from "@powerpipe/hooks/useTheme";
import { useDashboardTheme } from "@powerpipe/hooks/useDashboardTheme";

const ThemeToggle = () => {
  const { theme, setTheme } = useDashboardTheme();

  return (
    <button
      className={classNames("flex items-center h-5 w-5 text-gray-500")}
      onClick={() =>
        setTheme(
          theme.name === ThemeNames.STEAMPIPE_DEFAULT
            ? ThemeNames.STEAMPIPE_DARK
            : ThemeNames.STEAMPIPE_DEFAULT,
        )
      }
    >
      <DashboardIcon
        icon={
          theme.name === ThemeNames.STEAMPIPE_DARK
            ? "materialsymbols-solid:light_mode"
            : "materialsymbols-solid:dark_mode"
        }
        title={
          theme.name === ThemeNames.STEAMPIPE_DARK
            ? "Switch to light theme"
            : "Switch to dark theme"
        }
      />
    </button>
  );
};

export default ThemeToggle;
