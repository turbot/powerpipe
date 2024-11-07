import DetectionResultNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionResultNode";
import DetectionNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionNode";
import DetectionBenchmarkNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionBenchmarkNode";
import DetectionEmptyResultNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionEmptyResultNode";
import DetectionErrorNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionErrorNode";
import DetectionKeyValuePairNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionKeyValuePairNode";
import DetectionRootNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionRootNode";
import DetectionRunningNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionRunningNode";
import useGroupingFilterConfig from "./useGroupingFilterConfig";
import useDetectionGroupingConfig from "./useDetectionGroupingConfig";
import usePrevious from "./usePrevious";
import {
  DetectionDisplayGroup,
  DetectionFilter,
  DetectionNode as DetectionNodeType,
  DetectionResult,
  DetectionResultDimension,
  DetectionSeverity,
  DetectionSummary,
  DetectionTags,
  findDimension,
} from "@powerpipe/components/dashboards/grouping/common";
import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useReducer,
} from "react";
import { default as DetectionBenchmarkType } from "@powerpipe/components/dashboards/grouping/common/DetectionBenchmark";
import {
  ElementType,
  IActions,
  PanelDefinition,
  PanelsMap,
} from "@powerpipe/types";
import { useDashboard } from "./useDashboard";
import { useDashboardControls } from "@powerpipe/components/dashboards/layout/Dashboard/DashboardControlsProvider";

type CheckGroupingActionType = ElementType<typeof checkGroupingActions>;

export type CheckGroupNodeState = {
  expanded: boolean;
};

export type CheckGroupNodeStates = {
  [name: string]: CheckGroupNodeState;
};

export type CheckGroupingAction = {
  type: CheckGroupingActionType;
  [key: string]: any;
};

export type CheckGroupFilterValues = {
  control_tag: { key: {}; value: {} };
  dimension: { key: {}; value: {} };
};

type ICheckGroupingContext = {
  benchmark: DetectionBenchmarkType | null;
  definition: PanelDefinition;
  grouping: DetectionNodeType | null;
  groupingsConfig: DetectionDisplayGroup[];
  firstChildSummaries: DetectionSummary[];
  diffFirstChildSummaries?: DetectionSummary[];
  diffGrouping: DetectionNodeType | null;
  nodeStates: CheckGroupNodeStates;
  filterValues: CheckGroupFilterValues;
  dispatch(action: CheckGroupingAction): void;
};

const GroupingActions: IActions = {
  COLLAPSE_ALL_NODES: "collapse_all_nodes",
  COLLAPSE_NODE: "collapse_node",
  EXPAND_ALL_NODES: "expand_all_nodes",
  EXPAND_NODE: "expand_node",
  UPDATE_NODES: "update_nodes",
};

const checkGroupingActions = Object.values(GroupingActions);

const GroupingContext = createContext<ICheckGroupingContext | null>(null);

const addBenchmarkTrunkNode = (
  benchmark_trunk: DetectionBenchmarkType[],
  children: DetectionNodeType[],
  benchmarkChildrenLookup: { [name: string]: DetectionNodeType[] },
  groupingKeysBeforeBenchmark: string[],
  parentGroupType: string | null,
): DetectionNodeType => {
  let newChildren: DetectionNodeType[];
  if (benchmark_trunk.length > 1) {
    newChildren = [
      addBenchmarkTrunkNode(
        benchmark_trunk.slice(1),
        children,
        benchmarkChildrenLookup,
        groupingKeysBeforeBenchmark,
        parentGroupType,
      ),
    ];
  } else {
    newChildren = children;
  }
  const currentNode = benchmark_trunk.length > 0 ? benchmark_trunk[0] : null;
  if (!!currentNode?.name) {
    const lookupKey =
      groupingKeysBeforeBenchmark.length > 0
        ? `${groupingKeysBeforeBenchmark.join("/")}/${
            currentNode?.name || "Other"
          }`
        : currentNode?.name || "Other";
    const existingChildren = benchmarkChildrenLookup[lookupKey];
    if (existingChildren) {
      // We only want to add children that are not already in the list,
      // else we end up with duplicate nodes in the tree
      for (const child of newChildren) {
        if (
          existingChildren &&
          existingChildren.find((c) => c.name === child.name)
        ) {
          continue;
        }
        existingChildren.push(child);
      }
    } else {
      benchmarkChildrenLookup[lookupKey] = newChildren;
    }
  }
  return new DetectionBenchmarkNode(
    !!parentGroupType
      ? currentNode?.title || "Other"
      : currentNode?.sort || "Other",
    currentNode?.name || "Other",
    currentNode?.title || "Other",
    newChildren,
  );
};

