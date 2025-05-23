import { BaseChartProps } from "@powerpipe/components/dashboards/charts/types";
import { CardProps } from "@powerpipe/components/dashboards/Card";
import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import {
  DashboardInputs,
  DashboardPanelType,
  DashboardRunState,
  PanelDefinition,
  PanelDependenciesByStatus,
  PanelsMap,
} from "@powerpipe/types";
import { FlowProps } from "@powerpipe/components/dashboards/flows/types";
import { getNodeAndEdgeDataFormat } from "@powerpipe/components/dashboards/common/useNodeAndEdgeData";
import { GraphProps } from "@powerpipe/components/dashboards/graphs/types";
import { HierarchyProps } from "@powerpipe/components/dashboards/hierarchies/types";
import { ImageProps } from "@powerpipe/components/dashboards/Image";
import {
  InputProperties,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { NodeAndEdgeProperties } from "@powerpipe/components/dashboards/common/types";
import { PanelControlsProvider } from "@powerpipe/hooks/usePanelControls";
import { TableProps } from "@powerpipe/components/dashboards/Table";
import { TextProps } from "@powerpipe/components/dashboards/Text";
import { useContainer } from "@powerpipe/hooks/useContainer";
import { useDashboardInputs } from "@powerpipe/hooks/useDashboardInputs";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

type IPanelContext = {
  definition:
    | BaseChartProps
    | CardProps
    | FlowProps
    | GraphProps
    | HierarchyProps
    | ImageProps
    | InputProps
    | PanelDefinition
    | TableProps
    | TextProps;
  dependencies: PanelDefinition[];
  dependenciesByStatus: PanelDependenciesByStatus;
  inputPanelsAwaitingValue: PanelDefinition[];
  panelInformation: ReactNode | null;
  showPanelInformation: boolean;
  setPanelInformation: (information: ReactNode) => void;
  setShowPanelInformation: (show: boolean) => void;
};

type PanelProviderProps = {
  children: ReactNode;
  definition:
    | BaseChartProps
    | CardProps
    | FlowProps
    | GraphProps
    | HierarchyProps
    | ImageProps
    | InputProps
    | PanelDefinition
    | TableProps
    | TextProps;
  parentType: DashboardPanelType;
  showControls?: boolean;
};

const PanelContext = createContext<IPanelContext | null>(null);

const recordDependency = (
  definition: PanelDefinition,
  panelsMap: PanelsMap,
  inputs: DashboardInputs,
  dependencies: PanelDefinition[],
  dependenciesByStatus: PanelDependenciesByStatus,
  inputPanelsAwaitingValue: PanelDefinition[],
  recordedInputPanels: {},
) => {
  // Record this panel as a dependency
  dependencies.push(definition);

  // Keep track of this panel by its status
  const statuses =
    dependenciesByStatus[definition.status as DashboardRunState] || [];
  statuses.push(definition);
  dependenciesByStatus[definition.status as DashboardRunState] = statuses;

  // Is this panel an input? If so, does it have a value?
  const isInput = definition.panel_type === "input";
  const inputProperties = isInput
    ? (definition.properties as InputProperties)
    : null;
  const hasInputValue =
    isInput &&
    inputProperties?.unqualified_name &&
    !!inputs[inputProperties?.unqualified_name];
  if (isInput && !hasInputValue && !recordedInputPanels[definition.name]) {
    inputPanelsAwaitingValue.push(definition);
    recordedInputPanels[definition.name] = definition;
  }

  for (const dependency of definition?.dependencies || []) {
    const dependencyPanel = panelsMap[dependency];
    if (!dependencyPanel || !dependencyPanel.status) {
      continue;
    }
    recordDependency(
      dependencyPanel,
      panelsMap,
      inputs,
      dependencies,
      dependenciesByStatus,
      inputPanelsAwaitingValue,
      recordedInputPanels,
    );
  }
};

const PanelProvider = ({
  children,
  definition,
  parentType,
  showControls,
}: PanelProviderProps) => {
  const { updateChildStatus } = useContainer();
  const { panelsMap } = useDashboardState();
  const { inputs } = useDashboardInputs();
  const [showPanelInformation, setShowPanelInformation] = useState(false);
  const [panelInformation, setPanelInformation] = useState<ReactNode | null>(
    null,
  );
  const { dependencies, dependenciesByStatus, inputPanelsAwaitingValue } =
    useMemo(() => {
      if (!definition) {
        return {
          dependencies: [],
          dependenciesByStatus: {},
          inputPanelsAwaitingValue: [],
        };
      }
      const dataFormat = getNodeAndEdgeDataFormat(
        definition.properties as NodeAndEdgeProperties,
      );
      if (
        dataFormat === "LEGACY" &&
        (!definition.dependencies || definition.dependencies.length === 0)
      ) {
        return {
          dependencies: [],
          dependenciesByStatus: {},
          inputPanelsAwaitingValue: [],
        };
      }
      const dependencies: PanelDefinition[] = [];
      const dependenciesByStatus: PanelDependenciesByStatus = {};
      const inputPanelsAwaitingValue: PanelDefinition[] = [];
      const recordedInputPanels = {};

      if (dataFormat === "NODE_AND_EDGE") {
        const nodeAndEdgeProperties =
          definition.properties as NodeAndEdgeProperties;
        for (const node of nodeAndEdgeProperties.nodes || []) {
          const nodePanel = panelsMap[node];
          if (!nodePanel || !nodePanel.status) {
            continue;
          }
          recordDependency(
            nodePanel,
            panelsMap,
            inputs,
            dependencies,
            dependenciesByStatus,
            inputPanelsAwaitingValue,
            recordedInputPanels,
          );
        }
        for (const edge of nodeAndEdgeProperties.edges || []) {
          const edgePanel = panelsMap[edge];
          if (!edgePanel || !edgePanel.status) {
            continue;
          }
          recordDependency(
            edgePanel,
            panelsMap,
            inputs,
            dependencies,
            dependenciesByStatus,
            inputPanelsAwaitingValue,
            recordedInputPanels,
          );
        }
      }

      for (const dependency of definition.dependencies || []) {
        const dependencyPanel = panelsMap[dependency];
        if (!dependencyPanel || !dependencyPanel.status) {
          continue;
        }
        recordDependency(
          dependencyPanel,
          panelsMap,
          inputs,
          dependencies,
          dependenciesByStatus,
          inputPanelsAwaitingValue,
          recordedInputPanels,
        );
      }

      return { dependencies, dependenciesByStatus, inputPanelsAwaitingValue };
    }, [definition, panelsMap, inputs]);

  useEffect(() => {
    if (parentType !== "container") {
      return;
    }
    updateChildStatus(
      definition as PanelDefinition,
      inputPanelsAwaitingValue.length === 0 ? "visible" : "hidden",
    );
  }, [definition, inputPanelsAwaitingValue, parentType, updateChildStatus]);

  return (
    <PanelControlsProvider definition={definition} enabled={showControls}>
      <PanelContext.Provider
        value={{
          definition,
          dependencies,
          dependenciesByStatus,
          inputPanelsAwaitingValue,
          panelInformation,
          showPanelInformation,
          setPanelInformation,
          setShowPanelInformation,
        }}
      >
        {children}
      </PanelContext.Provider>
    </PanelControlsProvider>
  );
};

const usePanel = () => {
  const context = useContext(PanelContext);
  if (context === undefined) {
    throw new Error("usePanel must be used within a PanelContext");
  }
  return context as IPanelContext;
};

export { PanelContext, PanelProvider, usePanel };
