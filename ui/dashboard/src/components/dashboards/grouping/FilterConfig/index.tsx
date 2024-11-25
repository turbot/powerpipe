import FilterEditor, { validateFilter } from "@powerpipe/components/dashboards/grouping/FilterEditor";
import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import { Filter } from "../common";
import { useSearchParams } from "react-router-dom";

// const filtersToText = (filter: Filter) => {
//   if (filter.operator === "and") {
//     // And filter group
//     return filter.expressions?.map((item, index) => (
//       <Fragment key={index}>
//         {!!index && <span className="text-foreground-lighter">and</span>}
//         {filtersToText(item)}
//       </Fragment>
//     ));
//   }
//
//   if (filter.operator === "equal") {
//     // Convert filter to text
//     let textParts: ReactNode[] = [];
//     if (filter.key) {
//       textParts.push(<span>{filter.key}</span>);
//     } else {
//       textParts.push(<span className="capitalize">{filter.type}</span>);
//     }
//     textParts.push(<span>{filter.title || filter.value}</span>);
//
//     return (
//       <span className="space-x-1">
//         {textParts.map((item, index) => (
//           <Fragment key={index}>
//             {!!index && <span className="text-foreground-lighter">=</span>}
//             {item}
//           </Fragment>
//         ))}
//       </span>
//     );
//   }
//
//   return "<unsupported>";
// };

const FilterConfig = () => {
  const [, setSearchParams] = useSearchParams();
  const filterConfig = useFilterConfig();

  const saveFilterConfig = (toSave: Filter) => {
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

  return <FilterEditor filter={filterConfig} onApply={saveFilterConfig} />;
};

export default FilterConfig;