const getDetectionSeverityGroupingKey = (
  severity: DetectionSeverity | undefined,
): string => {
  switch (severity) {
    case "low":
      return "Low";
    case "medium":
      return "Medium";
    case "high":
      return "High";
    case "critical":
      return "Critical";
    default:
      return "Unspecified";
  }
};

const getDetectionSeveritySortKey = (
  severity: DetectionSeverity | undefined,
): string => {
  switch (severity) {
    case "critical":
      return "0";
    case "high":
      return "1";
    case "medium":
      return "2";
    case "low":
      return "3";
    default:
      return "4";
  }
};

const getCheckDimensionGroupingKey = (
  dimensionKey: string | undefined,
  dimensions: DetectionResultDimension[],
): string => {
  if (!dimensionKey) {
    return "<not set>";
  }
  const foundDimension = findDimension(dimensions, dimensionKey);
  return foundDimension ? foundDimension.value : `<not set>`;
};

function getCheckTagGroupingKey(
  tagKey: string | undefined,
  tags: DetectionTags,
) {
  if (!tagKey) {
    return "Tag key not set";
  }
  return tags[tagKey] || `<not set>`;
}

const getCheckGroupingKey = (
  checkResult: DetectionResult,
  group: DetectionDisplayGroup,
) => {
  switch (group.type) {
    case "detection_benchmark":
      if (checkResult.detection_benchmark_trunk.length <= 1) {
        return null;
      }
      return checkResult.detection_benchmark_trunk[
        checkResult.detection_benchmark_trunk.length - 1
      ].name;
    case "detection":
      return checkResult.detection.name;
    case "detection_tag":
      return getCheckTagGroupingKey(group.value, checkResult.tags);
    case "dimension":
      return getCheckDimensionGroupingKey(group.value, checkResult.dimensions);
    case "severity":
      return getDetectionSeverityGroupingKey(checkResult.detection.severity);
    default:
      return "Other";
  }
};

const getCheckGroupingNode = (
  detectionResult: DetectionResult,
  group: DetectionDisplayGroup,
  children: DetectionNodeType[],
  benchmarkChildrenLookup: { [name: string]: DetectionNodeType[] },
  groupingKeysBeforeBenchmark: string[] = [],
  parentGroupType: string | null,
): DetectionNodeType => {
  switch (group.type) {
    case "detection":
      return new DetectionNode(
        parentGroupType === "detection_benchmark"
          ? detectionResult.detection.sort
          : detectionResult.detection.title || detectionResult.detection.name,
        detectionResult.detection.name,
        detectionResult.detection.title,
        children,
      );
    case "detection_benchmark":
      return detectionResult.detection_benchmark_trunk.length > 1
        ? addBenchmarkTrunkNode(
            detectionResult.detection_benchmark_trunk.slice(1),
            children,
            benchmarkChildrenLookup,
            groupingKeysBeforeBenchmark,
            parentGroupType,
          )
        : children[0];
    case "detection_tag":
      const value = getCheckTagGroupingKey(group.value, detectionResult.tags);
      return new DetectionKeyValuePairNode(
        value,
        "control_tag",
        group.value || "Tag key not set",
        value,
        children,
      );
    case "dimension":
      const dimensionValue = getCheckDimensionGroupingKey(
        group.value,
        detectionResult.dimensions,
      );
      return new DetectionKeyValuePairNode(
        dimensionValue,
        "dimension",
        group.value || "Dimension key not set",
        dimensionValue,
        children,
      );
    case "severity":
      return new DetectionKeyValuePairNode(
        getDetectionSeveritySortKey(detectionResult.detection.severity),
        "severity",
        "severity",
        getDetectionSeverityGroupingKey(detectionResult.detection.severity),
        children,
      );
    default:
      throw new Error(`Unknown group type ${group.type}`);
  }
};

const addBenchmarkGroupingNode = (
  existingGroups: DetectionNodeType[],
  groupingNode: DetectionNodeType,
) => {
  const existingGroup = existingGroups.find(
    (existingGroup) => existingGroup.name === groupingNode.name,
  );
  if (existingGroup) {
    (existingGroup as DetectionBenchmarkNode).merge(groupingNode);
  } else {
    existingGroups.push(groupingNode);
  }
};

