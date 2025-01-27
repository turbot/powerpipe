import {
  DetectionDisplayGroupType,
  DetectionDisplayGroup,
} from "@powerpipe/components/dashboards/grouping/common";
import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const groupingKeys = [
  "benchmark",
  "detection",
  "detection_tag",
  "dimension",
  "reason",
  "resource",
  "result",
  "severity",
  "status",
];

const useDetectionGroupingConfig = () => {
  const [searchParams] = useSearchParams();
  return useMemo(() => {
    const rawGrouping = searchParams.get("grouping");
    if (rawGrouping) {
      const groupings: DetectionDisplayGroup[] = [];
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
            type: typeValueParts[0] as DetectionDisplayGroupType,
            value: typeValueParts[1],
          });
        } else {
          groupings.push({
            type: typeValueParts[0] as DetectionDisplayGroupType,
          });
        }
      }
      return groupings;
    } else {
      return [
        { type: "benchmark" },
        { type: "detection" },
        { type: "result" },
      ] as DetectionDisplayGroup[];
    }
  }, [searchParams]);
};

export default useDetectionGroupingConfig;
