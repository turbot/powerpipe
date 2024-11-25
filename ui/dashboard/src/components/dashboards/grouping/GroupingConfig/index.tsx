import GroupingEditor from "../GroupingEditor";
import useCheckGroupingConfig from "@powerpipe/hooks/useCheckGroupingConfig";
// import { CheckDisplayGroup } from "../common";
import { Noop } from "@powerpipe/types/func";
import { useSearchParams } from "react-router-dom";

type CheckGroupingConfigProps = {
  onClose: Noop;
};

// type CheckGroupingTitleLabelProps = {
//   item: CheckDisplayGroup;
// };

// const CheckGroupingTitleLabel = ({ item }: CheckGroupingTitleLabelProps) => {
//   switch (item.type) {
//     case "control_tag":
//     case "dimension":
//       return (
//         <div className="space-x-1">
//           <span className="capitalize">{item.type}</span>
//           <span className="text-foreground-lighter">=</span>
//           <span className="font-medium">{item.value}</span>
//         </div>
//       );
//     default:
//       return (
//         <div>
//           <span className="capitalize font-medium">{item.type}</span>
//         </div>
//       );
//   }
// };

const GroupingConfig = ({ onClose }: CheckGroupingConfigProps) => {
  const [, setSearchParams] = useSearchParams();
  const groupingConfig = useCheckGroupingConfig();

  const saveGroupingConfig = (toSave) => {
    setSearchParams((previous) => {
      const newParams = new URLSearchParams(previous);
      if (!toSave || !toSave.length) {
        newParams.delete("grouping");
      } else {
        newParams.set(
          "grouping",
          toSave
            .map((c) =>
              c.type === "dimension" || c.type === "control_tag"
                ? `${c.type}|${c.value}`
                : c.type,
            )
            .join(","),
        );
      }
      return newParams;
    });
  };

  return (
    <GroupingEditor
      config={groupingConfig}
      onCancel={onClose}
      onApply={saveGroupingConfig}
    />
  );
};

export default GroupingConfig;
