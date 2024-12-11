import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { DashboardActions } from "@powerpipe/types";
import { GlobalHotKeys } from "react-hotkeys";
import { noop } from "@powerpipe/utils/func";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

interface IDashboardPanelDetailContext {
  closePanelDetail: () => void;
}

interface DashboardPanelDetailProviderProps {
  children: ReactNode;
}

const DashboardPanelDetailContext =
  createContext<IDashboardPanelDetailContext | null>(null);

export const DashboardPanelDetailProvider = ({
  children,
}: DashboardPanelDetailProviderProps) => {
  const { dispatch } = useDashboardState();

  const [hotKeysHandlers, setHotKeysHandlers] = useState({
    CLOSE_PANEL_DETAIL: noop,
  });

  const hotKeysMap = {
    CLOSE_PANEL_DETAIL: ["esc"],
  };

  const closePanelDetail = useCallback(() => {
    dispatch({
      type: DashboardActions.SELECT_PANEL,
      panel: null,
    });
  }, [dispatch]);

  useEffect(() => {
    setHotKeysHandlers({
      CLOSE_PANEL_DETAIL: closePanelDetail,
    });
  }, [closePanelDetail]);

  return (
    <DashboardPanelDetailContext.Provider
      value={{
        closePanelDetail,
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
