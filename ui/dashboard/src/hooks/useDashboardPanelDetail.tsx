import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { GlobalHotKeys } from "react-hotkeys";
import { LeafNodeData } from "@powerpipe/components/dashboards/common";
import { noop } from "@powerpipe/utils/func";
import { PanelDefinition } from "@powerpipe/types";

export interface SidePanelInfo {
  panel: PanelDefinition;
  context?: any;
}

interface IDashboardPanelDetailContext {
  selectedPanel: PanelDefinition | null;
  selectPanel: (panelName: PanelDefinition | null, data?: LeafNodeData) => void;
  closePanel: () => void;
  panelOverrideData: LeafNodeData | null;
  selectedSidePanel: SidePanelInfo | null;
  selectSidePanel: (SidePanelInfo: SidePanelInfo | null) => void;
  closeSidePanel: () => void;
}

interface DashboardPanelDetailProviderProps {
  children: ReactNode;
}

const DashboardPanelDetailContext =
  createContext<IDashboardPanelDetailContext | null>(null);

export const DashboardPanelDetailProvider = ({
  children,
}: DashboardPanelDetailProviderProps) => {
  const [selectedPanel, setSelectedPanel] = useState<PanelDefinition | null>(
    null,
  );
  const [selectedSidePanel, setSelectedSidePanel] =
    useState<SidePanelInfo | null>(null);
  const [panelOverrideData, setPanelOverrideData] =
    useState<LeafNodeData | null>(null);

  const [hotKeysHandlers, setHotKeysHandlers] = useState({
    CLOSE_PANEL_DETAIL: noop,
  });

  const hotKeysMap = {
    CLOSE_PANEL_DETAIL: ["esc"],
  };

  const selectPanel = (
    panelName: PanelDefinition | null,
    data?: LeafNodeData,
  ) => {
    setSelectedPanel(panelName);
    if (data) {
      setPanelOverrideData(data);
    }
  };

  const closePanel = useCallback(() => {
    setSelectedPanel(null);
  }, []);

  useEffect(() => {
    setHotKeysHandlers({
      CLOSE_PANEL_DETAIL: closePanel,
    });
  }, [closePanel]);

  return (
    <DashboardPanelDetailContext.Provider
      value={{
        selectedPanel,
        selectPanel,
        closePanel,
        panelOverrideData,
        selectedSidePanel,
        selectSidePanel: (sidePanelInfo: SidePanelInfo | null) =>
          setSelectedSidePanel(sidePanelInfo),
        closeSidePanel: () => setSelectedSidePanel(null),
      }}
    >
      <GlobalHotKeys
        allowChanges
        keyMap={hotKeysMap}
        handlers={hotKeysHandlers}
      />
      {children}
    </DashboardPanelDetailContext.Provider>
  );
};

export const useDashboardPanelDetail = () => {
  const context = useContext(DashboardPanelDetailContext);
  if (!context) {
    throw new Error(
      "useDashboardPanelDetail must be used within a DashboardPanelDetailContext",
    );
  }
  return context;
};
