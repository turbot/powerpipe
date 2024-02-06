import Icon from "components/Icon";
import SearchPathEditor from "components/dashboards/SearchPathEditor";
import useDashboardSearchPath from "hooks/useDashboardSearchPath";
import { DashboardActions } from "types";
import { ReactNode, useEffect, useMemo, useState } from "react";
import { useDashboard } from "hooks/useDashboard";
import { useSearchParams } from "react-router-dom";

const SearchPathConfig = () => {
  const { dispatch, selectedDashboard, dashboardsMetadata } = useDashboard();
  const [showEditor, setShowEditor] = useState(false);
  const [isValid, setIsValid] = useState({ value: false, reason: "" });
  const [, setSearchParams] = useSearchParams();
  const searchPath = useDashboardSearchPath();

  const configuredSearchPath = useMemo(() => {
    if (!selectedDashboard || !dashboardsMetadata || !searchPath) {
      return [];
    }
    if (!!searchPath.length) {
      return searchPath;
    }
    if (!dashboardsMetadata[selectedDashboard.full_name]) {
      return [];
    }
    return (
      dashboardsMetadata[selectedDashboard.full_name].original_search_path || []
    );
  }, [dashboardsMetadata, searchPath, selectedDashboard]);

  const [modifiedSearchPath, setModifiedSearchPath] =
    useState<string[]>(configuredSearchPath);

  useEffect(
    () => setModifiedSearchPath(() => configuredSearchPath),
    [configuredSearchPath, setModifiedSearchPath],
  );

  useEffect(() => {
    dispatch({
      type: DashboardActions.SET_SELECTED_DASHBOARD_SEARCH_PATH,
      search_path: searchPath,
    });
  }, [dispatch, searchPath]);

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
        newParams.set("search_path", toSave.join(","));
      } else {
        newParams.delete("search_path");
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
                No search path set
              </span>
            )}
          </div>
          <Icon
            className="h-5 w-5 cursor-pointer shrink-0"
            icon="edit_square"
            onClick={() => setShowEditor(true)}
            title="Edit search path"
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
