import { getComponent } from "../dashboards";
// @ts-ignore
import { ReactComponent as Logo } from "./logos/powerpipe-logo.svg";
// @ts-ignore
import { ReactComponent as LogoDarkmode } from "./logos/powerpipe-logo-darkmode.svg";
// @ts-ignore
import { ReactComponent as LogoWordmark } from "./logos/powerpipe-logo-wordmark.svg";
// @ts-ignore
import { ReactComponent as LogoWordmarkDarkmode } from "./logos/powerpipe-logo-wordmark-darkmode.svg";
import { ThemeNames } from "@powerpipe/hooks/useTheme";
import { useDashboardSearchPath } from "@powerpipe/hooks/useDashboardSearchPath";
import { useDashboardTheme } from "@powerpipe/hooks/useDashboardTheme";

const PowerpipeLogo = () => {
  const { theme } = useDashboardTheme();
  const { searchPathPrefix } = useDashboardSearchPath();
  const ExternalLink = getComponent("external_link");

  return (
    <div className="mr-1 md:mr-4">
      <ExternalLink
        ignoreDataMode
        to={`/${!!searchPathPrefix.length ? `?search_path_prefix=${searchPathPrefix}` : ""}`}
      >
        <div className="block md:hidden w-8">
          {theme.name === ThemeNames.STEAMPIPE_DEFAULT && <Logo />}
          {theme.name === ThemeNames.STEAMPIPE_DARK && <LogoDarkmode />}
        </div>
        <div className="hidden md:block w-48">
          {theme.name === ThemeNames.STEAMPIPE_DEFAULT && <LogoWordmark />}
          {theme.name === ThemeNames.STEAMPIPE_DARK && <LogoWordmarkDarkmode />}
        </div>
      </ExternalLink>
    </div>
  );
};

export default PowerpipeLogo;
