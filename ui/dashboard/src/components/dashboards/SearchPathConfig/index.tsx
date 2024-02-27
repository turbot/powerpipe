import SearchPathEditor from "@powerpipe/components/dashboards/SearchPathEditor";
import useDashboardSearchPathPrefix from "@powerpipe/hooks/useDashboardSearchPathPrefix";
import { DashboardActions } from "@powerpipe/types";
import { Noop } from "@powerpipe/types/func";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useEffect, useMemo } from "react";
import { useSearchParams } from "react-router-dom";

type SearchPathConfigProps = {
  onClose: Noop;
};

const SearchPathConfig = ({ onClose }: SearchPathConfigProps) => {
  const { dispatch, selectedDashboard, dashboardsMetadata } = useDashboard();
  const [, setSearchParams] = useSearchParams();
  const searchPathPrefix = useDashboardSearchPathPrefix();

  const { configuredSearchPath, availableConnections } = useMemo(() => {
    if (!selectedDashboard || !dashboardsMetadata || !searchPathPrefix) {
      return { configuredSearchPath: [], availableConnections: [] };
    }
    const metadata = dashboardsMetadata[selectedDashboard?.full_name];
    return {
      configuredSearchPath: searchPathPrefix,
      availableConnections: metadata.original_search_path || [],
    };
  }, [dashboardsMetadata, searchPathPrefix, selectedDashboard]);

  useEffect(() => {
    dispatch({
      type: DashboardActions.SET_SELECTED_DASHBOARD_SEARCH_PATH_PREFIX,
      search_path_prefix: searchPathPrefix,
    });
  }, [dispatch, searchPathPrefix]);

  const saveSearchPath = (toSave: string[]) => {
    setSearchParams((previous) => {
      const newParams = new URLSearchParams(previous);
      if (!!toSave.length) {
        newParams.set("search_path_prefix", toSave.join(","));
      } else {
        newParams.delete("search_path_prefix");
      }
      return newParams;
    });
  };

  return (
    <SearchPathEditor
      availableConnections={availableConnections}
      searchPathPrefix={configuredSearchPath}
      onCancel={onClose}
      onApply={saveSearchPath}
      onSave={(toSave) => {
        saveSearchPath(toSave);
        onClose();
      }}
    />
  );
};

export default SearchPathConfig;
