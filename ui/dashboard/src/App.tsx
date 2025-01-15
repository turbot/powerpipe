import "./utils/registerComponents";
import Dashboard from "./components/dashboards/layout/Dashboard";
import DashboardHeader from "./components/DashboardHeader";
import DashboardList from "./components/DashboardList";
import useAnalytics from "./hooks/useAnalytics";
import WorkspaceErrorModal from "./components/dashboards/WorkspaceErrorModal";
import { DashboardProvider } from "./hooks/useDashboard";
import { FullHeightThemeWrapper, useTheme } from "./hooks/useTheme";
import { Route, Routes } from "react-router-dom";
import { useBreakpoint } from "./hooks/useBreakpoint";

const DashboardApp = ({
  analyticsContext,
  breakpointContext,
  themeContext,
}) => {
  const dashboards = (
    <DashboardProvider
      analyticsContext={analyticsContext}
      breakpointContext={breakpointContext}
      themeContext={themeContext}
      versionMismatchCheck={true}
    >
      <DashboardHeader />
      <WorkspaceErrorModal />
      <DashboardList wrapperClassName="p-4 h-full overflow-y-auto" />
      <Dashboard />
    </DashboardProvider>
  );

  return (
    <Routes>
      <Route path="/" element={dashboards} />
      <Route path="/snapshot/:dashboard_name" element={dashboards} />
      <Route path="/:dashboard_name" element={dashboards} />
    </Routes>
  );
};

const App = () => {
  const analyticsContext = useAnalytics();
  const breakpointContext = useBreakpoint();
  const themeContext = useTheme();

  return (
    <FullHeightThemeWrapper>
      <DashboardApp
        analyticsContext={analyticsContext}
        breakpointContext={breakpointContext}
        themeContext={themeContext}
      />
    </FullHeightThemeWrapper>
  );
};

export default App;

export { DashboardApp };
