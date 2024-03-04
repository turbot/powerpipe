import DashboardSearch from "@powerpipe/components/DashboardSearch";
import DashboardTagGroupSelect from "@powerpipe/components/DashboardTagGroupSelect";
import ManageSearchPathButton from "@powerpipe/components/ManageSearchPathButton";
import OpenSnapshotButton from "@powerpipe/components/OpenSnapshotButton";
import PowerpipeLogo from "@powerpipe/components/DashboardHeader/PowerpipeLogo";
import SaveSnapshotButton from "@powerpipe/components/SaveSnapshotButton";
import ThemeToggle from "@powerpipe/components/ThemeToggle";
import { classNames } from "@powerpipe/utils/styles";
import { getComponent } from "@powerpipe/components/dashboards";

const DashboardHeader = () => {
  const ExternalLink = getComponent("external_link");
  return (
    <>
      <div
        className={classNames(
          "flex w-screen px-4 py-3 items-center justify-between space-x-2 md:space-x-4 bg-dashboard-panel border-b border-divide print:hidden",
        )}
      >
        <PowerpipeLogo />
        <div className="flex flex-grow items-center space-x-2 md:space-x-4">
          <DashboardSearch />
          <ManageSearchPathButton />
          <DashboardTagGroupSelect />
          <SaveSnapshotButton />
          <OpenSnapshotButton />
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
