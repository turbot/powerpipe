import { classNames } from "@powerpipe/utils/styles";
import { ReactNode } from "react";

interface DashboardTitleProps {
  title: string | null | undefined;
  controls?: ReactNode;
}

const DashboardTitle = ({ title, controls }: DashboardTitleProps) => {
  if (!title) {
    return null;
  }
  const titleHeading = <h1 className={classNames("col-span-12")}>{title}</h1>;

  if (!controls) {
    return titleHeading;
  }

  return (
    <div className="col-span-12 flex items-center justify-between">
      {titleHeading}
      {controls}
    </div>
  );
};

export default DashboardTitle;
