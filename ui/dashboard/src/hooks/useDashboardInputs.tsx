import usePrefixedSearchParams from "@powerpipe/hooks/usePrefixedSearchParams";
import { createContext, ReactNode, useContext } from "react";

interface IDashboardInputsContext {
  inputs: Record<string, string>;
}

interface DashboardInputsProviderProps {
  children: ReactNode;
}

const DashboardInputsContext = createContext<IDashboardInputsContext | null>(
  null,
);

export const DashboardInputsProvider = ({
  children,
}: DashboardInputsProviderProps) => {
  const inputs = usePrefixedSearchParams("input.");

  return (
    <DashboardInputsContext.Provider value={{ inputs }}>
      {children}
    </DashboardInputsContext.Provider>
  );
};

export const useDashboardInputs = () => {
  const context = useContext(DashboardInputsContext);
  if (!context) {
    throw new Error(
      "useDashboardInputs must be used within a DashboardInputsContext",
    );
  }
  return context;
};
