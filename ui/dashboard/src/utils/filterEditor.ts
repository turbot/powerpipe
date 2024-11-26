const filterKeysSorter = (a, b) => {
  const aIsGrouped = Array.isArray(a?.options);
  const bIsGrouped = Array.isArray(b?.options);

  if (aIsGrouped && !bIsGrouped) {
    return 1; // a should come after b
  } else if (!aIsGrouped && bIsGrouped) {
    return -1; // a should come before b
  } else {
    return a?.label?.localeCompare(b?.label); // Alphabetical sort if both have colons or neither has
  }
};

const filterTypeMap = {
  benchmark: "Benchmark",
  detection_benchmark: "Benchmark",
  control: "Control",
  control_tag: "Control Tag",
  detection: "Detection",
  detection_tag: "Detection Tag",
  dimension: "Dimension",
  reason: "Reason",
  resource: "Resource",
  severity: "Severity",
  status: "Status",
};

export { filterKeysSorter, filterTypeMap };
