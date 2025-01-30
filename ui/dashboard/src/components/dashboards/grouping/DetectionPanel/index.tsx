import DetectionEmptyResultNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionEmptyResultNode";
import DetectionErrorNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionErrorNode";
import DetectionKeyValuePairNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionKeyValuePairNode";
import DetectionNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionNode";
import DetectionResultNode from "../common/node/DetectionResultNode";
import DetectionSummaryChart from "@powerpipe/components/dashboards/grouping/DetetctionSummaryChart";
import DocumentationView from "@powerpipe/components/dashboards/grouping/DocumentationView";
import PanelControls from "@powerpipe/components/dashboards/layout/Panel/PanelControls";
import sortBy from "lodash/sortBy";
import Table from "@powerpipe/components/dashboards/Table";
import useDownloadDetectionData from "@powerpipe/hooks/useDownloadDetectionData";
import {
  AlarmIcon,
  CollapseBenchmarkIcon,
  EmptyIcon,
  ErrorIcon,
  ExpandCheckNodeIcon,
} from "@powerpipe/constants/icons";
import { classNames } from "@powerpipe/utils/styles";
import {
  CheckSeverity,
  DetectionResult,
  DetectionSeveritySummary,
} from "@powerpipe/components/dashboards/grouping/common";
import { DashboardActions } from "@powerpipe/types";
import {
  GroupingActions,
  useDetectionGrouping,
} from "@powerpipe/hooks/useDetectionGrouping";
import {
  IPanelControl,
  PanelControlsProvider,
  usePanelControls,
} from "@powerpipe/hooks/usePanelControls";
import { noop } from "@powerpipe/utils/func";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect, useMemo, useState } from "react";

type DetectionChildrenProps = {
  depth: number;
  children: DetectionNode[];
};

type DetectionResultsProps = {
  empties: DetectionEmptyResultNode[];
  errors: DetectionErrorNode[];
  results: DetectionResultNode[];
};

type DetectionPanelProps = {
  depth: number;
  node: DetectionNode;
};

type DetectionPanelSeverityProps = {
  severity_summary: DetectionSeveritySummary;
};

type DetectionPanelSeverityBadgeProps = {
  severity: CheckSeverity;
  count: number;
};

type DetectionEmptyResultRowProps = {
  node: DetectionEmptyResultNode;
};

type DetectionResultRowProps = {
  result: DetectionResult;
};

type DetectionErrorRowProps = {
  error: string;
};

type DetectionResultRowStatusIconProps = {
  total: number;
};

const getMargin = (depth) => {
  switch (depth) {
    case 1:
      return "ml-[6px] md:ml-[24px]";
    case 2:
      return "ml-[12px] md:ml-[48px]";
    case 3:
      return "ml-[18px] md:ml-[72px]";
    case 4:
      return "ml-[24px] md:ml-[96px]";
    case 5:
      return "ml-[30px] md:ml-[120px]";
    case 6:
      return "ml-[36px] md:ml-[144px]";
    default:
      return "ml-0";
  }
};

const DetectionChildren = ({ children, depth }: DetectionChildrenProps) => {
  if (!children) {
    return null;
  }

  return (
    <>
      {children.map((child) => (
        <DetectionPanelWrapper key={child.name} depth={depth} node={child} />
      ))}
    </>
  );
};

const DetectionResultRowStatusIcon = ({
  total,
}: DetectionResultRowStatusIconProps) => {
  if (total > 0) {
    return <AlarmIcon className="h-5 w-5 text-alert" />;
  }
  return <EmptyIcon className="h-5 w-5 text-skip" />;
};

const getDetectionResultRowIconTitle = (total: number) => {
  if (total > 1) {
    return `${total} results.`;
  } else if (total === 1) {
    return "1 result.";
  }
  return "No results.";
};

