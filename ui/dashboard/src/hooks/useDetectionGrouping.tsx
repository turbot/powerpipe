import DetectionResultNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionResultNode";
import DetectionNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionNode";
import DetectionBenchmarkNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionBenchmarkNode";
import DetectionEmptyResultNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionEmptyResultNode";
import DetectionErrorNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionErrorNode";
import DetectionKeyValuePairNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionKeyValuePairNode";
import DetectionRootNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionRootNode";
import DetectionRunningNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionRunningNode";
import isObject from "lodash/isObject";
import useFilterConfig from "./useFilterConfig";
import useGroupingConfig from "@powerpipe/hooks/useGroupingConfig";
import usePrevious from "./usePrevious";
import {
  applyFilter,
  DetectionDisplayGroup,
  DetectionNode as DetectionNodeType,
  DetectionResult,
  DetectionResultDimension,
  DetectionSeverity,
  DetectionSummary,
  DetectionTags,
  Filter,
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
import { KeyValuePairs } from "@powerpipe/components/dashboards/common/types";
import {
  LeafNodeDataColumn,
  LeafNodeDataRow,
} from "@powerpipe/components/dashboards/common";
import { useDashboardState } from "./useDashboardState";
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
  detection_tag: { key: {}; value: {} };
  dimension: { key: {}; value: {} };
};

type ICheckGroupingContext = {
  benchmark: DetectionBenchmarkType | {} | null;
  definition: PanelDefinition;
  grouping: DetectionNodeType | null;
  groupingConfig: DetectionDisplayGroup[];
  firstChildSummaries: DetectionSummary[];
  hasSeverityResults: boolean;
  nodeStates: CheckGroupNodeStates;
  filterValues: CheckGroupFilterValues;
  dispatch(action: CheckGroupingAction): void;
};

const GroupingActions: IActions = {
  COLLAPSE_NODE: "collapse_node",
  EXPAND_ALL_NODES: "expand_all_nodes",
  EXPAND_NODE: "expand_node",
  UPDATE_NODES: "update_nodes",
};

const checkGroupingActions = Object.values(GroupingActions);

const GroupingContext = createContext<ICheckGroupingContext | null>({
  benchmark: null,
  definition: null,
  grouping: null,
  groupingConfig: [],
  firstChildSummaries: [],
  hasSeverityResults: false,
  nodeStates: {},
  filterValues: {
    detection_tag: { key: {}, value: {} },
    dimension: { key: {}, value: {} },
  },
  dispatch: () => {},
});

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
    currentNode?.documentation,
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
    case "benchmark":
      if (checkResult.benchmark_trunk.length <= 1) {
        return null;
      }
      return checkResult.benchmark_trunk[checkResult.benchmark_trunk.length - 1]
        .name;
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

