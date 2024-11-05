import {
  CheckDisplayGroup,
  CheckDisplayGroupType,
} from "@powerpipe/components/dashboards/grouping/common";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const groupingKeys = [
  "benchmark",
  "control",
  "control_tag",
  "dimension",
  "reason",
  "resource",
  "result",
  "severity",
  "status",
];

const useCheckGroupingConfig = () => {
  const [searchParams] = useSearchParams();
  return useMemo(() => {
    const rawGrouping = searchParams.get("grouping");
    if (rawGrouping) {
      const groupings: CheckDisplayGroup[] = [];
      const groupingParts = rawGrouping.split(",").filter((g) => !!g);
      for (const groupingPart of groupingParts) {
        const typeValueParts = groupingPart.split("|");
        const groupingKey = typeValueParts[0];

        // Is this a valid grouping key?
        const isValid = groupingKeys.includes(groupingKey);
        if (!isValid) {
          throw new Error(`Unsupported grouping key ${groupingKey}`);
        }

        if (typeValueParts.length > 1) {
          groupings.push({
            type: typeValueParts[0] as CheckDisplayGroupType,
            value: typeValueParts[1],
          });
        } else {
          groupings.push({
            type: typeValueParts[0] as CheckDisplayGroupType,
          });
        }
      }
      return groupings;
    } else {
      return [
        // { type: "status" },
        // { type: "reason" },
        // { type: "resource" },
        // { type: "severity" },
        // { type: "dimension", value: "account_id" },
        // { type: "dimension", value: "region" },
        // { type: "control_tag", value: "service" },
        // { type: "control_tag", value: "cis_type" },
        // { type: "control_tag", value: "cis_level" },
        { type: "benchmark" },
        { type: "control" },
        { type: "result" },
      ] as CheckDisplayGroup[];
    }
  }, [searchParams]);
};

export default useCheckGroupingConfig;
