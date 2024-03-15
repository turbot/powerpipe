import CheckFilterEditor, { validateFilter } from "../CheckFilterEditor";
import useCheckFilterConfig from "@powerpipe/hooks/useCheckFilterConfig";
import { CheckFilter } from "../common";
import { Fragment, ReactNode } from "react";
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

const CheckFilterConfig = () => {
  const [, setSearchParams] = useSearchParams();
  const [filterConfig] = useCheckFilterConfig();

  const saveFilterConfig = (toSave: CheckFilter) => {
    setSearchParams((previous) => {
      const newParams = new URLSearchParams(previous);
      if (!validateFilter(toSave)) {
        newParams.delete("where");
        return newParams;
      } else {
        const asJson = JSON.stringify(toSave);
        newParams.set("where", asJson);
        return newParams;
      }
    });
  };

  return <CheckFilterEditor filter={filterConfig} onApply={saveFilterConfig} />;
};

export default CheckFilterConfig;
