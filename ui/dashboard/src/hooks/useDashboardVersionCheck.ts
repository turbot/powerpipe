import { IDashboardContext } from "../types";
import { useEffect } from "react";

const useDashboardVersionCheck = (state: IDashboardContext) => {
  useEffect(() => {
    let cliVersion: string | null = "";
    let uiVersion: string | null = "";
    let mismatchedVersions = false;
    if (state.versionMismatchCheck) {
      console.group("Dashboard Version Check");
      const cliVersionRaw = state.metadata?.cli?.version;
      const uiVersionRaw = process.env.REACT_APP_VERSION;
      const hasVersionsSet = !!cliVersionRaw && !!uiVersionRaw;
      cliVersion = !!cliVersionRaw
        ? cliVersionRaw.startsWith("v")
          ? cliVersionRaw.substring(1)
          : cliVersionRaw
        : null;
      uiVersion = !!uiVersionRaw
        ? uiVersionRaw.startsWith("v")
          ? uiVersionRaw.substring(1)
          : uiVersionRaw
        : null;
      mismatchedVersions = hasVersionsSet && cliVersion !== uiVersion;
      console.log({
        state,
        raw_versions: { cliVersionRaw, uiVersionRaw },
        versions: { cliVersion, uiVersion },
        has_versions_set: hasVersionsSet,
        mismatched_versions: hasVersionsSet && cliVersion !== uiVersion,
      });
      console.groupEnd();

      // Add a version to force a reload with the new version to get the correct assets
      if (mismatchedVersions) {
        const searchParams = new URLSearchParams(window.location.search);
        // @ts-ignore this is always truthy as mismatchedVersions is only true if hasVersionsSet is true and hasVersionsSet is only true if cliVersionRaw is truthy
        searchParams.set("version", cliVersionRaw);
        window.location.replace(`${window.location.origin}?${searchParams}`);
      }
    }
  }, [state]);
};

export default useDashboardVersionCheck;