const DetectionResultRow = ({ result }: DetectionResultRowProps) => {
  const { panelsMap } = useDashboardState();
  return (
    <div className="flex bg-dashboard-panel print:bg-white last:rounded-b-md space-x-4 overflow-x-auto">
      <Table
        name={result.detection.name}
        panel_type="table"
        data={{ rows: result.rows, columns: result.columns }}
        properties={{
          display_columns:
            panelsMap[result.detection.name]?.properties?.display_columns || [],
        }}
        isDetectionTable
      />
    </div>
  );
};

const DetectionEmptyResultRow = ({ node }: DetectionEmptyResultRowProps) => {
  return (
    <div className="flex bg-dashboard-panel print:bg-white p-4 last:rounded-b-md space-x-4">
      <div className="flex-shrink-0" title={getDetectionResultRowIconTitle(0)}>
        <DetectionResultRowStatusIcon total={0} />
      </div>
      <div className="leading-4 mt-px">{node.title}</div>
    </div>
  );
};

const DetectionErrorRow = ({ error }: DetectionErrorRowProps) => {
  return (
    <div className="flex bg-dashboard-panel print:bg-white p-4 last:rounded-b-md space-x-4">
      <div className="flex-shrink-0" title="Error">
        <ErrorIcon className="h-5 w-5 text-alert" />
      </div>
      <div className="leading-4 mt-px">{error}</div>
    </div>
  );
};

const DetectionResults = ({
  empties,
  errors,
  results,
}: DetectionResultsProps) => {
  if (empties.length === 0 && errors.length === 0 && results.length === 0) {
    return null;
  }

  return (
    <div
      className={classNames(
        "border-t shadow-sm rounded-b-md divide-y divide-table-divide border-divide print:shadow-none print:border print:break-before-avoid-page print:break-after-avoid-page print:break-inside-auto",
      )}
    >
      {empties.map((emptyNode) => (
        <DetectionEmptyResultRow key={`${emptyNode.name}`} node={emptyNode} />
      ))}
      {errors.map((errorNode) => (
        <DetectionErrorRow key={`${errorNode.name}`} error={errorNode.error} />
      ))}
      {results.map((resultNode) => (
        <DetectionResultRow
          key={`${resultNode.result.detection.name}-${
            resultNode.result.resource
          }${
            resultNode.result.dimensions
              ? `-${resultNode.result.dimensions
                  .map((d) => `${d.key}=${d.value}`)
                  .join("-")}`
              : ""
          }`}
          result={resultNode.result}
        />
      ))}
    </div>
  );
};

const DetectionPanelSeverityBadge = ({
  severity,
  count,
}: DetectionPanelSeverityBadgeProps) => {
  const hasSeverity = count > 0;
  let label: string;
  switch (severity) {
    case "low":
      label = "Low";
      break;
    case "medium":
      label = "Medium";
      break;
    case "high":
      label = "High";
      break;
    case "critical":
      label = "Critical";
      break;
    default:
      label = "Unknown";
      break;
  }
  return (
    <div
      className={classNames(
        "border rounded-md text-sm divide-x",
        hasSeverity && severity === "low"
          ? "border-info bg-info text-white divide-white"
          : null,
        hasSeverity && severity === "medium"
          ? "border-severity bg-severity text-white divide-white"
          : null,
        hasSeverity && severity === "high"
          ? "border-orange bg-orange text-white divide-white"
          : null,
        hasSeverity && severity === "critical"
          ? "border-alert bg-alert text-white divide-white"
          : null,
        !hasSeverity ? "border-skip text-skip divide-skip" : null,
      )}
      title={`${count.toLocaleString()} ${severity} severity ${
        count === 1 ? "result" : "results"
      }`}
    >
      <span className={classNames("px-2 py-px")}>{label}</span>
      {count > 0 && <span className={classNames("px-2 py-px")}>{count}</span>}
    </div>
  );
};

