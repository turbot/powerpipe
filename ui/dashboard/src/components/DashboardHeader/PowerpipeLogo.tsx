import { getComponent } from "../dashboards";
// @ts-ignore
import Logo from "./logos/powerpipe-logo.svg?react";
// @ts-ignore
import LogoDarkmode from "./logos/powerpipe-logo-darkmode.svg?react";
// @ts-ignore
import LogoWordmark from "./logos/powerpipe-logo-wordmark.svg?react";
// @ts-ignore
import LogoWordmarkDarkmode from "./logos/powerpipe-logo-wordmark-darkmode.svg?react";
import { ThemeNames } from "@powerpipe/hooks/useTheme";
import { useDashboard } from "@powerpipe/hooks/useDashboard";

const PowerpipeLogo = () => {
  const {
    themeContext: { theme },
    searchPathPrefix,
  } = useDashboard();
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
