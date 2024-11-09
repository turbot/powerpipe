import DetectionEmptyResultNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionEmptyResultNode";
import DetectionErrorNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionErrorNode";
import DetectionKeyValuePairNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionKeyValuePairNode";
import DetectionNode from "@powerpipe/components/dashboards/grouping/common/node/DetectionNode";
import DetectionResultNode from "../common/node/DetectionResultNode";
import DetectionSummaryChart from "@powerpipe/components/dashboards/grouping/DetetctionSummaryChart";
import sortBy from "lodash/sortBy";
import Table from "@powerpipe/components/dashboards/Table";
import {
  AlarmIcon,
  CollapseBenchmarkIcon,
  EmptyIcon,
  ErrorIcon,
  ExpandCheckNodeIcon,
} from "@powerpipe/constants/icons";
import { classNames } from "@powerpipe/utils/styles";
import {
  DetectionResult,
  DetectionSeveritySummary,
} from "@powerpipe/components/dashboards/grouping/common";
import {
  GroupingActions,
  useDetectionGrouping,
} from "@powerpipe/hooks/useDetectionGrouping";
import { useMemo } from "react";

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
  label: string;
  count: number;
  title: string;
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
        <DetectionPanel key={child.name} depth={depth} node={child} />
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
  return (
    <div className="flex bg-dashboard-panel print:bg-white last:rounded-b-md space-x-4 overflow-x-auto">
      <div className="flex flex-col md:flex-row flex-grow">
        <Table
          name={`${result.detection.name}.table`}
          panel_type="table"
          data={{ rows: result.rows, columns: result.columns }}
          filterEnabled
          context={result.detection.name}
        />
      </div>
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
  count,
  label,
  title,
}: DetectionPanelSeverityBadgeProps) => {
  return (
    <div
      className={classNames(
        "border rounded-md text-sm divide-x",
        count > 0 ? "border-severity" : "border-skip",
        count > 0
          ? "bg-severity text-white divide-white"
          : "text-skip divide-skip",
      )}
      title={title}
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
      {critical !== undefined && (
        <DetectionPanelSeverityBadge
          label="Critical"
          count={critical}
          title={`${critical.toLocaleString()} critical severity ${
            critical === 1 ? "result" : "results"
          }`}
        />
      )}
      {high !== undefined && (
        <DetectionPanelSeverityBadge
          label="High"
          count={high}
          title={`${high.toLocaleString()} high severity ${
            high === 1 ? "result" : "results"
          }`}
        />
      )}{" "}
      {medium !== undefined && (
        <DetectionPanelSeverityBadge
          label="Medium"
          count={medium}
          title={`${medium.toLocaleString()} medium severity ${
            high === 1 ? "result" : "results"
          }`}
        />
      )}{" "}
      {low !== undefined && (
        <DetectionPanelSeverityBadge
          label="Low"
          count={low}
          title={`${high.toLocaleString()} low severity ${
            high === 1 ? "result" : "results"
          }`}
        />
      )}
    </>
  );
};

const DetectionPanel = ({ depth, node }: DetectionPanelProps) => {
  const { firstChildSummaries, dispatch, groupingsConfig, nodeStates } =
    useDetectionGrouping();
  const expanded = nodeStates[node.name]
    ? nodeStates[node.name].expanded
    : false;

  const [child_nodes, error_nodes, empty_nodes, result_nodes, can_be_expanded] =
    useMemo(() => {
      const children: DetectionNode[] = [];
      const errors: DetectionErrorNode[] = [];
      const empty: DetectionEmptyResultNode[] = [];
      const results: DetectionResultNode[] = [];
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
      }
      return [
        sortBy(children, "sort"),
        sortBy(errors, "sort"),
        sortBy(empty, "sort"),
        results,
        children.length > 0 ||
          (groupingsConfig &&
            groupingsConfig.length > 0 &&
            groupingsConfig[groupingsConfig.length - 1].type === "result" &&
            (errors.length > 0 || empty.length > 0 || results.length > 0)),
      ];
    }, [groupingsConfig, node]);

  return (
    <>
      <div
        id={node.name}
        className={classNames(
          getMargin(depth - 1),
          depth === 1 && node.type === "benchmark"
            ? "print:break-before-page"
            : null,
          node.type === "detection_benchmark" || node.type === "detection"
            ? "print:break-inside-avoid-page"
            : null,
        )}
      >
        <section
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
                <DetectionPanelSeverity
                  severity_summary={node.severity_summary}
                />
              </div>
              {/* <div className="TestName"> */}
              <div className="flex-shrink-0 w-40 md:w-72 lg:w-96">
                <DetectionSummaryChart
                  status={node.status}
                  summary={node.summary}
                  firstChildSummaries={firstChildSummaries}
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
        {can_be_expanded &&
          expanded &&
          groupingsConfig &&
          groupingsConfig[groupingsConfig.length - 1].type === "result" && (
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

export default DetectionPanel;
