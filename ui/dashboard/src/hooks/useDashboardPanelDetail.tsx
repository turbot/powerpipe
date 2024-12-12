import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { GlobalHotKeys } from "react-hotkeys";
import { noop } from "@powerpipe/utils/func";

interface IDashboardPanelDetailContext {
  selectedPanel: string | null;
  selectPanel: (panelName: string | null) => void;
  closePanel: () => void;
  selectedFilterAndGroupPanel: string | null;
  selectFilterAndGroupPanel: (panelName: string | null) => void;
  closeFilterAndGroupPanel: () => void;
}

interface DashboardPanelDetailProviderProps {
  children: ReactNode;
}

const DashboardPanelDetailContext =
  createContext<IDashboardPanelDetailContext | null>(null);

export const DashboardPanelDetailProvider = ({
  children,
}: DashboardPanelDetailProviderProps) => {
  const [selectedPanel, setSelectedPanel] = useState<string | null>(null);
  const [selectedFilterAndGroupPanel, setSelectedFilterAndGroupPanel] =
    useState<string | null>(null);

  const [hotKeysHandlers, setHotKeysHandlers] = useState({
    CLOSE_PANEL_DETAIL: noop,
  });

  const hotKeysMap = {
    CLOSE_PANEL_DETAIL: ["esc"],
  };

  const selectPanel = (panelName: string | null) => {
    setSelectedPanel(panelName);
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
        selectedFilterAndGroupPanel,
        selectFilterAndGroupPanel: (panelName: string | null) =>
          setSelectedFilterAndGroupPanel(panelName),
        closeFilterAndGroupPanel: () => setSelectedFilterAndGroupPanel(null),
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