const DetectionPanelSeverity = ({
  severity_summary,
}: DetectionPanelSeverityProps) => {
  const critical = severity_summary["critical"];
  const high = severity_summary["high"];
  const medium = severity_summary["medium"];
  const low = severity_summary["low"];

  if (
    critical === undefined &&
    high === undefined &&
    medium === undefined &&
    low === undefined
  ) {
    return null;
  }

  return (
    <>
      {critical !== undefined && critical !== null && (
        <DetectionPanelSeverityBadge severity="critical" count={critical} />
      )}
      {high !== undefined && high !== null && (
        <DetectionPanelSeverityBadge severity="high" count={high} />
      )}{" "}
      {medium !== undefined && medium !== null && (
        <DetectionPanelSeverityBadge severity="medium" count={medium} />
      )}{" "}
      {low !== undefined && low !== null && (
        <DetectionPanelSeverityBadge severity="low" count={low} />
      )}
    </>
  );
};

const recordChildResults = (
  node: DetectionNode,
  allResults: DetectionResultNode[],
): DetectionResultNode[] => {
  for (const child of node.children || []) {
    if (child.type === "result") {
      allResults.push(child as DetectionResultNode);
    } else if (child.children?.length) {
      return recordChildResults(child, allResults);
    }
  }
  return allResults;
};

