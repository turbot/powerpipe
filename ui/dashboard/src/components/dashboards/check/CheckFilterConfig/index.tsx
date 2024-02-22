import CheckFilterEditor from "../CheckFilterEditor";
import Icon from "../../../Icon";
import useCheckFilterConfig from "../../../../hooks/useCheckFilterConfig";
import { CheckFilter } from "../common";
import { Fragment, ReactNode, useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";

const filtersToText = (filter: CheckFilter) => {
  if (filter.operator === "and") {
    // And filter group
    return filter.expressions?.map((item, index) => (
      <Fragment key={index}>
        {!!index && <span className="text-foreground-lighter">and</span>}
        {filtersToText(item)}
      </Fragment>
    ));
  }

  if (filter.operator === "equal") {
    // Convert filter to text
    let textParts: ReactNode[] = [];
    if (filter.key) {
      textParts.push(<span>{filter.key}</span>);
    } else {
      textParts.push(<span className="capitalize">{filter.type}</span>);
    }
    textParts.push(<span>{filter.title || filter.value}</span>);

    return (
      <span className="space-x-1">
        {textParts.map((item, index) => (
          <Fragment key={index}>
            {!!index && <span className="text-foreground-lighter">=</span>}
            {item}
          </Fragment>
        ))}
      </span>
    );
  }

  return "<unsupported>";
};

const validateFilter = (filter: CheckFilter): boolean => {
  // Each and must have at least one expression
  if (
    filter.operator === "and" &&
    (!filter.expressions || filter.expressions.length < 1)
  ) {
    return false;
  }
  if (filter.operator === "and") {
    // @ts-ignore can't reach here if filter.expressions is not truthy
    return filter.expressions.every(validateFilter);
  }

  if (filter.operator === "equal") {
    return (
      !!filter.type && (filter.key !== undefined || filter.value !== undefined)
    );
  }

  return false;
};

const CheckFilterConfig = () => {
  const [showEditor, setShowEditor] = useState(false);
  const [isValid, setIsValid] = useState({ value: false, reason: "" });
  const [, setSearchParams] = useSearchParams();
  const filterConfig = useCheckFilterConfig();
  const [modifiedConfig, setModifiedConfig] =
    useState<CheckFilter>(filterConfig);

  useEffect(() => {
    if (!modifiedConfig) {
      setIsValid({ value: true, reason: "" });
      return;
    }

    setIsValid({ value: validateFilter(modifiedConfig), reason: "" });
  }, [modifiedConfig, setIsValid]);

  const saveFilterConfig = (toSave: CheckFilter) => {
    setSearchParams((previous) => {
      const filters =
        toSave.expressions?.filter((f) => validateFilter(f)) || [];
      const newParams = new URLSearchParams(previous);
      if (filters.length === 0) {
        newParams.delete("where");
        return newParams;
      } else {
        const asJson = JSON.stringify(toSave);
        newParams.set("where", asJson);
        return newParams;
      }
    });
  };

  return (
    <>
      <div className="flex items-center space-x-3 shrink-0">
        <Icon className="h-5 w-5" icon="filter_list" />
        {filterConfig.operator === "and" &&
          !!filterConfig.expressions &&
          filterConfig.expressions.length > 0 && (
            <div className="space-x-2">{filtersToText(filterConfig)}</div>
          )}
        {filterConfig.operator === "and" &&
          (!filterConfig.expressions ||
            filterConfig.expressions.length === 0) && (
            <span className="text-foreground-lighter">No filters</span>
          )}
        {!showEditor && (
          <Icon
            className="h-5 w-5 cursor-pointer shrink-0"
            icon="edit_square"
            onClick={() => setShowEditor(true)}
            title="Edit filter"
          />
        )}
      </div>
      {showEditor && (
        <>
          <CheckFilterEditor
            config={modifiedConfig}
            setConfig={setModifiedConfig}
            isValid={isValid}
            onCancel={() => setShowEditor(false)}
            onSave={() => {
              setShowEditor(false);
              saveFilterConfig(modifiedConfig);
            }}
          />
        </>
      )}
    </>
  );
};

export default CheckFilterConfig;

export { validateFilter };