function getBenchmarkChildrenLookupKey(
  groupingHierarchyKeys: string[],
  groupKey: string,
) {
  const groupingKeysBeforeBenchmark = groupingHierarchyKeys.slice(
    0,
    groupingHierarchyKeys.indexOf("detection_benchmark"),
  );
  const benchmarkChildrenLookupKey =
    groupingKeysBeforeBenchmark.length > 0
      ? `${groupingKeysBeforeBenchmark.join("/")}/${groupKey}`
      : groupKey;
  return { groupingKeysBeforeBenchmark, benchmarkChildrenLookupKey };
}

const groupCheckItems = (
  temp: { _: DetectionNodeType[] },
  checkResult: DetectionResult,
  groupingsConfig: DetectionDisplayGroup[],
  DetectionNodeStates: CheckGroupNodeStates,
  benchmarkChildrenLookup: { [name: string]: DetectionNodeType[] },
  groupingHierarchyKeys: string[],
) => {
  return groupingsConfig
    .filter((groupConfig) => groupConfig.type !== "result")
    .reduce(
      (
        cumulativeGrouping,
        currentGroupingConfig,
        currentIndex,
        filteredGroups,
      ) => {
        // We want to capture the parent group type to use later for sorting purposes.
        // If we're trying to decide how to sort a control node, we need to know if
        // we're under a benchmark or some other grouping type. If we're under a benchmark,
        // we'll sort by the order determined by the benchmark, else we'll sort by title
        const parentGroupType =
          currentIndex > 0 ? filteredGroups[currentIndex - 1].type : null;
        // Get this items grouping key - e.g. control or benchmark name
        const groupKey = getCheckGroupingKey(
          checkResult,
          currentGroupingConfig,
        );

        if (!groupKey) {
          return cumulativeGrouping;
        }

        groupingHierarchyKeys.push(groupKey);

        // Collapse all benchmark trunk nodes
        if (currentGroupingConfig.type === "detection_benchmark") {
          checkResult.detection_benchmark_trunk.forEach(
            (benchmark) =>
              (DetectionNodeStates[benchmark.name] = {
                expanded: false,
              }),
          );
        } else {
          DetectionNodeStates[groupKey] = {
            expanded: false,
          };
        }

        const { groupingKeysBeforeBenchmark, benchmarkChildrenLookupKey } =
          getBenchmarkChildrenLookupKey(groupingHierarchyKeys, groupKey);

        if (!cumulativeGrouping[groupKey]) {
          cumulativeGrouping[groupKey] = { _: [] };

          const groupingNode = getCheckGroupingNode(
            checkResult,
            currentGroupingConfig,
            cumulativeGrouping[groupKey]._,
            benchmarkChildrenLookup,
            groupingKeysBeforeBenchmark,
            parentGroupType,
          );

          if (groupingNode) {
            if (currentGroupingConfig.type === "detection_benchmark") {
              // For benchmarks, we need to get the benchmark nodes including the trunk
              addBenchmarkGroupingNode(cumulativeGrouping._, groupingNode);
            } else {
              cumulativeGrouping._.push(groupingNode);
            }
          }
        }

        // If the grouping key for this has already been logged by another result,
        // use the existing children from that - this covers cases where we may have
        // benchmark 1 -> benchmark 2 -> control 1
        // benchmark 1 -> control 2
        // ...when we build the benchmark grouping node for control 1, its key will be
        // for benchmark 2, but we'll add a hierarchical grouping node for benchmark 1 -> benchmark 2
        // When we come to get the benchmark grouping node for control 2, we'll need to add
        // the control to the existing children of benchmark 1
        if (
          currentGroupingConfig.type === "detection_benchmark" &&
          benchmarkChildrenLookup[benchmarkChildrenLookupKey]
        ) {
          const groupingEntry = cumulativeGrouping[groupKey];
          const { _, ...rest } = groupingEntry || {};
          cumulativeGrouping[groupKey] = {
            _: benchmarkChildrenLookup[benchmarkChildrenLookupKey],
            ...rest,
          };
        }

        return cumulativeGrouping[groupKey];
      },
      temp,
    );
};

