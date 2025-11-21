import CopyToClipboard from "@powerpipe/components/CopyToClipboard";
import ErrorMessage from "@powerpipe/components/ErrorMessage";
import Icon from "@powerpipe/components/Icon";
import LoadingIndicator from "@powerpipe/components/dashboards/LoadingIndicator";
import { DashboardRunState, PanelDefinition } from "@powerpipe/types";
import {
  EdgeStatus,
  GraphStatuses,
  NodeStatus,
  WithStatus,
} from "@powerpipe/components/dashboards/graphs/types";
import { Node } from "reactflow";

type NodeAndEdgePanelInformationProps = {
  panel: PanelDefinition;
  nodes: Node[];
  status: DashboardRunState;
  statuses: GraphStatuses;
};

const nodeOrEdgeTitle = (nodeOrEdge: NodeStatus | EdgeStatus) =>
  nodeOrEdge.title ||
  nodeOrEdge?.category?.title ||
  nodeOrEdge?.category?.name ||
  nodeOrEdge.id;

const formatNodeAndEdgeErrorsForCopy = (
  panel: PanelDefinition,
  statuses: GraphStatuses,
): string => {
  const {
    error: errorStatuses,
    running: runningStatuses,
    blocked: blockedStatuses,
    complete: completeStatuses,
  } = statuses;

  // Build markdown format
  let output = `# Graph Panel Error Report\n\n`;
  output += `**Panel:** ${panel.name}\n`;
  output += `**Type:** ${panel.panel_type || "graph"}\n`;
  output += `**Status:** ${panel.status || "unknown"}\n`;
  output += `**Timestamp:** ${new Date().toISOString()}\n\n`;

  // Summary section
  output += `## Summary\n`;
  output += `- Complete: ${completeStatuses.nodes.length} nodes, ${completeStatuses.edges.length} edges, ${completeStatuses.withs.length} withs\n`;
  output += `- Running: ${runningStatuses.nodes.length} nodes, ${runningStatuses.edges.length} edges, ${runningStatuses.withs.length} withs\n`;
  output += `- Waiting: ${blockedStatuses.nodes.length} nodes, ${blockedStatuses.edges.length} edges, ${blockedStatuses.withs.length} withs\n`;
  output += `- **Errors: ${errorStatuses.nodes.length} nodes, ${errorStatuses.edges.length} edges, ${errorStatuses.withs.length} withs**\n\n`;

  // Only include error details if there are errors
  if (errorStatuses.total === 0) {
    return output;
  }

  output += `## Error Details\n\n`;

  // With errors
  if (errorStatuses.withs.length > 0) {
    errorStatuses.withs.forEach((withBlock: WithStatus) => {
      output += `### With: ${withBlock.id}\n`;
      if (withBlock.title) {
        output += `**Title:** ${withBlock.title}\n`;
      }
      if (withBlock.error) {
        output += `**Error:**\n\`\`\`\n${withBlock.error}\n\`\`\`\n\n`;
      } else {
        output += `**Error:** (no error message provided)\n\n`;
      }
    });
  }

  // Node errors
  if (errorStatuses.nodes.length > 0) {
    errorStatuses.nodes.forEach((node: NodeStatus) => {
      output += `### Node: ${node.id}\n`;
      const title = nodeOrEdgeTitle(node);
      if (title && title !== node.id) {
        output += `**Title:** ${title}\n`;
      }
      if (node.category?.title) {
        output += `**Category:** ${node.category.title}\n`;
      }
      if (node.error) {
        output += `**Error:**\n\`\`\`\n${node.error}\n\`\`\`\n\n`;
      } else {
        output += `**Error:** (no error message provided)\n\n`;
      }
    });
  }

  // Edge errors
  if (errorStatuses.edges.length > 0) {
    errorStatuses.edges.forEach((edge: EdgeStatus) => {
      output += `### Edge: ${edge.id}\n`;
      const title = nodeOrEdgeTitle(edge);
      if (title && title !== edge.id) {
        output += `**Title:** ${title}\n`;
      }
      if (edge.category?.title) {
        output += `**Category:** ${edge.category.title}\n`;
      }
      if (edge.error) {
        output += `**Error:**\n\`\`\`\n${edge.error}\n\`\`\`\n\n`;
      } else {
        output += `**Error:** (no error message provided)\n\n`;
      }
    });
  }

  return output;
};

