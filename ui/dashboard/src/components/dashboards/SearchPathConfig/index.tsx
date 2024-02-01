import Icon from "components/Icon";
import SearchPathEditor from "components/dashboards/SearchPathEditor";
import useDashboardSearchPath from "hooks/useDashboardSearchPath";
import { DashboardActions } from "types";
import { ReactNode, useEffect, useState } from "react";
import { useDashboard } from "hooks/useDashboard";
import { useSearchParams } from "react-router-dom";

const SearchPathConfig = () => {
  const { dispatch } = useDashboard();
  const [showEditor, setShowEditor] = useState(false);
  const [isValid, setIsValid] = useState({ value: false, reason: "" });
  const [_, setSearchParams] = useSearchParams();
  const searchPath = useDashboardSearchPath();
  const [modifiedSearchPath, setModifiedSearchPath] =
    useState<string[]>(searchPath);

  useEffect(() => {
    dispatch({
      type: DashboardActions.SET_SELECTED_DASHBOARD_SEARCH_PATH,
      search_path: searchPath,
    });
  }, [searchPath]);

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
      newParams.set("search_path", toSave.join(","));
      return newParams;
    });
  };

  return (
    <>
      <div className="flex items-center space-x-3 shrink-0">
        <Icon className="h-5 w-5" icon="list" />
        <div className="space-x-0.5 truncate">
          {searchPath.length > 0 &&
            searchPath
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
          {searchPath.length === 0 && (
            <span className="text-foreground-lighter">No search path set</span>
          )}
        </div>
        {!showEditor && (
          <Icon
            className="h-5 w-5 cursor-pointer shrink-0"
            icon="edit_square"
            onClick={() => setShowEditor(true)}
            title="Edit search path"
          />
        )}
      </div>
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