const getDetectionResultNode = (detectionResult: DetectionResult) => {
  if (detectionResult.type === "loading") {
    return new DetectionRunningNode(detectionResult);
  } else if (detectionResult.type === "error") {
    return new DetectionErrorNode(detectionResult);
  } else if (detectionResult.type === "empty") {
    return new DetectionEmptyResultNode(detectionResult);
  }
  return new DetectionResultNode(detectionResult);
};

const reducer = (state: CheckGroupNodeStates, action) => {
  switch (action.type) {
    case GroupingActions.COLLAPSE_ALL_NODES: {
      const newNodes = {};
      for (const [name, node] of Object.entries(state)) {
        newNodes[name] = {
          ...node,
          expanded: false,
        };
      }
      return {
        ...state,
        nodes: newNodes,
      };
    }
    case GroupingActions.COLLAPSE_NODE:
      return {
        ...state,
        [action.name]: {
          ...(state[action.name] || {}),
          expanded: false,
        },
      };
    case GroupingActions.EXPAND_ALL_NODES: {
      const newNodes = {};
      Object.entries(state).forEach(([name, node]) => {
        newNodes[name] = {
          ...node,
          expanded: true,
        };
      });
      return newNodes;
    }
    case GroupingActions.EXPAND_NODE: {
      return {
        ...state,
        [action.name]: {
          ...(state[action.name] || {}),
          expanded: true,
        },
      };
    }
    case GroupingActions.UPDATE_NODES:
      return action.nodes;
    default:
      return state;
  }
};

type DetectionGroupingProviderProps = {
  children: null | JSX.Element | JSX.Element[];
  definition: PanelDefinition;
  diff_panels: PanelsMap | undefined;
};

function recordFilterValues(
  filterValues: {
    detection_benchmark: { value: {} };
    detection: { value: {} };
    detection_tag: { key: {}; value: {} };
    dimension: { key: {}; value: {} };
    severity: { value: {} };
    status: {
      total: number;
    };
  },
  detectionResult: DetectionResult,
) {
  // Record the benchmark of this check result to allow assisted filtering later
  if (
    !!detectionResult.detection_benchmark_trunk &&
    detectionResult.detection_benchmark_trunk.length > 0
  ) {
    for (const benchmark of detectionResult.detection_benchmark_trunk) {
      filterValues.detection_benchmark.value[benchmark.name] = filterValues
        .detection_benchmark.value[benchmark.name] || {
        title: benchmark.title,
        count: 0,
      };
      filterValues.detection_benchmark.value[benchmark.name].count += 1;
    }
  }

  // Record the control of this check result to allow assisted filtering later
  filterValues.detection.value[detectionResult.detection.name] = filterValues
    .detection.value[detectionResult.detection.name] || {
    title: detectionResult.detection.title,
    count: 0,
  };
  filterValues.detection.value[detectionResult.detection.name].count += 1;

  // Record the severity of this check result to allow assisted filtering later
  if (detectionResult.severity) {
    filterValues.severity.value[detectionResult.severity.toString()] =
      filterValues.severity.value[detectionResult.severity.toString()] || 0;
    filterValues.severity.value[detectionResult.severity.toString()] += 1;
  }

  // Record the status of this check result to allow assisted filtering later
  filterValues.status[detectionResult.status] =
    filterValues.status[detectionResult.status] || 0;
  filterValues.status[detectionResult.status] += 1;

  // Record the dimension keys/values + value/key counts of this check result to allow assisted filtering later
  // for (const dimension of checkResult.dimensions) {
  //   if (!(dimension.key in filterValues.dimension.key)) {
  //     filterValues.dimension.key[dimension.key] = {
  //       [dimension.value]: 0,
  //     };
  //   }
  //   if (!(dimension.value in filterValues.dimension.key[dimension.key])) {
  //     filterValues.dimension.key[dimension.key][dimension.value] = 0;
  //   }
  //   filterValues.dimension.key[dimension.key][dimension.value] += 1;
  //
  //   if (!(dimension.value in filterValues.dimension.value)) {
  //     filterValues.dimension.value[dimension.value] = {
  //       [dimension.key]: 0,
  //     };
  //   }
  //   if (!(dimension.key in filterValues.dimension.value[dimension.value])) {
  //     filterValues.dimension.value[dimension.value][dimension.key] = 0;
  //   }
  //   filterValues.dimension.value[dimension.value][dimension.key] += 1;
  // }

  // Record the dimension keys/values + value/key counts of this check result to allow assisted filtering later
  for (const [tagKey, tagValue] of Object.entries(detectionResult.tags || {})) {
    if (!(tagKey in filterValues.detection_tag.key)) {
      filterValues.detection_tag.key[tagKey] = {
        [tagValue]: 0,
      };
    }
    if (!(tagValue in filterValues.detection_tag.key[tagKey])) {
      filterValues.detection_tag.key[tagKey][tagValue] = 0;
    }
    filterValues.detection_tag.key[tagKey][tagValue] += 1;

    if (!(tagValue in filterValues.detection_tag.value)) {
      filterValues.detection_tag.value[tagValue] = {
        [tagKey]: 0,
      };
    }
    if (!(tagKey in filterValues.detection_tag.value[tagValue])) {
      filterValues.detection_tag.value[tagValue][tagKey] = 0;
    }
    filterValues.detection_tag.value[tagValue][tagKey] += 1;
  }
}

