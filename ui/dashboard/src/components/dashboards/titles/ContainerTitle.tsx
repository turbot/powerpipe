import { classNames } from "@powerpipe/utils/styles";
import { ReactNode } from "react";

interface ContainerTitleProps {
  title: string | null | undefined;
  controls?: ReactNode;
}

const ContainerTitle = ({ title, controls }: ContainerTitleProps) => {
  if (!title) {
    return null;
  }
  const titleHeading = (
    <h2 className={classNames("col-span-12 grow")}>{title}</h2>
  );

  if (!controls) {
    return titleHeading;
  }

  return (
    <div className="flex items-center justify-between">
      {titleHeading}
      {controls}
    </div>
  );
};

export default ContainerTitle;
