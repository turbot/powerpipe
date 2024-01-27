import { getComponent } from "../dashboards";

const items = [
  {
    title: "Find a Dashboard",
    description:
      "Steampipe Hub has hundreds of open source dashboards to get you started.",
    href: "https://hub.powerpipe.io/mods",
    withReferrer: true,
  },
  {
    title: "Build a Dashboard",
    description:
      "It's easy to create your own dashboard as code! Start with this tutorial.",
    href: "https://powerpipe.io/docs/mods/writing-dashboards",
    withReferrer: true,
  },
  {
    title: "Join our Community",
    description:
      "Connect directly with Steampipe users and the development team in Slack.",
    href: "https://steampipe.io/community/join",
    withReferrer: true,
  },
];

const CallToActions = () => {
  const ExternalLink = getComponent("external_link");
  return (
    <ul className="mt-4 md:mt-0 space-y-6">
      {items.map((item, itemIdx) => (
        <li key={itemIdx} className="flow-root">
          <div className="p-3 flex items-center space-x-4 rounded-md hover:bg-dashboard-panel focus-within:ring-2 focus-within:ring-blue-500">
            <ExternalLink
              className="focus:outline-none"
              ignoreDataMode
              to={item.href}
              withReferrer={item.withReferrer}
            >
              <span className="text-foreground">
                <>{item.title}</>
                <span aria-hidden="true" className="ml-1">
                  &rarr;
                </span>
              </span>
              <p className="mt-1 text-sm text-foreground-light">
                {item.description}
              </p>
            </ExternalLink>
          </div>
        </li>
      ))}
    </ul>
  );
};

export default CallToActions;
