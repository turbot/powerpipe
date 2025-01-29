import useDownloadPanelData from "./useDownloadPanelData";
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
import { noop } from "@powerpipe/utils/func";
import { PanelDefinition } from "@powerpipe/types";
import { TableProps } from "@powerpipe/components/dashboards/Table";
import { TextProps } from "@powerpipe/components/dashboards/Text";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";

export type IPanelControlsContext = {
  enabled: boolean;
  panelControls: IPanelControl[];
  customControls: IPanelControl[];
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
  disabled?: boolean;
  key: string;
  action?: (e: any) => Promise<void>;
  component?: ReactNode;
  icon?: string;
  title: string;
}

const PanelControlsContext = createContext<IPanelControlsContext | null>({
  enabled: false,
  panelControls: [],
  customControls: [],
  showPanelControls: false,
  setCustomControls: noop,
  setPanelData: noop,
  setShowPanelControls: noop,
});

const PanelControlsProvider = ({
  children,
  definition,
  enabled,
  panelDetailEnabled = true,
}: PanelControlsProviderProps) => {
  const [panelData, setPanelData] = useState<LeafNodeData | undefined>(
    definition.data,
  );
  const { download } = useDownloadPanelData(definition);
  const { selectPanel } = useDashboardPanelDetail();
  const [showPanelControls, setShowPanelControls] = useState(false);

  useEffect(() => setPanelData(() => definition.data), [definition.data]);

  const downloadPanelData = useCallback(
    async (e) => {
      e.stopPropagation();
      await download(panelData);
    },
    [definition, download, panelData],
  );

  const getBasePanelControls = () => {
    const controls: IPanelControl[] = [];
    if (!enabled || !definition) {
      return controls;
    }
    controls.push({
      key: "download-data",
      disabled: !panelData,
      title: panelData ? "Download data" : "No data to download",
      icon: "arrow-down-tray",
      action: downloadPanelData,
    });
    if (panelDetailEnabled) {
      controls.push({
        key: "view-panel-detail",
        title: "View detail",
        icon: "arrows-pointing-out",
        action: async () => selectPanel(definition, panelData),
      });
    }
    return controls;
  };

  const [panelControls, setPanelControls] = useState(getBasePanelControls());
  const [customControls, setCustomControls] = useState<IPanelControl[]>([]);

  useEffect(() => {
    const uniqueCustomControls: IPanelControl[] = [];
    let baseControls = getBasePanelControls();
    for (const control of customControls) {
      const existingIndex = baseControls.findIndex(
        (c) => c.key === control.key,
      );
      if (existingIndex === -1) {
        uniqueCustomControls.push(control);
      } else {
        baseControls = [
          ...baseControls.slice(0, existingIndex),
          control,
          ...baseControls.slice(existingIndex + 1),
        ];
      }
    }
    setPanelControls(() => [...uniqueCustomControls, ...baseControls]);
  }, [customControls, panelData]);

  return (
    <PanelControlsContext.Provider
      value={{
        enabled: enabled || false,
        panelControls,
        customControls,
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
