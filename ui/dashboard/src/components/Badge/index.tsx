import { ReactNode } from "react";
import { classNames } from "@powerpipe/utils/styles";

interface BadgeProps {
  children: ReactNode;
  type?: "info";
}

const Badge = ({ children, type = "info" }: BadgeProps) => {
  // TODO other badge types
  return (
    <span
      className={classNames(
        "text-sm px-1.5 py-0.5 rounded-md bg-opacity-20",
        type === "info" ? "bg-info text-info" : null,
      )}
    >
      {children}
    </span>
  );
};

export default Badge;