const escapeRegex = (string) => {
  if (!string) {
    return string;
  }
  return string.replace(/[-/\\^$*+?.()|[\]{}]/g, "\\$&");
};

const wildcardToRegex = (wildcard: string) => {
  const escaped = escapeRegex(wildcard);
  return escaped.replaceAll("\\*", ".*");
};

const includeResult = (
  checkResult: DetectionResult,
  checkFilterConfig: DetectionFilter,
): boolean => {
  if (
    !checkFilterConfig ||
    !checkFilterConfig.expressions ||
    checkFilterConfig.expressions.length === 0
  ) {
    return true;
  }
  let matches: boolean[] = [];
  for (const filter of checkFilterConfig.expressions) {
    if (!filter.type) {
      continue;
    }

    // @ts-ignore
    const valueRegex = new RegExp(`^${wildcardToRegex(filter.value)}$`);

    switch (filter.type) {
      case "detection_benchmark": {
        let matchesTrunk = false;
        for (const benchmark of checkResult.detection_benchmark_trunk || []) {
          const match = valueRegex.test(benchmark.name);
          if (match) {
            matchesTrunk = true;
            break;
          }
        }
        matches.push(matchesTrunk);
        break;
      }
      case "detection": {
        matches.push(valueRegex.test(checkResult.detection.name));
        break;
      }
      case "severity": {
        matches.push(valueRegex.test(checkResult.severity || ""));
        break;
      }
      case "status": {
        matches.push(valueRegex.test(checkResult.status.toString()));
        break;
      }
      case "dimension": {
        // @ts-ignore
        const keyRegex = new RegExp(`^${wildcardToRegex(filter.key)}$`);
        let matchesDimensions = false;
        for (const dimension of checkResult.dimensions || []) {
          if (
            keyRegex.test(dimension.key) &&
            valueRegex.test(dimension.value)
          ) {
            matchesDimensions = true;
            break;
          }
        }
        matches.push(matchesDimensions);
        break;
      }
      case "control_tag": {
        // @ts-ignore
        const keyRegex = new RegExp(`^${wildcardToRegex(filter.key)}$`);
        let matchesTags = false;
        for (const [tagKey, tagValue] of Object.entries(
          checkResult.tags || {},
        )) {
          if (keyRegex.test(tagKey) && valueRegex.test(tagValue)) {
            matchesTags = true;
            break;
          }
        }
        matches.push(matchesTags);
        break;
      }
      default:
        matches.push(true);
    }
  }
  return matches.every((m) => m);
};

