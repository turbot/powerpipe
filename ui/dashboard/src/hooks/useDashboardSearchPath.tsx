import { createContext, ReactNode, useContext, useMemo } from "react";
import { useSearchParams } from "react-router-dom";

interface IDashboardSearchPathContext {
  searchPathPrefix: string[];
}

interface DashboardSearchPathProviderProps {
  children: ReactNode;
}

const DashboardSearchPathContext =
  createContext<IDashboardSearchPathContext | null>(null);

export const DashboardSearchPathProvider = ({
  children,
}: DashboardSearchPathProviderProps) => {
  const [searchParams] = useSearchParams();
  const rawSearchPathPrefix = searchParams.get("search_path_prefix");

  const searchPathPrefix = useMemo(() => {
    if (!!rawSearchPathPrefix) {
      try {
        return rawSearchPathPrefix.split(",");
      } catch (error) {
        console.error("Error parsing search path prefix", error);
        return [];
      }
    } else {
      return [];
    }
  }, [rawSearchPathPrefix]);

  return (
    <DashboardSearchPathContext.Provider value={{ searchPathPrefix }}>
      {children}
    </DashboardSearchPathContext.Provider>
  );
};

export const useDashboardSearchPath = () => {
  const context = useContext(DashboardSearchPathContext);
  if (!context) {
    throw new Error(
      "useDashboardSearchPath must be used within a DashboardSearchPathContext",
    );
  }
  return context;
};