const DetectionPanel = ({ depth, node }: DetectionPanelProps) => {
  const {
    firstChildSummaries,
    hasSeverityResults,
    dispatch,
    groupingConfig,
    nodeStates,
  } = useDetectionGrouping();
  const {
    enabled: panelControlsEnabled,
    panelControls,
    showPanelControls,
    setCustomControls,
    setShowPanelControls,
  } = usePanelControls();
  const { dispatch: dispatchDashboardState, overlayVisible } =
    useDashboardState();
  const [referenceElement, setReferenceElement] = useState(null);
  const expanded = nodeStates[node.name]
    ? nodeStates[node.name].expanded
    : false;

  const [
    child_nodes,
    error_nodes,
    empty_nodes,
    result_nodes,
    descendant_result_nodes,
    can_be_expanded,
  ] = useMemo(() => {
    const children: DetectionNode[] = [];
    const errors: DetectionErrorNode[] = [];
    const empty: DetectionEmptyResultNode[] = [];
    const results: DetectionResultNode[] = [];
    let descendantResults: DetectionResultNode[] = [];
    for (const child of node.children || []) {
      if (child.type === "error") {
        errors.push(child as DetectionErrorNode);
      } else if (child.type === "result") {
        results.push(child as DetectionResultNode);
      } else if (child.type === "empty_result") {
        empty.push(child as DetectionEmptyResultNode);
      } else if (child.type !== "running") {
        children.push(child);
      }

      if (child.children?.length) {
        recordChildResults(child, descendantResults);
      }
    }
    return [
      sortBy(children, "sort"),
      sortBy(errors, "sort"),
      sortBy(empty, "sort"),
      results,
      descendantResults,
      children.length > 0 ||
        (groupingConfig &&
          groupingConfig.length > 0 &&
          groupingConfig[groupingConfig.length - 1].type === "result" &&
          (errors.length > 0 || empty.length > 0 || results.length > 0)),
    ];
  }, [groupingConfig, node]);

  const { download } = useDownloadDetectionData(
    node,
    node.type === "detection"
      ? node.children
      : descendant_result_nodes.length > 0
        ? descendant_result_nodes
        : undefined,
  );

  useEffect(() => {
    const controls: IPanelControl[] = [
      {
        key: "download-data",
        disabled:
          (node.type === "detection" && !node.children?.length) ||
          (node.type !== "detection" && descendant_result_nodes.length === 0),
        title: "Download data",
        icon: "arrow-down-tray",
        action: download,
      },
    ];
    if (node.type === "benchmark") {
      controls.push({
        key: "open-in-new-window",
        title: "Open in new window",
        icon: "open_in_new",
        action: async () => {
          window.open(window.location.origin + "/" + node.name, "_blank");
        },
      });
    }
    setCustomControls(controls);
  }, [
    download,
    node.name,
    node.type,
    descendant_result_nodes,
    setCustomControls,
  ]);

  const hasResults =
    can_be_expanded &&
    groupingConfig &&
    groupingConfig[groupingConfig.length - 1].type === "result";

  return (
    <>
      <div
        id={node.name}
        className={classNames(
          getMargin(depth - 1),
          depth === 1 && node.type === "benchmark"
            ? "print:break-before-page"
            : null,
          node.type === "benchmark" || node.type === "detection"
            ? "print:break-inside-avoid-page"
            : null,
        )}
        onMouseEnter={
          panelControlsEnabled ? () => setShowPanelControls(true) : noop
        }
        onMouseLeave={() => setShowPanelControls(false)}
      >
        <section
          ref={setReferenceElement}
          className={classNames(
            "bg-dashboard-panel shadow-sm rounded-md border-divide print:border print:bg-white print:shadow-none",
            can_be_expanded ? "cursor-pointer" : null,
            expanded &&
              (empty_nodes.length > 0 ||
                error_nodes.length > 0 ||
                result_nodes.length > 0)
              ? "rounded-b-none border-b-0"
              : null,
          )}
          onClick={() =>
            can_be_expanded
              ? dispatch({
                  type: expanded
                    ? GroupingActions.COLLAPSE_NODE
                    : GroupingActions.EXPAND_NODE,
                  name: node.name,
                })
              : null
          }
        >
          {showPanelControls && !overlayVisible && (
            <PanelControls
              referenceElement={referenceElement}
              controls={panelControls}
              withOffset
            />
          )}
          <div className="p-4 flex items-center space-x-6">
            <div className="flex flex-grow justify-between items-center space-x-6">
              <div className="flex items-center space-x-4">
                <h3
                  id={`${node.name}-title`}
                  className="mt-0"
                  title={node.title}
                >
                  {(node as DetectionKeyValuePairNode).type ===
                    "detection_tag" &&
                  !!(node as DetectionKeyValuePairNode).key ? (
                    <span>{(node as DetectionKeyValuePairNode).key}: </span>
                  ) : null}
                  {node.title}
                </h3>
                <DocumentationView
                  documentation={node.documentation}
                  onOpen={() => {
                    dispatchDashboardState({
                      type: DashboardActions.SET_OVERLAY_VISIBLE,
                      value: true,
                    });
                  }}
                  onClose={() => {
                    dispatchDashboardState({
                      type: DashboardActions.SET_OVERLAY_VISIBLE,
                      value: false,
                    });
                  }}
                />
                <DetectionPanelSeverity
                  severity_summary={node.severity_summary}
                />
              </div>
              {/* <div className="TestName"> */}
              <div className="flex-shrink-0 w-40 md:w-72 lg:w-96">
                <DetectionSummaryChart
                  status={node.status}
                  summary={node.summary}
                  severitySummary={node.severity_summary}
                  firstChildSummaries={firstChildSummaries}
                  hasSeverityResults={hasSeverityResults}
                />
              </div>
              {/* <div>{node.summary.total}</div> */}
              {/* </div> */}
            </div>
            {can_be_expanded && !expanded && (
              <ExpandCheckNodeIcon className="w-5 md:w-7 h-5 md:h-7 flex-shrink-0 text-foreground-lightest" />
            )}
            {expanded && (
              <CollapseBenchmarkIcon className="w-5 md:w-7 h-5 md:h-7 flex-shrink-0 text-foreground-lightest" />
            )}
            {!can_be_expanded && <div className="w-5 md:w-7 h-5 md:h-7" />}
          </div>
        </section>
        {hasResults && expanded && (
          <DetectionResults
            empties={empty_nodes}
            errors={error_nodes}
            results={result_nodes}
          />
        )}
      </div>
      {can_be_expanded && expanded && (
        <DetectionChildren children={child_nodes} depth={depth + 1} />
      )}
    </>
  );
};

const DetectionPanelWrapper = ({ node, depth }: DetectionPanelProps) => {
  const definition = useMemo(
    () => ({
      data: node.data,
      panel_type: node.type,
      dashboard: node.name,
    }),
    [node],
  );

  return (
    <PanelControlsProvider
      definition={definition}
      enabled={true}
      panelDetailEnabled={false}
    >
      <DetectionPanel node={node} depth={depth} />
    </PanelControlsProvider>
  );
};

export default DetectionPanelWrapper;
