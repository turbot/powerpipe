import SearchPathEditor from "@powerpipe/components/dashboards/SearchPathEditor";
import useDashboardSearchPathPrefix from "@powerpipe/hooks/useDashboardSearchPathPrefix";
import { SearchPathMetadata } from "@powerpipe/types";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const SearchPathConfig = ({ onClose }) => {
  const { selectedDashboard, metadata, dashboardsMetadata } = useDashboard();
  const [, setSearchParams] = useSearchParams();
  const searchPathPrefix = useDashboardSearchPathPrefix();

  const { configuredSearchPath, availableConnections } = useMemo(() => {
    const hasServerMetadataSearchPath =
      !!metadata?.search_path?.original_search_path &&
      !!metadata.search_path.original_search_path.length;
    const hasDashboardMetadataSearchPath =
      !!selectedDashboard &&
      !!dashboardsMetadata &&
      !!dashboardsMetadata[selectedDashboard.full_name] &&
      !!dashboardsMetadata[selectedDashboard.full_name]?.search_path &&
      !!dashboardsMetadata[selectedDashboard.full_name]?.search_path
        ?.original_search_path &&
      !!dashboardsMetadata[selectedDashboard.full_name]?.search_path
        ?.original_search_path?.length;

    if (!hasServerMetadataSearchPath && !hasDashboardMetadataSearchPath) {
      return { configuredSearchPath: [], availableConnections: [] };
    }

    let foundMetadata: SearchPathMetadata | null = null;
    if (selectedDashboard) {
      foundMetadata =
        dashboardsMetadata[selectedDashboard?.full_name]?.search_path;
    }

    if (!foundMetadata && !!metadata) {
      foundMetadata = metadata.search_path;
    }

    return {
      configuredSearchPath: searchPathPrefix,
      availableConnections: foundMetadata
        ? foundMetadata.original_search_path || []
        : [],
    };
  }, [metadata, dashboardsMetadata, searchPathPrefix, selectedDashboard]);

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
    onClose();
  };

  return (
    <SearchPathEditor
      availableConnections={availableConnections}
      searchPathPrefix={configuredSearchPath}
      onApply={saveSearchPath}
    />
  );
};

export default SearchPathConfig;
