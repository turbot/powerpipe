import BenchmarkNode from "@powerpipe/components/dashboards/grouping/common/node/BenchmarkNode";
import ControlEmptyResultNode from "@powerpipe/components/dashboards/grouping/common/node/ControlEmptyResultNode";
import ControlErrorNode from "@powerpipe/components/dashboards/grouping/common/node/ControlErrorNode";
import ControlNode from "@powerpipe/components/dashboards/grouping/common/node/ControlNode";
import ControlResultNode from "@powerpipe/components/dashboards/grouping/common/node/ControlResultNode";
import ControlRunningNode from "@powerpipe/components/dashboards/grouping/common/node/ControlRunningNode";
import KeyValuePairNode from "@powerpipe/components/dashboards/grouping/common/node/KeyValuePairNode";
import RootNode from "@powerpipe/components/dashboards/grouping/common/node/RootNode";
import useFilterConfig from "./useFilterConfig";
import usePrevious from "./usePrevious";
import {
  applyFilter,
  CheckDisplayGroup,
  CheckNode,
  CheckResult,
  CheckResultDimension,
  CheckResultStatus,
  CheckSeverity,
  CheckSummary,
  CheckTags,
  DisplayGroup,
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
import { default as BenchmarkType } from "@powerpipe/components/dashboards/grouping/common/Benchmark";
import {
  ElementType,
  IActions,
  PanelDefinition,
  PanelsMap,
} from "@powerpipe/types";

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

// Update CheckGroupFilterStatusValuesMap to include new states
type CheckGroupFilterStatusValuesMap = {
  [key in keyof typeof CheckResultStatus]: number;
} & {
  invalid: number;
  muted: number;
  tbd: number;
};

export type CheckGroupFilterValues = {
  status: CheckGroupFilterStatusValuesMap;
  control_tag: { key: {}; value: {} };
  dimension: { key: {}; value: {} };
  reason: { value: {} };
  resource: { value: {} };
  control: { value: {} };
  severity: { value: {} };
  benchmark: { value: {} };
};

// Extend CheckResultStatus type to include new states
export type ExtendedCheckResultStatus =
  | CheckResultStatus
  | "invalid"
  | "muted"
  | "tbd";

type ICheckGroupingContext = {
  benchmark: BenchmarkType | null;
  definition: PanelDefinition;
  grouping: CheckNode | null;
  groupingConfig: CheckDisplayGroup[];
  firstChildSummaries: CheckSummary[];
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

const GroupingContext = createContext<ICheckGroupingContext | null>(null);

const addBenchmarkTrunkNode = (
  benchmark_trunk: BenchmarkType[],
  children: CheckNode[],
  benchmarkChildrenLookup: { [name: string]: CheckNode[] },
  groupingKeysBeforeBenchmark: string[],
  parentGroupType: string | null,
): CheckNode => {
  let newChildren: CheckNode[];
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
  return new BenchmarkNode(
    !!parentGroupType
      ? currentNode?.title || "Other"
      : currentNode?.sort || "Other",
    currentNode?.name || "Other",
    currentNode?.title || "Other",
    currentNode?.documentation,
    newChildren,
  );
};

// Update function signatures to use ExtendedCheckResultStatus
const getCheckStatusGroupingKey = (
  status: ExtendedCheckResultStatus,
): string => {
  switch (status) {
    case CheckResultStatus.alarm:
      return "Alarm";
    case CheckResultStatus.invalid:
      return "Invalid";
    case CheckResultStatus.error:
      return "Error";
    case CheckResultStatus.ok:
      return "OK";
    case CheckResultStatus.info:
      return "Info";
    case CheckResultStatus.muted:
      return "Muted";
    case CheckResultStatus.skip:
    case CheckResultStatus.skipped:
      return "Skipped";
    case CheckResultStatus.empty:
      return "Empty";
    case CheckResultStatus.tbd:
      return "TBD";
    default:
      return status || "Unknown";
  }
};

const getCheckStatusSortKey = (status: ExtendedCheckResultStatus): string => {
  switch (status) {
    case CheckResultStatus.alarm:
      return "1";
    case CheckResultStatus.error:
      return "2";
    case CheckResultStatus.invalid:
      return "3";
    case CheckResultStatus.ok:
      return "4";
    case CheckResultStatus.info:
      return "5";
    case CheckResultStatus.muted:
      return "6";
    case CheckResultStatus.skip:
    case CheckResultStatus.skipped:
      return "7";
    case CheckResultStatus.empty:
      return "8";
    case CheckResultStatus.tbd:
      return "9";
    default:
      return "99";
  }
};

const getCheckSeverityGroupingKey = (
  severity: CheckSeverity | undefined,
): string => {
  switch (severity) {
    case "critical":
      return "Critical";
    case "high":
      return "High";
    case "low":
      return "Low";
    case "medium":
      return "Medium";
    default:
      return "Unspecified";
  }
};

const getCheckSeveritySortKey = (
  severity: CheckSeverity | undefined,
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
  dimensions: CheckResultDimension[],
): string => {
  if (!dimensionKey) {
    return "<not set>";
  }
  const foundDimension = findDimension(dimensions, dimensionKey);
  return foundDimension ? foundDimension.value : `<not set>`;
};

function getCheckTagGroupingKey(tagKey: string | undefined, tags: CheckTags) {
  if (!tagKey) {
    return "Tag key not set";
  }
  return tags[tagKey] || `<not set>`;
}

const getCheckReasonGroupingKey = (reason: string | undefined): string => {
  return reason || "<not set>";
};

const getCheckResourceGroupingKey = (resource: string | undefined): string => {
  return resource || "<not set>";
};

// const getCheckResultGroupingKey = (checkResult: CheckResult): string => {
//   return `${checkResult.control.name}-${checkResult.resource}`;
// };

const getCheckGroupingKey = (
  checkResult: CheckResult,
  group: CheckDisplayGroup,
) => {
  switch (group.type) {
    case "dimension":
      return getCheckDimensionGroupingKey(group.value, checkResult.dimensions);
    case "control_tag":
      return getCheckTagGroupingKey(group.value, checkResult.tags);
    case "reason":
      return getCheckReasonGroupingKey(checkResult.reason);
    case "resource":
      return getCheckResourceGroupingKey(checkResult.resource);
    // case "result":
    //   return getCheckResultGroupingKey(checkResult);
    case "severity":
      return getCheckSeverityGroupingKey(checkResult.control.severity);
    case "status":
      return getCheckStatusGroupingKey(checkResult.status);
    case "benchmark":
      if (checkResult.benchmark_trunk.length <= 1) {
        return null;
      }
      return checkResult.benchmark_trunk[checkResult.benchmark_trunk.length - 1]
        .name;
    case "control":
      return checkResult.control.name;
    default:
      return "Other";
  }
};

const getCheckGroupingNode = (
  checkResult: CheckResult,
  group: CheckDisplayGroup,
  children: CheckNode[],
  benchmarkChildrenLookup: { [name: string]: CheckNode[] },
  groupingKeysBeforeBenchmark: string[] = [],
  parentGroupType: string | null,
): CheckNode => {
  switch (group.type) {
    case "dimension":
      const dimensionValue = getCheckDimensionGroupingKey(
        group.value,
        checkResult.dimensions,
      );
      return new KeyValuePairNode(
        dimensionValue,
        "dimension",
        group.value || "Dimension key not set",
        dimensionValue,
        children,
      );
    case "control_tag":
      const value = getCheckTagGroupingKey(group.value, checkResult.tags);
      return new KeyValuePairNode(
        value,
        "control_tag",
        group.value || "Tag key not set",
        value,
        children,
      );
    case "reason":
      return new KeyValuePairNode(
        checkResult.reason || "𤭢", // U+24B62 - very high in sort order - will almost guarantee to put this to the end,
        "reason",
        "reason",
        getCheckReasonGroupingKey(checkResult.reason),
        children,
      );
    case "resource":
      return new KeyValuePairNode(
        checkResult.resource || "𤭢", // U+24B62 - very high in sort order - will almost guarantee to put this to the end
        "resource",
        "resource",
        getCheckResourceGroupingKey(checkResult.resource),
        children,
      );
    // case "result":
    //   return new ControlResultNode(
    //     checkResult,
    //     `${checkResult.control.name}-${checkResult.resource}`,
    //     "result",
    //     "result",
    //     getCheckResultGroupingKey(checkResult),
    //     children,
    //   );
    case "severity":
      return new KeyValuePairNode(
        getCheckSeveritySortKey(checkResult.control.severity),
        "severity",
        "severity",
        getCheckSeverityGroupingKey(checkResult.control.severity),
        children,
      );
    case "status":
      return new KeyValuePairNode(
        getCheckStatusSortKey(checkResult.status),
        "status",
        "status",
        getCheckStatusGroupingKey(checkResult.status),
        children,
      );
    case "benchmark":
      return checkResult.benchmark_trunk.length > 1
        ? addBenchmarkTrunkNode(
            checkResult.benchmark_trunk.slice(1),
            children,
            benchmarkChildrenLookup,
            groupingKeysBeforeBenchmark,
            parentGroupType,
          )
        : children[0];
    case "control":
      return new ControlNode(
        parentGroupType === "benchmark"
          ? checkResult.control.sort
          : checkResult.control.title || checkResult.control.name,
        checkResult.control.name,
        checkResult.control.title,
        checkResult.control.documentation,
        children,
      );
    default:
      throw new Error(`Unknown group type ${group.type}`);
  }
};

const addBenchmarkGroupingNode = (
  existingGroups: CheckNode[],
  groupingNode: CheckNode,
) => {
  const existingGroup = existingGroups.find(
    (existingGroup) => existingGroup.name === groupingNode.name,
  );
  if (existingGroup) {
    (existingGroup as BenchmarkNode).merge(groupingNode);
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

const groupCheckItems = (
  temp: { _: CheckNode[] },
  checkResult: CheckResult,
  groupingConfig: CheckDisplayGroup[],
  checkNodeStates: CheckGroupNodeStates,
  benchmarkChildrenLookup: { [name: string]: CheckNode[] },
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
        if (currentGroupingConfig.type === "benchmark") {
          checkResult.benchmark_trunk.forEach(
            (benchmark) =>
              (checkNodeStates[benchmark.name] = {
                expanded: false,
              }),
          );
        } else {
          checkNodeStates[groupKey] = {
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
        // benchmark 1 -> benchmark 2 -> control 1
        // benchmark 1 -> control 2
        // ...when we build the benchmark grouping node for control 1, its key will be
        // for benchmark 2, but we'll add a hierarchical grouping node for benchmark 1 -> benchmark 2
        // When we come to get the benchmark grouping node for control 2, we'll need to add
        // the control to the existing children of benchmark 1
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

const getCheckResultNode = (checkResult: CheckResult) => {
  if (checkResult.type === "loading") {
    return new ControlRunningNode(checkResult);
  } else if (checkResult.type === "error") {
    return new ControlErrorNode(checkResult);
  } else if (checkResult.type === "empty") {
    return new ControlEmptyResultNode(checkResult);
  }
  return new ControlResultNode(checkResult);
  // return new ControlResultNode(
  //   checkResult,
  //   `${checkResult.control.name}-${checkResult.resource}`,
  //   "result",
  //   "result",
  //   getCheckResultGroupingKey(checkResult),
  //   undefined,
  // );
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

type CheckGroupingProviderProps = {
  children: null | JSX.Element | JSX.Element[];
  definition: PanelDefinition;
  benchmarkChildren?: PanelDefinition[] | undefined;
  groupingConfig: DisplayGroup[];
  checkFilterConfig: Filter;
  panelsMap: PanelsMap;
  setDashboardControlsContext: (context: any) => void;
};

function recordFilterValues(
  filterValues: {
    severity: { value: {} };
    reason: { value: {} };
    resource: { value: {} };
    control: { value: {} };
    control_tag: { key: {}; value: {} };
    dimension: { key: {}; value: {} };
    benchmark: { value: {} };
    status: {
      alarm: number;
      skip: number;
      error: number;
      ok: number;
      empty: number;
      info: number;
    };
  },
  checkResult: CheckResult,
) {
  // Record the benchmark of this check result to allow assisted filtering later
  if (!!checkResult.benchmark_trunk && checkResult.benchmark_trunk.length > 0) {
    for (const benchmark of checkResult.benchmark_trunk) {
      filterValues.benchmark.value[benchmark.name] = filterValues.benchmark
        .value[benchmark.name] || { title: benchmark.title, count: 0 };
      filterValues.benchmark.value[benchmark.name].count += 1;
    }
  }

  // Record the control of this check result to allow assisted filtering later
  filterValues.control.value[checkResult.control.name] = filterValues.control
    .value[checkResult.control.name] || {
    title: checkResult.control.title,
    count: 0,
  };
  filterValues.control.value[checkResult.control.name].count += 1;

  // Record the status of this check result to allow assisted filtering later
  filterValues.status[checkResult.status] =
    filterValues.status[checkResult.status] || 0;
  filterValues.status[checkResult.status] += 1;

  // Record the reason of this check result to allow assisted filtering later
  if (checkResult.reason) {
    filterValues.reason.value[checkResult.reason] =
      filterValues.reason.value[checkResult.reason] || 0;
    filterValues.reason.value[checkResult.reason] += 1;
  }

  // Record the resource of this check result to allow assisted filtering later
  if (checkResult.resource) {
    filterValues.resource.value[checkResult.resource] =
      filterValues.resource.value[checkResult.resource] || 0;
    filterValues.resource.value[checkResult.resource] += 1;
  }

  // Record the severity of this check result to allow assisted filtering later
  if (checkResult.severity) {
    filterValues.severity.value[checkResult.severity.toString()] =
      filterValues.severity.value[checkResult.severity.toString()] || 0;
    filterValues.severity.value[checkResult.severity.toString()] += 1;
  }

  // Record the dimension keys/values + value/key counts of this check result to allow assisted filtering later
  for (const dimension of checkResult.dimensions) {
    if (!(dimension.key in filterValues.dimension.key)) {
      filterValues.dimension.key[dimension.key] = {
        [dimension.value]: 0,
      };
    }
    if (!(dimension.value in filterValues.dimension.key[dimension.key])) {
      filterValues.dimension.key[dimension.key][dimension.value] = 0;
    }
    filterValues.dimension.key[dimension.key][dimension.value] += 1;

    if (!(dimension.value in filterValues.dimension.value)) {
      filterValues.dimension.value[dimension.value] = {
        [dimension.key]: 0,
      };
    }
    if (!(dimension.key in filterValues.dimension.value[dimension.value])) {
      filterValues.dimension.value[dimension.value][dimension.key] = 0;
    }
    filterValues.dimension.value[dimension.value][dimension.key] += 1;
  }

  // Record the dimension keys/values + value/key counts of this check result to allow assisted filtering later
  for (const [tagKey, tagValue] of Object.entries(checkResult.tags || {})) {
    if (!(tagKey in filterValues.control_tag.key)) {
      filterValues.control_tag.key[tagKey] = {
        [tagValue]: 0,
      };
    }
    if (!(tagValue in filterValues.control_tag.key[tagKey])) {
      filterValues.control_tag.key[tagKey][tagValue] = 0;
    }
    filterValues.control_tag.key[tagKey][tagValue] += 1;

    if (!(tagValue in filterValues.control_tag.value)) {
      filterValues.control_tag.value[tagValue] = {
        [tagKey]: 0,
      };
    }
    if (!(tagKey in filterValues.control_tag.value[tagValue])) {
      filterValues.control_tag.value[tagValue][tagKey] = 0;
    }
    filterValues.control_tag.value[tagValue][tagKey] += 1;
  }
}

const includeResult = (result: CheckResult, filterConfig: Filter): boolean => {
  if (
    !filterConfig ||
    !filterConfig.expressions ||
    filterConfig.expressions.length === 0
  ) {
    return true;
  }
  let matches: boolean[] = [];
  for (const filter of filterConfig.expressions) {
    if (!filter.type) {
      continue;
    }

    switch (filter.type) {
      case "benchmark": {
        let matchesTrunk = false;
        for (const benchmark of result.benchmark_trunk || []) {
          const match = applyFilter(filter, benchmark.name);
          if (match) {
            matchesTrunk = true;
            break;
          }
        }
        matches.push(matchesTrunk);
        break;
      }
      case "control": {
        matches.push(applyFilter(filter, result.control.name));
        break;
      }
      case "reason": {
        matches.push(applyFilter(filter, result.reason));
        break;
      }
      case "resource": {
        matches.push(applyFilter(filter, result.resource));
        break;
      }
      case "severity": {
        matches.push(applyFilter(filter, result.severity || ""));
        break;
      }
      case "status": {
        matches.push(applyFilter(filter, result.status.toString()));
        break;
      }
      case "dimension": {
        let matchesDimensions = false;
        for (const dimension of result.dimensions || []) {
          if (
            filter.key === dimension.key &&
            applyFilter(filter, dimension.value)
          ) {
            matchesDimensions = true;
            break;
          }
        }
        matches.push(matchesDimensions);
        break;
      }
      case "control_tag": {
        let matchesTags = false;
        for (const [tagKey, tagValue] of Object.entries(result.tags || {})) {
          if (filter.key === tagKey && applyFilter(filter, tagValue)) {
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
  benchmarkChildren: PanelDefinition[] | undefined,
  panelsMap: PanelsMap | undefined,
  groupingConfig: CheckDisplayGroup[],
  checkFilterConfig: Filter,
  skip = false,
) => {
  return useMemo(() => {
    const filterValues: CheckGroupFilterValues = {
      benchmark: { value: {} },
      control: { value: {} },
      control_tag: { key: {}, value: {} },
      dimension: { key: {}, value: {} },
      reason: { value: {} },
      resource: { value: {} },
      severity: { value: {} },
      status: {
        alarm: 0,
        empty: 0,
        error: 0,
        info: 0,
        ok: 0,
        skip: 0,
        skipped: 0,
        invalid: 0,
        muted: 0,
        tbd: 0,
      },
    };

    if (!definition || skip || !panelsMap) {
      // Return empty but valid types
      return [null, null, [], null, {}, filterValues] as [
        BenchmarkType | null,
        CheckNode | null,
        CheckSummary[],
        any,
        any,
        CheckGroupFilterValues,
      ];
    }

    // @ts-ignore
    const nestedBenchmarks = benchmarkChildren?.filter(
      (child) => child.panel_type === "benchmark",
    );
    const nestedControls =
      definition.panel_type === "control"
        ? [definition]
        : // @ts-ignore
          benchmarkChildren?.filter((child) => child.panel_type === "control");

    const b = new BenchmarkType(
      "0",
      definition.name,
      definition.title,
      definition.description,
      definition.documentation,
      nestedBenchmarks,
      nestedControls,
      panelsMap,
      [],
    );

    const checkNodeStates: CheckGroupNodeStates = {};
    const result: CheckNode[] = [];
    const temp = { _: result };
    const benchmarkChildrenLookup = {};

    // We'll loop over each control result and build up the grouped nodes from there
    b.all_control_results.forEach((checkResult) => {
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
        groupingConfig,
        checkNodeStates,
        benchmarkChildrenLookup,
        [],
      );
      // Build and add a check result node to the children of the trailing group.
      // This will be used to calculate totals and severity, amongst other things.
      const node = getCheckResultNode(checkResult);
      grouping._.push(node);
    });

    const results = new RootNode(result);

    const firstChildSummaries: CheckSummary[] = [];
    for (const child of results.children) {
      firstChildSummaries.push(child.summary);
    }

    return [
      b,
      results,
      firstChildSummaries,
      checkNodeStates,
      filterValues,
    ] as const;
  }, [checkFilterConfig, definition, groupingConfig, panelsMap, skip]);
};

const GroupingProvider = ({
  children,
  definition,
  benchmarkChildren,
  groupingConfig,
  checkFilterConfig,
  panelsMap,
  setDashboardControlsContext,
}: CheckGroupingProviderProps) => {
  const [nodeStates, dispatch] = useReducer(reducer, { nodes: {} });

  const [
    benchmark,
    grouping,
    firstChildSummaries,
    tempNodeStates,
    filterValues,
  ] = useGroupingInternal(
    definition,
    benchmarkChildren,
    panelsMap,
    groupingConfig,
    checkFilterConfig,
  );

  const previousGroupings = usePrevious({ groupingConfig });

  useEffect(() => {
    if (
      previousGroupings &&
      // @ts-ignore
      JSON.stringify(previousGroupings.groupingConfig) ===
        JSON.stringify(groupingConfig)
    ) {
      return;
    }
    dispatch({
      type: GroupingActions.UPDATE_NODES,
      nodes: tempNodeStates,
    });
  }, [previousGroupings, groupingConfig, tempNodeStates]);

  useEffect(() => {
    setDashboardControlsContext(filterValues);
  }, [filterValues, setDashboardControlsContext]);

  return (
    <GroupingContext.Provider
      value={{
        benchmark: benchmark as BenchmarkType | null,
        definition,
        dispatch,
        firstChildSummaries: firstChildSummaries as CheckSummary[],
        grouping: grouping as CheckNode | null,
        groupingConfig,
        nodeStates,
        filterValues: filterValues as CheckGroupFilterValues,
      }}
    >
      {children}
    </GroupingContext.Provider>
  );
};

const useBenchmarkGrouping = () => {
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
  useBenchmarkGrouping,
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