const WaitingRow = ({ title }) => (
  <div className="flex items-center space-x-1">
    <Icon
      className="w-3.5 h-3.5 text-foreground-light shrink-0"
      icon="pending"
    />
    <span className="block truncate">{title}</span>
  </div>
);

const RunningRow = ({ title }) => (
  <div className="flex items-center space-x-1">
    <LoadingIndicator className="w-3.5 h-3.5 shrink-0" />
    <span className="block truncate">{title}</span>
  </div>
);

const ErrorRow = ({ title, error }: { title: string; error?: string }) => (
  <>
    <div className="flex items-center space-x-1">
      <Icon
        className="w-3.5 h-3.5 text-alert shrink-0"
        icon="materialsymbols-solid:error"
      />
      <span className="block">{title}</span>
    </div>
    {error && (
      <span className="block">
        <ErrorMessage error={error} />
      </span>
    )}
  </>
);

const NodeAndEdgePanelInformation = ({
  panel,
  nodes,
  status,
  statuses,
}: NodeAndEdgePanelInformationProps) => {
  const copyData = formatNodeAndEdgeErrorsForCopy(panel, statuses);

  return (
    <div className="space-y-2 overflow-y-scroll">
      {/* Copy button - minimal, just the icon */}
      <div className="flex justify-end px-4 pt-2" title="Copy to clipboard">
        <CopyToClipboard
          data={copyData}
          className="text-foreground-light hover:text-foreground transition-colors"
        />
      </div>

      <div className="space-y-1 px-4">
        <div>
          {statuses.complete.total} complete, {statuses.running.total} running,{" "}
          {statuses.blocked.total} waiting, {statuses.error.total}{" "}
          {statuses.error.total === 1 ? "error" : "errors"}.
        </div>
      {statuses.initialized.total === 0 &&
        statuses.blocked.total === 0 &&
        statuses.running.total === 0 &&
        statuses.complete.total === 0 &&
        status === "complete" &&
        nodes.length === 0 && (
          <span className="block text-foreground-light italic">
            No nodes or edges
          </span>
        )}
      {statuses.running.withs.map((withStatus, idx) => (
        <RunningRow
          key={`with:${withStatus.id}-${idx}`}
          title={`with: ${withStatus.title || withStatus.id}`}
        />
      ))}
      {statuses.running.nodes.map((node, idx) => (
        <RunningRow
          key={`node:${node.id}-${idx}`}
          title={`node: ${nodeOrEdgeTitle(node)}`}
        />
      ))}
      {statuses.running.edges.map((edge, idx) => (
        <RunningRow
          key={`edge:${edge.id}-${idx}`}
          title={`edge: ${nodeOrEdgeTitle(edge)}`}
        />
      ))}
      {statuses.blocked.withs.map((withStatus, idx) => (
        <WaitingRow
          key={`with:${withStatus.id}-${idx}`}
          title={`with: ${withStatus.title || withStatus.id}`}
        />
      ))}
      {statuses.blocked.nodes.map((node, idx) => (
        <WaitingRow
          key={`node:${node.id}-${idx}`}
          title={`node: ${nodeOrEdgeTitle(node)}`}
        />
      ))}
      {statuses.blocked.edges.map((edge, idx) => (
        <WaitingRow
          key={`edge:${edge.id}-${idx}`}
          title={`edge: ${nodeOrEdgeTitle(edge)}`}
        />
      ))}
      {statuses.error.withs.map((withStatus, idx) => (
        <ErrorRow
          key={`with:${withStatus.id}-${idx}`}
          title={`with: ${withStatus.title || withStatus.id}`}
          error={withStatus.error}
        />
      ))}
      {statuses.error.nodes.map((node, idx) => (
        <ErrorRow
          key={`node:${node.id}-${idx}`}
          title={`node: ${nodeOrEdgeTitle(node)}`}
          error={node.error}
        />
      ))}
      {statuses.error.edges.map((edge, idx) => (
        <ErrorRow
          key={`edge:${edge.id}-${idx}`}
          title={`edge: ${nodeOrEdgeTitle(edge)}`}
          error={edge.error}
        />
      ))}
      </div>
    </div>
  );
};

export default NodeAndEdgePanelInformation;
