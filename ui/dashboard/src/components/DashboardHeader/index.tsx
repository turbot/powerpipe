import ManageDatetimeRangeButton from "@powerpipe/components/dashboards/DatetimeRange/ManageDatetimeRangeButton";
import ManageSearchPathButton from "@powerpipe/components/dashboards/SearchPath/ManageSearchPathButton";
import PowerpipeLogo from "@powerpipe/components/DashboardHeader/PowerpipeLogo";
import SplitSnapshotButton from "@powerpipe/components/SplitSnapshotButton";
import ThemeToggle from "@powerpipe/components/ThemeToggle";
import { classNames } from "@powerpipe/utils/styles";
import { getComponent } from "@powerpipe/components/dashboards";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

const DashboardHeader = () => {
  const { selectedDashboard } = useDashboardState();
  const ExternalLink = getComponent("external_link");

  return (
    <>
      <div
        className={classNames(
          "flex w-screen px-4 py-3 items-center justify-between space-x-2 md:space-x-4 bg-dashboard-panel border-b border-divide print:hidden",
        )}
      >
        <PowerpipeLogo />
        <div
          className={classNames(
            "flex flex-grow items-center space-x-2 md:space-x-4",
            // Maintain height between dashboard list and dashboard
            selectedDashboard ? "" : "my-[1.75px]",
          )}
        >
          <ManageDatetimeRangeButton />
          <ManageSearchPathButton />
          <SplitSnapshotButton header />
        </div>
        <div className="space-x-2 sm:space-x-4 md:space-x-8 flex items-center justify-end">
          <ExternalLink
            className="text-base text-foreground-lighter hover:text-foreground"
            ignoreDataMode
            to="https://hub.powerpipe.io"
            withReferrer={true}
          >
            <>Hub</>
          </ExternalLink>
          <ExternalLink
            className="text-base text-foreground-lighter hover:text-foreground"
            ignoreDataMode
            to="https://powerpipe.io/docs"
            withReferrer={true}
          >
            <>Docs</>
          </ExternalLink>
          <ThemeToggle />
        </div>
      </div>
    </>
  );
};

export default DashboardHeader;
