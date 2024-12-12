import { createContext, ReactNode, useContext } from "react";
import { IThemeContext } from "@powerpipe/types";

interface DashboardThemeProviderProps {
  children: ReactNode;
  themeContext: IThemeContext;
}

const DashboardThemeContext = createContext<IThemeContext | null>(null);

export const DashboardThemeProvider = ({
  children,
  themeContext,
}: DashboardThemeProviderProps) => {
  return (
    <DashboardThemeContext.Provider value={themeContext}>
      {children}
    </DashboardThemeContext.Provider>
  );
};

export const useDashboardTheme = () => {
  const context = useContext(DashboardThemeContext);
  if (!context) {
    throw new Error(
      "useDashboardTheme must be used within a DashboardThemeProvider",
    );
  }
  return context;
};