const getDetectionGroupingNode = (
  detectionResult: DetectionResult,
  group: DetectionDisplayGroup,
  children: DetectionNodeType[],
  benchmarkChildrenLookup: { [name: string]: DetectionNodeType[] },
  groupingKeysBeforeBenchmark: string[] = [],
  parentGroupType: string | null,
): DetectionNodeType => {
  switch (group.type) {
    case "benchmark":
      return detectionResult.benchmark_trunk.length > 1
        ? addBenchmarkTrunkNode(
            detectionResult.benchmark_trunk.slice(1),
            children,
            benchmarkChildrenLookup,
            groupingKeysBeforeBenchmark,
            parentGroupType,
          )
        : children[0];
    case "detection":
      return new DetectionNode(
        parentGroupType === "benchmark"
          ? detectionResult.detection.sort
          : detectionResult.detection.title || detectionResult.detection.name,
        detectionResult.detection.name,
        detectionResult.detection.title,
        detectionResult.detection.documentation,
        detectionResult.detection._results?.[0],
        children,
      );
    case "detection_tag":
      const value = getCheckTagGroupingKey(group.value, detectionResult.tags);
      return new DetectionKeyValuePairNode(
        value,
        "detection_tag",
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
    groupingHierarchyKeys.indexOf("benchmark"),
  );
  const benchmarkChildrenLookupKey =
    groupingKeysBeforeBenchmark.length > 0
      ? `${groupingKeysBeforeBenchmark.join("/")}/${groupKey}`
      : groupKey;
  return { groupingKeysBeforeBenchmark, benchmarkChildrenLookupKey };
}

const groupDetectionItems = (
  temp: { _: DetectionNodeType[] },
  checkResult: DetectionResult,
  groupingConfig: DetectionDisplayGroup[],
  detectionNodeStates: CheckGroupNodeStates,
  benchmarkChildrenLookup: { [name: string]: DetectionNodeType[] },
  groupingHierarchyKeys: string[],
) => {
  return groupingConfig
    .filter((groupConfig) => groupConfig.type !== "result")
    .reduce(
      (
        cumulativeGrouping,
        currentGroupingConfig,
        currentIndex,
        filteredGroups,
      ) => {
        // We want to capture the parent group type to use later for sorting purposes.
        // If we're trying to decide how to sort a detection node, we need to know if
        // we're under a benchmark or some other grouping type. If we're under a benchmark,
        // we'll sort by the order determined by the benchmark, else we'll sort by title
        const parentGroupType =
          currentIndex > 0 ? filteredGroups[currentIndex - 1].type : null;
        // Get this items grouping key - e.g. detection or benchmark name
        const groupKey = getCheckGroupingKey(
          checkResult,
          currentGroupingConfig,
        );

        if (!groupKey) {
          return cumulativeGrouping;
        }

        groupingHierarchyKeys.push(groupKey);

        // Collapse all benchmark trunk nodes
        if (currentGroupingConfig.type === "benchmark") {
          checkResult.benchmark_trunk.forEach(
            (benchmark) =>
              (detectionNodeStates[benchmark.name] = {
                expanded: false,
              }),
          );
        } else {
          detectionNodeStates[groupKey] = {
            expanded: false,
          };
        }

        const { groupingKeysBeforeBenchmark, benchmarkChildrenLookupKey } =
          getBenchmarkChildrenLookupKey(groupingHierarchyKeys, groupKey);

        if (!cumulativeGrouping[groupKey]) {
          cumulativeGrouping[groupKey] = { _: [] };

          const groupingNode = getDetectionGroupingNode(
            checkResult,
            currentGroupingConfig,
            cumulativeGrouping[groupKey]._,
            benchmarkChildrenLookup,
            groupingKeysBeforeBenchmark,
            parentGroupType,
          );

          if (groupingNode) {
            if (currentGroupingConfig.type === "benchmark") {
              // For benchmarks, we need to get the benchmark nodes including the trunk
              addBenchmarkGroupingNode(cumulativeGrouping._, groupingNode);
            } else {
              cumulativeGrouping._.push(groupingNode);
            }
          }
        }

        // If the grouping key for this has already been logged by another result,
        // use the existing children from that - this covers cases where we may have
        // benchmark 1 -> benchmark 2 -> detection 1
        // benchmark 1 -> detection 2
        // ...when we build the benchmark grouping node for detection 1, its key will be
        // for benchmark 2, but we'll add a hierarchical grouping node for benchmark 1 -> benchmark 2
        // When we come to get the benchmark grouping node for detection 2, we'll need to add
        // the detection to the existing children of benchmark 1
        if (
          currentGroupingConfig.type === "benchmark" &&
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
  benchmarkChildren: PanelDefinition[] | undefined;
};

function recordFilterValues(
  filterValues: {
    benchmark: { value: {} };
    detection: { value: {} };
    detection_tag: { key: {}; value: {} };
    dimension: { key: {}; value: {} };
    severity: { value: {} };
  },
  detectionResult: DetectionResult,
) {
  // Record the benchmark of this check result to allow assisted filtering later
  if (
    !!detectionResult.benchmark_trunk &&
    detectionResult.benchmark_trunk.length > 0
  ) {
    for (const benchmark of detectionResult.benchmark_trunk) {
      filterValues.benchmark.value[benchmark.name] = filterValues.benchmark
        .value[benchmark.name] || {
        title: benchmark.title,
        count: 0,
      };
      filterValues.benchmark.value[benchmark.name].count +=
        detectionResult.rows?.length || 0;
    }
  }

  // Record the detection of this detection result to allow assisted filtering later
  filterValues.detection.value[detectionResult.detection.name] = filterValues
    .detection.value[detectionResult.detection.name] || {
    title: detectionResult.detection.title,
    count: 0,
  };
  filterValues.detection.value[detectionResult.detection.name].count +=
    detectionResult.rows?.length || 0;

  // Record the severity of this check result to allow assisted filtering later
  if (detectionResult.severity && detectionResult.rows?.length) {
    filterValues.severity.value[detectionResult.severity.toString()] =
      filterValues.severity.value[detectionResult.severity.toString()] || 0;
    filterValues.severity.value[detectionResult.severity.toString()] +=
      detectionResult.rows.length;
  }

  // Record the dimension keys/values + value/key counts of this check result to allow assisted filtering later
  for (const dimension of detectionResult.dimensions) {
    if (!(dimension.key in filterValues.dimension.key)) {
      filterValues.dimension.key[dimension.key] = {
        [dimension.value]: 0,
      };
    }
    if (!(dimension.value in filterValues.dimension.key[dimension.key])) {
      filterValues.dimension.key[dimension.key][dimension.value] = 0;
    }
    filterValues.dimension.key[dimension.key][dimension.value] +=
      detectionResult.rows?.length || 0;

    if (!(dimension.value in filterValues.dimension.value)) {
      filterValues.dimension.value[dimension.value] = {
        [dimension.key]: 0,
      };
    }
    if (!(dimension.key in filterValues.dimension.value[dimension.value])) {
      filterValues.dimension.value[dimension.value][dimension.key] = 0;
    }
    filterValues.dimension.value[dimension.value][dimension.key] +=
      detectionResult.rows?.length || 0;
  }

  // Record the dimension keys/values + value/key counts of this detection result to allow assisted filtering later
  for (const [tagKey, tagValue] of Object.entries(detectionResult.tags || {})) {
    if (!(tagKey in filterValues.detection_tag.key)) {
      filterValues.detection_tag.key[tagKey] = {
        [tagValue]: 0,
      };
    }
    if (!(tagValue in filterValues.detection_tag.key[tagKey])) {
      filterValues.detection_tag.key[tagKey][tagValue] = 0;
    }
    filterValues.detection_tag.key[tagKey][tagValue] +=
      detectionResult.rows?.length || 0;

    if (!(tagValue in filterValues.detection_tag.value)) {
      filterValues.detection_tag.value[tagValue] = {
        [tagKey]: 0,
      };
    }
    if (!(tagKey in filterValues.detection_tag.value[tagValue])) {
      filterValues.detection_tag.value[tagValue][tagKey] = 0;
    }
    filterValues.detection_tag.value[tagValue][tagKey] +=
      detectionResult.rows?.length || 0;
  }
}

const includeResult = (
  result: DetectionResult,
  panel: PanelDefinition,
  allFilters: KeyValuePairs<Filter>,
): boolean => {
  // If no filters, include this
  if (Object.keys(allFilters).length === 0) {
    return true;
  }

  const filterForRootPanel = allFilters[panel.name];
  const filterForDetection = allFilters[result.detection.name];

  // If no filters for the parent panel, or this panel, include
  if (!filterForRootPanel && !filterForDetection) {
    return true;
  }

  const columnLookup: KeyValuePairs<LeafNodeDataColumn> = {};

  let matches: boolean[] = [];
  for (const filter of [filterForRootPanel, filterForDetection].filter(
    (f) => !!f && !!f.expressions?.length,
  )) {
    for (const expression of filter.expressions || []) {
      if (!expression.type) {
        continue;
      }

      switch (expression.type) {
        case "benchmark": {
          let matchesTrunk = false;
          for (const benchmark of result.benchmark_trunk || []) {
            const match = applyFilter(expression, benchmark.name);
            if (match) {
              matchesTrunk = true;
              break;
            }
          }
          matches.push(matchesTrunk);
          break;
        }
        case "detection": {
          matches.push(applyFilter(expression, result.detection.name));
          break;
        }
        case "dimension": {
          let newRows: LeafNodeDataRow[] = [];
          let includeRow = false;
          let column: LeafNodeDataColumn | undefined;
          if (!(expression.key in columnLookup)) {
            column = result?.columns?.find((c) => c.name === expression.key);
            if (!column) {
              matches.push(false);
              break;
            }
            columnLookup[expression.key] = column;
          } else {
            column = columnLookup[expression.key];
          }

          for (const row of result.rows || []) {
            const rowValue = row[expression.key];
            includeRow =
              !!expression.key &&
              expression.key in row &&
              applyFilter(
                expression,
                rowValue,
                column.data_type === "jsonb" ||
                  column.data_type === "varchar[]" ||
                  isObject(rowValue),
              );
            if (includeRow) {
              newRows.push(row);
            }
          }
          result.rows = newRows;
          matches.push(true);
          break;
        }
        case "detection_tag": {
          let matchesTags = false;
          for (const [tagKey, tagValue] of Object.entries(result.tags || {})) {
            if (
              expression.key === tagKey &&
              applyFilter(expression, tagValue)
            ) {
              matchesTags = true;
              break;
            }
          }
          matches.push(matchesTags);
          break;
        }
        case "severity": {
          matches.push(applyFilter(expression, result.severity || ""));
          break;
        }
        default:
          matches.push(true);
      }
    }
  }
  return matches.every((m) => m);
};

const useGroupingInternal = (
  definition: PanelDefinition | null,
  benchmarkChildren: PanelDefinition[] | undefined,
  panelsMap: PanelsMap | undefined,
  groupingConfig: DetectionDisplayGroup[],
  skip = false,
) => {
  const { allFilters } = useFilterConfig(definition?.name);

  return useMemo(() => {
    const filterValues = {
      benchmark: { value: {} },
      detection: { value: {} },
      detection_tag: { key: {}, value: {} },
      dimension: { key: {}, value: {} },
      severity: { value: {} },
    };

    if (!definition || skip || !panelsMap) {
      return [null, null, null, [], false, {}, filterValues];
    }

    // @ts-ignore
    const nestedDetectionBenchmarks = benchmarkChildren?.filter(
      (child) => child.panel_type === "benchmark",
    );
    const nestedDetections =
      definition.panel_type === "detection"
        ? [definition]
        : // @ts-ignore
          benchmarkChildren?.filter(
            (child) => child.panel_type === "detection",
          );

    const b = new DetectionBenchmarkType(
      "0",
      definition.name,
      definition.title,
      definition.description,
      definition.documentation,
      nestedDetectionBenchmarks,
      nestedDetections,
      panelsMap,
      [],
    );

    const detectionNodeStates: CheckGroupNodeStates = {};
    const result: DetectionNodeType[] = [];
    const temp = { _: result };
    const benchmarkChildrenLookup = {};

    // We'll loop over each detection result and build up the grouped nodes from there
    b.all_detection_results.forEach((detectionResult) => {
      // Record values pre-filter so we can expand out from filtered states with all values later on
      recordFilterValues(filterValues, detectionResult);

      // See if the result needs to be filtered
      if (!includeResult(detectionResult, definition, allFilters)) {
        return;
      }

      // Build a grouping node - this will be the leaf node down from the root group
      // e.g. benchmark -> detection (where detection is the leaf)
      const grouping = groupDetectionItems(
        temp,
        detectionResult,
        groupingConfig,
        detectionNodeStates,
        benchmarkChildrenLookup,
        [],
      );
      // Build and add a check result node to the children of the trailing group.
      // This will be used to calculate totals and severity, amongst other things.
      const node = getDetectionResultNode(detectionResult);
      grouping._.push(node);
    });

    const results = new DetectionRootNode(result);

    const firstChildSummaries: DetectionSummary[] = [];
    let hasSeverityResults: Boolean = false;
    for (const child of results.children) {
      firstChildSummaries.push(child.summary);
      if (
        child.severity_summary.critical !== undefined ||
        child.severity_summary.high !== undefined ||
        child.severity_summary.medium !== undefined ||
        child.severity_summary.low !== undefined
      ) {
        hasSeverityResults = true;
      }
    }

    return [
      b,
      results,
      firstChildSummaries,
      hasSeverityResults,
      detectionNodeStates,
      filterValues,
    ] as const;
  }, [allFilters, definition, groupingConfig, panelsMap, skip]);
};

const GroupingProvider = ({
  children,
  definition,
  benchmarkChildren,
}: DetectionGroupingProviderProps) => {
  const { panelsMap } = useDashboardState();
  const { setContext: setDashboardControlsContext } = useDashboardControls();
  const [nodeStates, dispatch] = useReducer(reducer, { nodes: {} });
  const { grouping: groupingConfig } = useGroupingConfig(definition.name);

  const [
    benchmark,
    grouping,
    firstChildSummaries,
    hasSeverityResults,
    tempNodeStates,
    filterValues,
  ] = useGroupingInternal(
    definition,
    benchmarkChildren,
    panelsMap,
    groupingConfig,
  );

  const previousGroupings = usePrevious({ groupingConfig });

  useEffect(() => {
    if (
      !previousGroupings ||
      (!!previousGroupings &&
        JSON.stringify(previousGroupings.groupingConfig) ===
          JSON.stringify(groupingConfig))
    ) {
      return;
    }
    dispatch({
      type: GroupingActions.UPDATE_NODES,
      nodes: tempNodeStates,
    });
  }, [previousGroupings, groupingConfig, tempNodeStates]);

  useEffect(() => {
    setDashboardControlsContext(() => filterValues);
  }, [filterValues]);

  return (
    <GroupingContext.Provider
      value={{
        benchmark,
        definition,
        dispatch,
        firstChildSummaries,
        hasSeverityResults,
        grouping,
        groupingConfig,
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