const useGroupingInternal = (
  definition: PanelDefinition | null,
  panelsMap: PanelsMap | undefined,
  groupingsConfig: DetectionDisplayGroup[],
  skip = false,
) => {
  const checkFilterConfig = useGroupingFilterConfig();

  return useMemo(() => {
    const filterValues = {
      detection_benchmark: { value: {} },
      detection: { value: {} },
      detection_tag: { key: {}, value: {} },
      dimension: { key: {}, value: {} },
      severity: { value: {} },
      status: { total: 0 },
    };

    if (!definition || skip || !panelsMap) {
      return [null, null, null, [], {}, filterValues];
    }

    // @ts-ignore
    const nestedDetectionBenchmarks = definition.children?.filter(
      (child) => child.panel_type === "detection_benchmark",
    );
    const nestedDetections =
      definition.panel_type === "detection"
        ? [definition]
        : // @ts-ignore
          definition.children?.filter(
            (child) => child.panel_type === "detection",
          );

    const rootBenchmarkPanel = panelsMap[definition.name];
    const b = new DetectionBenchmarkType(
      "0",
      rootBenchmarkPanel.name,
      rootBenchmarkPanel.title,
      rootBenchmarkPanel.description,
      nestedDetectionBenchmarks,
      nestedDetections,
      panelsMap,
      [],
    );

    const DetectionNodeStates: CheckGroupNodeStates = {};
    const result: DetectionNodeType[] = [];
    const temp = { _: result };
    const benchmarkChildrenLookup = {};

    // We'll loop over each control result and build up the grouped nodes from there
    b.all_detection_results.forEach((checkResult) => {
      // Record values pre-filter so we can expand out from filtered states with all values later on
      recordFilterValues(filterValues, checkResult);

      // See if the result needs to be filtered
      if (!includeResult(checkResult, checkFilterConfig)) {
        return;
      }

      // Build a grouping node - this will be the leaf node down from the root group
      // e.g. benchmark -> control (where control is the leaf)
      const grouping = groupCheckItems(
        temp,
        checkResult,
        groupingsConfig,
        DetectionNodeStates,
        benchmarkChildrenLookup,
        [],
      );
      // Build and add a check result node to the children of the trailing group.
      // This will be used to calculate totals and severity, amongst other things.
      const node = getDetectionResultNode(checkResult);
      grouping._.push(node);
    });

    const results = new DetectionRootNode(result);

    const firstChildSummaries: DetectionSummary[] = [];
    for (const child of results.children) {
      firstChildSummaries.push(child.summary);
    }

    return [
      b,
      { ...rootBenchmarkPanel, children: definition.children },
      results,
      firstChildSummaries,
      DetectionNodeStates,
      filterValues,
    ] as const;
  }, [checkFilterConfig, definition, groupingsConfig, panelsMap, skip]);
};

const GroupingProvider = ({
  children,
  definition,
  diff_panels,
}: DetectionGroupingProviderProps) => {
  const { panelsMap } = useDashboard();
  const { setContext: setDashboardControlsContext } = useDashboardControls();
  const [nodeStates, dispatch] = useReducer(reducer, { nodes: {} });
  const groupingsConfig = useDetectionGroupingConfig();

  const [
    benchmark,
    panelDefinition,
    grouping,
    firstChildSummaries,
    tempNodeStates,
    filterValues,
  ] = useGroupingInternal(definition, panelsMap, groupingsConfig);

  const [, , diffGrouping, diffFirstChildSummaries] = useGroupingInternal(
    definition,
    diff_panels,
    groupingsConfig,
    !diff_panels,
  );

  const previousGroupings = usePrevious({ groupingsConfig });

  useEffect(() => {
    if (
      previousGroupings &&
      // @ts-ignore
      previousGroupings.groupingsConfig === groupingsConfig
    ) {
      return;
    }
    dispatch({
      type: GroupingActions.UPDATE_NODES,
      nodes: tempNodeStates,
    });
  }, [previousGroupings, groupingsConfig, tempNodeStates]);

  useEffect(() => {
    setDashboardControlsContext(filterValues);
  }, [filterValues, setDashboardControlsContext]);

  return (
    <GroupingContext.Provider
      value={{
        benchmark,
        // @ts-ignore
        definition: panelDefinition,
        dispatch,
        firstChildSummaries,
        diffFirstChildSummaries,
        diffGrouping,
        grouping,
        groupingsConfig,
        nodeStates,
        filterValues,
      }}
    >
      {children}
    </GroupingContext.Provider>
  );
};

const useDetectionGrouping = () => {
  const context = useContext(GroupingContext);
  if (context === undefined) {
    throw new Error("useCheckGrouping must be used within a GroupingContext");
  }
  return context as ICheckGroupingContext;
};

export {
  GroupingActions,
  GroupingContext,
  GroupingProvider,
  useDetectionGrouping,
};

// https://stackoverflow.com/questions/50737098/multi-level-grouping-in-javascript
// keys = ['level1', 'level2'],
//     result = [],
//     temp = { _: result };
//
// data.forEach(function (a) {
//   keys.reduce(function (r, k) {
//     if (!r[a[k]]) {
//       r[a[k]] = { _: [] };
//       r._.push({ [k]: a[k], [k + 'list']: r[a[k]]._ });
//     }
//     return r[a[k]];
//   }, temp)._.push({ Id: a.Id });
// });
