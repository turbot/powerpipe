import Icon from "@powerpipe/components/Icon";
import SearchPathEditor from "@powerpipe/components/dashboards/SearchPathEditor";
import useDashboardSearchPathPrefix from "@powerpipe/hooks/useDashboardSearchPathPrefix";
import { DashboardActions } from "@powerpipe/types";
import { ReactNode, useEffect, useMemo, useState } from "react";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useSearchParams } from "react-router-dom";

const SearchPathConfig = () => {
  const { dispatch, selectedDashboard, dashboardsMetadata } = useDashboard();
  const [showEditor, setShowEditor] = useState(false);
  const [isValid, setIsValid] = useState({ value: false, reason: "" });
  const [, setSearchParams] = useSearchParams();
  const searchPathPrefix = useDashboardSearchPathPrefix();

  const configuredSearchPath = useMemo(() => {
    if (!selectedDashboard || !dashboardsMetadata || !searchPathPrefix) {
      return [];
    }
    if (!!searchPathPrefix.length) {
      return searchPathPrefix;
    }
    return [];
    // if (!dashboardsMetadata[selectedDashboard.full_name]) {
    //   return [];
    // }
    // return (
    //   dashboardsMetadata[selectedDashboard.full_name].original_search_path || []
    // );
  }, [dashboardsMetadata, searchPathPrefix, selectedDashboard]);

  const [modifiedSearchPath, setModifiedSearchPath] =
    useState<string[]>(configuredSearchPath);

  useEffect(
    () => setModifiedSearchPath(() => configuredSearchPath),
    [configuredSearchPath, setModifiedSearchPath],
  );

  useEffect(() => {
    dispatch({
      type: DashboardActions.SET_SELECTED_DASHBOARD_SEARCH_PATH_PREFIX,
      search_path_prefix: searchPathPrefix,
    });
  }, [dispatch, searchPathPrefix]);

  useEffect(() => {
    const isValid = modifiedSearchPath.every((c) => !!c);
    setIsValid({
      value: isValid,
      reason: !isValid ? "Search path contains empty connection" : "",
    });
  }, [modifiedSearchPath, setIsValid]);

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
    <>
      {!showEditor && (
        <div className="flex items-center space-x-3 shrink-0">
          <Icon className="h-5 w-5 shrink-0" icon="list" />
          <div className="space-x-0.5 truncate">
            {modifiedSearchPath.length > 0 &&
              modifiedSearchPath
                .map<ReactNode>((item, i) => (
                  <span key={`${item}-${i}`} className="font-medium">
                    {item}
                  </span>
                ))
                .reduce((prev, curr, idx) => [
                  prev,
                  <span key={idx} className="text-foreground-lighter">
                    ,
                  </span>,
                  curr,
                ])}
            {modifiedSearchPath.length === 0 && (
              <span className="text-foreground-lighter">
                No search path prefix set
              </span>
            )}
          </div>
          <Icon
            className="h-5 w-5 cursor-pointer shrink-0"
            icon="edit_square"
            onClick={() => setShowEditor(true)}
            title="Edit search path prefix"
          />
        </div>
      )}
      {showEditor && (
        <>
          <SearchPathEditor
            searchPath={modifiedSearchPath}
            setSearchPath={setModifiedSearchPath}
            isValid={isValid}
            onCancel={() => setShowEditor(false)}
            onSave={() => {
              setShowEditor(false);
              saveSearchPath(modifiedSearchPath);
            }}
          />
        </>
      )}
    </>
  );
};

export default SearchPathConfig;
