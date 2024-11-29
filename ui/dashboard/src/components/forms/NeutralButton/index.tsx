import Button, { ButtonProps } from "../Button";
import { classNames } from "@powerpipe/utils/styles";
import { forwardRef } from "react";

const NeutralButton = forwardRef(
  (
    {
      children,
      className = "",
      disabled = false,
      onClick,
      size = "md",
      title,
      type,
    }: ButtonProps,
    ref,
  ) => {
    return (
      <Button
        ref={ref}
        className={classNames(
          className,
          "bg-dashboard-panel border border-black-scale-2 text-light hover:bg-black-scale-2 hover:border-black-scale-2 disabled:bg-dashboard disabled:text-light",
        )}
        disabled={disabled}
        onClick={onClick}
        size={size}
        title={title}
        type={type}
      >
        {children}
      </Button>
    );
  },
);

export default NeutralButton;
