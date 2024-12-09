import useDownloadPanelData from "./useDownloadPanelData";
import useSelectPanel from "./useSelectPanel";
import { BaseChartProps } from "@powerpipe/components/dashboards/charts/types";
import { CardProps } from "@powerpipe/components/dashboards/Card";
import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { FlowProps } from "@powerpipe/components/dashboards/flows/types";
import { GraphProps } from "@powerpipe/components/dashboards/graphs/types";
import { HierarchyProps } from "@powerpipe/components/dashboards/hierarchies/types";
import { ImageProps } from "@powerpipe/components/dashboards/Image";
import { InputProps } from "@powerpipe/components/dashboards/inputs/types";
import { LeafNodeData } from "@powerpipe/components/dashboards/common";
import { PanelDefinition } from "@powerpipe/types";
import { TableProps } from "@powerpipe/components/dashboards/Table";
import { TextProps } from "@powerpipe/components/dashboards/Text";

export type IPanelControlsContext = {
  enabled: boolean;
  panelControls: IPanelControl[];
  showPanelControls: boolean;
  setShowPanelControls: (show: boolean) => void;
  setCustomControls: (controls: IPanelControl[]) => void;
  setPanelData: (data: LeafNodeData) => void;
};

type PanelControlsProviderProps = {
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
  enabled?: boolean;
  panelDetailEnabled?: boolean;
};

export interface IPanelControl {
  action: (e: any) => Promise<void>;
  component?: ReactNode;
  icon?: string;
  title: string;
}

const PanelControlsContext = createContext<IPanelControlsContext | null>(null);

const PanelControlsProvider = ({
  children,
  definition,
  enabled,
  panelDetailEnabled = true,
}: PanelControlsProviderProps) => {
  const [panelData, setPanelData] = useState<LeafNodeData | undefined>(
    definition.data,
  );
  const { download } = useDownloadPanelData(definition, panelData);
  const { select } = useSelectPanel(definition);
  const [showPanelControls, setShowPanelControls] = useState(false);

  useEffect(() => setPanelData(() => definition.data), [definition.data]);

  const downloadPanelData = useCallback(
    async (e) => {
      e.stopPropagation();
      await download();
    },
    [download],
  );

  const getBasePanelControls = useCallback(() => {
    const controls: IPanelControl[] = [];
    if (!enabled || !definition) {
      return controls;
    }
    if (panelData) {
      controls.push({
        action: downloadPanelData,
        icon: "arrow-down-tray",
        title: "Download data",
      });
    }
    if (panelDetailEnabled) {
      controls.push({
        action: select,
        icon: "arrows-pointing-out",
        title: "View detail",
      });
    }
    return controls;
  }, [definition, downloadPanelData, panelDetailEnabled, select, enabled]);

  const [panelControls, setPanelControls] = useState(getBasePanelControls());
  const [customControls, setCustomControls] = useState<IPanelControl[]>([]);

  useEffect(() => {
    setPanelControls(() => [...customControls, ...getBasePanelControls()]);
  }, [customControls, getBasePanelControls]);

  return (
    <PanelControlsContext.Provider
      value={{
        enabled: enabled || false,
        panelControls,
        showPanelControls,
        setShowPanelControls,
        setCustomControls,
        setPanelData,
      }}
    >
      {children}
    </PanelControlsContext.Provider>
  );
};

const usePanelControls = () => {
  const context = useContext(PanelControlsContext);
  if (context === undefined) {
    throw new Error(
      "usePanelControls must be used within a PanelControlsContext",
    );
  }
  return context as IPanelControlsContext;
};

export { PanelControlsProvider, usePanelControls };
